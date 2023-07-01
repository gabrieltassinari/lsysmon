package runtime

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rprobaina/lpfs"
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
