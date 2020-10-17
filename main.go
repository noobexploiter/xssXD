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
	var s string
	flag.IntVar(&c, "c", 50, "Set the Concurrency ")
	flag.StringVar(&s, "s", "none", "Specify the payload to use")
	flag.Parse()
	inputs := make(chan string)
	var wg sync.WaitGroup
	input := bufio.NewScanner(os.Stdin)
	go func() {
		for input.Scan() {
			ur, err := url.Parse(input.Text())
			if err != nil {
				continue
			}
			x := ur.Query()
			if len(x) == 0 {
				continue
			}
			inputs <- input.Text()
		}
		close(inputs)
	}()
	for i := 0; i < c; i++ {
		wg.Add(1)
		go workers(inputs, &wg, s)
	}
	wg.Wait()
}

func buildurl(s string, st string) {
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
	if st != "none" {
		for i := range x {
			params.Add(i, st)
		}
		finalurl := baseurl + params.Encode()
		if specifiedpayload(finalurl, st) {
			fmt.Println(s, "is reflecting", st)
		}
	} else {
		for i := range x {
			params.Add(i, "ab1<ab3\"ab4>")
		}
		finalurl := baseurl + params.Encode()
		chars := checkxss(finalurl)
		if len(chars) == 0 {
			return
		}
		fmt.Println(s, "is reflecting", strings.Join(chars, ", "))
	}

}

func checkErr(e error) {
	if e != nil {
		fmt.Println(e)
	}

}

func checkxss(s string) []string {
	//fmt.Println("TESTING", s)
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)   // <- do not forget to release
	defer fasthttp.ReleaseResponse(resp) // <- do not forget to release

	req.SetRequestURI(s)

	fasthttp.Do(req, resp)

	allowedchars := []string{}

	bodyBytes := resp.Body()
	if strings.Contains(string(bodyBytes), "ab1<") {
		allowedchars = append(allowedchars, "<")
	}
	if strings.Contains(string(bodyBytes), "ab3\"") {
		allowedchars = append(allowedchars, "\"")
	}
	if strings.Contains(string(bodyBytes), "ab4>") {
		allowedchars = append(allowedchars, ">")
	}
	return allowedchars
}

func workers(cha chan string, wg *sync.WaitGroup, s string) {
	for i := range cha {
		buildurl(i, s)
	}
	wg.Done()
}

func specifiedpayload(s string, st string) bool {
	//fmt.Println("TESTING", s)
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)   // <- do not forget to release
	defer fasthttp.ReleaseResponse(resp) // <- do not forget to release

	req.SetRequestURI(s)

	fasthttp.Do(req, resp)

	bodyBytes := resp.Body()
	if strings.Contains(string(bodyBytes), st) {
		return true
	}
	return false
}
