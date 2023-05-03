package data

import (
	"encoding/json"
	"net/http"

	"github.com/rprobaina/lpfs"
)

type swapLabels struct {
	Name string
	Size int
}

type memoryLabels struct {
	Total int
}

type jsonLabels struct {
	SwapLabel swapLabels	 `json:"swap"`
	MemoryLabel memoryLabels `json:"memory"`
}

func Labels(w http.ResponseWriter) error {
	mt, err := lpfs.GetMemTotal()
	if err != nil {
		return err
	}

	sn, err := lpfs.GetSwapFilename()
	if err != nil {
		return err
	}

	ss, err := lpfs.GetSwapSize()
	if err != nil {
		return err
	}

	msg := jsonLabels {
		swapLabels {
			Name: sn,
			Size: ss,
		},
		memoryLabels {
			Total: mt,
		},
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	w.Write(b)
	return nil
}
