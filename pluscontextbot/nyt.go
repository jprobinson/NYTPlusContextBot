package pluscontextbot

import (
	"encoding/json"
	"errors"
	"fmt"
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

	// get the datas
	var nytResp SearchResponse
	err = json.NewDecoder(resp.Body).Decode(&nytResp)
	if err != nil {
		return nil, err
	}

	if len(nytResp.Response.Docs) == 0 {
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
		Meta struct {
			Hits   int `json:"hits"`
			Time   int `json:"time"`
			Offset int `json:"offset"`
		} `json:"meta"`
		Docs []struct {
			WebURL        string        `json:"web_url"`
			Snippet       string        `json:"snippet"`
			LeadParagraph interface{}   `json:"lead_paragraph"`
			Abstract      string        `json:"abstract"`
			PrintPage     interface{}   `json:"print_page"`
			Blog          []interface{} `json:"blog"`
			Source        string        `json:"source"`
			Multimedia    []struct {
				Width   int    `json:"width"`
				URL     string `json:"url"`
				Height  int    `json:"height"`
				Subtype string `json:"subtype"`
				Legacy  struct {
					Wide       string `json:"wide"`
					Wideheight string `json:"wideheight"`
					Widewidth  string `json:"widewidth"`
				} `json:"legacy"`
				Type string `json:"type"`
			} `json:"multimedia"`
			Headline struct {
				Main   string `json:"main"`
				Kicker string `json:"kicker"`
			} `json:"headline"`
			Keywords []struct {
				Rank  string `json:"rank"`
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"keywords"`
			PubDate        time.Time `json:"pub_date"`
			DocumentType   string    `json:"document_type"`
			NewsDesk       string    `json:"news_desk"`
			SectionName    string    `json:"section_name"`
			SubsectionName string    `json:"subsection_name"`
			Byline         struct {
				Person []struct {
					Organization string `json:"organization"`
					Role         string `json:"role"`
					Firstname    string `json:"firstname"`
					Rank         int    `json:"rank"`
					Lastname     string `json:"lastname"`
				} `json:"person"`
				Original string `json:"original"`
			} `json:"byline"`
			TypeOfMaterial   string      `json:"type_of_material"`
			ID               string      `json:"_id"`
			WordCount        string      `json:"word_count"`
			SlideshowCredits interface{} `json:"slideshow_credits"`
		} `json:"docs"`
	} `json:"response"`
	Status    string `json:"status"`
	Copyright string `json:"copyright"`
}
