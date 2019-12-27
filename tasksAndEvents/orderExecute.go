package tasksAndEvents

import "encoding/json"

type orderExecute struct {
	Class  string        `json:"class"`
	Method string        `json:"method"`
	Params *OrderExecute `json:"params"`
}

func (o *orderExecute) ToBytes() ([]byte, error) {
	o.Class = "OrderExecute"
	o.Method = "handle"

	data, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	return data, nil
}

type OrderExecute struct {
	OrderId    int    `json:"order_id"`
	BranchType string `json:"branch_type"`
}

/**
这个主要是更新用户课程的异步操作
*/
func (o *OrderExecute) Update() (int, error) {
	data := new(orderExecute)
	data.Params = o

	return operateDB(data)
}
