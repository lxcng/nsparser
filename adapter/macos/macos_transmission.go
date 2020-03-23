package macos

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	def *MacosTransmissionAdapter = &MacosTransmissionAdapter{alive: false}
)

type MacosTransmissionAdapter struct {
	files []string
	mx    sync.RWMutex
	alive bool
}

func NewMacosTransmissionAdapter() *MacosTransmissionAdapter {
	def.mx.Lock()
	defer def.mx.Unlock()
	if !def.alive {
		def.alive = true
		go def.updateJob()
	}
	return def
}

func (ad *MacosTransmissionAdapter) updateJob() {
	log.Println("started scanning")
	for {
		ad.update()
		time.Sleep(time.Second * 10)
	}
}

func (ad *MacosTransmissionAdapter) GetLocal(title, translator string, episodeNumberRule *regexp.Regexp) map[int]struct{} {
	rp := strings.NewReplacer(
		" ", ".*",
		"-", ".*",
		":", ".*",
		";", ".*",
		".", ".*",
		",", ".*",
		"–", ".*",
	)
	rx := regexp.MustCompile(`.*` + rp.Replace(title) + `.*`)
	rx1 := regexp.MustCompile(translator)
	res := make(map[int]struct{})
	ad.mx.RLock()
	defer ad.mx.RUnlock()
	for _, n := range ad.files {
		if rx.MatchString(n) && rx1.MatchString(n) {
			numStr := episodeNumberRule.FindStringSubmatch(n)
			if len(numStr) < 2 {
				continue
			}
			numTmp, err := strconv.ParseInt(numStr[1], 10, 32)
			if err != nil {
				continue
			}
			res[int(numTmp)] = struct{}{}
		}
	}
	log.Printf("transmission scanner: found %d files for %s\n", len(res), title)
	return res
}

func (ad *MacosTransmissionAdapter) update() {
	files, err := ioutil.ReadDir(os.Getenv("HOME") + `/Library/Application Support/Transmission/Torrents`)
	if err != nil {
		log.Println("can't open directory:", err)
	}
	names := make([]string, len(files))
	for i, f := range files {
		names[i] = f.Name()
	}
	ad.mx.Lock()
	ad.files = names
	ad.mx.Unlock()
	log.Printf("transmission scanner: loaded %d files\n", len(files))
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
