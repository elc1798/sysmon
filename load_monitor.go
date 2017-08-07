package sysmon

import (
	"time"

	linuxproc "github.com/c9s/goprocinfo/linux"
)

const (
	LOADINFO_PATH = "/proc/loadavg"
)

var (
	LOAD_MON_FIELDS = []string{
		"last_1_min",
		"last_5_min",
		"last_15_min",
		"num_curr_proc",
		"num_total_proc",
	}
)

type LoadMonitor struct {
	startTime time.Time
	values    map[string]float64
}

func (m *LoadMonitor) Name() string {
	return "load_monitor"
}

func (m *LoadMonitor) Init() {
	m.startTime = time.Now()
	m.values = make(map[string]float64)
}

func (m *LoadMonitor) GetFields() []string {
	return LOAD_MON_FIELDS
}

func (m *LoadMonitor) GetUptime() time.Duration {
	return time.Since(m.startTime)
}

func (m *LoadMonitor) GetValue(field string) float64 {
	return m.values[field]
}

func (m *LoadMonitor) UpdateValues() error {
	info, err := linuxproc.ReadLoadAvg(LOADINFO_PATH)
	if err != nil {
		return err
	}

	m.values["last_1_min"] = float64(info.Last1Min)
	m.values["last_5_min"] = float64(info.Last5Min)
	m.values["last_15_min"] = float64(info.Last15Min)
	m.values["num_curr_proc"] = float64(info.ProcessRunning)
	m.values["num_total_proc"] = float64(info.ProcessTotal)

	return nil
}
