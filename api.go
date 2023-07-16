package mobilelegendapi

import (
	"context"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

const (
	Miya     = "1"
	Balmond  = "2"
	Saber    = "3"
	Alucard  = "4"
	Chou     = "5"
	Fanny    = "6"
	Karina   = "7"
	Johnson  = "8"
	Minotaur = "9"
	Franco   = "10" // TODO: add more
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
		log.Panicln(err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(result))

	if err != nil {
		log.Panicln(err)
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

func (core *Core) GetHeroById(ctx context.Context, id string) Hero {
	chromectx, cancel := chromedp.NewContext(core.Context)
	defer cancel()

	url := ENDPOINT + "/hero/" + id + "/skill"
	var result string
	err := chromedp.Run(chromectx, chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("Page loaded")
			return nil
		}),
		// wait until ready
		chromedp.WaitReady("div.skilllist"),
		chromedp.OuterHTML(`html`, &result),
	})

	if err != nil {
		log.Panicln(err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(result))

	if err != nil {
		log.Panicln(err)
	}

	log.Println("already get document")
	var hero Hero
	var skills []Skill
	doc.Find("#skill > .skilllist").Each(func(i int, s *goquery.Selection) {
		// get skill
		var skill Skill
		s.Find("ul > li").Each(func(i int, sc *goquery.Selection) { // TODO: fix duplicate skill
			// name
			skill.Name = sc.Find("p").Text()

			// image url
			img, _ := sc.Find("img").Attr("data-src")
			skill.ImageUrl = img[2:]
		})

		s.Find(".skilldesc").Each(func(i int, sc *goquery.Selection) {
			skill.Description = sc.Find("p").Text()
			skill.Tips = sc.Find(".tips").Text()
		})

		skills = append(skills, skill)
	})

	doc.Find(".name > h3").Each(func(i int, s *goquery.Selection) {
		hero.Name = s.Text()
	})

	hero.Skills = skills[:len(skills)/2]
	hero.Id = id

	return hero
}
