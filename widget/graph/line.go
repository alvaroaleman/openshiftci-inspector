package graph

import (
	"github.com/janoszen/openshiftci_inspector/widget"
)

func NewLineGraphWidget(Series []Serie) (widget.Widget, error) {
	return &lineGraphWidget{
		Series: Series,
	}, nil
}

type lineGraphWidget struct {
	Series []Serie `json:"series"`
}

func (g *lineGraphWidget) Type() string {
	return "line"
}

func (g *lineGraphWidget) Validate() error {
	return nil
}

func (g *lineGraphWidget) JSON() interface{} {
	return g
}

type Serie struct {
	Legend string            `json:"legend"`
	Points map[int64]float64 `json:"points"`
}
