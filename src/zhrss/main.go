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
	Refreshing  bool
}

type sysInfo struct {
	resultInfo
	Link     string
	Location *time.Location
}

func (i *sysInfo) fetchDoc() (*goquery.Document, error) {
	doc, err := goquery.NewDocument(i.Link)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (i *sysInfo) parseFeed(doc *goquery.Document) *feeds.Feed {
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

	for x := 0; x < itemsQ.Length(); x++ {
		item = itemsQ.Eq(x)

		linkQ := item.Find(".zm-profile-section-main a").Last()
		title := linkQ.Text()
		if title == "" {
			continue
		}

		link, successed := linkQ.Attr("href")
		if !successed {
			continue
		}
		if strings.HasPrefix(link, "/") {
			link = fmt.Sprintf("%s://%s%s", linkInfo.Scheme, linkInfo.Host, link)
		}
		id := stringToUUID(link)

		dataTime, successed := item.Attr("data-time")
		if !successed {
			continue
		}
		timestamp, err := strconv.ParseInt(dataTime, 10, 64)
		if err != nil {
			continue
		}

		created := time.Unix(timestamp, 0).In(i.Location)
		author := item.Find(".author-link").Last().Text()
		content := item.Find("textarea.content").Text()

		items = append(items, &feeds.Item{
			Id:          id,
			Created:     created,
			Title:       title,
			Author:      &feeds.Author{Name: author},
			Link:        &feeds.Link{Href: link},
			Description: content,
		})
	}

	return &feeds.Feed{
		Id:    stringToUUID(i.Link),
		Title: doc.Find("title").First().Text(),
		Link:  &feeds.Link{Href: i.Link},
		Description: doc.Find(
			"div.zm-profile-header-description span.content",
		).First().Text(),
		Author: &feeds.Author{Name: doc.Find(
			"div.title-section span.name",
		).First().Text()},
		Created: time.Now().In(i.Location),
		Items:   items,
	}
}

func (i *sysInfo) refreshRSSResult() (string, error) {
	log.Print("ready to refresh result")

	doc, err := i.fetchDoc()
	if err != nil {
		return "", err
	}

	feed := i.parseFeed(doc)
	if feed == nil {
		return "", err
	}

	rss, err := feed.ToRss()
	if err != nil {
		return "", err
	}
	return rss, nil
}

func (i *sysInfo) handle(w http.ResponseWriter, r *http.Request) {

	now := time.Now().In(i.Location)
	if i.Refreshing || now.Before(i.NextCheckAt) {
		log.Print("return result from cache")
		fmt.Fprint(w, i.Result)
		return
	}

	i.Refreshing = true
	rss, err := i.refreshRSSResult()
	if err != nil {
		log.Print("return result from cache due to :", err)
	} else if i.NextCheckAt.Before(now) {
		i.Result = rss
		log.Print("next check at ", i.NextCheckAt)
	}
	i.Refreshing = false
	i.NextCheckAt = now.Add(time.Duration(i.CacheTTL) * time.Second)

	fmt.Fprint(w, i.Result)
}

func stringToUUID(s string) string {
	var uuid [16]byte
	for i := 0; i < len(s); i++ {
		uuid[i%16] ^= s[i]
	}
	uuid[8] = (uuid[8] | 0x80) & 0xBf
	uuid[6] = (uuid[6] | 0x40) & 0x4f
	return fmt.Sprintf(
		"%x-%x-%x-%x-%x", uuid[:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:],
	)
}

func main() {
	var serverAddr string
	var locationName string
	var urlPath string
	info := new(sysInfo)

	flag.StringVar(
		&info.Link, "url", "https://www.zhihu.com/people/mr_lyc",
		"user time line page url",
	)
	flag.StringVar(&serverAddr, "addr", ":8080", "address to listen")
	flag.StringVar(&locationName, "location", "UTC", "location name")
	flag.Int64Var(&info.CacheTTL, "cache", 600, "result cache ttl")
	flag.StringVar(&urlPath, "path", "/", "url path")
	flag.Parse()

	location, err := time.LoadLocation(locationName)
	if err != nil {
		log.Fatal("load location failed: ", err)
	}
	info.Refreshing = false
	info.Location = location
	info.Result, err = info.refreshRSSResult()
	if err != nil {
		log.Fatal("refresh rss result failed: ", err)
	}

	http.HandleFunc(urlPath, info.handle)
	log.Printf("feed server for %s listen on %s", info.Link, serverAddr)
	err = http.ListenAndServe(serverAddr, nil)
	if err != nil {
		log.Fatal("server start failed: ", err)
	}
}
