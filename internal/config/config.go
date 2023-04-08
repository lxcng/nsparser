package config

import (
	"bytes"
	"encoding/json"
	"log"

	"nsparser/internal/parser"
	"os"
	"sync"
)

type config struct {
	parser.ParserOpts

	Translators []*translator `json:","`
	path        string
}

func NewConfig(path string) *Config {
	c := &config{path: path}
	c.load()
	c.compile(nil)
	return &Config{c: c}
}

func (c *config) load() {
	bt, err := os.ReadFile(c.path)
	if err != nil {
		log.Println("can't load config file", err)
		return
	}
	err = json.Unmarshal(bt, c)
	if err != nil {
		log.Println("can't unmarshal config", err)
		return
	}
	// fmt.Println("config loaded")
}

func (c *config) Save() {
	buff := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buff)
	enc.SetIndent("", "\t")
	err := enc.Encode(c)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(c.path, buff.Bytes(), 0777)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println("config saved")
}

// int funcs

func (c *config) getParent() Entry {
	return nil
}

func (c *config) getChilds() []Entry {
	res := make([]Entry, len(c.Translators))
	for i, t := range c.Translators {
		res[i] = Entry(t)
	}
	return res
}

func (c *config) getParserOpts() *parser.ParserOpts {
	return &c.ParserOpts
}

func (c *config) start() error {
	for _, e := range c.getChilds() {
		err := e.start()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *config) compile(Entry) {
	for _, e := range c.getChilds() {
		e.compile(c)
	}
}

func (c *config) parse(wg *sync.WaitGroup) error {
	nwg := sync.WaitGroup{}
	for _, e := range c.getChilds() {
		nwg.Add(1)
		go e.parse(&nwg)
	}
	nwg.Wait()
	// log.Println("parsed")
	return nil
}
