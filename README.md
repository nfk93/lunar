### Case for Lunar interview
This project contains two runnable microservices and a quickly made CLI to interact with them

In `/counter` there is a restful microservice that tracks counters on the paths `/counter/{id}`. It increments on a POST request and returns its value on a GET request. To decrement, use a POST request on `/counter/{id}/decr`
To run it, use `go run counter.go` in the counter folder

In `/service` there is a microservice for a business that wants to track its inventory. It can receive orders, check inventory for a given item id, and add to its inventory all via RPC. It logs all its request to a log file that can be inspected to see which orders have been processed. To run it use `go run service.go` in the service folder.

The CLI to interact with the service can be run by using `go run client.go ARGS`, where args specify what RPC to make. You can use the following
* `go run client.go CheckInventory ID` where `ID` is the id of the item to check
* `go run client.go ProcessOrder ORDERID ITEMID1 AMOUNT1 ITEMID2 AMOUNT2...` where `ORDERID` is the id of the order and `ITEMIDn AMOUNTn` is an item to add and its amount. Specify as many as you'd like. Disclaimer, there is no graceful exit in the CLI if the order can't be processed 
* `go run client.go AddInventory ITEMID1 AMOUNT1 ITEMID2 AMOUNT2...` where `ITEMIDn AMOUNTn` specify the amount of which item to add
