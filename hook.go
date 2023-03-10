package main

import (
	"encoding/json"
	"github.com/google/go-github/github"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type HookOptions struct {
	App    *App
	Secret string
}

func NewHookHandler(o *HookOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		evName := r.Header.Get("X-Github-Event")
		if evName != "push" {
			log.Printf("Ignoring '%s' event", evName)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		/*
			if o.Secret != "" {
				ok := false
				for _, sig := range strings.Fields(r.Header.Get("X-Hub-Signature")) {
					if !strings.HasPrefix(sig, "sha1=") {
						continue
					}
					sig = strings.TrimPrefix(sig, "sha1=")
					mac := hmac.New(sha1.New, []byte(o.Secret))
					mac.Write(body)
					if sig == hex.EncodeToString(mac.Sum(nil)) {
						ok = true
						break
					}
				}
				if !ok {
					log.Printf("Ignoring '%s' event with incorrect signature", evName)
					return
				}
			}
		*/
		ev := github.PushEvent{}
		err = json.Unmarshal(body, &ev)
		if err != nil {
			log.Printf("Ignoring '%s' event with invalid payload", evName)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if ev.Repo.FullName == nil || *ev.Repo.SSHURL != o.App.Repo {
			log.Printf("Ignoring '%s' event with incorrect repository name '%s'", evName, *ev.Repo.SSHURL)
			return
		}

		ref := ev.GetRef()
		branchName := strings.TrimPrefix(ref, "refs/heads/")

		if branchName != o.App.Branch {
			log.Printf("Ignoring '%s' event with incorrect branch name '%s'", evName, branchName)
			return
		}
		log.Printf("Handling '%s' event for %s", evName, o.App.Repo)

		err = o.App.fetchChanges()
	})
}
