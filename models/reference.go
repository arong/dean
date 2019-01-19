package models

type Reference interface {
	IncreaseRef(interface{})
	DecreaseRef(interface{})
}
