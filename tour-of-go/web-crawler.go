package main

import (
	"fmt"
	"sync"
	"runtime"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)


}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.

func Crawl(url string, depth int, fetcher Fetcher) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:

	var resultMap map[string][]string
	resultMap = make(map[string][]string)

	var fet func(u string, dep int, f Fetcher, r map[string][]string, w *sync.WaitGroup)

	var cache func(r map[string][]string, ch chan []string)

	cache = func(r map[string][]string, ch chan []string){

		for strs := range ch {
			key := strs[0]
			val := strs[1:]
			//fmt.Println(key,"--",val)
			r[key] = val
		}
	}

	ch := make(chan []string)
	go cache(resultMap, ch)

	fet = func(u string, dep int, f Fetcher, r map[string][]string, w *sync.WaitGroup){

		if dep <= 0 {
			if w != nil{
				w.Done()
			}
			return
		}
		var body_urls []string
		var body string
		var urls []string
		var err error
		var ok bool

		body_urls , ok = r[u]

		if ! ok {
			body, urls, err = f.Fetch(u)
			if err != nil {
				fmt.Println(err)
				if w != nil{
					w.Done()
				}
				return
			} else {
				strs := []string{}
				strs = append(strs, u)
				strs = append(strs, body)
				for _ , val:= range urls {
					strs = append(strs, val)
				}

				ch <- strs
				fmt.Printf("found: %s %q\n", u, body)
			}

		} else {
			body = body_urls[0]
			urls = body_urls[1:]
			fmt.Printf("cached found: %s %q\n", u, body)
		}

		var wg sync.WaitGroup
		wg.Add(len(urls))
		for _, ur := range urls {

			go fet(ur, dep-1, f, r, &wg)

		}
		wg.Wait()
		if w != nil{
			w.Done()
		}
	}

	fet(url, depth, fetcher, resultMap, nil)
	return
}

func main() {
	runtime.GOMAXPROCS(4)
	Crawl("http://golang.org/", 4, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult


type fakeResult struct {
	body string
	urls []string
}


func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}


// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/":   &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}

