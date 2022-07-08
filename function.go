package shan3

import "context"

type Function struct {
	Name          string
	Fn            func(context.Context, UserDTO, []byte, map[string]string) (interface{}, error)
	Transcation   bool
	ContentType   string
	Authorization bool
	Logged        bool
}

type MethodFunction map[string]Function

func (mf MethodFunction) GetMethod(name string) (interface{}, bool) {
	f, ok := mf[name]
	return f, ok
}
