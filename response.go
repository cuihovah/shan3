package shan3

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
)

type Retval struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}


func compress(data []byte) []byte {
	var b bytes.Buffer
	write, _ := gzip.NewWriterLevel(&b, gzip.BestCompression)
	defer write.Close()
	write.Write(data)
	write.Flush()
	return b.Bytes()
}

func ResponseWapperSucc(w http.ResponseWriter, data interface{}) {
	var ret Retval
	ret.Code = 0
	ret.Msg = "OK"
	ret.Data = data
	result, err := json.Marshal(ret)
	if err != nil {
		fmt.Println(err.Error())
	}
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Content-Encoding", "gzip")
	result = compress(result)
	w.Write(result)
}

func ResponseHandleError(w http.ResponseWriter, msg string, data interface{}) {
	var ret Retval
	ret.Code = 3
	ret.Msg = msg
	ret.Data = data
	result, _ := json.Marshal(ret)
	w.Header().Add("Content-Type", "application/json")
	w.Write(result)
}

