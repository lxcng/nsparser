package adapter

import (
	"nsparser/adapter/macos"
	"nsparser/adapter/ubuntu"
)

type Adapter interface {
	AddMagnet(magnet, dir string) error
}

var (
	scanners = map[int]int{}
)

func NewScanner(co ClientOpts) Adapter {
	switch *co.Os {
	case "macos":
		return macos.NewMacosTransmissionAdapter()
	case "ubuntu":
		return ubuntu.NewUbuntuTransmissionAdapter()
	default:
		return nil
	}
}
