package main

import (
	"encoding/binary"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type counterStorage struct {
	counters map[string]uint64
	lock     sync.RWMutex
}

func (c *counterStorage) getCounterByIdHandler(w http.ResponseWriter, r *http.Request) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	vars := mux.Vars(r)

	// Default value of map[string]uint64 is 0 when looking up an unitialized key,
	// so we don't need to condition on the id's existence in the map
	v := c.counters[vars["id"]]
	w.WriteHeader(http.StatusOK)
	binary.Write(w, binary.LittleEndian, v)
}

func (c *counterStorage) incrementCounterByIdHandler(w http.ResponseWriter, r *http.Request) {
	c.lock.Lock()
	defer c.lock.Unlock()

	vars := mux.Vars(r)

	// Default value of map[string]uint64 is 0 when looking up an unitialized key,
	// so we don't need to condition on the id's existence in the map
	w.WriteHeader(http.StatusOK)
	c.counters[vars["id"]] += 1
}

func (c *counterStorage) decrementCounterByIdHandler(w http.ResponseWriter, r *http.Request) {
	c.lock.Lock()
	defer c.lock.Unlock()

	vars := mux.Vars(r)
	id := c.counters[vars["id"]]

	// Default value of map[string]uint64 is 0 when looking up an unitialized key,
	// so we don't need to condition on the id's existence in the map
	if id > 0 {
		w.WriteHeader(http.StatusOK)
		c.counters[vars["id"]] -= 1
	} else {
		http.Error(w, fmt.Sprintf("counter %v is already at 0", id), http.StatusConflict)
	}
}

func startListening() {
	// initialize the counter
	counters := counterStorage{
		counters: make(map[string]uint64),
	}

	r := mux.NewRouter()
	r.HandleFunc("/counter/{id}", counters.getCounterByIdHandler).Methods("GET")
	r.HandleFunc("/counter/{id}", counters.incrementCounterByIdHandler).Methods("POST")
	r.HandleFunc("/counter/{id}/decr", counters.decrementCounterByIdHandler).Methods("POST")

	http.ListenAndServe(":10002", r)
}

func main() {
	startListening()
}
