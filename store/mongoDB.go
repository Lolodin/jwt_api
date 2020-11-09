package store

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
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
// Добавить юзера в БД, возвращает ошибку или nil T
func (s *MongoStore) AddUser(name, password string) error {
	bpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u := User{Name: name, Password: bpass}

	s.Client.Connect(context.Background())


	
	sess, err:=s.Client.StartSession()
	if err != nil {
		return err
	}
	err = sess.StartTransaction()
	defer sess.EndSession(context.Background())

	if err != nil {
		return err
	}
	mongo.WithSession(context.Background(), sess, func(sessionContext mongo.SessionContext) error {
		_, err = s.Client.Database("user").Collection("users").InsertOne(context.Background(), &u)
		if err != nil {
			return sess.AbortTransaction(context.Background())
		}
		return nil
	})
	return sess.CommitTransaction(context.Background())



}
// Возвращает user или nil в случае неудачи
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
func (s *MongoStore) AddSession(name, uuid string, refToken []byte, exp int64) error {
	s.Client.Connect(context.Background())
	session := Session{}
	session.Name = name
	session.Create = time.Now().Unix()
	session.Expr = exp
	session.RefToken = refToken
	session.UUID = uuid
	sess, err:=s.Client.StartSession()
	if err != nil {
		return err
	}
	err = sess.StartTransaction()
	if err !=nil {
		return err
	}
return 	mongo.WithSession(context.Background(), sess, func(sessionContext mongo.SessionContext) error {
		_, err = s.Client.Database("user").Collection("sessions").InsertOne(context.Background(), &session)
		if err != nil {
			sess.AbortTransaction(context.Background())
			return err
		}
		return nil
	})


}
func (s *MongoStore) GetSession(uuid string) *Session {
	filter := bson.D{{"uuid", uuid}}
	u := Session{}
	s.Client.Connect(context.Background())
	opt:= options.Session().SetDefaultReadConcern(readconcern.Majority())
	sess, _:=s.Client.StartSession(opt)
	defer sess.EndSession(context.Background())
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
	opt:= options.Session().SetDefaultReadConcern(readconcern.Majority())
	sess, _:=s.Client.StartSession(opt)
	defer sess.EndSession(context.Background())
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
