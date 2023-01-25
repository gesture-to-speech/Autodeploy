# Autodeploy

## Run on server

### Add SSH key for deployment
You need to create a default SSH key by running `ssh-keygen -t rsa -b 4096`, then click enter to create a key
with the default name `id_rsa`. Click enter again to create a key without a passphrase and click enter again to confirm.
A random SSH key will be generated. You need to add it in GitHub to the repository you are trying to autodeploy. To
do that, go to the repository's page, click `Settings -> Deploy keys (under Security on the left)`, and click button
`Add deploy key`. Set the title to `Autodeploy`, `Key` to the output of this command
`cat ~/.ssh/id_rsa.pub`, and click `Add key`.

### Autodeploy setup
First you need to add `github.com` to the list of known hosts. The easiest way to do it is by running 
`ssh -T git@github.com` and writing `yes` in the prompt. `github.com` should be added to the list of known hosts.

Next, create a directory where all the files will be stored and enter it. Next download the Autodeploy executable
by running `rm -f Autodeploy && wget --no-cache --no-cookies --no-check-certificate
https://github.com/gesture-to-speech/Autodeploy/raw/main/bin/Autodeploy -O Autodeploy && sudo chmod -R 0777 Autodeploy`.
Create `config.tml` based on the `config-template.tml` file. See the explanation of `config.tml` file below.

To start the server run `sudo rm -f nohup.out && sudo nohup ./Autodeploy &`. This will start the server listening to GitHub
webhook in the background. Logs will be located in the file `nohup.out`. To stop the server run
`sudo kill $(pgrep Autodeploy)`.

### GitHub Webhook setup
Next you need to set up webhook. Go to `Webhooks` (under `Code and automation` on the left in `Settings`) and click
button `Add webhook`. Set `Payload URL` to `<address of VM>:4550/hook` (webhook is listening on port 4550),
`Content type` to `application/json`, leave secret blank, ensure that
`Just the push event` and `Active` are set, and click button `Add webhook`.

## How it works
When the server is running, whenever anything is pushed to the repository specified in the `config.tml`, 
GitHub webhook will send a request to the `Autodeploy` server. If the secret sent through the request and branch,
to which the push event happened, is the same as specified in the `config.tml`, then the server will first run the
`stop.sh` script specified in the repository (if exists), then fetch and pull the changes, and run the `start.sh` script
specified in the repository (if exists).

`start.sh` and `stop.sh` are defined for a given repository. `start.sh` is responsible for setting up the environment
and starting the required process (or processes). `stop.sh` is responsible for cleanly stopping the running process.

## Config file
`config.tml` defines the behaviour of the `Autodeploy` script. It consists of two sections: `app` and `hook`.

In the `app` we define all variables related to the server:
- repo -> SSH to the repository that should be autodeployed
- branch -> name of the branch that should be autodeployed
- mainFolder -> absolute path to which the repo should be cloned, it should end with `/`

In the `hook` we define all variables related to the GitHub hook:
- key -> secret value that was set during the hook setup - leave it blank

## Development
To install all dependencies for development run `sudo bash install.sh` and add `bin` folder to the path
`go env -w GOBIN=<absolute path to the folder>/Autodeploy/bin`.

To compile code to an executable run `go install Autodeploy`.
