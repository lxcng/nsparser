package config

import (
	"nsparser/internal/adapter"
	"nsparser/internal/parser"
	"sync"
)

type Entry interface {
	getParent() Entry
	getChilds() []Entry
	getClientOpts() *adapter.ClientOpts
	getParserOpts() *parser.ParserOpts

	compile(Entry)
	start() error
	parse(*sync.WaitGroup) error
}

func wg() *sync.WaitGroup {
	wg := sync.WaitGroup{}
	wg.Add(1)
	return &wg
}
