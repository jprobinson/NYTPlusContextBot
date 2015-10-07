package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/ChimeraCoder/anaconda"
	"github.com/jprobinson/NYTPlusContextBot/pluscontextbot"
)

var credsFile = flag.String("creds", "/opt/nyt/etc/creds.json", "path to creds.json file")

const NYTMinusContextID = "2189503302"

func main() {
	flag.Parse()

	// get credentials from file
	f, err := os.Open(*credsFile)
	if err != nil {
		log.Fatal("unable to open creds file ("+*credsFile+"): ", err)
	}
	defer f.Close()

	var creds Creds
	err = json.NewDecoder(f).Decode(&creds)
	if err != nil {
		log.Fatal("unable to parse creds file ("+*credsFile+"): ", err)
	}

	// setup the Twitter client
	anaconda.SetConsumerKey(creds.TwtConsumer)
	anaconda.SetConsumerSecret(creds.TwtConsumerSecret)
	api := anaconda.NewTwitterApi(creds.TwtAccess, creds.TwtAccessSecret)
	stream := api.PublicStreamFilter(url.Values{
		// @NYTMinusContext's ID
		"follow": []string{NYTMinusContextID},
	})

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
		log.Printf("Received kill signal %s", <-ch)
		stream.Interrupt()
		stream.End()
	}()

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
		msg, err = pluscontextbot.PublishTweet(article, api)
		if err != nil {
			log.Print("unable to publish tweet: ", err)
			continue
		}

		log.Print("published tweet: ", msg)
	}

	log.Print("shutting down")
}

type Creds struct {
	TwtAccess, TwtAccessSecret     string
	TwtConsumer, TwtConsumerSecret string

	NYTToken string
}
