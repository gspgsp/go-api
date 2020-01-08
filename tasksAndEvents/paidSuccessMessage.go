package tasksAndEvents

import "encoding/json"

type paidSuccessMessage struct {
	Class  string              `json:"class"`
	Method string              `json:"method"`
	Params *PaidSuccessMessage `json:"params"`
}

func (o *paidSuccessMessage) ToBytes() ([]byte, error) {
	o.Class = "PaidSuccessMessage"
	o.Method = "paidSuccessMessage"

	data, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	return data, nil
}

type PaidSuccessMessage struct {
	OrderId    int    `json:"order_id"`
	BranchType string `json:"branch_type"`
	PaySource  string `json:"pay_source"`
	EventType  string `json:"event_type"`
}

/**
这个主要是更新用户课程的异步操作
*/
func (o *PaidSuccessMessage) Send() (int, error) {
	data := new(paidSuccessMessage)
	data.Params = o

	return operateDB(data)
}
