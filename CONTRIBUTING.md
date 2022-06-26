# Development environment

---
Most instructions for running the system from source can be found in the [README](README.md). This section adds to that, with recommended IDE's, build configurations nd other instructions.

## Back-end

---
The back-end server can be developed in any IDE that supports Go. It was mostly development using [Goland](https://www.jetbrains.com/go/).

### API Documentation
The API documentation is generated from annotations int the Go code using [swag](https://github.com/swaggo/swag).
If annotations change in the code, the docs must be generated again with a command. First install `swag` once.

    cd ./mainServer
    go install github.com/swaggo/swag/cmd/swag@latest

After every update to the comment annotations in the code, update the documentation. (from the `/mainServer` folder).

    swag init -g server/router.go

Make sure to restart the Go server after this, before trying to access the updated API documentation. So, terminate the server if it's running (Ctr+C) and use `build` and `run` again as described in [README.md](README.md).

### Run configuration
The full steps can be incorporated in a single command if necessary. Below is an example of a Powershell command.

    cd ./mainServer; swag init -g server/router.go; if($?) {go build mainServer}; if($?) {go run mainServer}

## Front-end

---
The front-end server can be developed in any IDE that supports node.js. It was mostly developed using [WebStorm](https://www.jetbrains.com/webstorm/). As mentioned in the [README](README.md), make sure to use the right `run` command based on your platform.


