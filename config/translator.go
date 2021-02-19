package config

import (
	"nsparser/adapter"
	"nsparser/parser"
	"sync"
)

type translator struct {
	parser.ParserOpts
	adapter.ClientOpts

	Shows []*show `json:","`

	mx     sync.RWMutex
	parent Entry
}

func (t *translator) getParent() Entry {
	return t.parent
}

func (t *translator) getChilds() []Entry {
	res := make([]Entry, len(t.Shows))
	for i, s := range t.Shows {
		res[i] = Entry(s)
	}
	return res
}

func (t *translator) getParserOpts() *parser.ParserOpts {

	return &t.ParserOpts
}

func (t *translator) getClientOpts() *adapter.ClientOpts {

	return &t.ClientOpts
}

func (t *translator) start(id string) (bool, error) {
	t.mx.RLock()
	defer t.mx.RUnlock()
	for _, e := range t.getChilds() {
		ok, err := e.start(id)
		if ok {
			return ok, err
		}
	}
	return false, nil
}

func (t *translator) compile(e Entry, sf func()) {
	t.mx.Lock()
	defer t.mx.Unlock()
	t.parent = e
	for _, e := range t.getChilds() {
		e.compile(t, sf)
	}
}

func (t *translator) parse(wg *sync.WaitGroup) error {
	t.mx.Lock()
	defer t.mx.Unlock()
	nwg := sync.WaitGroup{}
	for _, e := range t.getChilds() {
		nwg.Add(1)
		go e.parse(&nwg)
	}
	nwg.Wait()
	wg.Done()
	return nil
}

func (t *translator) getViews() []*View {
	t.mx.RLock()
	defer t.mx.RUnlock()
	res := []*View{}
	for _, e := range t.getChilds() {
		res = append(res, e.getViews()...)
	}
	return res
}
