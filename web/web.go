package web

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type People struct {
	FirstName string `json:"name"`
	Surname   string `json:"surname"`
}

type joke struct {
	ID   int    `json:"id"`
	Joke string `json:"joke"`
}

type Response struct {
	T     string `json:"type"`
	Value joke   `json:"value"`
}

func FetchName(path chan string, person People) string {
	req, err := fetch("http://uinames.com/api/?region=England")
	if err != nil {
		log.Fatal(err)
	} else {
		if err != nil {
			log.Fatal(err)
		} else {

			err = json.Unmarshal(req, &person)
			if err != nil {
				log.Fatal(err)
			}
			jokepath := "http://api.icndb.com/jokes/random?firstName=" + person.FirstName + "&lastName=" + person.Surname + "&limitTo=[nerdy]"
			path <- jokepath
			return (person.FirstName + " " + person.Surname)
		}

	}
	return "ERR"

}

func FetchJoke(path chan string, r Response) string {
	urlsArray := []string{}

	select {
	case path := <-path:
		urlsArray = append(urlsArray, path)
	}

	resc, errc := make(chan string), make(chan error)

	for _, url := range urlsArray {
		go func(url string) {
			body, err := fetch(url)
			if err != nil {
				errc <- err
				return
			}
			err = json.Unmarshal(body, &r)
			if err != nil {
				log.Fatal(err)
			}
			joke := r.Value.Joke
			resc <- joke
		}(url)
	}

	for i := 0; i < len(urlsArray); i++ {
		select {
		case res := <-resc:
			return res
			log.Print(res)
		case err := <-errc:
			log.Fatal(err)
		}
	}
	return "ERR"
}

func fetch(url string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "Chrome: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}
