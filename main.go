package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/ChimeraCoder/anaconda"
	"github.com/NYTimes/gizmo/config"
	"github.com/jprobinson/NYTPlusContextBot/pluscontextbot"
)

const NYTMinusContextID = "2189503302"

func main() {
	// pull creds out of envirnment
	var creds Creds
	config.LoadEnvConfig(&creds)

	// setup the Twitter client
	anaconda.SetConsumerKey(creds.TwtConsumer)
	anaconda.SetConsumerSecret(creds.TwtConsumerSecret)
	api := anaconda.NewTwitterApi(creds.TwtAccess, creds.TwtAccessSecret)
	api.SetLogger(anaconda.BasicLogger)
	stream := api.PublicStreamFilter(url.Values{
		// @NYTMinusContext's ID
		"follow": []string{NYTMinusContextID},
	})

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
		log.Printf("Received kill signal %s", <-ch)
		stream.Stop()
	}()

	log.Print("listening for tweets")

	for res := range stream.C {
		tweet, ok := res.(anaconda.Tweet)
		if !ok {
			log.Printf("encountered a non-tweet: %#v", res)
			continue
		}

		if tweet.User.IdStr != NYTMinusContextID {
			log.Print("ignoring non-NYTMinusContext tweet from "+tweet.User.ScreenName, ": "+tweet.Text)
			continue
		}

		log.Printf("got a tweet: %#v", tweet)

		article, err := pluscontextbot.SearchNYT(tweet.Text, creds.NYTToken)
		if err != nil {
			log.Printf(`unable to complete search of tweet ("%s"): %s`, tweet.Text, err)
			continue
		}

		log.Printf("found the article: %#v", article)

		var msg string
		msg, err = pluscontextbot.PublishTweet(article, tweet.IdStr, api)
		if err != nil {
			log.Print("unable to publish tweet: ", err)
			continue
		}

		log.Print("published tweet: ", msg)
	}

	log.Print("shutting down")
}

type Creds struct {
	TwtAccess         string `envconfig:"TWT_ACCESS"`
	TwtAccessSecret   string `envconfig:"TWT_ACCESS_SECRET"`
	TwtConsumer       string `envconfig:"TWT_CONSUMER"`
	TwtConsumerSecret string `envconfig:"TWT_CONSUMER_SECRET"`

	NYTToken string `envconfig:"NYT_TOKEN"`
}
