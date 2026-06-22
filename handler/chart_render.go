//go:build !no_chart

package handler

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/kenshaw/rasterm"
	"github.com/xo/echartsgoja"
	"github.com/xo/resvg"
	"github.com/xo/usql/metacmd/charts"
)

func renderChart(ctx context.Context, stdout io.Writer, typ rasterm.TermType, cfg charts.ChartConfig, data string) error {
	echarts := echartsgoja.New(echartsgoja.WithWidthHeight(cfg.W, cfg.H))
	res, err := echarts.RenderOptions(ctx, data)
	if err != nil {
		return err
	}
	if cfg.File != "" {
		fmt.Println("writing to", cfg.File)
		return os.WriteFile(cfg.File, []byte(res), 0o644)
	}
	img, err := resvg.Render([]byte(res), resvg.WithBackground(cfg.Background))
	if err != nil {
		return err
	}
	return typ.Encode(stdout, img)
}
