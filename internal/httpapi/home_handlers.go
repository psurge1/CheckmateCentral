package httpapi

import (
	"log"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Home Page"))
	if err != nil {
		log.Println("Failed to write to response")
	}
}
