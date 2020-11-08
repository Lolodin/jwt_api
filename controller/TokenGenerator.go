package controller

import (
	"github.com/Lolodin/jwt_api/store"
	"github.com/dgrijalva/jwt-go"
	"log"

	"golang.org/x/crypto/bcrypt"
	"time"
)

//Возвращает мап с двумя токенами, делаем маршал в жсон и отправляем клиенту
func GenerateTokens(uuid, name string, s *store.MongoStore) map[string]string {
	//Access Token
	a, err := GenerateAccessToken(name,uuid)
	if err != nil {
		log.Println(err)
	}
	//Refresh
	r, bref,exp, err:=GenerateRefreshToken(uuid)
	if err != nil {
		log.Println(err)
	}
	// Запись сессии в БД
	s.AddSession(name, uuid,bref, exp)
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
func GenerateRefreshToken(uuid string) (string, []byte, int64,error) {
	expRef := time.Now().Add(time.Hour * 48).Unix()
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["uuid"] = uuid
	claims["exp"] = expRef
	t, err := token.SignedString([]byte(SECRET))
	if err != nil {
		return "",nil,0, err
	}
	bhash, err :=bcrypt.GenerateFromPassword([]byte(t), bcrypt.DefaultCost)
	if err != nil {
		return "",nil,0, err
	}
	return t,bhash,expRef, nil
}
