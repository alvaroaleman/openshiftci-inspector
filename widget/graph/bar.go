package graph

import (
	"github.com/janoszen/openshiftci_inspector/widget"
)

func NewBarGraph(Values []Value) (widget.Widget, error) {
	return &barGraphWidget{
		Values: Values,
	}, nil
}

type barGraphWidget struct {
	Values []Value `json:"values"`
}

func (g *barGraphWidget) Type() string {
	return "bar"
}

func (g *barGraphWidget) Validate() error {
	return nil
}

func (g *barGraphWidget) JSON() interface{} {
	return g
}

type Value struct {
	Legend string  `json:"legend"`
	Point  float64 `json:"point"`
	Time   int64   `json:"time"`
}

