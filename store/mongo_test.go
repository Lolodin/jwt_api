package store

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

var connectMongo = "mongodb://user:root@localhost:27017"

func TestMongoStore_AddUser(t *testing.T) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectMongo))
	if err != nil {
		t.Error(err)
	}
	s := NewMongoStore(client)
	s.AddUser("testuser", "123456")
}
func TestMongoStore_GetUser(t *testing.T) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectMongo))
	if err != nil {
		t.Error(err)
	}
	s := NewMongoStore(client)
	user := s.GetUser("testuser")
	res := bcrypt.CompareHashAndPassword(user.Password, []byte("123456"))
	t.Log(user.Name == "testuser", res == nil)
}
func TestMongoStore_CheckSession(t *testing.T) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectMongo))
	if err != nil {
		t.Error(err)
	}
	s := NewMongoStore(client)
	user := s.CheckSession("ca8fdae648fc8ec128e9e404b026ac85")
	fmt.Println(user)
}
func TestMongoStore_DeleteSession(t *testing.T) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectMongo))
	if err != nil {
		t.Error(err)
	}
	s := NewMongoStore(client)
	user := s.DeleteSession("ca8fdae648fc8ec128e9e404b026ac85")
	fmt.Println(user)
}
