package util

import (
	"context"
	"os"
	"time"

	"github.com/elisha1995/magic-stream-movies/server/database"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
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
func UpdateAllTokens(userId, token, refreshToken string) (err error) {
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

	var userCollection = database.OpenCollection("users")

	_, err = userCollection.UpdateOne(ctx, bson.M{"user_id": userId}, updateData)

	if err != nil {
		return err
	}
	return nil
}
