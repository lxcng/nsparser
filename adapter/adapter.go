package adapter

import (
	"nsparser/adapter/macos"
	"nsparser/adapter/ubuntu"
	"regexp"
)

type Adapter interface {
	GetLocal(title, translator string, episodeNumberRule *regexp.Regexp) map[int]struct{}
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
