package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/rszewczyk/gotalks/intro/xkcd-server/xkcd"
)

const doctypeTmpl = `<!doctype html>`

var xclient = xkcd.Client{
	URL: "http://xkcd.com",
}

func main() {
	// Fetch multiple xkcd comics by ids and display them (in correct order)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var (
			ids   []int
			chans []chan io.WriterTo
		)
		for _, s := range strings.Split(r.URL.Path[1:], ".") {
			id, err := strconv.ParseInt(s, 10, 0)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			ids = append(ids, int(id))
			chans = append(chans, make(chan io.WriterTo))
		}

		for i, id := range ids {
			go worker(id, chans[i])
		}

		w.Write([]byte(doctypeTmpl))

		for _, ch := range chans {
			comic := <-ch
			comic.WriteTo(w)
		}
	})

	log.Fatalln(http.ListenAndServe(":7776", nil))
}

func worker(id int, result chan io.WriterTo) {
	comic, _ := xclient.GetComic(id)
	result <- comic
}
