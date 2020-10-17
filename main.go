package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/valyala/fasthttp"
)

func main() {
	var c int
	flag.IntVar(&c, "c", 50, "Set the Concurrency ")
	flag.Parse()
	inputs := make(chan string)
	var wg sync.WaitGroup
	input := bufio.NewScanner(os.Stdin)
	go func() {
		for input.Scan() {
			inputs <- input.Text()
		}
		close(inputs)
	}()
	for i := 0; i < c; i++ {
		wg.Add(1)
		go workers(inputs, &wg)
	}
	wg.Wait()
}

func buildurl(s string) {
	ur, err := url.Parse(s)
	if err != nil {
		return
	}
	x := ur.Query()
	if len(x) == 0 {
		return
	}
	baseurl := ur.Scheme + "://" + ur.Host + ur.Path + "?"
	params := url.Values{}
	for i := range x {
		params.Add(i, "<'\">")
	}
	finalurl := baseurl + params.Encode()
	//fmt.Printf("Testing %s \t", finalurl)
	if checkxss(finalurl) {
		fmt.Println(s, "might be vulnerable to xss")
	}
}

func checkErr(e error) {
	if e != nil {
		fmt.Println(e)
	}

}

func checkxss(s string) bool {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)   // <- do not forget to release
	defer fasthttp.ReleaseResponse(resp) // <- do not forget to release

	req.SetRequestURI(s)

	fasthttp.Do(req, resp)

	bodyBytes := resp.Body()
	if strings.Contains(string(bodyBytes), "<'\">") {
		return true
	}
	return false
}

func workers(cha chan string, wg *sync.WaitGroup) {
	for i := range cha {
		buildurl(i)
	}
	wg.Done()
}
