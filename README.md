# go-matrix
A Go library for Matrix.org Clients

## Works
- Login to HomeServer
- Joining Rooms
- Sending Events
- Receiving Events

## Example
```
package main

import (
	"os"
	"github.com/ericevenchick/go-matrix"
)

func main() {
	c := matrix.NewClient("https://matrix.org")
	c.PasswordLogin(os.Args[1], os.Args[2])
	c.JoinRoom(os.Args[3])
	msg := matrix.MessageEvent{MessageType: "m.text", Body: "Hello World!"}
	err := c.SendEvent("!kRvllZVWbwFkpOwidF:matrix.org", "m.room.message", msg)
	if err != nil {
		panic(err)
	}
}
```
