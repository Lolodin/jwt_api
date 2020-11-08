package main

import (
	"fmt"
	"github.com/Lolodin/jwt_api/controller"
	store2 "github.com/Lolodin/jwt_api/store"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
)

func main() {
	//БД
	fmt.Println(os.Getenv("DATABASE_URL"), os.Getenv("PORT"))
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("DATABASE_URL")))
	if err != nil {
		log.Fatal(err)
	}
	store := store2.NewMongoStore(client)

	mux := http.NewServeMux()
	mux.HandleFunc("/getTokens", controller.GetTokens(&store))
	mux.HandleFunc("/", controller.Index())
	mux.HandleFunc("/reg", controller.Reg())
	mux.HandleFunc("/register", controller.Register(&store))
	mux.HandleFunc("/refresh", controller.RefreshTokens(&store))
	mux.HandleFunc("/deleteRef", controller.DeleteRefreshToken(&store))
	mux.HandleFunc("/deleteAll", controller.DeleteAllUserTokens(&store))

	http.ListenAndServe(":"+os.Getenv("PORT"), mux)
}
