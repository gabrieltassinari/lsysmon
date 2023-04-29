package data

import (
    "fmt"
    "encoding/json"

    "github.com/rprobaina/lpfs"
)

type MemoryJSON struct{
	Total int 
	Free int 
}

type SwapJSON struct{
	Filename string
	Size int
	Used int
}

func Memory() (string, error) {
	total, err := lpfs.GetMemTotal()
	if err != nil {
		return "", err
	}

	free, err := lpfs.GetMemFree()
	if err != nil {
		return "", err
	}

	msg := MemoryJSON {
		Total: total, 
		Free: free, 
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	data := fmt.Sprintf("event: memory\ndata: %v\n\n", string(b))

	return data, nil
}

func Swap() (string, error) {
	filename, err := lpfs.GetSwapFilename()
	if err != nil {
		return "", err
	}

	size, err := lpfs.GetSwapSize()
	if err != nil {
		return "", err
	}

	used, err := lpfs.GetSwapUsed()
	if err != nil {
		return "", err
	}

	msg := SwapJSON {
		Filename: filename,
		Size: size,
		Used: used,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	data := fmt.Sprintf("event: swap\ndata: %v\n\n", string(b))

	return data, nil
}	

func Uptime() (string, error) {
	uptime, err := lpfs.GetUptimeSystem()
	if err != nil {
		return "", err
	}

	data := fmt.Sprintf("event: uptime\ndata: %f\n\n", uptime)

	return data, nil
}
