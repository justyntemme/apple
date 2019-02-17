package webserver

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/justyntemme/apple/web"
)

func Start() {
	http.HandleFunc("/", serve)
	http.ListenAndServe(":8080", nil)

}

func serve(w http.ResponseWriter, r *http.Request) {
	res := web.Response{}
	p := web.People{}
	c := make(chan string)
	var wg sync.WaitGroup
	name := ""
	joke := ""

	wg.Add(2)
	go func() {
		defer wg.Done()
		name = web.FetchName(c, p)
	}()
	go func() {
		defer wg.Done()
		joke = web.FetchJoke(c, res)
	}()

	wg.Wait()

	fmt.Fprintln(w, " "+joke)

}
