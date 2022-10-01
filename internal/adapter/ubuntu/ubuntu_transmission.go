package ubuntu

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

var (
	def *UbuntuTransmissionAdapter = &UbuntuTransmissionAdapter{}
)

type UbuntuTransmissionAdapter struct {
}

func NewUbuntuTransmissionAdapter() *UbuntuTransmissionAdapter {

	return def
}

func (tr *UbuntuTransmissionAdapter) AddTorrent(url, dir string) error {
	cmd := exec.Command(
		"/bin/sh", "-c",
		fmt.Sprintf(
			"%s; %s; %s;",
			tr.wget(url),
			tr.addTorrent(trim(url), dir),
			tr.rm(trim(url)),
		),
	)
	bt, err := cmd.Output()
	if err != nil {
		log.Println(err.Error(), string(bt))
		return err
	}
	return nil
}

func (tr *UbuntuTransmissionAdapter) wget(url string) string {
	return fmt.Sprintf("wget %s", url)
}

func (tr *UbuntuTransmissionAdapter) addTorrent(path, dir string) string {
	return fmt.Sprintf("transmission-remote -a %s", path)
}
func (tr *UbuntuTransmissionAdapter) rm(path string) string {
	return fmt.Sprintf("rm %s", path)
}

func trim(url string) string {
	s := strings.Split(url, "/")
	return s[len(s)-1]
}
