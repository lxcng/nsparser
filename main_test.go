package main

import (
	"fmt"
	"nsparser/config"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	config.NewConf("config.json")
	time.Sleep(time.Second * 5)
	all := config.GetAll()
	fmt.Println(all)
	err := config.Start("Sword Art Online Alternative - Gun Gale Online")
	if err != nil {
		t.Fatal(err)
	}
}
