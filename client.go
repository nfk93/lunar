package main

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
)

type Order struct {
	Id    string
	Items map[string]uint64
}

func makeProcessOrderCall(order Order, client *rpc.Client) {
	var void struct{}

	err := client.Call("Rpc.ProcessOrder", order, &void)
	if err != nil {
		fmt.Printf(err.Error())
	} else {
		fmt.Printf("succesfully processed order %v", order)
	}
}

func makeCheckInventoryCall(itemId string, client *rpc.Client) {
	var result uint64

	err := client.Call("Rpc.CheckInventory", itemId, &result)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("CheckInventory %s: %v", itemId, result)
	}
}

func makeAddInventoryCall(items map[string]uint64, client *rpc.Client) {
	var result struct{}

	err := client.Call("Rpc.AddInventory", items, &result)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Added inventory %v", items)
	}
}

func runClient() {
	var err error
	var reply int

	// Create a TCP connection to localhost on port 1234
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	for {
		client.Call("Task.Increment", struct{}{}, &reply)
		fmt.Printf("Called Task.Increment remotely and received: %d\n", reply)
		time.Sleep(2 * time.Second)
	}
}
