A simplistic build application for any virtual machine

## Background

Some developers don't want to, or can, use cloud providers to host and deploy our applications. Instead we use on-premise
virtual machines where we might have root access, but these servers don't usually have any fancy deployment tools and we don't want to allow
other applications to directly connect to it to deploy something. Sometimes you also might want to avoid writing an automated pipeline for your builds.

This app serves a very simple REST API, with a single action endpoint, that can be triggered via a HTTP POST request to execute a script on the server. 
This script can contain instructions for rebuilding your application each time you commit to your repository, or make changes to your website content,
for example. This is particularly helpful if you have something like a front-end application with Jamstack architecture, or just want to enable
continuous deployment for a backend application.

This application has only been tested on a Linux server. Error scenarios have not been tested.

## Getting started

1. Clone this repo
2. Run `go install ./src`
3. Copy `.env.dist` as `.env` and set your variable values.
4. Create a bash script to be executed somewhere on your server. Copy it's full path and allow it to be executed with `chmod +x build.sh` for example.
5. Build the app with `GOOS=<operating-system> GOARCH=<architecture> go build -o bin/main cmd/main.go` and run it with something like pm2: `pm2 start bin/main --name builder-app`.
6. Set up a reverse proxy in Apache / Nginx for localhost:8081 (default port in env)
7. Do a POST call to your server. Don't forget to include the additional header, if you set one in the `BUILDER_WEBHOOK_SECRET`. For example `POST http://localhost:8080/build` (with headers).
8. (optional) Set your Github repository webhook to point to the same server, to `/build` path, and set the same secret as the content (or chnge the server implementation if you want different secrets) a

## Example .env file

```bash
# Secret / password that needs to exist in the header or body for request validation
BUILDER_WEBHOOK_SECRET=barFoo123

# Name for the header where secret is set in.
# Only fill if you decide to put the key in the request headers.
BUILDER_WEBHOOK_SECRET_HEADER=x-build-secret

# Secret for Github webhook. This needs to be added to Github.
# It is used to validate the request signature.
BUILDER_GITHUB_SECRET=fooBar123

# Port where HTTP server is exposed in
BUILDER_PORT=8081

# Path to executable shell script
BUILDER_EXEC_PATH="/var/www/build.sh"
```

## Execution on the server

You can write a shell script, or bash script, anywhere on the server and point the `BUILDER_EXEC_PATH` to that file.
Currently providing any arguments from this application is not possible. 

In case you run into permission issues with a shell script, remember to run `chmod +x <file>`

An example bash script for a JavaScript front-end project with Node.js:

```sh
#!/bin/bash

cd /var/www/html;
git pull;
rm -rm node_modules;
npm install;
npm run build && npm run generate && pm2 restart <ProcessName>;
```

A more advanced case would be to create an empty directory to which you clone the entire repo, install packages, copy a base environment variable file to, then build it, move that build over to the actual host folder and restart your process manager. This way you won't hit any conflicts with `git pull`.

You can also take this repo and extend it by providing a request body to the `/build` endpoint where you tell your server where it can download the build artifact and simply unpack it on the server, to avoid running Git, which might end up being problematic.
