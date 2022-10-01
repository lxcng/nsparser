package adapter

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func AddTorrent(url string) error {
	cmd := exec.Command(
		"/bin/sh", "-c",
		fmt.Sprintf(
			"%s; %s; %s;",
			wget(url),
			addTorrent(trim(url)),
			rm(trim(url)),
		),
	)
	bt, err := cmd.Output()
	if err != nil {
		log.Println(err.Error(), string(bt))
		return err
	}
	return nil
}

func Flush() error {
	cmd := exec.Command(
		"transmission-remote", "-l",
	)
	bt, err := cmd.Output()
	if err != nil {
		log.Println(err.Error(), string(bt))
		return err
	}
	//
	lines := strings.Split(string(bt), "\n")
	ids := make([]int64, 0)
	for i := 1; i < len(lines)-2; i++ {
		l := strings.TrimLeft(lines[i], " ")
		ts := strings.Split(l, " ")
		id, err := strconv.ParseInt(ts[0], 10, 64)
		if err != nil {
			return err
		}
		ids = append(ids, id)
	}
	//
	cmd = exec.Command(
		"/bin/sh", "-c",
		removeTorrents(ids),
	)
	bt, err = cmd.Output()
	if err != nil {
		log.Println(err.Error(), string(bt))
		return err
	}
	return nil
}

func wget(url string) string {
	return fmt.Sprintf("wget %s", url)
}

func addTorrent(path string) string {
	return fmt.Sprintf("transmission-remote -a %s", path)
}

func removeTorrents(ids []int64) string {
	res := ""
	for _, id := range ids {
		res += fmt.Sprintf("transmission-remote -t %v -r; ", id)
	}
	return res
}
func rm(path string) string {
	return fmt.Sprintf("rm %s", path)
}

func trim(url string) string {
	s := strings.Split(url, "/")
	return s[len(s)-1]
}
