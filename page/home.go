package page

import (
	"bytes"
	"log"
	"nsparser/config"
	"sync"
	"time"

	"golang.org/x/net/html"
)

var (
	mx sync.RWMutex
)

func GetHome() []byte {
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
	log.Printf("home rendered in %v nsec\n", time.Now().Sub(start).Nanoseconds())
	return res
}

func render(n *html.Node) []byte {
	buff := bytes.NewBuffer(nil)
	err := html.Render(buff, n)
	if err != nil {
		log.Fatal(err)
	}
	return buff.Bytes()
}

func makeBody() (*html.Node, *html.Node) {
	root := &html.Node{
		Type: html.DocumentNode,
	}
	root.AppendChild(&html.Node{
		Type: html.DoctypeNode,
		Data: "doctype",
	})
	htm := &html.Node{
		Type: html.ElementNode,
		Data: "html",
		Attr: []html.Attribute{{"", "lang", "en"}},
	}
	root.AppendChild(htm)
	htm.AppendChild(&html.Node{
		Type: html.ElementNode,
		Data: "head",
	})
	body := &html.Node{
		Type: html.ElementNode,
		Data: "body",
	}
	htm.AppendChild(body)
	return root, body
}

func makeThead(labels ...string) *html.Node {
	thead := &html.Node{
		Type: html.ElementNode,
		Data: "thead",
	}
	tr := &html.Node{
		Type: html.ElementNode,
		Data: "tr",
	}
	thead.AppendChild(tr)
	for _, l := range labels {
		th := &html.Node{
			Type: html.ElementNode,
			Data: "th",
		}
		th.AppendChild(makeTextNode(l))
		tr.AppendChild(th)
	}
	return thead
}

func makeTbody() *html.Node {
	tbody := &html.Node{
		Type: html.ElementNode,
		Data: "tbody",
	}
	for _, d := range config.GetAll() {
		tbody.AppendChild(makeDisplay(d))
	}
	return tbody
}

func makeDisplay(d *config.View) *html.Node {
	res := &html.Node{
		Type: html.ElementNode,
		Data: "tr",
	}
	title := &html.Node{
		Type: html.ElementNode,
		Data: "td",
	}
	title.AppendChild(makeTextNode(d.Title))
	res.AppendChild(title)

	loaded := &html.Node{
		Type: html.ElementNode,
		Data: "td",
	}
	loaded.AppendChild(makeTextNode(d.Present))
	res.AppendChild(loaded)

	missing := &html.Node{
		Type: html.ElementNode,
		Data: "td",
	}
	missing.AppendChild(makeRefNode("api/load/"+d.Id, d.Episodes))
	res.AppendChild(missing)

	return res
}

func makeTextNode(t string) *html.Node {
	return &html.Node{
		Type: html.TextNode,
		Data: t,
	}
}

func makeRefNode(url, text string) *html.Node {
	res := &html.Node{
		Type: html.ElementNode,
		Data: "a",
		Attr: []html.Attribute{
			html.Attribute{
				Namespace: "",
				Key:       "href",
				Val:       url,
			},
			html.Attribute{
				Namespace: "",
				Key:       "title",
				Val:       text,
			},
		},
	}
	res.AppendChild(makeTextNode(text))
	return res
}
