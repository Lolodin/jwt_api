package controller

import (
	"encoding/json"
	"fmt"
	"github.com/Lolodin/jwt_api/store"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const SECRET = "TESTTESTTESTTEST"
type erroranswer struct {
	Error string `json:"error"`
}

type UserData struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Password string `json:"password"`
}
type Refresh struct {
	Ref string `json:"ref"`

}
type Token struct {
	Exp int64 `json:"exp"`
	Name string `json:"name"`
	UUID string `json:"uuid"`
}
//Получение токена доступа и обновления
func GetTokens(s *store.MongoStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GetTokens")
		user := UserData{}
		buff, err :=ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
		}
		json.Unmarshal(buff, &user)
		fmt.Println(user)


		//check pass
		DBuser:=s.GetUser(user.Name)
		if DBuser == nil {
			w.Write([]byte("Wrong User"))
			return
		}

		res := bcrypt.CompareHashAndPassword(DBuser.Password,[]byte(user.Password))
		if res != nil {
			http.Error(w, "User non exist", 404 )
			return
		}
		//checkSession uuid
		s.DeleteSession(user.UUID)
		//Access Token

		m:= GenerateTokens(user.UUID, user.Name, s)
		jst, err:=json.Marshal(m)
		if err != nil {
			log.Print(err)
		}
		w.Write(jst)
	}


}
//Страница для теста запросов
func Index()  func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("indexAction")
		t, _ := template.ParseFiles("html/index.html")
		err := t.Execute(w, "index")
		if err != nil {
			fmt.Println(err.Error())
		}
	}


}
//Страница для регистрации юзера
func Reg()  func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("regAction")
		t, _ := template.ParseFiles("html/reg.html")
		err := t.Execute(w, "reg")
		if err != nil {
			fmt.Println(err.Error())
		}
	}


}
//Маршрут для регистрации
func Register(s *store.MongoStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Register")
		user := UserData{}
		buff, err :=ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
		}
		json.Unmarshal(buff, &user)
		fmt.Println(user)
		s.AddUser(user.Name, user.Password)

	}


}

//Обновляем Токены
func RefreshTokens(s *store.MongoStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("RefTokens")
		ref := Refresh{}
		buff, err :=ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
		}
		json.Unmarshal(buff, &ref)
		fmt.Println(ref)

		//RefreshToken сравниваем с токеном из БД

		token, err:=jwt.Parse(ref.Ref, func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET), nil
		})
		if err != nil {
			fmt.Println(err)
		}
		if token.Valid {
			claims:=token.Claims.(jwt.MapClaims)
			uuid := claims["uuid"].(string)
			exp  := claims["exp"].(int64)
			// Получаем сессию из БД и удаляем в любом случае

			session:=s.GetSession(uuid)
			s.DeleteSession(uuid)
			check:=bcrypt.CompareHashAndPassword(session.RefToken, []byte(ref.Ref))
			if check != nil {
				err, _:=json.Marshal(erroranswer{Error: check.Error()})
				w.Write(err)
				return
			}


			if exp<= time.Now().Unix() {
				err, _:=json.Marshal(erroranswer{Error: "refresh token not valid"})
				w.Write(err)
				return
			}
			newtoken:=GenerateTokens(session.UUID, session.Name, s)
			

		}

	}


}