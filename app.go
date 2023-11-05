package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
)

type App struct {
	Repo       string
	Branch     string
	MainFolder string
	RepoFolder string
}

func NewApp(repo string, branch string, mainFolder string) *App {
	repoSSHParts := strings.Split(repo, "/")
	repoName := repoSSHParts[len(repoSSHParts)-1]
	repoName = strings.TrimSuffix(repoName, ".git")

	a := &App{
		Repo:       repo,
		Branch:     branch,
		MainFolder: mainFolder,
		RepoFolder: mainFolder + repoName + "/",
	}

	return a
}

func (a *App) initRepo() error {
	log.Printf("Making the directory %s safe", a.RepoFolder)
	err := executeCommand(a.RepoFolder, "git", "config", "--global", "--add", "safe.directory", "'*'")

	_, err = os.Stat(a.RepoFolder)
	if !os.IsNotExist(err) {
		log.Print("Repository already initialized, fetching new changes")
		return a.fetchChanges()
	}

	log.Print("Initializing repository")
	err = executeCommand(a.MainFolder, "git", "clone", a.Repo)
	if err != nil {
		return err
	}

	err = executeCommand(a.RepoFolder, "git", "checkout", a.Branch)
	if err != nil {
		return err
	}
	log.Printf("Repository initialized")

	return nil
}

func (a *App) fetchChanges() error {
	_, err := os.Stat(a.RepoFolder + "stop.sh")
	if os.IsNotExist(err) {
		log.Printf("No stop.sh file found in repository folder %s", a.RepoFolder)
	} else {
		log.Printf("Make stop.sh executable")
		err = executeCommand(a.RepoFolder, "chmod", "+x", "stop.sh")
		if err != nil {
			log.Printf("an error occurred while trying to make stop.sh executable; %s", err)
		} else {
			log.Printf("Running stop.sh in repo %s", a.RepoFolder)
			err = executeCommand(a.RepoFolder, "sudo", "-E", "./stop.sh")
			if err != nil {
				log.Printf("an error occurred while running stop.sh; %s", err)
			}
		}
	}

	log.Print("Reset changes")
	err = executeCommand(a.RepoFolder, "git", "reset", "--hard")
	if err != nil {
		return err
	}

	log.Print("Make sure that git always rebases")
	err = executeCommand(a.RepoFolder, "git", "config", "--global", "pull.rebase", "true")
	if err != nil {
		return err
	}

	log.Print("Fetch new changes")
	err = executeCommand(a.RepoFolder, "git", "fetch", "origin")
	if err != nil {
		return err
	}

	err = executeCommand(a.RepoFolder, "git", "checkout", a.Branch)
	if err != nil {
		return err
	}

	err = executeCommand(a.RepoFolder, "git", "pull")
	if err != nil {
		return err
	}
	log.Print("Finished fetching")

	_, err = os.Stat(a.RepoFolder + "start.sh")
	if os.IsNotExist(err) {
		log.Printf("No start.sh file found in repository folder %s", a.RepoFolder)
	} else {
		log.Printf("Make start.sh executable")
		err = executeCommand(a.RepoFolder, "chmod", "+x", "start.sh")
		if err != nil {
			log.Printf("an error occurred while trying to make start.sh executable; %s", err)
		} else {
			log.Printf("Running start.sh in repo %s", a.RepoFolder)
			err = executeCommand(a.RepoFolder, "sudo", "-E", "./start.sh")
			if err != nil {
				log.Printf("an error occurred while running start.sh; %s", err)
			}
		}
	}

	return nil
}

func executeCommand(dir string, commandName string, arg ...string) error {
	var stdoutBuf bytes.Buffer

	cmd := exec.Command(commandName, arg...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	log.Printf("stdout: %s\n", stdoutBuf.String())

	if err != nil {
		return err
	}
	return nil
}
