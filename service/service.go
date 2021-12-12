package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"github.com/nfk93/lunar/order"
)

const counterServerAddr string = "http://127.0.0.1:10002"

type Rpc struct{}

func (r *Rpc) AddInventory(inventoryToAdd map[string]uint64, _ *struct{}) error {
	err := addInventory(inventoryToAdd)
	return err
}

func (r *Rpc) ProcessOrder(order order.Order, _ *struct{}) error {
	err := processOrder(order)
	return err
}

func (r *Rpc) CheckInventory(itemId string, value *uint64) error {
	v, err := checkInventory(itemId)
	if err != nil {
		return err
	}

	*value = v
	return nil
}

func addInventory(items map[string]uint64) error {
	log.Printf("adding inventory: %v", items)

	for k, v := range items {
		path := fmt.Sprintf("/counter/%v", k)
		for i := 0; i < int(v); i++ {
			// TODO handle response
			_, err := http.Post(counterServerAddr+path, "", nil)
			if err != nil {
				// TODO: handle error
				return err
			}
		}
	}

	return nil
}

func processOrder(order order.Order) error {
	log.Printf("received order %s, checking if processable...", order.ToString())

	// check if our current supply is sufficient enough to process the order
	for item, amount := range order.Items {
		inventoryAmount, err := checkInventory(item)
		if err != nil {
			return err
		}

		if inventoryAmount < amount {
			log.Printf("order %s, can't be processed due to insufficient supply of %s", order.Id, item)
			return fmt.Errorf("insufficient supply: need %v item:%v, supply only has %v",
				amount, item, inventoryAmount)
		}
	}
	log.Printf("order %s can be processed, starting processing", order.Id)

	// begin processing the order, stopping if we fail to decrement the counter,
	// rolling back the part of the order we have already processed
	// Failing to decrement means that somebody else has used up the supply while
	// we were processing
	var processed map[string]uint64 = make(map[string]uint64)
	for item, amount := range order.Items {
		path := fmt.Sprintf("/counter/%v/decr", item)
		for i := 0; i < int(amount); i++ {
			resp, _ := http.Post(counterServerAddr+path, "", nil)
			if resp.StatusCode != 200 {
				log.Printf("aborting order %s due to supply being depleted. Starting rollback", order.Id)
				rollback(processed)
				return fmt.Errorf("couldn't process order because supply was depleted during processing")
			}
			processed[item] += 1
		}
	}

	log.Printf("succesfully procesed order %s", order.Id)
	return nil
}

func checkInventory(itemId string) (uint64, error) {
	path := fmt.Sprintf("/counter/%v", itemId)
	resp, err := http.Get(counterServerAddr + path)

	if err != nil {
		return 0, err
	}

	var inventoryAmount uint64
	err = binary.Read(resp.Body, binary.LittleEndian, &inventoryAmount)
	if err != nil {
		return 0, err
	}
	log.Printf("lookup counter for id %s", itemId)
	return inventoryAmount, nil
}

func rollback(m map[string]uint64) {
	for item, amount := range m {
		path := fmt.Sprintf("/counter/%v", item)
		for i := 0; i < int(amount); i++ {
			// TODO handle response
			_, err := http.Post(counterServerAddr+path, "", nil)
			if err != nil {
				// TODO handle error by retrying until correct status is achieved
			}
		}
	}
}

func main() {
	rpc_ := new(Rpc)
	err := rpc.Register(rpc_)
	if err != nil {
		log.Fatal("Format of service Task isn't correct. ", err)
	}
	// Register a HTTP handler
	rpc.HandleHTTP()
	listener, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("Listen error: ", e)
	}
	log.Printf("Serving RPC server on port %d", 1234)

	// setup logging file
	f, err := os.OpenFile("log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// Start accept incoming HTTP connections
	http.Serve(listener, nil)
	if err != nil {
		log.Fatal("Error serving: ", err)
	}
}
