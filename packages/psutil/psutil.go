package psutil

import (
	"fmt"

	i3 "github.com/denbeigh2000/goi3bar"
	"github.com/denbeigh2000/goi3bar/config"
	"github.com/shirou/gopsutil/cpu"

	"time"
)

type psutilCPUPercConfig struct {
	Interval      string  `json:"interval"`
	Name          string  `json:"name"`
	WarnThreshold float64 `json:"warn_threshold"`
	CritThreshold float64 `json:"crit_threshold"`
	PerCPU        bool    `json:"per_cpu"`
}

type psutilCPUPercBuilder struct{}

type psutilCPU struct {
	Name          string  `json:"name"`
	WarnThreshold float64 `json:"warn_threshold"`
	CritThreshold float64 `json:"crit_threshold"`
	PerCPU        bool    `json:"per_cpu"`
}

func (b psutilCPUPercBuilder) Build(c config.Config) (i3.Producer, error) {
	conf := psutilCPUPercConfig{}
	err := c.ParseConfig(&conf)

	interval, err := time.ParseDuration(conf.Interval)
	if err != nil {
		return nil, err
	}

	return &i3.BaseProducer{
		Generator: &psutilCPU{
			Name:          conf.Name,
			WarnThreshold: conf.WarnThreshold,
			CritThreshold: conf.CritThreshold,
			PerCPU:        conf.PerCPU,
		},
		Interval: interval,
	}, nil
}

func init() {
	config.Register("psutil_cpu", psutilCPUPercBuilder{})
}

func (p psutilCPU) Generate() (out []i3.Output, err error) {

	percs, err := cpu.Percent(0, p.PerCPU)
	if err != nil {
		return
	}
	out = make([]i3.Output, len(percs)+1)
	out[0].FullText = "CPU:"
	out[0].Color = i3.DefaultColors.OK
	out[0].Separator = false

	for i, perc := range percs {
		var o i3.Output
		o.FullText = fmt.Sprintf("%02.0f", perc)
		if i < len(percs)-1 {
			o.Separator = false
		} else {
			o.Separator = true
		}
		switch {
		case perc >= p.CritThreshold:
			o.Color = i3.DefaultColors.Crit
		case perc >= p.WarnThreshold:
			o.Color = i3.DefaultColors.Warn
		default:
			o.Color = i3.DefaultColors.OK
		}
		out[i+1] = o
	}

	return
}
