package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	url            string
	payload        string
	threads        int
	requestCounter int
	printedMsgs    []string
	waitGroup      sync.WaitGroup
)

func printMsg(msg string) {
	if !contains(printedMsgs, msg) {
		fmt.Printf("\n%s after %d requests\n", msg, requestCounter)
		printedMsgs = append(printedMsgs, msg)
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func handleStatusCodes(statusCode int) {
	requestCounter++
	fmt.Printf("\r%d requests have been sent", requestCounter)

	if statusCode == 429 {
		printMsg("You have been throttled")
	}
	if statusCode == 500 {
		printMsg("Status code 500 received")
	}
}

func sendGET() {
	defer waitGroup.Done()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", randomUserAgent())

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	handleStatusCodes(resp.StatusCode)
}

func sendPOST() {
	defer waitGroup.Done()

	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", randomUserAgent())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	handleStatusCodes(resp.StatusCode)
}

func randomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:54.0) Gecko/20100101 Firefox/54.0",
		"Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:40.0) Gecko/20100101 Firefox/40.1",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) Gecko/20100101 Firefox/34.0",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/30.0.1599.134 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.94 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/23.0.1271.97 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/19.0.1084.52 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/17.0.963.79 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/15.0.874.121 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/14.0.835.202 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/13.0.782.220 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/9.0.597.98 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/8.0.622.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/7.0.517.44 Safari/537.36",
	}

	return userAgents[rand.Intn(len(userAgents))]
}

func main() {
	flag.StringVar(&url, "g", "", "Specify GET request. Usage: -g '<url>'")
	flag.StringVar(&url, "p", "", "Specify POST request. Usage: -p '<url>'")
	flag.StringVar(&payload, "d", "", "Specify data payload for POST request")
	flag.IntVar(&threads, "t", 500, "Specify number of threads to be used")
	flag.Parse()

	if url == "" {
		flag.Usage()
		return
	}

	if flag.NFlag() != 4 && flag.NFlag() != 5 {
		fmt.Println("Incorrect number of flags provided.")
		flag.Usage()
		return
	}

	waitGroup.Add(threads)

	for i := 0; i < threads; i++ {
		if url != "" {
			if payload != "" {
				go sendPOST()
			} else {
				go sendGET()
			}
		}
	}
	waitGroup.Wait()
}
