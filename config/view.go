package config

type View struct {
	Id       string
	Episodes string `json:","`
	Present  string `json:","`
	Title    string `json:","`
}
