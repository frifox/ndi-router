package main

import "encoding/json"

type Emit struct {
	Type string
}

// Requests

type UpdateOutput struct {
	Output string
	Input  string
}

// Responses

type RouterMatrix struct {
	Emit
	Matrix map[string]string
}

func (e *RouterMatrix) MarshalJSON() (data []byte, err error) {
	e.Type = "RouterMatrix"
	return json.Marshal(*e)
}
