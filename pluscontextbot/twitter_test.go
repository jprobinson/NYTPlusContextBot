package pluscontextbot

import "testing"

func TestTwitterize(t *testing.T) {

	tests := []struct {
		name  string
		given Article
		want  string
	}{
		{
			"No Truncate",
			Article{
				Headline: "Man Bites Dog",
				URL:      "http://www.nytimes.com/man-bites-dog",
			},
			"@NYTMinusContext Man Bites Dog http://www.nytimes.com/man-bites-dog",
		},
		{
			"URL in Headline",
			Article{
				Headline: "Man Bites Dog.com",
				URL:      "http://www.nytimes.com/man-bites-dog",
			},
			"@NYTMinusContext Man Bites Dog http://www.nytimes.com/man-bites-dog",
		},
		{
			"Truncate",
			Article{
				Headline: "Man Bites Dog And It's Really, Really Gory and Stuff. Top Notch Headline, Here. All full of words and whatnot.",
				URL:      "http://www.nytimes.com/man-bites-dog",
			},
			"@NYTMinusContext Man Bites Dog And It's Really, Really Gory and Stuff. Top Notch Headline, Here. All full of word... http://www.nytimes.com/man-bites-dog",
		},
		{
			"Exact Length",
			Article{
				Headline: "Man Bites Dog And It's Really, Really Gory and Stuff. Top Notch Headline, Here. A lot of words 1234",
				URL:      "http://www.nytimes.com/man-bites-dog",
			},
			"@NYTMinusContext Man Bites Dog And It's Really, Really Gory and Stuff. Top Notch Headline, Here. A lot of words 1234 http://www.nytimes.com/man-bites-dog",
		},
	}

	for _, test := range tests {
		got := twitterize(&test.given)

		if got != test.want {
			t.Errorf("twitterize[%s] got %s, expected %s", test.name, got, test.want)
		}
	}
}
