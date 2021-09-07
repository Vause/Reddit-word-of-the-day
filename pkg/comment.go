package pkg

import (
	"encoding/json"
	"fmt"
	"sync"
)

const CommentURL = "https://api.pushshift.io/reddit/search/comment/?link_id=t3_%s&sort_type=score&sort=desc&score=%3E10&fields=body"

//Struct for all data that is returned from pushshift for a comment query
type CommentResponse struct {
	Data []Comment `json:"data"`
}

//Struct for comment body
type Comment struct {
	Body string `json:"body"`
}

func getTopCommentsHelper(topSubRedditLinkIds []Submission, comments chan []Comment) {
	var wg sync.WaitGroup

	for _, subRedditLinkId := range topSubRedditLinkIds {
		wg.Add(1)
		go getTopComments(subRedditLinkId, comments, &wg)
	}

	go func() {
		wg.Wait()
		close(comments)
	}()
}

func getTopComments(subRedditLinkId Submission, c1 chan<- []Comment, wg *sync.WaitGroup) {
	defer (*wg).Done()

	responseData := Get(fmt.Sprintf(CommentURL, subRedditLinkId.Id))

	var responseObject CommentResponse
	json.Unmarshal(responseData, &responseObject)

	for len(responseObject.Data) < 1 {
		newResponse := retryGetComment(subRedditLinkId)
		responseObject.Data = newResponse
	}

	c1 <- responseObject.Data
}

func retryGetComment(subRedditLinkId Submission) (newResponse []Comment) {
	responseData := Get(fmt.Sprintf(CommentURL, subRedditLinkId.Id))

	var responseObject CommentResponse
	json.Unmarshal(responseData, &responseObject)

	return responseObject.Data
}
