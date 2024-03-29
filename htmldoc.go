// Package exhtml provides functions to extract `*html.Node`, raw `[]byte` contents, links from a website.
package exhtml

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

func request(src string) (*http.Response, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", src, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36")
	return client.Do(req)
}

// GetRawAndDoc can get html raw bytes and html.Node by rawurl.
func GetRawAndDoc(url *url.URL, retryTimeout time.Duration) ([]byte, *html.Node, error) {
	// Get response form url
	deadline := time.Now().Add(retryTimeout)
	for tries := 0; time.Now().Before(deadline); tries++ {
		resp, err := request(url.String())
		if err == nil { // success
			defer resp.Body.Close()
			raw, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, nil, errors.WithMessage(
					err, "exhtml: GetRawAndDoc: ReadAll",
				)
			}
			doc, err := html.Parse(bytes.NewBuffer(raw))
			return raw, doc, nil
		}
		log.SetPrefix("[wait]")
		log.SetFlags(0)
		log.Printf("server not responding (%s); retrying...", err)
		log.SetPrefix("")
		log.SetFlags(3)
		time.Sleep(time.Second << uint(tries)) // exponential back-off
	}
	return nil, nil, nil
}

// ExtractRssGuids get value from <guid>
func ExtractRssGuids(weburl string) ([]string, error) {
	resp, err := request(weburl)
	if err != nil {
		return nil, err
	}
	gf := gofeed.NewParser()
	feed, err := gf.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	ls := []string{}
	for _, e := range feed.Items {
		ls = append(ls, e.GUID)
	}
	return ls, nil
}

func ExtractRss(weburl string) ([]string, error) {
	resp, err := request(weburl)
	if err != nil {
		return nil, err
	}
	gf := gofeed.NewParser()
	feed, err := gf.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	ls := []string{}
	for _, e := range feed.Items {
		ls = append(ls, e.Link)
	}
	return ls, nil
}

// ExtractLinks makes an HTTP GET request to the specified URL, parses
// the response as HTML, and returns the links in the HTML document.
func ExtractLinks(weburl string) ([]string, error) {
	u, err := url.Parse(weburl)
	if err != nil {
		return nil, err
	}
	resp, err := request(weburl)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", weburl, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", weburl, err)
	}
	var links []string
	visitNode := func(n *html.Node) {
		// TODO: compress layers
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				// append only the target website
				if strings.HasPrefix(a.Val, "http") && strings.Contains(a.Val, u.Hostname()) {
					links = append(links, link.String())
				} else if strings.HasPrefix(a.Val, "/") {
					links = append(links, link.String())
				}

			}
		}
	}
	ForEachNode(doc, visitNode, nil)
	return links, nil
}

func ForEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ForEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

func ElementsByTagAttr(doc *html.Node, tagName, attrName, attrValue string) []*html.Node {
	var nodes []*html.Node
	if tagName == "" || doc == nil {
		return nil
	}
	if doc.Type == html.ElementNode {
		if tagName == doc.Data {
			for _, a := range doc.Attr {
				if attrName == "" {
					nodes = append(nodes, doc)
				}
				if attrName != "" && attrValue == "" {
					if a.Key == attrName {
						nodes = append(nodes, doc)
					}
				}
				if attrName != "" && attrValue != "" {
					if a.Key == attrName && a.Val == attrValue {
						nodes = append(nodes, doc)
					}
				}
			}
		}
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, ElementsByTagAttr(c, tagName, attrName, attrValue)...)
	}
	return nodes
}

func DivWithAttr(doc *html.Node, attrName, attrValue string) []*html.Node {
	var nodes []*html.Node
	if attrName == "" || attrValue == "" || doc == nil {
		return nil
	}
	if doc.Type == html.ElementNode {
		if "div" == doc.Data {
			for _, a := range doc.Attr {
				if a.Key == attrName && a.Val == attrValue {
					nodes = append(nodes, doc)
				}
			}
		}
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, DivWithAttr(c, attrName, attrValue)...)
	}
	return nodes
}

func DivWithAttr2(raw []byte, attrName, attrValue string) []byte {
	if attrName == "" || attrValue == "" || raw == nil {
		return nil
	}
	z := html.NewTokenizer(bytes.NewReader(raw))
	for {
		tt := z.Next()
		t := z.Token()
		if err := z.Err(); err != nil && err == io.EOF {
			break
		}
		switch tt {
		case html.StartTagToken:
			if "div" == t.Data {
				for _, a := range t.Attr {
					if a.Key == attrName && a.Val == attrValue {
						return z.Buffered()
					}
				}
			}
		}
	}
	return nil
}

func ElementsNext(doc *html.Node) []*html.Node {
	nodes := []*html.Node{}
	if doc == nil {
		return nil
	}
	visitNode := func(n *html.Node) {
		if n.NextSibling != nil {
			nodes = append(nodes, n)
		}
	}
	ForEachNode(doc, visitNode, nil)
	return nodes
}

