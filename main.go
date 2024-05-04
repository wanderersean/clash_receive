package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type StringSaver struct {
	mu             sync.Mutex
	data           string
	lastUpdateTime time.Time
}

type Resp struct {
	Data           string
	LastUpdateTime time.Time
}

func (s *StringSaver) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if request.Method == "POST" {
		bs, err := io.ReadAll(request.Body)
		if err != nil {
			log.Default().Println("error when process post", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		s.data = string(bs)
		s.lastUpdateTime = time.Now()
		writer.WriteHeader(http.StatusOK)
		log.Default().Println("succeeded save data")
	} else if request.Method == "GET" {
		bs, err := json.Marshal(Resp{Data: s.data, LastUpdateTime: s.lastUpdateTime})
		if err != nil {
			log.Default().Println("error when process marshal", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = writer.Write(bs)
		if err != nil {
			log.Default().Println("error when process get", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}

func main() {
	_ = http.ListenAndServe("0.0.0.0:8888", &StringSaver{})
}
