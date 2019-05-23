package main

import (
	"log"
	"net/http"
)

func main() {
	log.Print("Port listened on 8009")
	http.ListenAndServe(":8009", nil)
}
