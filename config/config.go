package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"nsparser/adapter"
	"nsparser/parser"
	"sync"
)

type config struct {
	parser.ParserOpts
	adapter.ClientOpts

	Translators []*translator `json:","`

	mx   sync.RWMutex
	path string

	parser *parser.Parser
}

func (c *config) load() {
	bt, err := ioutil.ReadFile(c.path)
	if err != nil {
		log.Println("can't load config file", err)
		return
	}
	c.mx.Lock()
	defer c.mx.Unlock()
	err = json.Unmarshal(bt, c)
	if err != nil {
		log.Println("can't unmarshal config", err)
		return
	}
	log.Println("config loaded")
}

func (c *config) save() {
	c.mx.RLock()
	defer c.mx.RUnlock()
	buff := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buff)
	enc.SetIndent("", "\t")
	err := enc.Encode(c)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(c.path, buff.Bytes(), 0777)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("config saved")
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

func (c *config) getClientOpts() *adapter.ClientOpts {
	return &c.ClientOpts
}

func (c *config) start(id string) (bool, error) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	for _, e := range c.getChilds() {
		ok, err := e.start(id)
		if ok {
			return ok, err
		}
	}
	return false, nil
}

func (c *config) compile(Entry, func()) {
	c.mx.Lock()
	defer c.mx.Unlock()
	for _, e := range c.getChilds() {
		e.compile(c, c.save)
	}
}

func (c *config) parse(wg *sync.WaitGroup) error {
	c.mx.Lock()
	defer c.mx.Unlock()
	nwg := sync.WaitGroup{}
	for _, e := range c.getChilds() {
		nwg.Add(1)
		go e.parse(&nwg)
	}
	nwg.Wait()
	log.Println("parsed")
	return nil
}

func (c *config) getViews() []*View {
	c.mx.RLock()
	defer c.mx.RUnlock()
	res := []*View{}
	for _, e := range c.getChilds() {
		res = append(res, e.getViews()...)
	}
	return res
}
