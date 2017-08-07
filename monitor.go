package sysmon

import (
	"time"
)

// System resource monitor
type Monitor interface {
	Name() string
	Init()
	GetFields() []string
	GetUptime() time.Duration
	GetValue(string) float64
	UpdateValues() error
}
