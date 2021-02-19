package config

import (
	"nsparser/adapter"
	"nsparser/parser"
	"sync"
)

type Entry interface {
	getParent() Entry
	getChilds() []Entry
	getClientOpts() *adapter.ClientOpts
	getParserOpts() *parser.ParserOpts

	compile(Entry, func())
	start(string) (bool, error)
	parse(*sync.WaitGroup) error

	getViews() []*View
}
