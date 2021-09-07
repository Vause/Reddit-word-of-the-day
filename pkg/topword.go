package pkg

import "fmt"

func GetTopWord(topSubReddits []string, submissionsChannel chan []Submission, commentsChannel chan []Comment) {
	getTopSubmissionsHelper(topSubReddits, submissionsChannel)
	for submissionLinkIds := range submissionsChannel {
		getTopCommentsHelper(submissionLinkIds, commentsChannel)
	}

	for topComments := range commentsChannel {
		for i := 0; i < len(topComments); i++ {
			fmt.Println(topComments[i].Body)
		}
	}
}
