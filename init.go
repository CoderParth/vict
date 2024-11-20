package main

import (
	"fmt"
	"log"
	"os"
)

func handleInit(commands []string) {
	repoName := ""
	if len(commands) > 1 { // a repo name was provided
		repoName = commands[1]
	} // Else, no repo name was provided by user, during initialization
	// so initialize in the current dir

	currDir, errInWd := os.Getwd()
	if errInWd != nil {
		log.Fatal(errInWd)
	}

	if len(repoName) != 0 {
		err := os.Mkdir(repoName, 0755)
		if err != nil {
			if os.IsExist(err) {
				fmt.Printf("Repo: %v  already exists in %v \n", repoName, currDir)
				return
			}
			log.Fatal(err)
		}
		currDir += "/" // Adding / to make the directory right for the new repo.
	}

	currDir += repoName
	currDir += "/.vict"

	// create .vict directory
	err := os.Mkdir(currDir, 0755)
	if err != nil {
		if os.IsExist(err) {
			fmt.Printf("Reinitialized existing Vict repository in %v \n", currDir)
			return
		}
		log.Fatal(err)
	}

	// Create essential files and dirs
	createFilesAndDirs(currDir)
	fmt.Printf("Initialized empty Vict repository in %v \n", currDir)
}

func createFilesAndDirs(currDir string) {
	files := [3]string{"/HEAD", "/config", "/description"}
	dirs := [4]string{"/hooks", "/info", "/objects", "/refs"}
	for _, f := range files {
		createFile(currDir+f, currDir)
	}
	for _, d := range dirs {
		createDir(currDir + d)
	}
}

func createFile(filePath string, currDir string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	if filePath == currDir+"/HEAD" {
		file.WriteString("ref: refs/heads/main\n")
	}
	defer file.Close()
}

func createDir(dir string) {
	err := os.Mkdir(dir, 0755)
	if err != nil {
		log.Fatal(err)
	}
}
