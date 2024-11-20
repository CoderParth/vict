# Vict - A clone of Git (or at least a toy version of it)

Vict was my attempt to create a clone of git - at least a toy version of it – while trying to learn the mysterious ways that git works under the hood.
The commands that have been implemented so far are:

`init`

`init <repo-name>`

`add <file-name-1> <file-name-2> … <file-name-n>`

`<commit -m “Your-commit-message”`

`log`

A hidden “.vict” directory is created when you run “init” which is similar to how git creates a hidden “.git” directory, when you initialize a repository.

It is written in Go and makes fair use of its powerful concurrency.

In case you want to give it a try, please follow the steps below.

**_The following steps assume that you have go installed on your machine._**

## Start by running the go install command:

`go install github.com/coderparth/vict@latest`

## On Linux (Bash):

1. Open the `.bashrc` file in a text editor, then type the following:

   `nano ~/.bashrc`

2. Scroll to the bottom of the file and add the following line:

   `export PATH=$PATH:$HOME/go/bin`

3. In `nano`, press `CTRL + O` to save the file.

4. Reload the shell to apply the changes:

   `source ~/.bashrc`

5. DONE!!! Run the binary from any location:

   `vict`

---

## On macOS:

1. Determine your shell by running:
   `echo $SHELL`

2. Add `export PATH=$PATH:$HOME/go/bin` to `~/.zshrc` (or `~/.bash_profile` if using Bash).

3. Reload the shell:
   `source ~/.zshrc` (or `source ~/.bash_profile`)

4. DONE!!! Run the binary from any location:
   `vict`

---

## On Windows:

1. Go to **Environment Variables** from the Start Menu.

2. Add `C:\Users\your-username\go\bin` to the `PATH`.

3. Open a new **Command Prompt** or **PowerShell** and verify with:
   `echo %PATH%`

4. Run the binary from any location:
   `vict`

# Below are some examples of `vict` commands

```bash
codep@CODP:~/projects/victclone$ vict
These are Vict commands that are currently available for use:

start a working area
   init      Create an empty Vict repository or reinitialize an existing one

work on the current change
   add       Add file contents to the index

examine the history and state (see also: git help revisions)
   log       Show commit logs
   status    Show the working tree status

grow, mark and tweak your common history
   commit    Record changes to the repository

```

```

codep@CODP:~/projects/victclone$ vict init
Initialized empty Vict repository in /home/codep/projects/victclone/.vict
```

```
codep@CODP:~/projects/victclone$ vict add haha.txt
2024/11/20 20:56:40 open haha.txt: no such file or directory
```

```
codep@CODP:~/projects/victclone$ vict add a.txt
codep@CODP:~/projects/victclone$ vict commit -m "haha - first commit"
Committed as 4cbe2515b731392190e99b187032f54d7a3bb524
```

```
codep@CODP:~/projects/victclone$ vict log
commit 4cbe2515b731392190e99b187032f54d7a3bb524
    haha - first commit
```

```
codep@CODP:~/projects/victclone$ vict status
On branch main
Untracked files: [b.txt]
Changed files: []
Deleted files: []
```