func ElementsRmByTag(doc *html.Node, name ...string) {
	if len(name) == 0 || doc == nil {
		return
	}
	visitNode := func(n *html.Node) {
		if n.NextSibling != nil && n.NextSibling.Type == html.ElementNode {
			for _, tag := range name {
				if tag == n.NextSibling.Data {
					n.Parent.RemoveChild(n.NextSibling)
				}
			}
		}
	}
	ForEachNode(doc, visitNode, nil)
	rmFirstTag := func(n *html.Node) {
		if n.FirstChild != nil && n.FirstChild.Type == html.ElementNode {
			for _, tag := range name {
				if tag == n.FirstChild.Data {
					n.RemoveChild(n.FirstChild)
				}
			}
		}
	}
	ForEachNode(doc, rmFirstTag, nil)
}

// ElementsRmByTagClass rm nodes
// if class != "" rm by tag and class else rm just by tag.
func ElementsRmByTagClass(doc *html.Node, tag, class string) {
	if tag == "" || doc == nil {
		return
	}
	visitNode := func(n *html.Node) {
		if n.NextSibling != nil && n.NextSibling.Type == html.ElementNode {
			if tag == n.NextSibling.Data {
				if class != "" {
					for _, a := range n.NextSibling.Attr {
						if a.Key == "class" && a.Val == class {
							n.Parent.RemoveChild(n.NextSibling)
						}
					}
				} else {
					n.Parent.RemoveChild(n.NextSibling)
				}
			}
		}
	}
	ForEachNode(doc, visitNode, nil)
	rmFirstTag := func(n *html.Node) {
		if n.FirstChild != nil && n.FirstChild.Type == html.ElementNode {
			if tag == n.FirstChild.Data {
				if class != "" {
					for _, a := range n.FirstChild.Attr {
						if a.Key == "class" && a.Val == class {
							n.RemoveChild(n.FirstChild)
						}
					}
				} else {
					n.RemoveChild(n.FirstChild)
				}
			}
		}
	}
	ForEachNode(doc, rmFirstTag, nil)
}

// ElementsRmByTagAttr rm nodes
// if attrName != "" rm by tag and attr else rm just by tag.
func ElementsRmByTagAttr(doc *html.Node, tag, attrName, attrValue string) {
	if tag == "" || doc == nil {
		return
	}
	visitNode := func(n *html.Node) {
		if n.NextSibling != nil && n.NextSibling.Type == html.ElementNode {
			if tag == n.NextSibling.Data {
				if attrName != "" {
					for _, a := range n.NextSibling.Attr {
						if a.Key == attrName && a.Val == attrValue {
							n.Parent.RemoveChild(n.NextSibling)
						}
					}
				} else {
					n.Parent.RemoveChild(n.NextSibling)
				}
			}
		}
	}
	ForEachNode(doc, visitNode, nil)
	rmFirstTag := func(n *html.Node) {
		if n.FirstChild != nil && n.FirstChild.Type == html.ElementNode {
			if tag == n.FirstChild.Data {
				if attrName != "" {
					for _, a := range n.FirstChild.Attr {
						if a.Key == attrName && a.Val == attrValue {
							n.RemoveChild(n.FirstChild)
						}
					}
				} else {
					n.RemoveChild(n.FirstChild)
				}
			}
		}
	}
	ForEachNode(doc, rmFirstTag, nil)
}

func ElementsByTag(doc *html.Node, name ...string) []*html.Node {
	var nodes []*html.Node
	if len(name) == 0 || doc == nil {
		return nil
	}
	if doc.Type == html.ElementNode {
		for _, tag := range name {
			if tag == doc.Data {
				nodes = append(nodes, doc)
			}
		}
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, ElementsByTag(c, name...)...)
	}
	return nodes
}

func ElementsByTag2(raw []byte, tags ...string) []byte {
	if raw == nil || tags == nil {
		return nil
	}
	z := html.NewTokenizer(bytes.NewReader(raw))
	var b bytes.Buffer
	depth := 0
	for {
		tt := z.Next()
		if err := z.Err(); err == io.EOF {
			break
		}
		switch tt {
		case html.ErrorToken:
			return []byte(z.Err().Error())
		case html.TextToken:
			if depth > 0 {
				b.Write(z.Text())
			}
		case html.StartTagToken, html.EndTagToken:
			tn, _ := z.TagName()
			for _, tag := range tags {
				if string(tn) == tag {
					b.Write(z.Text())
				}
			}
			if tt == html.StartTagToken {
				depth++
			} else {
				depth--
			}
		}
	}
	return b.Bytes()
}

func TagWithAttr(doc *html.Node, tag, attr string) []*html.Node {
	var nodes []*html.Node
	if attr == "" || tag == "" || doc == nil {
		return nil
	}
	if doc.Type == html.ElementNode {
		if tag == doc.Data {
			for _, a := range doc.Attr {
				if a.Key == attr {
					nodes = append(nodes, doc)
				}
			}
		}
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, TagWithAttr(c, tag, attr)...)
	}
	return nodes
}

