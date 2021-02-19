package macos

import (
	"log"
	"os/exec"
)

var (
	def *MacosTransmissionAdapter = &MacosTransmissionAdapter{}
)

type MacosTransmissionAdapter struct {
}

func NewMacosTransmissionAdapter() *MacosTransmissionAdapter {

	return def
}

func (tr *MacosTransmissionAdapter) AddMagnet(magnet, dir string) error {
	cmd := exec.Command("transmission-remote", "-a", magnet, "-w", dir)
	bt, err := cmd.Output()
	if err != nil {
		log.Println(err.Error(), string(bt))
		return err
	}
	log.Printf("started magnet: %s\n", magnet)
	return nil
}
