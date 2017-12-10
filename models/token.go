package models

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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
	var privateKey *rsa.PrivateKey
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
	return token.SignedString(privateKey)
}

//GetPublicPem gets the public for the passed in private key.
//Returns the bytes of the key in pem format.
func GetPublicPem() ([]byte, error) {
	var privateKey *rsa.PrivateKey
	//We need to get the private key from S3!

	privateKey.Public()

	pubDer, err := x509.MarshalPKIXPublicKey(privateKey.Public())
	if err != nil {
		return nil, err
	}
	pubPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubDer,
	})
	return pubPem, nil
}
