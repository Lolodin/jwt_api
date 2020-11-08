package main

import (
	"github.com/Lolodin/jwt_api/controller"
	store2 "github.com/Lolodin/jwt_api/store"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)
var connectMongo = "mongodb://user:root@localhost:27017"
func main() {
	//БД
	client, err:=mongo.NewClient(options.Client().ApplyURI(connectMongo))
	if err != nil {
		log.Fatal(err)
	}
	store := store2.NewMongoStore(client)

	mux := http.NewServeMux()
	mux.HandleFunc("/getTokens", controller.GetTokens(&store))
	mux.HandleFunc("/", controller.Index())
	mux.HandleFunc("/refresh", controller.RefreshTokens(&store))


	http.ListenAndServe(":8080", mux)
}
