package main

import (
	"github.com/Vause/Reddit-word-of-the-day/pkg"
)

func main() {
	topSubReddits := []string{
		"funny",
	}

	submissionsChannel := make(chan []pkg.Submission)
	commentsChannel := make(chan []pkg.Comment)

	pkg.GetTopWord(topSubReddits, submissionsChannel, commentsChannel)
}
