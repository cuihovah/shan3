package shan3

import "time"

type Log struct {
	UserId string `json:"user_id" bson:"user_id"`
	UserName string `json:"user_name" bson:"user_name"`
	ResourceId string `json:"resource_id" bson:"resource_id"`
	Name string `json:"name" bson:"name"`
	Method string `json:"method" bson:"method"`
	CreateTime time.Time `json:"create_time" bson:"create_time"`
}

