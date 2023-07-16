package mobilelegendapi

import "strings"

type NewsType string

const (
	DefaultType NewsType = "news"
	ArticleType NewsType = "event"
	GuideType   NewsType = "guide"
	UnknownType NewsType = "unknown"
)

func StringToNewsType(s string) NewsType {
	switch strings.ToLower(s) {
	case "news":
		return DefaultType
	case "event":
		return ArticleType
	case "guide":
		return GuideType
	default:
		return UnknownType
	}
}

type News struct {
	Id          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Thumbnail   string   `json:"thumbnail"`
	Type        NewsType `json:"type"`
}
