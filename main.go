package main

import (
	"log"
	"net/http"

	"github.com/pelletier/go-toml"
)

func main() {
	cfg, err := toml.LoadFile("config.tml")
	catch(err)

	repo, ok := cfg.Get("app.repo").(string)
	if !ok {
		log.Fatal("app.repo not defined")
	}

	branch, ok := cfg.Get("app.branch").(string)
	if !ok {
		log.Fatal("app.branch not defined")
	}

	mainFolder, ok := cfg.Get("app.mainFolder").(string)
	if !ok {
		log.Fatal("app.mainFolder not defined")
	}

	key, ok := cfg.Get("hook.key").(string)
	if !ok {
		log.Fatal("hook.key not defined")
	}

	app := NewApp(repo, branch, mainFolder)

	err = app.initRepo()
	catch(err)

	http.Handle("/hook", NewHookHandler(&HookOptions{
		App:    app,
		Secret: key,
	}))

	err = http.ListenAndServe("0.0.0.0:4550", nil)
	catch(err)
}

func catch(err error) {
	if err != nil {
		panic(err)
	}
}
