package main

import "fmt"

func handleLog() {
	currentCommitHash := getCurrentCommitHash()
	for currentCommitHash != "" && currentCommitHash != "0000000000000000000000000000000000000000" {
		commitData, err := readCommitObject(currentCommitHash)
		if err != nil {
			fmt.Println("Error reading commit:", err)
			break
		}
		fmt.Printf("commit %s\n", currentCommitHash)
		fmt.Printf("    %s\n", commitData.Message)
		// Move to the parent commit
		currentCommitHash = commitData.ParentCommit
	}
}
