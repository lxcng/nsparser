package config

import (
	"nsparser/internal/parser"
)

type translator struct {
	parser.ParserOpts

	Shows  []*show `json:","`
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

func (t *translator) start() error {
	for _, e := range t.getChilds() {
		err := e.start()
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *translator) compile(e Entry) {
	t.parent = e
	for _, e := range t.getChilds() {
		e.compile(t)
	}
}

func (t *translator) parse() error {
	for _, e := range t.getChilds() {
		e.parse()
	}
	return nil
}

func (x *translator) getShows() []string {
	res := make([]string, 0, len(x.Shows))
	for _, s := range x.Shows {
		res = append(res, s.Title)
	}
	return res
}

func (x *translator) add(title, present string) error {
	for _, s := range x.Shows {
		if s.Title == title {
			s.AppendPresent(present)
			return s.parse()
		}
	}
	s := &show{
		Title:   title,
		Present: present,
	}
	s.compile(x)
	err := s.parse()
	if err != nil {
		return err
	}

	x.Shows = append(x.Shows, s)
	return nil
}

func (x *translator) delete(i int) {
	copy(x.Shows[i:], x.Shows[i+1:])
	x.Shows = x.Shows[:len(x.Shows)-1]
}
