# Builder App with a REST API for webhooks

Some developers don't want to, or can, use cloud providers to host and deploy our applications. Instead we have to use on-premise
virtual machines where we might have root access, but these servers don't usually have any fancy deployment tools. 

This app serves a webhook, with a user specified path that can be called via HTTP POST to execute a shell script on the server. 
This script can contain instructions for rebuilding your application each time you commit to your repository, or make changes to your
website content.

This application has only been tested on a Linux server. Error scenarios have not been tested.

## Getting started

1. Clone this repo
2. Run `go install ./src`
3. Copy `.env.dist` as `.env` and set your variable values.
4. Build the app, copy it to your web server and host it (WIP)
5. Do a POST call to your server's address with the specified `BUILDER_PORT` and `BUILDER_CONTENT_WEBHOOK_PATH`. Don't forget to include the additional header, if you set one in the `BUILDER_REQ_HEADER`. For example `POST http://localhost:8082/gxrjg4y6s6kjshznb1a5` (with headers).
6. (optional) Set your Github repository webhook to point to the same server,to `BUILDER_GITHUB_WEBHOOK_PATH`, and set the same secret as the content (or change the server implementation if you want different secrets) 

## Example .env file

```bash
# Secret / password that needs to exist in the header or body for request validation
BUILDER_WEBHOOK_SECRET=barFoo123

# Name for the header where secret is set in.
# Only fill if you decide to put the key in the request headers.
BUILDER_SECRET_HEADER_KEY=x-webhook-secret

# Name for the property in the request body where secret is set in.
# Only fill if you decide to put the key in the request body.
BUILDER_SECRET_BODY_KEY=

# Secret for Github webhook. This needs to be added to Github.
# It is used to validate the request signature.
BUILDER_GITHUB_SECRET=fooBar123

# Port where HTTP server is exposed in
BUILDER_PORT=8008
```

## Execution on the server

You can write a shell script, or bash script, anywhere on the server and point the `BUILDER_EXEC_PATH` to that file.
Currently providing any arguments from this application is not possible. 

In case you run into permission issues with a shell script, remember to run `chmod +x <file>`

An example bash script:

```sh
#!/bin/bash

cd /var/www/html;
git pull;
rm -rm node_modules;
npm install;
npm run build && npm run generate && pm2 restart <ProcessName>;
```