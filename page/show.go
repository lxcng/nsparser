package page

import (
	"log"
	"time"

	"golang.org/x/net/html"
)

func GetShow(title string) []byte {
	start := time.Now()
	root, body := makeBody()
	table := &html.Node{
		Type: html.ElementNode,
		Data: "table",
	}
	body.AppendChild(table)

	table.AppendChild(makeThead("Title", "Present", "Missing"))
	table.AppendChild(makeTbody())

	res := render(root)
	log.Printf("show '%s' rendered in %v nsec\n", title, time.Now().Sub(start).Nanoseconds())
	return res
}
