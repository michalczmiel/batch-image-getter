package internal

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/net/html"
)

var userAgents = []string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 Mobile/14F89 Safari/602.1",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) FxiOS/8.1.1b4948 Mobile/14F89 Safari/603.2.4",
	"Mozilla/5.0 (iPad; CPU OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 Mobile/14F89 Safari/602.1",
	"Mozilla/5.0 (Linux; Android 4.3; GT-I9300 Build/JSS15J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.125 Mobile Safari/537.36",
	"Mozilla/5.0 (Android 4.3; Mobile; rv:54.0) Gecko/54.0 Firefox/54.0",
	"Mozilla/5.0 (Linux; Android 4.3; GT-I9300 Build/JSS15J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.91 Mobile Safari/537.36 OPR/42.9.2246.119956",
	"Opera/9.80 (Android; Opera Mini/28.0.2254/66.318; U; en) Presto/2.12.423 Version/12.16",
}

func getRandomUserAgent() string {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s) // initialize local pseudorandom generator
	i := r.Intn(len(userAgents))
	return userAgents[i]
}

const DefaultReferer = "https://www.google.com"

func GetHtmlDocFromUrl(url string, httpClient HttpClent, parameters *Parameters) (*html.Node, error) {
	response, err := httpClient.Request(url, map[string]string{
		"User-Agent": parameters.UserAgent,
		"Referer":    parameters.Referer,
	})
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	doc, err := html.Parse(response.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

type HttpClent interface {
	Request(url string, headers map[string]string) (*http.Response, error)
}

type DefaultHttpClient struct {
	client *http.Client
}

const Timeout = 15 * time.Second

func NewHttpClient() HttpClent {
	return &DefaultHttpClient{
		client: &http.Client{
			Timeout: Timeout,
		},
	}
}

func (c *DefaultHttpClient) Request(url string, headers map[string]string) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if headers["User-Agent"] == "" {
		headers["User-Agent"] = getRandomUserAgent()
	}

	if headers["Referer"] == "" {
		headers["Referer"] = DefaultReferer
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status code %d", response.StatusCode)
	}

	return response, nil
}
