package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"net/url"
	"time"
)

const SECRET = "TEST"
//Получение токена доступа и обновления
func GetTokens(w http.ResponseWriter, r *http.Request) {
	param:=r.URL.RawQuery
	id, err:=url.ParseQuery(param)
	if err != nil {
		log.Print("Ошибка получения параметров запроса")
	}
	uuid:= id["uuid"]
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["uuid"] = uuid
	claims["exp"] = time.Now().Add(time.Hour * 8).Unix()
    t, err := token.SignedString([]byte(SECRET))
    if err != nil {
    	log.Print("Ошибка подписи токена")
	}
	m := make(map[string]string)
	m["token"] = t
	jst, err:=json.Marshal(m)
	if err != nil {
		log.Print(err)
	}
	w.Write(jst)

}
