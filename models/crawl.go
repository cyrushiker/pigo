package models

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/proxy"
)

var httpClient *http.Client

func init() {
	// create a socks5 dialer
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:1080", nil, proxy.Direct)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		os.Exit(1)
	}
	// set our socks5 as the dialer
	httpTransport := &http.Transport{Dial: dialer.Dial}
	// setup a http client
	httpClient = &http.Client{Transport: httpTransport}
}

func Extract(url string) ([]string, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, a := range n.Attr {
				if a.Key != "src" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

func downloadPic(base, src string) error {
	logger.Printf("Get picture from #%s#", src)
	ps := strings.Split(src, "/")
	name := strings.Join(ps[3:], "_")
	if _, err := os.Stat(base); os.IsNotExist(err) {
		os.Mkdir(base, os.ModePerm)
	}
	out, err := os.Create(base + name)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("GET", src, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func DownloadPics(ctx context.Context, base string, srcs []string, w int) error {
	var wg sync.WaitGroup
	errc := make(chan error)
	sc := make(chan string)

	go func() {
		for _, src := range srcs {
			sc <- src
		}
		close(sc)
	}()

	for i := 0; i < w; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for s := range sc {
				err := downloadPic(base, s)
				if err != nil {
					errc <- err
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	logger.Println(strings.Repeat("**", 32))
	for e := range errc {
		logger.Println(e)
	}

	return nil
}
