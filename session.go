package shan3

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session interface {
	GetId() interface{}
	GetUserId() string
	GetUserName() string
}

type SessionServer interface {
	SessionGetId(context.Context, string) (UserDTO, error)
	UserCreateInstance(context.Context, UserDTO) (UserDTO, error)
	SessionSet(context.Context, UserDTO) (Session, error)
	UserFindOne(context.Context, bson.M) (UserDTO, error)
	GetSessionId(context.Context, interface{}) (string, error)
	ParseToken(context.Context, string, interface{}) (UserDTO, error)
}

type DeafultSession struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	UserId    string             `json:"user_id" bson:"user_id"`
	UserName  string             `json:"user_name" bson:"user_name"`
	Role      string             `json:"role" bson:"role"`
	UUID      string             `json:"uuid" bson:"uuid"`
	Timestamp int64              `json:"timestamp" bson:"timestamp"`
}

func (d DeafultSession) GetId() interface{} {
	return d.Id
}

func (d DeafultSession) GetUserId() string {
	return d.UserId
}

func (d DeafultSession) GetUserName() string {
	return d.UserName
}

func DeafultSessionDecodeCookie(ss string) (interface{}, error) {
	result := DeafultSession{}
	buf, err := base64.StdEncoding.DecodeString(ss)
	if err != nil {
		return result, err
	}
	b, err := base64.StdEncoding.DecodeString(string(buf))
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(b, &result)
	return result, err
}

func DeafultSessionEncodeCookie(session interface{}) (string, error) {
	buf, err := json.Marshal(session)
	if err != nil {
		return "", err
	}
	result := base64.StdEncoding.EncodeToString([]byte(base64.StdEncoding.EncodeToString(buf)))
	return result, nil
}

func NewDefaultSession() DeafultSession {
	return DeafultSession{}
}