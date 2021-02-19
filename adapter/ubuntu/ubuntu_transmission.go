package ubuntu

import (
	"log"
	"os/exec"
)

var (
	def *UbuntuTransmissionAdapter = &UbuntuTransmissionAdapter{}
)

type UbuntuTransmissionAdapter struct {
}

func NewUbuntuTransmissionAdapter() *UbuntuTransmissionAdapter {

	return def
}

func (tr *UbuntuTransmissionAdapter) AddMagnet(magnet, dir string) error {
	cmd := exec.Command("transmission-gtk", "-m", magnet)
	bt, err := cmd.Output()
	if err != nil {
		log.Println(err.Error(), string(bt))
		return err
	}
	log.Printf("started magnet: %s\n", magnet)
	return nil
}
