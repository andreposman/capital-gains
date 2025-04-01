package json

import "encoding/json"

type Operation struct {
	Operation string  `json:"operation"`
	UnitCost  float64 `json:"unit-cost"`
	Quantity  int     `json:"quantity"`
}

func ParseInput(input []byte) ([]Operation, error) {
	var operations []Operation
	err := json.Unmarshal(input, &operations)
	return operations, err
}
