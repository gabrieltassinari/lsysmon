package data

import (
	"encoding/json"
	"net/http"

	"github.com/rprobaina/lpfs"
)

type SwapLabel struct {
	Name string
	Size int
}

type MemoryLabel struct {
	Total int
}

type Json struct {
	SwapLabel SwapLabel
	MemoryLabel MemoryLabel
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

	msg := Json {
		SwapLabel {
			Name: sn,
			Size: ss,
		},
		MemoryLabel {
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
