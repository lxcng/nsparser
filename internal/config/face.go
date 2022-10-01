package config

type Config struct {
	c *config
}

func (x *Config) Save() {
	x.c.Save()
}

func (x *Config) GetTranslators() []string {
	res := make([]string, 0, len(x.c.Translators))
	for _, tr := range x.c.Translators {
		res = append(res, *tr.Translator)
	}
	return res
}

func (x *Config) GetShows(tr int) []string {
	return x.findTranslator(tr).getShows()
}

func (x *Config) AddShow(i int, title, present string) error {
	return x.findTranslator(i).add(title, present)
}

func (x *Config) DeleteShow(i, title int) error {
	x.findTranslator(i).delete(title)
	return nil
}

func (x *Config) StartAll() error {
	return x.c.start()
}

func (x *Config) findTranslator(i int) *translator {
	return x.c.Translators[i]
}
