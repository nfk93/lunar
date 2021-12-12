package order

import "fmt"

type Order struct {
	Id    string
	Items map[string]uint64
}

func (o *Order) ToString() string {
	return fmt.Sprintf("order:%v{%v}", o.Id, o.Items)
}
