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

	"github.com/PuerkitoBio/goquery"
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

type MeizituCrawl struct {
	sync.Mutex
	worker  int
	base    string
	url     string
	pages   int
	limit   int
	links   chan string
	nexts   chan string
	pics    chan string
	picsSet map[string]struct{}
}

func NewMeizituCrawl(url string, limit int) *MeizituCrawl {
	return &MeizituCrawl{
		base:    "/tmp/meizitu/",
		worker:  10,
		url:     url,
		limit:   limit,
		pages:   0,
		links:   make(chan string, 20),
		nexts:   make(chan string, 20),
		pics:    make(chan string, 20),
		picsSet: make(map[string]struct{}),
	}
}

func (mc *MeizituCrawl) addNext(url string) {
	mc.Lock()
	defer mc.Unlock()
	if mc.pages < mc.limit {
		mc.pages++
		mc.nexts <- url
		if mc.pages == mc.limit {
			logger.Println(strings.Repeat("*", 32), "close channal >> nexts", strings.Repeat("*", 32))
			close(mc.nexts)
		}
	}
}

func (mc *MeizituCrawl) addLinks(urls []string) {
	for _, url := range urls {
		mc.links <- url
	}
}

func (mc *MeizituCrawl) addPics(pics []string) {
	for _, pic := range pics {
		var ok bool
		mc.Lock()
		if _, ok = mc.picsSet[pic]; !ok {
			mc.picsSet[pic] = struct{}{}
		}
		mc.Unlock()
		if !ok {
			mc.pics <- pic
		}
	}
}

func (mc *MeizituCrawl) Crawl() {
	mc.addNext(mc.url)
	go func() {
		for url := range mc.nexts {
			next, links, _, err := Extract2(url)
			if err != nil {
				logger.Printf("Error: %v", err)
				continue
			}
			mc.addLinks(links)
			mc.addNext(next)
		}
		logger.Println(strings.Repeat("*", 32), "close channal >> links", strings.Repeat("*", 32))
		close(mc.links)
	}()

	go func() {
		for url := range mc.links {
			_, _, pics, err := Extract2(url)
			if err != nil {
				logger.Printf("Error: %v", err)
				continue
			}
			mc.addPics(pics)
		}
		logger.Println(strings.Repeat("*", 32), "close channal >> pics", strings.Repeat("*", 32))
		close(mc.pics)
	}()

	var wg sync.WaitGroup
	wg.Add(mc.worker)
	for i := 0; i < mc.worker; i++ {
		go func() {
			defer wg.Done()
			for pic := range mc.pics {
				err := downloadPic(context.TODO(), mc.base, pic)
				if err != nil {
					logger.Println(err)
				}
			}
		}()
	}
	wg.Wait()

	logger.Println("crawl done.")
}

func Extract2(url string) (next string, links []string, pics []string, err error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		err = fmt.Errorf("getting %s: %s", url, resp.Status)
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("xmlpath: parsing %s with error: %v", url, err)
		return
	}

	// list
	doc.Find("#maincontent .inWrap ul li .pic").Each(func(i int, s *goquery.Selection) {
		if link, ok := s.Find("a").Attr("href"); ok {
			logger.Println(link)
			links = append(links, link)
		}
	})
	// next
	doc.Find(".navigation #wp_page_numbers ul li:nth-last-child(2)").Each(func(i int, s *goquery.Selection) {
		if n, ok := s.Find("a").Attr("href"); ok {
			n, err := resp.Request.URL.Parse(n)
			if err != nil {
				return
			}
			next = n.String()
			logger.Println("next page is --> ", next)
		}
	})
	// pics
	doc.Find(".postContent #picture p img").Each(func(i int, s *goquery.Selection) {
		if m, ok := s.Attr("src"); ok {
			// logger.Println(m)
			pics = append(pics, m)
		}
	})
	return
}

func downloadPic(ctx context.Context, base, src string) error {
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
				err := downloadPic(ctx, base, s)
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

func longOpDo(ctx context.Context) (string, error) {
	s := make(chan string, 1)
	c := make(chan error, 1)
	go func() {
		s1, c1 := longOp(ctx)
		s <- s1
		c <- c1
	}()
	select {
	case <-ctx.Done():
		return "false 1", ctx.Err()
	case r := <-s:
		err := <-c
		return r, err
	}
}

func longOp(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "false 2", ctx.Err()
	default:
	}

	logger.Println("sleep 4 second")
	time.Sleep(4 * time.Second)
	return "true", nil
}
