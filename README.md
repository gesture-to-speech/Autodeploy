# Autodeploy

## Run autodeploy
To download and run autodeploy server at your machine you need to add `config.tml`,
download the executable file `rm -f Autodeploy && wget https://github.com/GestureToSpeech/Autodeploy/raw/master/bin/Autodeploy -O Autodeploy && sudo chmod -R 0777 Autodeploy`
and run it `rm -f nohup.out && nohup ./Autodeploy &`. This will start the server in the background. Logs will be in 
`nohup.out`. To stop it, run `kill $(pgrep Autodeploy)`.

Go to settings in your repository, click on `Webhooks`, and click `Add webhook`. Set payload URL to `http://<address of VM>:4550/hook`,
content type to `application/json`, secret to whatever value you want (make sure to set the same value in `config.tml`),
set `Just the push event`, set `Active` to true, and click `Add webhook`.

Now the webhook is set up. Your server will receive info whenever anything is pushed.

You can also add `start.sh` and `stop.sh` bash scripts to you repository to set up and start the repository code, 
and safely stop it before updating repository.

## Development
To install all dependencies for development: `sudo bash install.sh`.

Add bin path: `go env -w GOBIN=/path/to/folder/Autodeploy/bin`

Compile to executable: `go install Autodeploy`
