package main

import (
	"encoding/binary"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gorilla/mux"
)

func TestCounter(t *testing.T) {
	counters := counterStorage{
		counters: make(map[string]uint64),
	}
	router := mux.NewRouter()
	router.HandleFunc("/counter/{id}", counters.getCounterByIdHandler).Methods("GET")
	router.HandleFunc("/counter/{id}", counters.incrementCounterByIdHandler).Methods("POST")

	testcounters := []string{"a", "b", "c", "asdalsmkd"}
	var wg sync.WaitGroup

	// Spawn threads that make POST requests
	for _, id := range testcounters {
		path := "/counter/" + id
		for i := 0; i < 20; i++ {
			wg.Add(1)

			go func() {
				req, _ := http.NewRequest("POST", path, nil)
				rr := httptest.NewRecorder()
				for j := 0; j < 1000; j++ {

					router.ServeHTTP(rr, req)
				}
				wg.Done()
			}()
		}
	}

	// wait until all threads making requests has finished
	wg.Wait()

	// make GET requests and check that they have the correct value
	rr := httptest.NewRecorder()
	for _, id := range testcounters {
		req, err := http.NewRequest("GET", "/counter/"+id, nil)
		if err != nil {
			t.Fatal(err)
		}
		router.ServeHTTP(rr, req)

		expected := uint64(20000)
		var value uint64
		err = binary.Read(rr.Body, binary.LittleEndian, &value)
		if err != nil {
			t.Fatal(err)
		}
		if value != expected {
			t.Errorf("unexpected value at counter %s: got %v want %v", id, value, expected)
		}
	}
}
