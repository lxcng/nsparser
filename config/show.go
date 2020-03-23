package config

import (
	"fmt"
	"log"
	"math"
	"nsparser/adapter"
	"nsparser/parser"
	"strconv"
	"strings"
	"sync"
)

type show struct {
	parser.ParserOpts  `json:",omitempty"`
	adapter.ClientOpts `json:",omitempty"`

	pOpts parser.ParserOpts
	cOpts adapter.ClientOpts

	id string

	episodes map[int]string
	eps      string

	Present string `json:","`
	present map[int]struct{}

	Title string `json:","`

	mx sync.RWMutex

	parent Entry

	parser  *parser.Parser
	adapter adapter.Adapter
}

func (s *show) getParent() Entry {
	return s.parent
}

func (s *show) getChilds() []Entry {
	return nil
}

func (s *show) getParserOpts() *parser.ParserOpts {
	return &s.ParserOpts
}

func (s *show) getClientOpts() *adapter.ClientOpts {
	return &s.ClientOpts
}

func (s *show) start(id string) (bool, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	if s.id == id {
		for _, ep := range s.episodes {
			err := s.adapter.AddMagnet(ep, *s.cOpts.DownloadFolder)
			if err != nil {
				return true, err
			}
		}
		return true, nil
	}
	return false, nil
}

func (s *show) compile(e Entry) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.parent = e
	s.id = s.Title
	s.episodes = map[int]string{}
	for current := Entry(s); current != nil; current = current.getParent() {
		s.pOpts.Merge(current.getParserOpts())
		s.cOpts.Merge(current.getClientOpts())
	}
	s.parser = parser.NewParser(s.pOpts)
	s.adapter = adapter.NewScanner(s.cOpts)
}

func (s *show) parse(wg *sync.WaitGroup) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	epsTmp := s.parser.ParsePage(*s.pOpts.Translator, s.Title)
	if epsTmp == nil {
		return fmt.Errorf("can't load episodes of \"%s\"", s.Title)
	}

	s.unfoldPresent()
	tmp := s.adapter.GetLocal(s.Title, *s.pOpts.Translator, s.parser.MagnetRule)
	for k, v := range tmp {
		s.present[k] = v
	}
	s.Present = s.compileEps(s.present)

	for n, m := range epsTmp {
		if _, ok := s.present[n]; !ok {
			s.episodes[n] = m
		} else {
			delete(s.episodes, n)
		}
	}
	eps := make(map[int]struct{})
	for k, _ := range s.episodes {
		eps[k] = struct{}{}
	}
	s.eps = s.compileEps(eps)
	wg.Done()
	log.Printf("%s: present: %s, ready: %s\n", s.Title, s.Present, s.eps)
	return nil
}

func (s *show) getViews() []*View {
	eps := make(map[int]struct{})
	for k, _ := range s.episodes {
		eps[k] = struct{}{}
	}
	res := View{
		Id:       s.id,
		Episodes: s.eps,
		Present:  s.Present,
		Title:    s.Title,
	}
	return []*View{&res}
}

func (s *show) unfoldPresent() {
	res := map[int]struct{}{}
	if s.Present == "" {
		return
	}
	ranges := strings.Split(s.Present, ", ")
	for _, r := range ranges {
		num, err := strconv.ParseInt(r, 10, 32)
		if err == nil {
			res[int(num)] = struct{}{}
		} else {
			borders := strings.Split(r, "-")
			if len(borders) != 2 {
				log.Fatalf("invalid Present for %s", s.Title)
			}
			min, err := strconv.ParseInt(borders[0], 10, 32)
			if err != nil {
				log.Fatalf("invalid Present for %s", s.Title)
			}
			max, err := strconv.ParseInt(borders[1], 10, 32)
			if err != nil {
				log.Fatalf("invalid Present for %s", s.Title)
			}
			for i := min; i <= max; i++ {
				res[int(i)] = struct{}{}
			}
		}
	}
	s.present = res
}

func (s *show) compileEps(eps map[int]struct{}) string {
	res := ""
	min, max := math.MaxInt32, 0
	for n := range eps {
		if n > max {
			max = n
		}
		if n < min {
			min = n
		}
	}
	curr := ""
	start := 0
	for i := min; i <= max+1; i++ {
		if _, ok := eps[i]; ok {
			if curr == "" {
				curr += fmt.Sprint(i)
				start = i
			}
		} else {
			if curr != "" {
				if start != i-1 {
					curr += fmt.Sprintf("-%d", i-1)
				}
				if res != "" {
					res += ", "
				}
				res += curr
				curr = ""
			}
		}
	}
	return res
}
