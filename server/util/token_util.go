package util

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/elisha1995/magic-stream-movies/server/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Role      string
	UserId    string
	jwt.RegisteredClaims
}

var SecretKey = os.Getenv("SECRET_KEY")
var SecretRefreshKey = os.Getenv("SECRET_REFRESH_KEY")

// GenerateAllTokens generates a pair of JSON Web Tokens (access token and refresh token) given user details.
//
// The access token is a short-lived token that is used to authenticate and authorize
// requests to the MagicStream API. The access token is valid for 24 hours.
//
// The refresh token is a long-lived token that is used to obtain a new access token
// when the current access token expires. The refresh token is valid for 7 days.
//
// The function takes in the user's email, first name, last name, role, and user ID, and
// returns the signed access token and refresh token as strings. If there is an error during
// token generation, the function returns an empty string and an error.
func GenerateAllTokens(email, firstName, lastName, role, userId string) (string, string, error) {
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "MagicStream",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SecretKey))

	if err != nil {
		return "", "", err
	}

	refreshClaims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "MagicStream",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(SecretRefreshKey))

	if err != nil {
		return "", "", err
	}

	return signedToken, signedRefreshToken, nil
}

// UpdateAllTokens updates the access token and refresh token for a user in the database.
//
// The function takes in the user's ID, the new access token, and the new refresh token, and
// returns an error if there is an issue updating the tokens. If the update is successful, the
// function returns nil.
//
// The function uses the WithTimeout function from the context package to set a timeout of 100
// seconds for the update operation. If the update operation takes longer than 100 seconds to
// complete, the function will return a context.DeadlineExceeded error.
//
// The function uses the UpdateOne function from the mongo-driver package to update the user's document in
// the database. The function sets the "token", "refresh_token", and "update_at" fields in the
// user's document to the new access token, new refresh token, and the current time, respectively.
//
// If there is an error during the update operation, the function returns the error. Otherwise, the
// function returns nil.
func UpdateAllTokens(userId, token, refreshToken string, client *mongo.Client) (err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	updateAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updateData := bson.M{
		"$set": bson.M{
			"token":         token,
			"refresh_token": refreshToken,
			"update_at":     updateAt,
		},
	}

	var userCollection = database.OpenCollection("users", client)

	_, err = userCollection.UpdateOne(ctx, bson.M{"user_id": userId}, updateData)

	if err != nil {
		return err
	}
	return nil
}

func GetAccessToken(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}
	tokenString := authHeader[len("Bearer "):]

	if tokenString == "" {
		return "", errors.New("bearer token is required")
	}

	return tokenString, nil

}

func ValidateToken(tokenString string) (*SignedDetails, error) {
	claims := &SignedDetails{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, err
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	return claims, nil

}

// GetUserIdFromContext retrieves the user ID from the Gin context.
//
// The function returns the user ID as a string if it is found in the context.
// If the user ID is not found in the context, the function returns an empty string
// and an error indicating that the user ID does not exist in the context.
//
// If the user ID found in the context is not a string, the function returns an empty string
// and an error indicating that the user ID is not a string.
func GetUserIdFromContext(c *gin.Context) (string, error) {
	userId, exists := c.Get("userId")

	if !exists {
		return "", errors.New("userId does not exists in this context")
	}

	id, ok := userId.(string)

	if !ok {
		return "", errors.New("unable to retrieve userId")
	}

	return id, nil

}

func GetRoleFromContext(c *gin.Context) (string, error) {
	role, exists := c.Get("role")

	if !exists {
		return "", errors.New("role does not exists in this context")
	}

	memberRole, ok := role.(string)

	if !ok {
		return "", errors.New("unable to retrieve userId")
	}

	return memberRole, nil

}

func ValidateRefreshToken(tokenString string) (*SignedDetails, error) {
	claims := &SignedDetails{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

		return []byte(SecretRefreshKey), nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, err
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("refresh token has expired")
	}

	return claims, nil
}
