package cinc

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	jwt "github.com/dgrijalva/jwt-go"
	gin "github.com/gin-gonic/gin"
	"github.com/jasonmoore30/CinC-API/models"
	//uuid "github.com/satori/go.uuid"
)

//ErrorRes is when something in the API goes wrong
type ErrorRes struct {
	Message string `json:"message"`
}

//ErrorsRes is for when many errors can be returned
type ErrorsRes struct {
	Errors []ErrorRes `json:"errors,omitempty"`
}

//RegisterUserReq is a request from a client to create a new user
type RegisterUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Confirm  string `json:"confirm"`
	//SiteID   string `json:"siteId"`
}

//AuthReq is a request to authenticate a user
type AuthReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	//SiteID   string `json:"siteId"`
}

//AuthRes is a response from the API for authenticating a user
type AuthRes struct {
	Token string `json:"token"`
}

//RegisterUser is an HTTP Router Handle for registering new users
// nolint: gocyclo
func registerUser(c *gin.Context) {

	var req RegisterUserReq

	var errs []ErrorRes
	if req.Email == "" {
		errs = append(errs, ErrorRes{Message: "Must include an email address"})
	}

	if req.Password == "" {
		errs = append(errs, ErrorRes{Message: "Must supply a password"})
	}
	if req.Password != req.Confirm {
		errs = append(errs, ErrorRes{Message: "password and confirm must match"})
	}
	if len(errs) > 0 {
		//writeResponse(w, http.StatusBadRequest, &ErrorsRes{Errors: errs})
		c.JSON(http.StatusBadRequest, &ErrorsRes{Errors: errs})
		return
	}
	_, err := models.FindUser(req.Email) //models.findUser(email string)
	if err != nil {
		if err == models.ErrNoUserFound { //might need to remove this check
			var cryptPass []byte
			cryptPass, err = bcrypt.GenerateFromPassword([]byte(req.Password), 10)
			if err != nil {
				//writeResponse(w, http.StatusInternalServerError, &ErrorsRes{Errors: []ErrorRes{ErrorRes{Message: err.Error()}}})
				c.JSON(http.StatusInternalServerError, &ErrorsRes{Errors: []ErrorRes{ErrorRes{Message: err.Error()}}})
				return
			}
			err = models.AddUser(&models.User{ //models.addUser
				Email:    req.Email,
				Password: base64.StdEncoding.EncodeToString(cryptPass),
				//Sites:    []string{req.SiteID},
			})
			if err != nil {
				//writeResponse(w, http.StatusInternalServerError, &ErrorsRes{Errors: []ErrorRes{ErrorRes{Message: err.Error()}}})
				c.JSON(http.StatusInternalServerError, &ErrorsRes{Errors: []ErrorRes{ErrorRes{Message: err.Error()}}})
				return
			}
			//w.WriteHeader(http.StatusCreated)
			c.AbortWithStatus(http.StatusCreated)
			return
		}
		c.JSON(http.StatusInternalServerError, &ErrorsRes{Errors: []ErrorRes{ErrorRes{Message: err.Error()}}})
		return
	}
	//writeResponse(w, http.StatusBadRequest, &ErrorsRes{Errors: []ErrorRes{ErrorRes{Message: "Email is already registered"}}})
	c.JSON(http.StatusBadRequest, &ErrorsRes{Errors: []ErrorRes{ErrorRes{Message: "Email is already registered"}}})
}

//AuthUser is an HTTP Router Handle for Authentication new users and return tokens
// nolint: gocyclo
func authUser(c *gin.Context) {
	var req AuthReq
	err := c.Bind(&req)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, &ErrorsRes{Errors: []ErrorRes{ErrorRes{Message: "Request must be in JSON format"}}})
		return
	}
	var errs []ErrorRes
	if req.Email == "" {
		errs = append(errs, ErrorRes{Message: "Must include an email address"})
	}
	if req.Password == "" {
		errs = append(errs, ErrorRes{Message: "Must supply a password"})
	}

	if len(errs) > 0 {
		c.JSON(http.StatusBadRequest, &ErrorsRes{Errors: errs})
		return
	}

	user, err := models.FindUser(req.Email) //models.findUser()
	if err != nil {
		if err == models.ErrNoUserFound {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.JSON(http.StatusInternalServerError, &ErrorsRes{Errors: []ErrorRes{ErrorRes{Message: err.Error()}}})
		return
	}
	cryptPass, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorsRes{Errors: []ErrorRes{ErrorRes{Message: err.Error()}}})
		return
	}
	//TODO: find a better way to check as this is taking about 2 seconds to check the passwords
	err = bcrypt.CompareHashAndPassword(cryptPass, []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}
	token, err := models.GenerateUserToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorsRes{Errors: []ErrorRes{ErrorRes{Message: err.Error()}}})
		return
	}
	c.JSON(http.StatusOK, &AuthRes{
		Token: token,
	})
}

//PublicKey is a way to get the public key to verify tokens
func publicKey(c *gin.Context) {
	//LETS HOPE THAT THIS WILL WORK AS IS
	pem, err := models.GetPublicPem()
	if err != nil {
		c.JSON(http.StatusInternalServerError, &ErrorsRes{Errors: []ErrorRes{ErrorRes{Message: err.Error()}}})
		return
	}
	c.Header("Content-Type", "application/x-pem-file")
	//w.Header().Set("Content-Type", "application/x-pem-file")
	c.Writer.WriteHeader(http.StatusOK)
	_, err = c.Writer.Write(pem)
	if err != nil {
		fmt.Println("Could not write to response: ", err)
	}

}

//AuthToken Middleware
func AuthToken(c *gin.Context) {
	var claims models.CustomClaims
	var parse jwt.Parser
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	authPieces := strings.Split(authHeader, " ")
	var rawToken string
	if authPieces[0] != "Bearer" {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	} else if authPieces[0] == "Bearer" {
		rawToken = authPieces[1]
	}

	token, error := parse.ParseWithClaims(rawToken, &claims, func(_ *jwt.Token) (interface{}, error) {

		pubPem, err := models.GetPublicPem()
		if err != nil {
			return nil, err
		}
		pubBlock, _ := pem.Decode(pubPem)
		return x509.ParsePKIXPublicKey(pubBlock.Bytes)
	})

	if error != nil {
		log.Println(error)
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	if token.Valid {
		c.Next()
	} else {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
}
