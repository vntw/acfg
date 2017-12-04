package user

import (
	"errors"
	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	Username string
	Password string
}

var tokenSecret []byte

// TODO: Find better way
func SetTokenSecret(s []byte) {
	tokenSecret = s
}

var users []*User

func addUser(username, password string) *User {
	u := &User{username, password}
	users = append(users, u)
	return u
}

func AddConfigUsers(users map[string]string) {
	for username, password := range users {
		addUser(username, password)
	}
}

func MatchUser(username string, password string) (*User, error) {
	u := findUser(username)

	if u == nil {
		return nil, errors.New(fmt.Sprintf("user %s not found", username))
	}
	if u.Password != password {
		return nil, errors.New(fmt.Sprintf("got invalid password for user %s", username))
	}

	return u, nil
}

func CreateToken(u *User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"usr": u.Username,
	})

	tokenString, err := token.SignedString(tokenSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func AuthToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return tokenSecret, nil
	})

	if err != nil {
		log.Println(err)
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if usr, ok := claims["usr"]; ok {
			return findUser(usr.(string)) != nil
		}
	}

	return false
}

func findUser(username string) *User {
	for _, u := range users {
		if u.Username == username {
			return u
		}
	}

	return nil
}
