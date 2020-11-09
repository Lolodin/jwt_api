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
	"strconv"
	"time"
)

const SECRET = "TEST"

type status struct {
	Status string `json:"status"`
}

type UserData struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
type Refresh struct {
	Ref string `json:"ref"`
}
type Token struct {
	Exp  int64  `json:"exp"`
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

//Получение токена доступа и обновления
func Login(s *store.MongoStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Login")
		user := UserData{}
		buff, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
		}
		json.Unmarshal(buff, &user)
		fmt.Println(user)

		//check pass
		DBuser := s.GetUser(user.Name)
		if DBuser == nil {

			err, _ := json.Marshal(status{Status: "user not found"})
			w.Write(err)
			return
		}

		res := bcrypt.CompareHashAndPassword(DBuser.Password, []byte(user.Password))
		if res != nil {

			err, _ := json.Marshal(status{Status: "Password not correct"})
			w.Write(err)
			return
		}
		//checkSession uuid
		s.DeleteSession(user.UUID)
		//Access Token

		m := GenerateTokens(user.UUID, user.Name, s)
		if m == nil {
			json.Marshal(status{Status: ""})
		}
		m["status"] = "Tokens send"
		jst, err := json.Marshal(m)
		if err != nil {
			log.Print(err)
		}
		w.Write(jst)
		log.Println("TOKENS SEND USER")
	}

}

//Страница для теста запросов
func Index() func(w http.ResponseWriter, r *http.Request) {
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
func Reg() func(w http.ResponseWriter, r *http.Request) {
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
		buff, err := ioutil.ReadAll(r.Body)
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
		buff, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
		}
		json.Unmarshal(buff, &ref)
		ref.Ref = DecodeBase64(ref.Ref)
		fmt.Println(ref)

		//RefreshToken сравниваем с токеном из БД

		token, err := jwt.Parse(ref.Ref, func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET), nil
		})
		if err != nil {
			fmt.Println(err)
		}
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			uuid := claims["uuid"].(string)
			exp := claims["exp"].(string)
			expint, err := strconv.Atoi(exp)
			if err != nil {
				err, _ := json.Marshal(status{Status: "refresh token not valid"})
				w.Write(err)
				return
			}

			// Получаем сессию из БД и удаляем в любом случае
			fmt.Println(uuid)
			session := s.GetSession(uuid)
			s.DeleteSession(uuid)
			if session == nil {
				err, _ := json.Marshal(status{Status: "refresh token not valid"})
				w.Write(err)
				return
			}
			check := bcrypt.CompareHashAndPassword(session.RefToken, []byte(ref.Ref))
			if check != nil {

				err, _ := json.Marshal(status{Status: check.Error()})
				w.Write(err)
				return
			}

			if int64(expint) <= time.Now().Unix() {

				err, _ := json.Marshal(status{Status: "refresh token not valid"})
				w.Write(err)
				return
			}
			newtoken := GenerateTokens(session.UUID, session.Name, s)
			if newtoken == nil {
				err, _ := json.Marshal(status{Status: "Error created token"})
				w.Write(err)
				return
			}
			newtoken["status"] = "Tokens updated"
			js, err := json.Marshal(newtoken)
			if err != nil {

				err, _ := json.Marshal(status{Status: "Error created token"})
				w.Write(err)
				return
			}
			w.Write(js)
			log.Println("Tokens updated")

		}

	}

}

func LogOut(s *store.MongoStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ref := Refresh{}
		buff, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
		}
		json.Unmarshal(buff, &ref)
		ref.Ref = DecodeBase64(ref.Ref)
		fmt.Println(ref)

		//RefreshToken сравниваем с токеном из БД

		token, err := jwt.Parse(ref.Ref, func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET), nil
		})
		if err != nil {
			fmt.Println(err)
		}
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			uuid := claims["uuid"].(string)
			exp := claims["exp"].(string)
			expint, err := strconv.Atoi(exp)
			if err != nil {
				err, _ := json.Marshal(status{Status: "refresh token not valid"})
				w.Write(err)
				return
			}

			// Получаем сессию из БД и удаляем в любом случае
			fmt.Println(uuid)
			session := s.GetSession(uuid)
			s.DeleteSession(uuid)
			if session == nil {
				err, _ := json.Marshal(status{Status: "refresh token not valid"})
				w.Write(err)
				return
			}
			check := bcrypt.CompareHashAndPassword(session.RefToken, []byte(ref.Ref))
			if check != nil {

				err, _ := json.Marshal(status{Status: check.Error()})
				w.Write(err)
				return
			}

			if int64(expint) <= time.Now().Unix() {

				err, _ := json.Marshal(status{Status: "refresh token not valid"})
				w.Write(err)
				return
			}

			answer, _ := json.Marshal(status{Status: "Refresh Token Delete"})
			w.Write(answer)
			log.Println("TOKEN DELETE")

		}

	}

}
func LogOutAll(s *store.MongoStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ref := Refresh{}
		buff, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
		}
		json.Unmarshal(buff, &ref)
		ref.Ref = DecodeBase64(ref.Ref)
		fmt.Println(ref)

		//RefreshToken сравниваем с токеном из БД

		token, err := jwt.Parse(ref.Ref, func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET), nil
		})
		if err != nil {
			fmt.Println(err)
		}
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			uuid := claims["uuid"].(string)
			exp := claims["exp"].(string)
			expint, err := strconv.Atoi(exp)
			if err != nil {
				err, _ := json.Marshal(status{Status: "refresh token not valid"})
				w.Write(err)
				return
			}

			// Получаем сессию из БД и удаляем в любом случае
			fmt.Println(uuid)
			session := s.GetSession(uuid)
			s.DeleteSession(uuid)
			if session == nil {
				err, _ := json.Marshal(status{Status: "refresh token not valid"})
				w.Write(err)
				return
			}
			check := bcrypt.CompareHashAndPassword(session.RefToken, []byte(ref.Ref))
			if check != nil {

				err, _ := json.Marshal(status{Status: check.Error()})
				w.Write(err)
				return
			}

			if int64(expint) <= time.Now().Unix() {

				err, _ := json.Marshal(status{Status: "refresh token not valid"})
				w.Write(err)
				return
			}
			s.DeleteAllSessions(session.Name)
			answer, _ := json.Marshal(status{Status: "Delete All Tokens"})
			w.Write(answer)
			log.Println("TOKEN DELETE")

		}

	}

}
