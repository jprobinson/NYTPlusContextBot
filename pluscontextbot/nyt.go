package pluscontextbot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func SearchNYT(text string, APIToken string) (*Article, error) {
	query := url.QueryEscape(fmt.Sprintf(`body:%s`, strconv.Quote(text)))
	urlStr := `http://api.nytimes.com/svc/search/v2/articlesearch.json?fq=%s&api-key=%s`
	urlStr = fmt.Sprintf(urlStr, query, APIToken)

	// create request with timeout
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	req.Close = true
	// make request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// get the datas
	var nytResp SearchResponse
	err = json.Unmarshal(body, &nytResp)
	if err != nil {
		log.Printf("unable to unmarshal JSON (err: %s)\n Response:\n %s", err, string(body))
		return nil, err
	}

	if len(nytResp.Response.Docs) == 0 {
		log.Printf("unable to find results?\n Response:\n %s", string(body))
		return nil, errors.New("not found")
	}

	// extract the datas
	data := nytResp.Response.Docs[0]
	// if a kicker exists, toss it on
	headline := data.Headline.Main
	if data.Headline.Kicker != "" {
		headline = data.Headline.Kicker + " | " + headline
	}
	return &Article{data.WebURL, headline}, nil
}

type Article struct {
	URL      string
	Headline string
}

type SearchResponse struct {
	Response struct {
		Docs []struct {
			WebURL   string `json:"web_url"`
			Headline struct {
				Main   string `json:"main"`
				Kicker string `json:"kicker"`
			} `json:"headline"`
		} `json:"docs"`
	} `json:"response"`
}
