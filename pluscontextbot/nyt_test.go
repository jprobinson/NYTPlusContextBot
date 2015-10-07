package pluscontextbot

import "testing"

func TestNYTSearch_integration(t *testing.T) {
	text := `a hyper-realistic baby astronaut with communist and capitalist patches on its suit and a dove sitting atop his head`
	got, err := SearchNYT(text, "sample-key")
	if err != nil {
		t.Fatal("unable to perform search: ", err)
	}

	if got.Headline != "T Magazine | Past, Future or Alien?" ||
		got.URL != "http://tmagazine.blogs.nytimes.com/2014/02/07/on-view-past-future-or-alien/" {
		t.Errorf("Did not get the article we expected: %#v", got)
	}
}
