package main

import (
	"log"
	"net/http"
	"sync"
)

var sum int64 = 0

var productNum int64 = 100000

var mutex sync.Mutex

func GetOneProduct() bool {
	mutex.Lock()
	defer mutex.Unlock()
	if sum < productNum {
		sum += 1
		return true
	}
	return false
}

func GetProduct(w http.ResponseWriter, req *http.Request) {
	if GetOneProduct() {
		w.Write([]byte("true"))
		return
	}
	w.Write([]byte("false"))
	return
}

func main() {
	http.HandleFunc("/getOne", GetProduct)
	err := http.ListenAndServe(":8084", nil)
	if err != nil {
		log.Fatal("Err: ", err)
	}
}
