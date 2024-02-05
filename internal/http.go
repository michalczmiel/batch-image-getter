package internal

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/net/html"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.3",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Safari/605.1.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.3",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:122.0) Gecko/20100101 Firefox/122.",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.3",
}

func getRandomUserAgent() string {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s) // initialize local pseudorandom generator
	i := r.Intn(len(userAgents))
	return userAgents[i]
}

const DefaultReferer = "https://www.google.com"

func GetHtmlDocFromUrl(url string, httpClient HttpClient, parameters *Parameters) (*html.Node, error) {
	response, err := httpClient.Request(url, map[string]string{
		"Referer": parameters.Referer,
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

type HttpClient interface {
	Request(url string, headers map[string]string) (*http.Response, error)
}

type DefaultHttpClient struct {
	client    *http.Client
	userAgent string
}

const Timeout = 15 * time.Second

func NewHttpClient(userAgent string) HttpClient {
	if userAgent == "" {
		userAgent = getRandomUserAgent()
	}

	return &DefaultHttpClient{
		client: &http.Client{
			Timeout: Timeout,
		},
		userAgent: userAgent,
	}
}

func (c *DefaultHttpClient) Request(url string, headers map[string]string) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if headers["Referer"] == "" {
		headers["Referer"] = DefaultReferer
	}

	headers["User-Agent"] = c.userAgent

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
