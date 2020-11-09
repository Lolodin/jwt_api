package controller

import (
	"encoding/base64"
	"fmt"
	"github.com/Lolodin/jwt_api/store"
	"github.com/dgrijalva/jwt-go"
	"log"

	"golang.org/x/crypto/bcrypt"
	"time"
)

//Возвращает мап с двумя токенами, делаем маршал в жсон и отправляем клиенту
func GenerateTokens(uuid, name string, s *store.MongoStore) map[string]string {
	//Access Token
	a, err := GenerateAccessToken(name, uuid)
	if err != nil {
		log.Println(err)
	}
	//Refresh
	r, bref, exp, err := GenerateRefreshToken(uuid)
	if err != nil {
		log.Println(err)
	}
	// Запись сессии в БД
	if err:=s.AddSession(name, uuid, bref, exp); err != nil{
		return nil
	}
	m := make(map[string]string)
	m["token"] = a
	m["refresh_token"] = r
	return m

}

// return JWT
func GenerateAccessToken(name, uuid string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = name
	claims["uuid"] = uuid
	claims["exp"] = time.Now().Add(time.Hour * 8).Unix()
	t, err := token.SignedString([]byte(SECRET))
	if err != nil {
		return "", err
	}
	return t, nil
}

//return refresh token, bryptToken, error
func GenerateRefreshToken(uuid string) (string, []byte, int64, error) {
	expRef := time.Now().Add(time.Hour * 48).Unix()
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["uuid"] = uuid
	claims["exp"] = fmt.Sprint(expRef)
	t, err := token.SignedString([]byte(SECRET))
	if err != nil {
		return "", nil, 0, err
	}
	bhash, err := bcrypt.GenerateFromPassword([]byte(t), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, 0, err
	}
	t = EncodeBase64(t)
	return t, bhash, expRef, nil
}

func EncodeBase64(token string) string {
	endcod := base64.RawStdEncoding
	return endcod.EncodeToString([]byte(token))
}
func DecodeBase64(token string) string {
	decode := base64.RawStdEncoding
	b, e := decode.DecodeString(token)
	if e != nil {
		return ""
	}
	return string(b)
}
