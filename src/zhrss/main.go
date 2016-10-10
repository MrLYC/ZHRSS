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

type zhihuInfo struct {
	Link        string
	LastCheckAt time.Time
}

func (i *zhihuInfo) checkAndFetch() (*goquery.Document, error) {
	doc, err := goquery.NewDocument(i.Link)
	if err != nil {
		return nil, err
	}

	i.LastCheckAt = time.Now()
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
		created := time.Unix(timestamp, 0)
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
		Created: time.Now(),
		Items:   items,
	}
}

func (i *zhihuInfo) handle(w http.ResponseWriter, r *http.Request) {
	doc, err := i.checkAndFetch()
	if err != nil {
		return
	}
	feed := i.parseFeed(doc)
	rss, err := feed.ToRss()
	if err != nil {
		return
	}
	fmt.Fprint(w, rss)
}

func main() {
	var serverAddr string
	info := new(zhihuInfo)

	flag.StringVar(
		&info.Link, "u", "https://www.zhihu.com/people/mr_lyc",
		"user time line page url",
	)
	flag.StringVar(&serverAddr, "a", ":8080", "address to listen")

	http.HandleFunc("/", info.handle)
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		log.Fatal("server start failed: ", err)
	}

}
