package shan3

import (
	"context"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"
)

type Server interface{
	GetMgo() MgoClient
	GetLogger() LogClient
	GetFunction() MethodFunction
	ParseToken(http.ResponseWriter, *http.Request) (UserDTO, error)
}

func transaction(
	s Server,
	ff Function,
	w http.ResponseWriter,
	r *http.Request) (interface{}, error) {
	fn := ff.Fn
	token, err := s.ParseToken(w, r)
	if err != nil {
		return nil, errors.New("请登录")
	}
	query := QueryParse(r)
	body, err := BodyBuffer(r)
	if err != nil {
		return nil, err
	}
	var ret interface{}

	ctx, cancel := context.WithCancel(WithValue(context.TODO(), w, r))
	go s.GetMgo().UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		defer func() {
			if e := recover(); e != nil {
				ret = nil
				err = errors.New(fmt.Sprintf("%v", e))
				fmt.Printf("%s", debug.Stack())
				err = sessionContext.AbortTransaction(sessionContext)
				cancel()
			}
		}()
		err = sessionContext.StartTransaction()
		if err != nil {
			ret = nil
			return err
		}
		ret, err = fn(sessionContext, token, body, query)
		if err != nil {
			sessionContext.AbortTransaction(sessionContext)
		} else {
			if ff.Logged == true {
				s.GetLogger().AppendLog(context.TODO(), Log{
					token.GetId(),
					token.GetName(),
					query["id"],
					ff.Name,
					query["method"],
					time.Now(),
				})
			}
			err = sessionContext.CommitTransaction(sessionContext)
		}
		cancel()
		return err
	})
	<-ctx.Done()
	return ret, err
}

func process(
	s Server,
	ff Function,
	w http.ResponseWriter,
	r *http.Request) (interface{}, error) {
	fn := ff.Fn
	token, err := s.ParseToken(w, r)
	if err != nil {
		return nil, errors.New("请登录")
	}
	query := QueryParse(r)
	body, err := BodyBuffer(r)
	if err != nil {
		return nil, err
	}
	var ret interface{}
	ctxCancel, cancel := context.WithCancel(WithValue(context.TODO(), w, r))
	go func() {
		defer func() {
			if e := recover(); e != nil {
				ret = nil
				err = errors.New(fmt.Sprintf("%v", e))
				fmt.Printf("%s", debug.Stack())
				cancel()
			}
		}()
		ret, err = fn(ctxCancel, token, body, query)
		if ff.Logged== true {
			s.GetLogger().AppendLog(context.TODO(), Log{
				token.GetId(),
				token.GetName(),
				query["id"],
				ff.Name,
				query["method"],
				time.Now(),
			})
		}
		cancel()
	}()
	select {
	case <-ctxCancel.Done():
		return ret, err
		//case <- timeout.Done():
		//	return nil, errors.New("请求超时")
	}
}

func openProcess(
	_ Server,
	ff Function,
	w http.ResponseWriter,
	r *http.Request) (interface{}, error) {
	fn := ff.Fn
	var token UserDTO
	query := QueryParse(r)
	body, err := BodyBuffer(r)
	if err != nil {
		return nil, err
	}
	ctxCancel, _ := context.WithCancel(WithValue(context.TODO(), w, r))
	ret, err := fn(ctxCancel, token, body, query)
	return ret, err
}

func run(s Server) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		method := GetMethodName(r)
		fn, ok := s.GetFunction().GetMethod(method).(Function)
		if ok == true {
			var ret interface{}
			var err error
			if fn.Authorization == false {
				ret, err = openProcess(s, fn, w, r)
			} else if fn.Transcation == true {
				ret, err = transaction(s, fn, w, r)
			} else {
				ret, err = process(s, fn, w, r)
			}
			if err != nil {
				ResponseHandleError(w, err.Error(), ret)
			} else if fn.ContentType == "" {
				ResponseWapperSucc(w, ret)
			} else {
				w.Header().Add("Content-Type", fn.ContentType)
				w.Write(ret.([]byte))
			}
			return
		}
		w.WriteHeader(404)
	}
}

func Serv(s Server, port, apiversion string) error {
	router := httprouter.New()
	router.POST(apiversion, run(s))
	router.GET(apiversion, run(s))
	log.Printf("OPENUSERV is started. The PID is %d and listening on port %s\n", os.Getpid(), port)
	return http.ListenAndServe(port, router)
}