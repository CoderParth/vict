package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func handleStatus() {
	readIndexFile()
	walkDir()

	fmt.Println("On branch main")
	fmt.Printf("Untracked files: %v \n", statusData.UntrackedFiles)
	fmt.Printf("Changed files: %v \n", statusData.ChangedFiles)

	var deletedFiles []string
	for _, data := range idxFileData {
		if _, ok := filesInDir.files[data.FilePath]; !ok {
			deletedFiles = append(deletedFiles, data.FilePath)
		}
	}

	fmt.Printf("Deleted files: %v \n", deletedFiles)
}

func walkDir() {
	if err := filepath.Walk(".", scanPath); err != nil {
		fmt.Printf("Error scanning the directory %v:\n", err)
	}
	wg.Wait()
}

func scanPath(path string, info os.FileInfo, err error) error {
	// Get the name of the current file path
	_, name := filepath.Split(path)

	if name == ".git" || name == ".vict" {
		if info.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}

	if !info.IsDir() {
		wg.Add(1)
		go compareFileToIndex(path) // Compare it to the .vict/index file
	}
	return nil
}

func compareFileToIndex(path string) {
	defer wg.Done()

	// Add the current path to the "filesInDir.files" map.
	// This is used later to track the deleted files.
	filesInDir.mu.Lock()
	filesInDir.files[path] = 0
	filesInDir.mu.Unlock()

	for _, data := range idxFileData {
		if data.FilePath == path {
			// File present, check the hash
			fileContent, err := os.ReadFile(path)
			if err != nil {
				log.Fatal(err)
			}

			sha1_hash := calculateHash(fileContent)

			if sha1_hash != data.Hash {
				// hash has changed, add it to changed files.
				statusData.mu.Lock()
				statusData.ChangedFiles = append(statusData.ChangedFiles, path)
				statusData.mu.Unlock()
			}

			return
		}
	}

	// Else if the file is not in the index, it is an untracked file.
	statusData.mu.Lock()
	statusData.UntrackedFiles = append(statusData.UntrackedFiles, path)
	statusData.mu.Unlock()
}

func readIndexFile() {
	currDir, errInWd := os.Getwd()
	if errInWd != nil {
		log.Fatal(errInWd)
	}
	currDir += "/.vict"

	idxFile, err := os.Open(currDir + "/index")
	if err != nil {
		log.Fatal("No files have been added yet. Please run vict add <filename>")
		return
	}
	defer idxFile.Close()
	readIdxInChunks(idxFile)
}

func readIdxInChunks(idxFile *os.File) {
	for {
		// Read the file mode (6 bytes)
		var mode [6]byte
		_, err := idxFile.Read(mode[:])
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		// Read the SHA-1 hash (20 bytes)
		hash := make([]byte, 20)
		readBytes(idxFile, hash)
		hashStr := hex.EncodeToString(hash)

		// Read the path length (4 bytes, uint32)
		var lengthBytes [4]byte
		readBytes(idxFile, lengthBytes[:])
		pathLength := binary.BigEndian.Uint32(lengthBytes[:])

		// Read the file path (variable length)
		filePath := make([]byte, pathLength)
		readBytes(idxFile, filePath)

		data := IndexFileData{
			Hash:     hashStr,
			FilePath: string(filePath),
		}

		idxFileData = append(idxFileData, data)
	}
}

func readBytes(idxFile *os.File, byteLen []byte) {
	_, err := idxFile.Read(byteLen)
	if err != nil {
		log.Fatal(err)
	}
}
