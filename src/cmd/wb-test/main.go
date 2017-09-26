package main

import (
	"bufio"
	"log"
	"os"
	"net/http"
	"io/ioutil"
	"strings"
)

func init() {
	log.SetFlags(0)
}

func main() {
	k := 2;
	processingChannel := make(chan string, k); // Channel for URLs
	doneChan := make(chan bool);               // Used when all urls handled
	results := make(map[string]int);
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		url := scanner.Text();

		// Run handler
		processingChannel <- url

		go handleUrl(url, results, processingChannel, doneChan)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	// waiting for all handlers
	<-doneChan;

	total := 0

	// Printing results for each url
	for key, value := range results {
		log.Println(key, ':', value)
		total += value;
	}

	log.Printf("Total: %v", total)
}

func handleUrl(url string, results map[string]int, processingChannel chan string, doneChan chan bool) {
	// Go and get count
	res := countGoEntries(request(url));

	// Store result
	results[url] = res;

	// Pull processed url so next can be pushed
	<-processingChannel;
	// if no url pushed we are done, call fot it
	if len(processingChannel) == 0 {
		doneChan <- true;
	};
}

func request(url string) string {
	res, err := http.Get(url)

	if err != nil {
		log.Fatalln(err);
	}

	bodyBytes, _ := ioutil.ReadAll(res.Body);
	return string(bodyBytes);
}

func countGoEntries(s string) int {
	return strings.Count(s, "Go")
}
