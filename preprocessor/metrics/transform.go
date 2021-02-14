package metrics

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/promql"

	"github.com/janoszen/openshiftci_inspector/widget"
	"github.com/janoszen/openshiftci_inspector/widget/graph"
)

func transformScalar(s promql.Scalar) (widget.Widget, error) {
	return nil, fmt.Errorf("scalar results are not supported")
}

func transformMatrix(m promql.Matrix, query Query) (widget.Widget, error) {
	var series []graph.Serie
	for _, s := range m {
		legend, err := transformLabels(s.Metric, query.LegendFormat)
		if err != nil {
			return nil, err
		}
		series = append(series, graph.Serie{
			Legend: legend,
			Points: transformPoints(s.Points),
		})
	}

	return graph.NewLineGraphWidget(series)
}

func transformPoints(points []promql.Point) map[int64]float64 {
	result := map[int64]float64{}
	for _, p := range points {
		result[p.T] = p.V
	}
	return result
}

func transformVector(v promql.Vector, query Query) (widget.Widget, error) {
	var values []graph.Value
	for _, s := range v {
		legend, err := transformLabels(s.Metric, query.LegendFormat)
		if err != nil {
			return nil, err
		}
		values = append(values, graph.Value{
			Legend: legend,
			Time: s.T,
			Point: s.V,
		})
	}

	return graph.NewBarGraph(values)
}

func transformLabels(metric labels.Labels, format string) (string, error) {
	tpl, err := template.New("legend").Parse(format)
	if err != nil {
		return "", fmt.Errorf("failed to parse legend template %s (%w)", format, err)
	}

	metricLabels := map[string]string{}
	for _, l := range metric {
		metricLabels[l.Name] = l.Value
	}
	writer := &bytes.Buffer{}
	if err := tpl.Execute(writer, metricLabels); err != nil {
		return "", fmt.Errorf("failed to execute legend template %s (%w)", format, err)
	}
	return string(writer.Bytes()), nil
}

