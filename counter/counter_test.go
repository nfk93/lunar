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
	router.HandleFunc("/counter/{id}/decr", counters.decrementCounterByIdHandler).Methods("POST")

	testcounters := []string{"a", "b", "c", "d"}
	var wg sync.WaitGroup

	// Spawn threads that make POST requests to increment the counters
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

	// Spawn threads that make POST requests that decrement the counter a
	path := "/counter/" + testcounters[0] + "/decr"
	for i := 0; i < 10; i++ {
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

	wg.Wait()

	// make GET requests and check that they have the correct value
	rr := httptest.NewRecorder()
	var m = map[string]uint64{
		"a": 10000,
		"b": 20000,
		"c": 20000,
		"d": 20000,
	}
	for _, id := range testcounters {
		req, err := http.NewRequest("GET", "/counter/"+id, nil)
		if err != nil {
			t.Fatal(err)
		}
		router.ServeHTTP(rr, req)

		var value uint64
		err = binary.Read(rr.Body, binary.LittleEndian, &value)
		if err != nil {
			t.Fatal(err)
		}
		if value != m[id] {
			t.Errorf("unexpected value at counter %s: got %v want %v", id, value, m[id])
		}
	}
}
