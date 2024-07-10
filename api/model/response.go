package model

type Response struct {
	StatusCode  int
	Description string
	Data        interface{}
}
