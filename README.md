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
    $ go run main.go
* Connect to the server:
    $ nc localhost 8080
* Using commands:
    $ /help
    $ /list
    $ /log