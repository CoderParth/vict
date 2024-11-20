package main

import "sync"

type IndexFileData struct {
	Hash     string
	FilePath string
}

type Status struct {
	UntrackedFiles []string
	ChangedFiles   []string
	DeletedFiles   []string
	mu             sync.Mutex
}

type CommitData struct {
	Message      string
	ParentCommit string
	TreeHash     string
}

type DirFiles struct {
	files map[string]int
	mu    sync.Mutex
}

var (
	wg          sync.WaitGroup
	idxFileData []IndexFileData
	statusData  Status
	filesInDir  DirFiles = DirFiles{files: make(map[string]int)}
)

// initialize files map

var availableCommands string = `These are Vict commands that are currently availabe for use:

start a working area 
   init      Create an empty Vict repository or reinitialize an existing one

work on the current change 
   add       Add file contents to the index

examine the history and state (see also: git help revisions)
   log       Show commit logs
   status    Show the working tree status

grow, mark and tweak your common history
   commit    Record changes to the repository
  `
