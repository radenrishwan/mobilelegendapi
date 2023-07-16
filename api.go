package mobilelegendapi

import (
	"context"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

const ENDPOINT = "https://m.mobilelegends.com/id"

type Core struct {
	context.Context
}

func NewCore(ctx context.Context) *Core {
	alloc, _ := chromedp.NewContext(ctx)

	return &Core{
		Context: alloc,
	}
}

func (core *Core) GetNews(ctx context.Context) []News {
	chromectx, cancel := chromedp.NewContext(core.Context)
	defer cancel()

	url := ENDPOINT + "/news/EVENTS/list"
	var result string
	err := chromedp.Run(chromectx, chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("Page loaded")
			return nil
		}),
		chromedp.WaitVisible(`div.midnewsitem`),
		chromedp.OuterHTML(`html`, &result),
	})

	if err != nil {
		panic(err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(result))

	if err != nil {
		panic(err)
	}

	log.Println("already get document")

	var news []News
	doc.Find(".midnews").Each(func(i int, s *goquery.Selection) {
		var res News
		// get href attribute
		id, _ := s.Find("a").Attr("href")
		res.Id = id

		// get title
		title := s.Find(".desc > h3").Text()
		res.Title = title

		// get tags
		tags := s.Find(".tags > .tag").Text()
		res.Type = StringToNewsType(tags)

		// get description
		desc := s.Find(".desc > p").Text()
		res.Description = desc

		// get thumbnail
		thumbnail, _ := s.Find("img").Attr("data-src")
		res.Thumbnail = thumbnail

		// append to news
		news = append(news, res)
	})

	return news
}
