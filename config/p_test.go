package config

import (
	"fmt"
	"regexp"
	"testing"
)

func TestEp(t *testing.T) {
	r := ".*- ([0-9]+.*[0-9]*) "
	rg := regexp.MustCompile(r)
	res := rg.FindStringSubmatch("[HorribleSubs] Kakushigoto - 09 [1080p].mkv")
	fmt.Println(res)
}
