package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
)

type resultInfo struct {
	CacheTTL    int64
	NextCheckAt time.Time
	Result      string
}

type zhihuInfo struct {
	resultInfo
	Link string
}

func (i *zhihuInfo) fetchDoc() (*goquery.Document, error) {
	doc, err := goquery.NewDocument(i.Link)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (i *zhihuInfo) parseFeed(doc *goquery.Document) *feeds.Feed {
	itemsQ := doc.Find(
		"#zh-profile-activity-page-list .zm-profile-section-item",
	)
	items := make([]*feeds.Item, 0, itemsQ.Length())
	var item *goquery.Selection
	linkInfo, err := url.Parse(i.Link)
	if err != nil {
		log.Print("parse feed failed due to: ", err)
		return nil
	}

	for i := 0; i < itemsQ.Length(); i++ {
		item = itemsQ.Eq(i)
		dataTime, successed := item.Attr("data-time")
		if !successed {
			continue
		}
		timestamp, err := strconv.ParseInt(dataTime, 10, 64)
		if err != nil {
			continue
		}
		created := time.Unix(timestamp, 0).UTC()
		linkQ := item.Find(".zm-profile-section-main a").Last()

		link, successed := linkQ.Attr("href")
		if !successed {
			continue
		}
		if strings.HasPrefix(link, "/") {
			link = fmt.Sprintf("%s://%s%s", linkInfo.Scheme, linkInfo.Host, link)
		}

		title := linkQ.Text()
		if title == "" {
			continue
		}
		author := item.Find(".author-link").Last().Text()

		content := item.Find("textarea.content").Text()

		items = append(items, &feeds.Item{
			Created:     created,
			Title:       title,
			Author:      &feeds.Author{Name: author},
			Link:        &feeds.Link{Href: link},
			Description: content,
		})
	}

	return &feeds.Feed{
		Title: doc.Find("title").First().Text(),
		Link:  &feeds.Link{Href: i.Link},
		Description: doc.Find(
			"div.zm-profile-header-description span.content",
		).First().Text(),
		Author: &feeds.Author{Name: doc.Find(
			"div.title-section span.name",
		).First().Text()},
		Created: time.Now().UTC(),
		Items:   items,
	}
}

func (i *zhihuInfo) handle(w http.ResponseWriter, r *http.Request) {

	now := time.Now().UTC()
	if now.Before(i.NextCheckAt) {
		log.Print("return result from cache")
		fmt.Fprint(w, i.Result)
		return
	}
	log.Print("ready to refresh result")

	doc, err := i.fetchDoc()
	if err != nil {
		log.Print("return result from cache due to: ", err)
		fmt.Fprint(w, i.Result)
		return
	}
	feed := i.parseFeed(doc)
	if feed == nil {
		log.Print("return result from cache due to empty feed")
		fmt.Fprint(w, i.Result)
		return
	}

	rss, err := feed.ToRss()
	if err != nil {
		log.Print("return result from cache due to: ", err)
		fmt.Fprint(w, i.Result)
		return
	}
	if i.NextCheckAt.Before(now) {
		i.Result = rss
		i.NextCheckAt = now.Add(time.Duration(i.CacheTTL) * time.Second)
		log.Print("next check at ", i.NextCheckAt)
	}

	fmt.Fprint(w, i.Result)
}

func main() {
	var serverAddr string
	info := new(zhihuInfo)
	info.NextCheckAt = time.Now().UTC()

	flag.StringVar(
		&info.Link, "url", "https://www.zhihu.com/people/mr_lyc",
		"user time line page url",
	)
	flag.StringVar(&serverAddr, "addr", ":8080", "address to listen")
	flag.Int64Var(&info.CacheTTL, "cache", 600, "result cache ttl")
	flag.Parse()

	http.HandleFunc("/", info.handle)
	log.Printf("feed server for %s listen on %s", info.Link, serverAddr)
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		log.Fatal("server start failed: ", err)
	}
}
