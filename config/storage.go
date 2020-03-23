package config

import (
	"sync"
	"time"
)

type storage struct {
	mx   sync.RWMutex
	conf *config
}

var st *storage

func NewConf(path string) {
	st = &storage{conf: &config{path: path}}
	go st.reloadJob()
}

func Start(id string) error {
	return st.start(id)
}

func Save() {
	st.save()
}

func GetAll() []*View {
	return st.getAll()
}

func (s *storage) save() {
	s.mx.RLock()
	defer s.mx.RUnlock()
	s.conf.save()
}

func (s *storage) start(id string) error {
	s.mx.RLock()
	defer s.mx.RUnlock()
	_, err := s.conf.start(id)
	return err
}

func (s *storage) getAll() []*View {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.conf.getViews()
}

func (s *storage) reloadJob() {
	for {
		s.mx.RLock()
		tmp := config{path: s.conf.path}
		s.mx.RUnlock()
		tmp.load()
		tmp.compile(nil)
		tmp.parse(nil)
		s.mx.Lock()
		s.conf = &tmp
		s.mx.Unlock()
		time.Sleep(time.Second * 10)
	}
}
