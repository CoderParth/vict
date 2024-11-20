package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func handleCommit(commands []string) {
	if len(commands) == 1 {
		fmt.Printf("vict: '%v' is not a vict command. See 'vict --help'. \n", commands)
		return
	}
	if commands[1] != "-m" {
		fmt.Printf("vict: '%v' is not a vict command. See 'vict --help'. \n", commands)
		return
	}

	msg := strings.Join(commands[2:], " ")
	if len(msg) == 0 {
		fmt.Println("error: switch `m' requires a value ")
		return
	}

	parentCommitHash := getCurrentCommitHash()
	treeHash := createTreeObject()                                    // represents the current file state
	commitHash := createCommitObject(msg, treeHash, parentCommitHash) // new commit
	updateHead(commitHash)                                            // update the head file to point to new commit
	fmt.Printf("Committed as %s\n", commitHash)
}

func createTreeObject() string {
	var buffer bytes.Buffer
	for _, data := range idxFileData {
		// entry = (file mode + hash + path)
		entry := fmt.Sprintf("100644 %s %s", data.Hash, data.FilePath)
		buffer.WriteString(entry)
	}
	// Create a tree object (zlib compression)
	treeData := buffer.String()
	compressedTreeData := compress(treeData)

	treeHash := calculateShaOfObject(compressedTreeData)
	writeObjectToFile(treeHash, compressedTreeData)
	return treeHash
}

func getCurrentCommitHash() string {
	currDir, errInWd := os.Getwd()
	if errInWd != nil {
		log.Fatal(errInWd)
	}
	currDir += "/.vict/HEAD"

	headFile, err := os.Open(currDir)
	if err != nil {
		log.Fatal(err)
	}
	defer headFile.Close()

	var ref string
	_, err = fmt.Fscanf(headFile, "ref: refs/heads/main\n%s", &ref)
	if err != nil {
		return ""
	}

	return ref // This is the current commit hash
}

func createCommitObject(msg string, treeHash string, parentCommitHash string) string {
	// Prepare the commit data with parent reference (if any)
	commitData := fmt.Sprintf("tree %s\n", treeHash)
	// If this is not the first commit, include the parent commit hash
	if parentCommitHash != "" {
		commitData += fmt.Sprintf("parent %s\n", parentCommitHash)
	}

	commitData += fmt.Sprintf("\n%s\n", msg)
	compressedCommitData := compress(commitData)
	commitHash := calculateShaOfObject(compressedCommitData)
	writeObjectToFile(commitHash, compressedCommitData)
	return commitHash
}

// Save the tree object to disk
func writeObjectToFile(hash string, data []byte) {
	objectPath := fmt.Sprintf(".vict/objects/%s", hash[:2]) // First two characters of the hash
	err := os.MkdirAll(objectPath, 0755)
	if err != nil {
		log.Fatal(err)
	}
	// Write the object to the file in objects folder
	filePath := fmt.Sprintf("%s/%s", objectPath, hash[2:])
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.Write(data)
}

// Open HEAD file and write the new commit reference
func updateHead(commitHash string) {
	currDir, errInWd := os.Getwd()
	if errInWd != nil {
		log.Fatal(errInWd)
	}
	currDir += "/.vict/HEAD"
	headFile, err := os.Create(currDir)
	if err != nil {
		log.Fatal(err)
	}
	defer headFile.Close()
	_, err = fmt.Fprintf(headFile, "ref: refs/heads/main\n%s\n", commitHash)
	if err != nil {
		log.Fatal(err)
	}
}

func readCommitObject(commitHash string) (CommitData, error) {
	currDir, errInWd := os.Getwd()
	if errInWd != nil {
		log.Fatal(errInWd)
	}
	currDir += "/.vict/objects/" + commitHash[:2] // First 2 characters form a directory

	// Read the compressed commit file
	commitFilePath := currDir + "/" + commitHash[2:]
	commitFile, err := os.Open(commitFilePath)
	if err != nil {
		return CommitData{}, fmt.Errorf("failed to open commit object: %w", err)
	}
	defer commitFile.Close()

	// Decompress the commit data
	var decompressedData bytes.Buffer
	zlibReader, err := zlib.NewReader(commitFile)
	if err != nil {
		return CommitData{}, fmt.Errorf("failed to create zlib reader: %w", err)
	}
	_, err = io.Copy(&decompressedData, zlibReader)
	if err != nil {
		return CommitData{}, fmt.Errorf("failed to decompress commit data: %w", err)
	}

	// Read the commit data (tree, parent, message)
	commitData := decompressedData.String()
	var parentCommit, treeHash, message string

	_, err = fmt.Sscanf(commitData, "tree %s\nparent %s\n", &treeHash, &parentCommit)
	if err != nil {
		// Handle first commit (no parent)
		_, err = fmt.Sscanf(commitData, "tree %s\n", &treeHash)
		if err != nil {
			return CommitData{}, fmt.Errorf("failed to parse commit data: %w", err)
		}
		parentCommit = "" // First commit has no parent
	}

	// read the commit message from the rest of the data (after tree and parent)
	messageStartIndex := strings.Index(commitData, "\n\n") + 2 // Find where the message starts
	message = commitData[messageStartIndex:]
	message = strings.TrimSpace(message)

	return CommitData{
		Message:      message,
		ParentCommit: parentCommit,
		TreeHash:     treeHash,
	}, nil
}

func calculateShaOfObject(data []byte) string {
	h := sha1.New()
	h.Write(data)
	commitHash := hex.EncodeToString(h.Sum(nil))
	return commitHash
}
