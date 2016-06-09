package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/rszewczyk/gotalks/intro/xkcd-server/xkcd"
)

const doctypeTmpl = `<!doctype html>`

var xclient = xkcd.Client{
	URL: "http://xkcd.com",
}

func main() {
	// Fetch a single xkcd comic by id and display it
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.URL.Path[1:], 10, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		comic, err := xclient.GetComic(int(id))
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		w.Write([]byte(doctypeTmpl))
		comic.WriteTo(w)
	})

	log.Fatalln(http.ListenAndServe(":7777", nil))
}
