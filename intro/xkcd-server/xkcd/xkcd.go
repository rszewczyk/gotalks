package xkcd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const comicTemplate = `<div>
<h3>%s</h3>
<img src="%s">
</div>
`

// Response represents the data from an xkcd API call
type apiResult struct {
	Title string `json:"safe_title"`
	Image string `json:"img"`
}

// WriteTo implements the io.WriterTo interface
func (r *apiResult) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprintf(w, comicTemplate, r.Title, r.Image)
	return int64(n), err
}

// Client fetches comics from the configured xkcd URL
type Client struct {
	URL string

	mu    sync.RWMutex //protects the cache
	cache map[int]apiResult
	once  sync.Once
}

func (client *Client) init() {
	client.once.Do(func() {
		client.cache = make(map[int]apiResult)
	})
}

func (client *Client) buildURL(id int) string {
	url := client.URL
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}

	// 0 fetches current comic
	if id != 0 {
		url = fmt.Sprintf("%s%d/", url, id)
	}

	return url + "info.0.json"
}

// GetComic fetches the xkcd comic with the given id
func (client *Client) GetComic(id int) (io.WriterTo, error) {
	defer func(start time.Time) {
		log.Printf("GetComic %d took %d ms\n", id, time.Since(start)/time.Millisecond)
	}(time.Now())

	client.init()

	if comic := client.getFromCache(id); comic != nil {
		return comic, nil
	}

	apiResponse, err := http.Get(client.buildURL(id))
	if err != nil {
		return nil, fmt.Errorf("Error connecting to server: %v", err)
	}
	defer apiResponse.Body.Close()

	data, err := ioutil.ReadAll(io.LimitReader(apiResponse.Body, 64e3))
	if err != nil {
		return nil, fmt.Errorf("Could not read response: %v", err)
	}

	comic := &apiResult{}
	err = json.Unmarshal(data, comic)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal reponse: %v", err)
	}

	client.putInCache(id, comic)

	return comic, nil
}

func (client *Client) getFromCache(id int) *apiResult {
	client.mu.RLock()
	defer client.mu.RUnlock()

	if r, ok := client.cache[id]; ok {
		return &r
	}

	return nil
}

func (client *Client) putInCache(id int, response *apiResult) {
	client.mu.Lock()
	defer client.mu.Unlock()

	client.cache[id] = *response
}
