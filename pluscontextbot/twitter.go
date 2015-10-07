package pluscontextbot

import (
	"strings"
	"unicode/utf8"

	"golang.org/x/exp/utf8string"

	"github.com/ChimeraCoder/anaconda"
)

func PublishTweet(article *Article, api *anaconda.TwitterApi) (string, error) {
	msg := twitterize(article)
	_, err := api.PostTweet(msg, nil)
	return msg, err
}

func twitterize(article *Article) string {
	msg := "@NYTMinusContext " + article.Headline
	// remove any URLs from headline so our link works
	msg = strings.Replace(msg, ".com", "", -1)
	// leave 24 for space and link!
	if utf8.RuneCountInString(msg) > 116 {
		msg = utf8string.NewString(msg).Slice(0, 113) + "..."
	}
	// tack on the URL
	msg += " " + article.URL
	return msg
}
