package adapter

import (
	"nsparser/internal/adapter/ubuntu"
)

type Adapter interface {
	AddTorrent(url, dir string) error
}

var (
	scanners = map[int]int{}
)

func NewScanner(co ClientOpts) Adapter {
	switch *co.Os {
	// case "macos":
	// return macos.NewMacosTransmissionAdapter()
	case "ubuntu":
		return ubuntu.NewUbuntuTransmissionAdapter()
	default:
		return nil
	}
}
