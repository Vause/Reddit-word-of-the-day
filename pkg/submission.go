package pkg

import (
	"encoding/json"
	"fmt"
	"sync"
)

const SubmissionURL = "https://api.pushshift.io/reddit/search/submission/?subreddit=%s&sort_type=num_comments&sort=desc&after=3d&fields=id"

//Struct for all data that is returned from pushshift for a submission query
type SubmissionResponse struct {
	Data []Submission `json:"data"`
}

//Struct for submission link IDs
type Submission struct {
	Id string `json:"id"`
}

func getTopSubmissionsHelper(topSubReddits []string, submissions chan []Submission) {
	var wg sync.WaitGroup

	for _, subReddit := range topSubReddits {
		wg.Add(1)
		go getTopSubmissions(subReddit, submissions, &wg)
	}

	go func() {
		wg.Wait()
		close(submissions)
	}()
}

func getTopSubmissions(topSubReddit string, c chan<- []Submission, wg *sync.WaitGroup) {
	defer (*wg).Done()
	responseData := Get(fmt.Sprintf(SubmissionURL, topSubReddit))

	var responseObject SubmissionResponse
	json.Unmarshal(responseData, &responseObject)

	for len(responseObject.Data) < 1 {
		newResponse := retryGetSubmission(topSubReddit)
		responseObject.Data = newResponse
	}

	c <- responseObject.Data
}

func retryGetSubmission(topSubReddit string) (newResponse []Submission) {
	responseData := Get(fmt.Sprintf(SubmissionURL, topSubReddit))

	var responseObject SubmissionResponse
	json.Unmarshal(responseData, &responseObject)

	return responseObject.Data
}
