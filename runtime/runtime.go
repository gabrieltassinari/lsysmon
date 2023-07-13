package runtime

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rprobaina/lpfs"
)

var (
	prevIdle  int
	prevTotal int
	cpuUsage  float64
)

type MemoryJSON struct {
	Buffers int
	Cached  int
	Free    int
}

type SwapJSON struct {
	Filename string
	Size     int
	Used     int
}

func ProcessesStat(w http.ResponseWriter) error {
	pps, err := lpfs.GetPerProcessStat()
	if err != nil {
		return err
	}

	b, err := json.Marshal(pps)
	if err != nil {
		return err
	}

	data := fmt.Sprintf("event: process\ndata: %v\n\n", string(b))

	fmt.Fprintf(w, data)

	return nil
}

func Memory(w http.ResponseWriter) error {
	buffers, err := lpfs.GetMemBuffers()
	if err != nil {
		return err
	}

	cached, err := lpfs.GetMemCached()
	if err != nil {
		return err
	}

	free, err := lpfs.GetMemFree()
	if err != nil {
		return err
	}

	msg := MemoryJSON{
		Buffers: buffers,
		Cached:  cached,
		Free:    free,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	data := fmt.Sprintf("event: memory\ndata: %v\n\n", string(b))

	fmt.Fprintf(w, data)

	return nil
}

func Swap(w http.ResponseWriter) error {
	filename, err := lpfs.GetSwapFilename()
	if err != nil {
		return err
	}

	size, err := lpfs.GetSwapSize()
	if err != nil {
		return err
	}

	used, err := lpfs.GetSwapUsed()
	if err != nil {
		return err
	}

	msg := SwapJSON{
		Filename: filename,
		Size:     size,
		Used:     used,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	data := fmt.Sprintf("event: swap\ndata: %v\n\n", string(b))

	fmt.Fprintf(w, data)

	return nil
}

func Uptime(w http.ResponseWriter) error {
	uptime, err := lpfs.GetUptimeSystem()
	if err != nil {
		return err
	}

	data := fmt.Sprintf("event: uptime\ndata: %f\n\n", uptime)

	fmt.Fprintf(w, data)

	return nil
}

func Cpu(w http.ResponseWriter) error {
	user, err := lpfs.GetCpuUserTime()
	if err != nil {
		return err
	}

	nice, err := lpfs.GetCpuNiceTime()
	if err != nil {
		return err
	}

	system, err := lpfs.GetCpuSystemTime()
	if err != nil {
		return err
	}

	idle, err := lpfs.GetCpuIdleTime()
	if err != nil {
		return err
	}

	iowait, err := lpfs.GetCpuIowaitTime()
	if err != nil {
		return err
	}

	irq, err := lpfs.GetCpuIrqTime()
	if err != nil {
		return err
	}

	softirq, err := lpfs.GetCpuSoftirqTime()
	if err != nil {
		return err
	}

	steal, err := lpfs.GetCpuStealTime()
	if err != nil {
		return err
	}

	total := user + nice + system + idle + iowait + irq + softirq + steal

	deltaIdle := idle - prevIdle
	deltaTotal := total - prevTotal
	cpuUsage = (1.0 - float64(deltaIdle)/float64(deltaTotal)) * 100.0

	prevIdle = idle
	prevTotal = total

	data := fmt.Sprintf("event: cpu\ndata: %f\n\n", cpuUsage)

	fmt.Fprintf(w, data)

	return nil
}