func ElementsByTagAndClass(doc *html.Node, tag, class string) []*html.Node {
	var nodes []*html.Node
	if tag == "" || class == "" || doc == nil {
		return nil
	}
	if doc.Type == html.ElementNode {
		if tag == doc.Data {
			for _, a := range doc.Attr {
				if a.Key == "class" && a.Val == class {
					nodes = append(nodes, doc)
				}
			}
		}
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, ElementsByTagAndClass(c, tag, class)...)
	}
	return nodes
}

func ElementsByTagAndClass2(raw []byte, tag, class string) []byte {
	if raw == nil || tag == "" || class == "" {
		return nil
	}
	z := html.NewTokenizer(bytes.NewReader(raw))
	var b bytes.Buffer
	for {
		tt := z.Next()
		t := z.Token()
		if err := z.Err(); err == io.EOF {
			break
		}
		switch tt {
		case html.StartTagToken:
			if tag == t.Data {
				for _, a := range t.Attr {
					if a.Key == "class" && a.Val == class {
						b.Write(z.Buffered())
					}
				}
			}
		}
	}
	return b.Bytes()
}

func ElementsByTagAndId(doc *html.Node, tag, id string) []*html.Node {
	var nodes []*html.Node
	if doc == nil || tag == "" || id == "" {
		return nil
	}
	if doc.Type == html.ElementNode {
		if tag == doc.Data {
			for _, a := range doc.Attr {
				if a.Key == "id" && a.Val == id {
					nodes = append(nodes, doc)
				}
			}
		}
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, ElementsByTagAndId(c, tag, id)...)
	}
	return nodes
}

func ElementsByTagAndId2(raw []byte, tag, id string) []byte {
	if raw == nil || tag == "" || id == "" {
		return nil
	}
	z := html.NewTokenizer(bytes.NewReader(raw))
	for {
		tt := z.Next()
		t := z.Token()
		if err := z.Err(); err != nil && err == io.EOF {
			break
		}
		switch tt {
		case html.StartTagToken:
			if tag == t.Data {
				for _, a := range t.Attr {
					if a.Key == "id" && a.Val == id {
						return z.Buffered()
					}
				}
			}
		}
	}
	return nil
}

func ElementsByTagAndType(doc *html.Node, tag, attrType string) []*html.Node {
	var nodes []*html.Node
	if tag == "" || attrType == "" || doc == nil {
		return nil
	}
	if doc.Type == html.ElementNode {
		if tag == doc.Data {
			for _, a := range doc.Attr {
				if a.Key == "type" && a.Val == attrType {
					nodes = append(nodes, doc)
				}
			}
		}
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, ElementsByTagAndType(c, tag, attrType)...)
	}
	return nodes
}

func ElementsNextByTag(doc *html.Node, tag string) []*html.Node {
	var nodes []*html.Node
	if tag == "" || doc == nil {
		return nil
	}
	if doc == nil || tag == "" {
		return nil
	}
	if doc.Type == html.ElementNode {
		if tag == doc.Data && doc.NextSibling != nil {
			nodes = append(nodes, doc.NextSibling)
		}
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, ElementsNextByTag(c, tag)...)
	}
	return nodes
}

// MetasByName focus on `<meta name="dateModified" content="2020/09/29 11:27" />`
func MetasByName(doc *html.Node, values ...string) []*html.Node {
	var nodes []*html.Node
	if doc == nil || values == nil {
		return nil
	}
	if doc.Type == html.ElementNode {
		if doc.Data == "meta" {
			for _, a := range doc.Attr {
				if a.Key == "name" {
					for _, v := range values {
						if v == a.Val {
							nodes = append(nodes, doc)
						}
					}
				}
			}
		}
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, MetasByName(c, values...)...)
	}
	return nodes
}

// MetasByItemprop focus on `<meta itemprop="dateModified" content="2020/09/29 11:27" />`
func MetasByItemprop(doc *html.Node, values ...string) []*html.Node {
	var nodes []*html.Node
	if doc == nil || values == nil {
		return nil
	}
	if doc.Type == html.ElementNode {
		if doc.Data == "meta" {
			for _, a := range doc.Attr {
				if a.Key == "itemprop" {
					for _, v := range values {
						if v == a.Val {
							nodes = append(nodes, doc)
						}
					}
				}
			}
		}
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, MetasByItemprop(c, values...)...)
	}
	return nodes
}

// MetasByProperty focus on `<meta property="dateModified" content="2020/09/29 11:27" />`
func MetasByProperty(doc *html.Node, values ...string) []*html.Node {
	var nodes []*html.Node
	if doc == nil || values == nil {
		return nil
	}
	if doc.Type == html.ElementNode {
		if doc.Data == "meta" {
			for _, a := range doc.Attr {
				if a.Key == "property" {
					for _, v := range values {
						if v == a.Val {
							nodes = append(nodes, doc)
						}
					}
				}
			}
		}
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, MetasByProperty(c, values...)...)
	}
	return nodes
}
