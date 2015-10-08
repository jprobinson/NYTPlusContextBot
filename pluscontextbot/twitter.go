package pluscontextbot

import (
	"html"
	"net/url"
	"regexp"
	"unicode/utf8"

	"golang.org/x/exp/utf8string"

	"github.com/ChimeraCoder/anaconda"
)

func PublishTweet(article *Article, replyToId string, api *anaconda.TwitterApi) (string, error) {
	msg := twitterize(article)
	_, err := api.PostTweet(msg, url.Values{
		"in_reply_to_status_id": []string{replyToId},
	})
	return msg, err
}

var cleanRegex = regexp.MustCompile(`\.(com|org|net|io)`)

func twitterize(article *Article) string {
	msg := "@NYTMinusContext " + html.UnescapeString(article.Headline)
	// remove any URLs from headline so our link works
	msg = cleanRegex.ReplaceAllString(msg, "")
	// leave 24 for space and link!
	if utf8.RuneCountInString(msg) > 116 {
		msg = utf8string.NewString(msg).Slice(0, 113) + "..."
	}
	// tack on the URL
	msg += " " + article.URL
	return msg
}
