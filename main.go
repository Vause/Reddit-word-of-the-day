package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

//Struct for all data that is returned from pushshift for a submission query
type SubmissionResponse struct {
	Data []Submission `json:"data"`
}

//Struct for all data that is returned from pushshift for a comment query
type CommentResponse struct {
	Data []Comment `json:"data"`
}

//Struct for submission link IDs
type Submission struct {
	Id string `json:"id"`
}

//Struct for comment body
type Comment struct {
	Body string `json:"body"`
}

func main() {
	topSubReddits := []string{
		"funny",
		"AskReddit",
		"gaming",
		"aww",
		"Music",
		"pics",
		"worldnews",
		"science",
		"todayilearned",
	}

	getTopLinksHelper(topSubReddits)
}

func getTopLinksHelper(topSubReddits []string) {
	c := make(chan []Submission)
	var wg sync.WaitGroup

	for _, subReddit := range topSubReddits {
		wg.Add(1)
		go getTopLinks(subReddit, c, &wg)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for msg := range c {
		getTopCommentsHelper(msg)
	}
}

func getTopCommentsHelper(topSubRedditLinkIds []Submission) {
	c1 := make(chan []Comment)
	var wg sync.WaitGroup

	for _, subRedditLinkId := range topSubRedditLinkIds {
		wg.Add(1)
		go getTopComments(subRedditLinkId, c1, &wg)
	}

	go func() {
		wg.Wait()
		close(c1)
	}()

	for msg := range c1 {
		for i := 0; i < len(msg); i++ {
			fmt.Println(msg[i].Body)
		}
	}
}

func getTopComments(subRedditLinkId Submission, c1 chan []Comment, wg *sync.WaitGroup) {
	defer (*wg).Done()
	response, err := http.Get(fmt.Sprintf("https://api.pushshift.io/reddit/search/comment/?link_id=t3_%s&sort_type=score&sort=desc&score=%3E100&fields=body", subRedditLinkId.Id))

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject CommentResponse
	json.Unmarshal(responseData, &responseObject)

	c1 <- responseObject.Data
}
func getTopLinks(topSubReddit string, c chan []Submission, wg *sync.WaitGroup) {
	defer (*wg).Done()
	response, err := http.Get(fmt.Sprintf("https://api.pushshift.io/reddit/search/submission/?subreddit=%s&sort_type=num_comments&sort=desc&after=3d&fields=id", topSubReddit))

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject SubmissionResponse
	json.Unmarshal(responseData, &responseObject)

	c <- responseObject.Data
}
