package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)
	log.Println("本地開發伺服器已啟動：http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
