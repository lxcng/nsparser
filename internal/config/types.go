package config

import (
	"nsparser/internal/parser"
)

type Entry interface {
	getParent() Entry
	getChilds() []Entry
	getParserOpts() *parser.ParserOpts

	compile(Entry)
	start() error
	parse() error
}
