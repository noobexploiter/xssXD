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
		params.Add(i, "ab1<ab2'ab3\"ab4>")
	}
	finalurl := baseurl + params.Encode()
	chars := checkxss(finalurl)
	if len(chars) == 0 {
		return
	}
	fmt.Println(s, "is reflecting", strings.Join(chars, ", "))
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
	defer fasthttp.ReleaseRequest(req)   
	defer fasthttp.ReleaseResponse(resp) 

	req.SetRequestURI(s)

	fasthttp.Do(req, resp)

	allowedchars := []string{}

	bodyBytes := resp.Body()
	if strings.Contains(string(bodyBytes), "ab1<") {
		allowedchars = append(allowedchars, "<")
	}
	if strings.Contains(string(bodyBytes), "ab2'") {
		allowedchars = append(allowedchars, "'")
	}
	if strings.Contains(string(bodyBytes), "ab3\"") {
		allowedchars = append(allowedchars, "\"")
	}
	if strings.Contains(string(bodyBytes), "ab4>") {
		allowedchars = append(allowedchars, ">")
	}
	return allowedchars
}

func workers(cha chan string, wg *sync.WaitGroup) {
	for i := range cha {
		buildurl(i)
	}
	wg.Done()
}
