COSC1170
========

COSC1170 Foundation Distributed Computing

## Assignment 1 notes
- **start_client.go**: initiates Client functionality
- **start_server.go**: initiates Server listener
- `go test ./...` runs the tests for the application, output can be found in tramservice.log
- check out commit history for a breakdown of changes that occured during development
- unable to ensure that a character is represented as 2 bytes in Go as a string is a slice of bytes rather than characters which are not guaranteed to be a fixed size. Generally it's utf8 so can vary depending on the char it is representing. A "Rune" type representing a char in Golang is an alias for uint32 (4 bytes).
- Based on instructions from lab assistant, TransactionID is just copied the client for validation of server response.


 
