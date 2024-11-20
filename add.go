package main

import (
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func handleAdd(commands []string) {
	if len(commands) <= 1 {
		fmt.Printf("vict: '%v' is not a vict command. See 'vict --help'. \n", commands)
		return
	}

	currDir, errInWd := os.Getwd()
	if errInWd != nil {
		log.Fatal(errInWd)
	}
	currDir += "/.vict"

	// check if index file already exists
	_, err := os.Stat(currDir + "/index")
	if os.IsNotExist(err) {
		indexFile, err := os.Create(currDir + "/index")
		if err != nil {
			log.Fatal(err)
		}
		defer indexFile.Close()
	}

	if commands[1] == "." { // TODO: Add a feature for adding all files
		fmt.Printf("This command is currently unavailable \n")
		return
	}
	// add specific files
	addSpecFiles(commands[1:], currDir)
}

func addSpecFiles(fileNames []string, currDir string) {
	// Open the index file for appending
	idxFile, err := os.OpenFile(currDir+"/index", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer idxFile.Close()

	for _, fileName := range fileNames {
		fileContent, err := os.ReadFile(fileName)
		if err != nil {
			log.Fatal(err)
		}

		sha1_hash := calculateHash(fileContent)
		pathLength := uint32(len(fileName))

		// Seek to the end of the file before writing
		_, err = idxFile.Seek(0, io.SeekEnd)
		if err != nil {
			log.Fatal(err)
		}

		idxFile.Write([]byte("100644")) // file mode (6 bytes)
		hashBytes, err := hex.DecodeString(sha1_hash)
		if err != nil {
			log.Fatal(err)
		}
		idxFile.Write(hashBytes)                            // hash (20 bytes)
		binary.Write(idxFile, binary.BigEndian, pathLength) // path length (4 bytes)
		idxFile.Write([]byte(fileName))                     // file path
	}
}

func calculateHash(fileContent []byte) string {
	contentSize := len(fileContent)
	// Follow the git standard
	strToBeHashed := "blob " + strconv.Itoa(contentSize) + "\x00" + string(fileContent)
	h := sha1.New()
	h.Write([]byte(strToBeHashed))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	return sha1_hash
}
