package sysmon

import (
	"time"

	linuxproc "github.com/c9s/goprocinfo/linux"
)

const (
	MEMINFO_PATH = "/proc/meminfo"
)

var (
	MEM_MON_FIELDS = []string{
		"mem_total",
		"mem_free",
		"mem_available",
		"mem_cached",
		"swap_total",
		"swap_free",
		"swap_cached",
	}
)

type MemoryMonitor struct {
	startTime time.Time
	values    map[string]float64
}

func (m *MemoryMonitor) Name() string {
	return "mem_monitor"
}

func (m *MemoryMonitor) Init() {
	m.startTime = time.Now()
	m.values = make(map[string]float64)
}

func (m *MemoryMonitor) GetFields() []string {
	return MEM_MON_FIELDS
}

func (m *MemoryMonitor) GetUptime() time.Duration {
	return time.Since(m.startTime)
}

func (m *MemoryMonitor) GetValue(field string) float64 {
	return m.values[field]
}

func (m *MemoryMonitor) UpdateValues() error {
	info, err := linuxproc.ReadMemInfo(MEMINFO_PATH)
	if err != nil {
		return err
	}

	m.values["mem_total"] = float64(info.MemTotal)
	m.values["mem_free"] = float64(info.MemFree)
	m.values["mem_available"] = float64(info.MemAvailable)
	m.values["mem_cached"] = float64(info.Cached)
	m.values["swap_total"] = float64(info.SwapTotal)
	m.values["swap_free"] = float64(info.SwapFree)
	m.values["swap_cached"] = float64(info.SwapCached)

	return nil
}
