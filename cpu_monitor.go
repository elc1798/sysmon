package sysmon

import (
	"fmt"
	"log"
	"time"

	linuxproc "github.com/c9s/goprocinfo/linux"
)

const (
	CPUINFO_PATH = "/proc/cpuinfo"
	CPUSTAT_PATH = "/proc/stat"
)

var CPU_MON_FIELDS []string

// CPU Usage monitor. Usage is reported as a percentage
type CPUMonitor struct {
	startTime    time.Time
	cpuStatCache []linuxproc.CPUStat
	numCores     int
	values       map[string]float64
}

func (m *CPUMonitor) Name() string {
	return "cpu_monitor"
}

func (m *CPUMonitor) Init() {
	m.startTime = time.Now()

	m.numCores = 32 // Use 32 cores by default, leave unused empty
	if info, err := linuxproc.ReadCPUInfo(CPUINFO_PATH); err != nil {
		log.Printf("Failed to read CPU Info: %v", err)
	} else {
		m.numCores = info.NumCore()
	}

	m.cpuStatCache = make([]linuxproc.CPUStat, m.numCores)
	CPU_MON_FIELDS = make([]string, m.numCores)
	for i := 0; i < m.numCores; i++ {
		CPU_MON_FIELDS[i] = fmt.Sprintf("cpu%d", i)
	}

	m.values = make(map[string]float64)
}

func (m *CPUMonitor) GetFields() []string {
	return CPU_MON_FIELDS
}

func (m *CPUMonitor) GetUptime() time.Duration {
	return time.Since(m.startTime)
}

func (m *CPUMonitor) GetValue(field string) float64 {
	return m.values[field]
}

func (m *CPUMonitor) UpdateValues() error {
	stats, err := linuxproc.ReadStat(CPUSTAT_PATH)
	if err != nil {
		return err
	}

	for i, s := range stats.CPUStats {
		m.values[fmt.Sprintf("cpu%d", i)] = calcSingleCoreUsage(s, m.cpuStatCache[i])

		// Cache old value
		m.cpuStatCache[i] = s
	}

	return nil
}

func calcSingleCoreUsage(curr, prev linuxproc.CPUStat) float64 {
	// https://stackoverflow.com/questions/11356330/getting-cpu-usage-with-golang

	prevIdle := prev.Idle + prev.IOWait
	prevNonIdle := prev.User + prev.Nice + prev.System + prev.IRQ + prev.SoftIRQ + prev.Steal
	prevTotal := prevIdle + prevNonIdle

	idle := curr.Idle + curr.IOWait
	nonIdle := curr.User + curr.Nice + curr.System + curr.IRQ + curr.SoftIRQ + curr.Steal
	total := idle + nonIdle

	totald := total - prevTotal
	idled := idle - prevIdle

	return (float64(totald) - float64(idled)) / float64(totald) * 100.0
}
