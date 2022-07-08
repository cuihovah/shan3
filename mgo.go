package shan3

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type MgoClient interface{
	UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error
}
