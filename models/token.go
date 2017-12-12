package models

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

//CustomClaims are the the claims that will be in the JWT token returned to the authenticated user
type CustomClaims struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
	jwt.StandardClaims
}

//GenerateUserToken ..
func GenerateUserToken(user *User) (string, error) {

	//NEED TO GET PRIVATE KEY FROM S3 Storage
	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("us-east-1")}))
	svc := s3.New(sess)

	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:    aws.String("id_rsa"),
	})
	defer result.Body.Close()

	var privKey = make([]byte, 1675)

	n, err := result.Body.Read(privKey)
	//fmt.Println("Raw Bytes: ", privKey)
	if err != nil {
		log.Println("Error reading the bytes from S3. Read ", n, " bytes before failure")
	}

	pemBlock, rest := pem.Decode(privKey)
	// if pemBlock == nil || pemBlock.Type != "PRIVATE KEY" {
	if pemBlock == nil {
		log.Fatal("failed to decode PEM block containing private key")
	}
	//fmt.Println("Making it here 2!")

	//THIS IS NEW!
	myKey := new(rsa.PrivateKey)
	myKey, err = x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if myKey == nil {
		myKey, err = x509.ParsePKCS1PrivateKey(rest)
		if err != nil {
			return "", err
		}
		fmt.Println("Inner key return.")
	}

	//We have private key from S3 Storage now

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &CustomClaims{
		Email: user.Email,
		ID:    user.ID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(36 * time.Hour).Unix(), //TODO: Come up with an expiration time
			Issuer:    "CinC",
		},
	})

	if myKey == nil {
		fmt.Println("Token being signed with a nil private key")
	}
	completeToken, err := token.SignedString(myKey)
	return completeToken, nil
}

//GetPublicPem gets the public for the passed in private key.
//Returns the bytes of the key in pem format.
func GetPublicPem() ([]byte, error) {

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Could not load env file. Defaulting to config vars?")
		//return nil, err
	}

	mySession, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	if err != nil {
		log.Println("Couldn't connect to the AWS platform")
		return nil, err
	}
	svc := s3.New(mySession)

	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:    aws.String("id_rsa"),
	})
	defer result.Body.Close()

	var privKey = make([]byte, 1675)

	_, err = result.Body.Read(privKey)
	//fmt.Println("Raw Bytes: ", privKey)
	if err != nil {
		//log.Println("Error reading the bytes from S3 in GetPublicPem. Read ", n, " bytes before failure")
	}

	pemBlock, rest := pem.Decode(privKey)
	// if pemBlock == nil || pemBlock.Type != "PRIVATE KEY" {
	if pemBlock == nil {
		log.Fatal("failed to decode PEM block containing private key")
	}
	//fmt.Println("Making it here 2!")

	var myKey = new(rsa.PrivateKey)
	myKey, err = x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
	if myKey == nil {
		myKey, err = x509.ParsePKCS1PrivateKey(rest)
		if err != nil {
			//return nil, err
		}
		fmt.Println("Inner key return.")
		//return myKey, nil
	}

	//now we should have the private key!
	//Pull out the public key

	pubDer, err := x509.MarshalPKIXPublicKey(myKey.Public())
	if err != nil {
		return nil, err
	}
	pubPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubDer,
	})
	return pubPem, nil
}
