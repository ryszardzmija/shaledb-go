package membership

import (
	"encoding/json"
	"errors"
	"io"
)

type File struct {
	Nodes []Node `json:"nodes"`
}

type Node struct {
	ID      string      `json:"id"`
	Address NodeAddress `json:"address"`
}

type NodeAddress struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
}

func Load(r io.Reader) (File, error) {
	var cfg File

	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cfg); err != nil {
		return File{}, err
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return File{}, errors.New("membership config must contain a single JSON object")
		}

		return File{}, err
	}

	return cfg, nil
}
