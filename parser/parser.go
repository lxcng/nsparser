package parser

import (
	"bytes"
	"log"
	"regexp"
	"strconv"
	"time"

	"golang.org/x/net/html"
)

type Parser struct {
	resolution          int
	episodeNumberRegexp string
	EpisodeNumberRule   *regexp.Regexp
	magnetRegexp        string
	MagnetRule          *regexp.Regexp
}

func NewParser(po ParserOpts) *Parser {
	res := &Parser{
		resolution:          *po.Resolution,
		episodeNumberRegexp: *po.EpisodeNumberRegex,
		magnetRegexp:        *po.MagnetRegex,
	}
	var err error
	res.EpisodeNumberRule, err = regexp.Compile(res.episodeNumberRegexp)
	if err != nil {
		log.Fatal(err)
	}
	res.MagnetRule, err = regexp.Compile(res.magnetRegexp)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func (p *Parser) ParsePage(user, title string) map[int]string {
	start := time.Now()
	res := map[int]string{}
	for i := 1; i < 100; i++ {
		eps := p.parsePage(user, title, i)
		if eps == nil || len(eps) == 0 {
			break
		}
		for n, m := range eps {
			res[n] = m
		}
	}
	log.Printf("parsed %d episodes of \"%s\" by %s in %v\n", len(res), title, user, time.Now().Sub(start).String())
	return res
}

func (p *Parser) parsePage(user, title string, page int) map[int]string {
	bt, err := getPage(user, title, p.resolution, page)
	if err != nil {
		log.Fatal(err)
	}
	doc, err := html.Parse(bytes.NewReader(bt))
	if err != nil {
		log.Fatal(err)
	}

	tbody := p.lookupTable(doc)
	if tbody == nil {
		return nil
	}
	res := p.lookupMagnets(tbody)
	return res
}

func (p *Parser) lookupTable(n *html.Node) *html.Node {
	htm := lookupElement(n, "html", html.ElementNode, [][2]string{{"lang", "en"}}, nil)
	if htm == nil {
		return nil
	}
	body := lookupElement(htm, "body", html.ElementNode, nil, nil)
	if body == nil {
		return nil
	}
	container := lookupElement(body, "div", html.ElementNode, [][2]string{{"class", "container"}}, nil)
	if container == nil {
		return nil
	}
	containerRow := lookupElement(container, "div", html.ElementNode, [][2]string{{"class", "row"}}, nil)
	if containerRow == nil {
		return nil
	}
	tableResponsive := lookupElement(containerRow, "div", html.ElementNode, [][2]string{{"class", "table-responsive"}}, nil)
	if tableResponsive == nil {
		return nil
	}
	table := lookupElement(tableResponsive, "table", html.ElementNode, nil, nil)
	if table == nil {
		return nil
	}
	tbody := lookupElement(table, "tbody", html.ElementNode, nil, nil)
	if tbody == nil {
		return nil
	}
	return tbody
}

func lookupElement(n *html.Node, data string, tp html.NodeType, attr [][2]string, attrBl [][2]string) *html.Node {
	var res *html.Node
	for i := n.FirstChild; i != nil; i = i.NextSibling {
		if i.Data == data && i.Type == tp {
			ok := true
			if attr != nil && i.Attr != nil {
				for _, at := range attr {
					present := false
					for _, iat := range i.Attr {
						if at[0] == iat.Key {
							if at[1] == "" || at[1] == iat.Val {
								present = true
								break
							}
						}
					}
					if !present {
						ok = false
					}
				}
				for _, at := range attrBl {
					present := false
					for _, iat := range i.Attr {
						if at[0] == iat.Key {
							if at[1] == "" || at[1] == iat.Val {
								present = true
								break
							}
						}
					}
					if present {
						ok = false
					}
				}
			}
			if ok {
				res = i
				break
			}
		}
	}
	return res
}

func (p *Parser) lookupMagnets(n *html.Node) map[int]string {
	res := map[int]string{}
	for i := n.FirstChild; i != nil; i = i.NextSibling {
		if i.Data == "tr" {
			if i.Attr != nil {
				for _, iat := range i.Attr {
					if iat.Key == "class" {
						num, mag, ok := p.lookupMagnet(i)
						if ok {
							res[num] = mag
						}
					}
				}
			}
		}
	}
	return res
}

func (p *Parser) lookupMagnet(n *html.Node) (int, string, bool) {
	num := 0
	mag := ""
	for i := n.FirstChild; i != nil; i = i.NextSibling {
		if i.Data == "td" {
			if i.Attr != nil {
				for _, iat := range i.Attr {
					if iat.Key == "colspan" {
						if iat.Val == "2" {
							numTmp, ok := p.lookupEpisodeNumber(i)
							if !ok {
								return 0, "", false
							}
							num = numTmp
						}
					}
					if iat.Key == "class" {
						if iat.Val == "text-center" {
							if i.FirstChild != nil {
								magTmp, ok := p.lookupMagnetUrl(i)
								if ok {
									mag = magTmp
								}
							}
						}
					}
				}
			}
		}
	}
	return num, mag, true
}

func (p *Parser) lookupEpisodeNumber(n *html.Node) (int, bool) {
	href := lookupElement(n, "a", html.ElementNode, [][2]string{{"title", ""}}, [][2]string{{"class", "comments"}})
	if href == nil {
		return 0, false
	}
	for _, attr := range href.Attr {
		if attr.Key == "title" {
			numStr := p.EpisodeNumberRule.FindStringSubmatch(attr.Val)
			if len(numStr) < 2 {
				return 0, false
			}
			numTmp, err := strconv.ParseInt(numStr[1], 10, 32)
			if err != nil {
				return 0, false
			}
			return int(numTmp), true
		}
	}
	return 0, true
}

func (p *Parser) lookupMagnetUrl(n *html.Node) (string, bool) {
	for i := n.FirstChild; i != nil; i = i.NextSibling {
		if i.Data == "a" {
			if i.Attr != nil {
				for _, iat := range i.Attr {
					if iat.Key == "href" {
						if p.MagnetRule.MatchString(iat.Val) {
							return iat.Val, true
						}
					}
				}
			}
		}
	}
	return "", false
}
