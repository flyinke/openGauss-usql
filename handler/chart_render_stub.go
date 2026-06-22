//go:build no_chart

package handler

import (
	"context"
	"errors"
	"io"

	"github.com/kenshaw/rasterm"
	"github.com/xo/usql/metacmd/charts"
)

var errChartDisabled = errors.New(`chart rendering is unavailable in "no_chart" builds`)

func renderChart(context.Context, io.Writer, rasterm.TermType, charts.ChartConfig, string) error {
	return errChartDisabled
}
