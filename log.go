package shan3

import "context"

type LogClient interface {
	AppendLog(context.Context, interface{})
}
