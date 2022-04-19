# Go Chatroom
## Description
Simple TCP chatroom written in Go.

## Features
* Low overhead and memory footprint
* No dependencies
* Multithreaded and thread-safe
* Supports multiple users with usernames
* Commands:
    * /help - displays this help message
    * /list - displays a list of all users
    * /log - displays the chat log

## Usage
* Change the `IP` (default: localhost) and `PORT` constants to match your server
* Run the server:
   ```bash
   $ go run main.go
   ```
* Build the server:
   ```bash
   $ go build main.go
   ```
* Connect to the server:
    ```bash
    $ nc localhost 8080
    ```