package server

import (
	"github.com/ant0ine/go-json-rest"
	"log"
)

func GetRecord(w *rest.ResponseWriter, req *rest.Request) {
	loc := req.PathParam("loc")
	mu.RLock()
	r, exists := records[loc]
	mu.RUnlock()
	if !exists {
		log.Printf("First request '%s'", loc)
		mu.Lock()
		r, _ = pm25(loc)
		records[loc] = r
		mu.Unlock()
	}
	w.WriteJson(r)
}
