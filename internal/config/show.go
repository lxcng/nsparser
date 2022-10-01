package config

import (
	"fmt"
	"log"
	"math"
	"nsparser/internal/adapter"
	"nsparser/internal/parser"
	"sort"
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

	Missing []episode `json:","`

	Title string `json:","`

	parent Entry

	parser  *parser.Parser
	adapter adapter.Adapter
}

type episode struct {
	N       int
	Torrent string
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

func (s *show) start() error {
	for key, ep := range s.episodes {
		err := s.adapter.AddTorrent(ep, *s.cOpts.DownloadFolder)
		if err != nil {
			return err
		}
		delete(s.episodes, key)
		s.present[key] = struct{}{}
	}
	s.Present = s.compileEps(s.present)
	eps := make(map[int]struct{})
	for k, _ := range s.episodes {
		eps[k] = struct{}{}
	}
	s.eps = s.compileEps(eps)
	return nil
}

func (s *show) compile(e Entry) {
	s.parent = e
	s.id = s.Title
	s.episodes = map[int]string{}
	s.present = decompileEps(s.Present, s.Title)
	s.unfoldMissing()
	for current := Entry(s); current != nil; current = current.getParent() {
		s.pOpts.Merge(current.getParserOpts())
		s.cOpts.Merge(current.getClientOpts())
	}
	s.parser = parser.NewParser(s.pOpts)
	s.adapter = adapter.NewScanner(s.cOpts)
}

func (s *show) parse(wg *sync.WaitGroup) error {
	epsTmp := s.parser.ParsePage(s.Title)
	if epsTmp == nil {
		return fmt.Errorf("can't load episodes of \"%s\"", s.Title)
	}

	s.present = decompileEps(s.Present, s.Title)
	s.unfoldMissing()

	s.Present = s.compileEps(s.present)

	for epNum, torr := range epsTmp {
		if _, ok := s.present[epNum]; !ok {
			s.episodes[epNum] = torr
		} else {
			delete(s.episodes, epNum)
		}
	}
	eps := make(map[int]struct{})
	for k, _ := range s.episodes {
		eps[k] = struct{}{}
	}
	s.eps = s.compileEps(eps)
	wg.Done()
	log.Printf("%s: present: %s, ready: %s\n", s.Title, s.Present, func(s string) string {
		if s == "" {
			return "0"
		}
		return s
	}(s.eps))
	s.compileMissing()
	return nil
}

func (s *show) AppendPresent(pr string) {
	n := decompileEps(pr, s.Title)
	for k := range n {
		s.present[k] = struct{}{}
	}
}

// func (s *show) unfoldPresent() {
// 	res := map[int]struct{}{}
// 	if s.Present == "" {
// 		s.present = res
// 		return
// 	}
// 	ranges := strings.Split(s.Present, ", ")
// 	for _, r := range ranges {
// 		num, err := strconv.ParseInt(r, 10, 32)
// 		if err == nil {
// 			res[int(num)] = struct{}{}
// 		} else {
// 			borders := strings.Split(r, "-")
// 			if len(borders) != 2 {
// 				log.Fatalf("invalid Present for %s", s.Title)
// 			}
// 			min, err := strconv.ParseInt(borders[0], 10, 32)
// 			if err != nil {
// 				log.Fatalf("invalid Present for %s", s.Title)
// 			}
// 			max, err := strconv.ParseInt(borders[1], 10, 32)
// 			if err != nil {
// 				log.Fatalf("invalid Present for %s", s.Title)
// 			}
// 			for i := min; i <= max; i++ {
// 				res[int(i)] = struct{}{}
// 			}
// 		}
// 	}
// 	s.present = res
// }

func (s *show) unfoldMissing() {
	s.episodes = make(map[int]string)
	for _, ep := range s.Missing {
		s.episodes[ep.N] = ep.Torrent
	}
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

func decompileEps(eps, title string) map[int]struct{} {
	res := map[int]struct{}{}
	if eps == "" {
		return res
	}
	ranges := strings.Split(eps, ", ")
	for _, r := range ranges {
		num, err := strconv.ParseInt(r, 10, 32)
		if err == nil {
			res[int(num)] = struct{}{}
		} else {
			borders := strings.Split(r, "-")
			if len(borders) != 2 {
				log.Fatalf("invalid Present for %s", title)
			}
			min, err := strconv.ParseInt(borders[0], 10, 32)
			if err != nil {
				log.Fatalf("invalid Present for %s", title)
			}
			max, err := strconv.ParseInt(borders[1], 10, 32)
			if err != nil {
				log.Fatalf("invalid Present for %s", title)
			}
			for i := min; i <= max; i++ {
				res[int(i)] = struct{}{}
			}
		}
	}
	return res
}

func (s *show) compileMissing() {
	tmp := make([]episode, 0, len(s.episodes))
	for numEp, torr := range s.episodes {
		tmp = append(tmp, episode{numEp, torr})
	}
	sort.Slice(tmp, func(i, j int) bool { return tmp[i].N < tmp[j].N })
	s.Missing = tmp
}
