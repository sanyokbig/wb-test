package main

import (
	"bufio"
	"log"
	"os"
	"net/http"
	"io/ioutil"
	"strings"
)

type Result struct {
	url   string
	count int
}

func init() {
	log.SetFlags(0)
}

func main() {
	k := 5;
	processingChannel := make(chan string, k); // Channel for URLs
	doneChan := make(chan bool);               // Used when all urls handled
	results := []Result{}                      // Results array
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		url := scanner.Text();

		// Waiting for slot in channel
		processingChannel <- url

		// Run handler
		go handleUrl(url, &results, processingChannel, doneChan)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	// waiting for all handlers
	<-doneChan;
	total := 0

	// Printing results for each url
	for _, res := range results {
		log.Println(res.url, ':', res.count)
		total += res.count;
	}

	log.Printf("Total: %v", total)
}

func handleUrl(url string, resultsPtr *[]Result, processingChannel chan string, doneChan chan bool) {
	// Go and get count
	var count int;
	body, err := request(url);

	if err != nil {
		// Instead of running counter, set counter to zero
		count = 0;
	} else {
		count = countGoEntries(body);
	}

	// Store result
	result := Result{url, count}
	*resultsPtr = append(*resultsPtr, result);

	// Pull processed url so next can be pushed
	<-processingChannel;
	// if no url pushed we are done, call fot it
	if len(processingChannel) == 0 {
		doneChan <- true;
	};
}

func request(url string) (string, error) {
	res, err := http.Get(url)

	if err != nil {
		log.Println(err);
		return "", err
	}

	bodyBytes, err := ioutil.ReadAll(res.Body);
	if err != nil {
		log.Println(err);
		return "", err
	}

	return string(bodyBytes), nil;
}

func countGoEntries(s string) int {
	return strings.Count(s, "Go")
}
