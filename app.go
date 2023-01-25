package main

import (
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
	_, err := os.Stat(a.RepoFolder)
	if !os.IsNotExist(err) {
		log.Print("Repository already initialized")
		return nil
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
	if os.IsExist(err) {
		log.Print("Running stop.sh")
		err = executeCommand(a.RepoFolder, "sudo", "-E", "sh", "stop.sh")
		if err != nil {
			return err
		}
	} else {
		log.Print("No stop.sh file found in repository folder")
	}

	log.Print("Fetching changes")

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
	if os.IsExist(err) {
		log.Print("Running start.sh")
		err = executeCommand(a.RepoFolder, "sudo", "-E", "sh", "start.sh")
		if err != nil {
			return err
		}
	} else {
		log.Print("No start.sh file found in repository folder")
	}

	return nil
}

func executeCommand(dir string, commandName string, arg ...string) error {
	cmd := exec.Command(commandName, arg...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
