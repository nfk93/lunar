package main

import (
	"encoding/binary"
	"fmt"
	"net/http"
	"sync"
)

func main() {

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)

		go func() {
			for j := 0; j < 1000; j++ {
				http.Post("http://127.0.0.1:10002/counter", "", nil)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	resp, err := http.Get("http://127.0.0.1:10002/counter")
	if err != nil {
		panic(err.Error())
	}

	var value uint64
	err = binary.Read(resp.Body, binary.LittleEndian, &value)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("count is ", value)
}
