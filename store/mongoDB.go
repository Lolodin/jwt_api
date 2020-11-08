package store

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"time"
)

//mongodb://user:root@localhost:27017
type MongoStore struct {
	*mongo.Client
}
type User struct {
	UUID     string
	Name     string
	Password []byte
}
type Session struct {
	UUID     string
	Name     string
	Expr     int64
	RefToken []byte
	Create   int64
}

func NewMongoStore(client *mongo.Client) MongoStore {
	s := MongoStore{client}
	return s
}
func (s *MongoStore) AddUser(name, password string) error {
	bpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u := User{Name: name, Password: bpass}

	s.Client.Connect(context.Background())
	_, err = s.Client.Database("user").Collection("users").InsertOne(context.Background(), &u)
	return err
}
func (s *MongoStore) GetUser(name string) *User {
	filter := bson.D{{"name", name}}
	u := User{}
	s.Client.Connect(context.Background())
	err := s.Client.Database("user").Collection("users").FindOne(context.Background(), filter).Decode(&u)
	if err != nil {
		return nil
	}
	return &u
}
func (s *MongoStore) AddSession(name, uuid string, refToken []byte, exp int64) {
	s.Client.Connect(context.Background())
	session := Session{}
	session.Name = name
	session.Create = time.Now().Unix()
	session.Expr = exp
	session.RefToken = refToken
	session.UUID = uuid

	_, err := s.Client.Database("user").Collection("sessions").InsertOne(context.Background(), &session)
	if err != nil {
		fmt.Println(err)
	}
}
func (s *MongoStore) GetSession(uuid string) *Session {
	filter := bson.D{{"uuid", uuid}}
	u := Session{}
	s.Client.Connect(context.Background())
	err := s.Client.Database("user").Collection("sessions").FindOne(context.Background(), filter).Decode(&u)
	if err != nil {
		return nil
	}
	return &u
}

// Return true if session exist
func (s *MongoStore) CheckSession(uuid string) bool {
	filter := bson.D{{"uuid", uuid}}
	u := Session{}
	s.Client.Connect(context.Background())
	err := s.Client.Database("user").Collection("sessions").FindOne(context.Background(), filter).Decode(&u)
	if err != nil {
		return false
	}
	if u.Name != "" {
		return true
	}
	return false
}

//delete Session, return true if session delete
func (s *MongoStore) DeleteSession(uuid string) bool {
	filter := bson.D{{"uuid", uuid}}

	s.Client.Connect(context.Background())
	res, err := s.Client.Database("user").Collection("sessions").DeleteOne(context.Background(), filter)
	if err != nil || res.DeletedCount == 0 {
		return false
	}
	return true

}

//delete Sessions, return true if all sessions delete
func (s *MongoStore) DeleteAllSessions(name string) bool {
	filter := bson.D{{"name", name}}
	s.Client.Connect(context.Background())
	res, err := s.Client.Database("user").Collection("sessions").DeleteMany(context.Background(), filter)
	if err != nil || res.DeletedCount == 0 {
		return false
	}
	return true

}
