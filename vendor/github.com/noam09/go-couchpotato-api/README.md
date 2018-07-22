# go-couchpotato-api

🥔 A simple wrapper for the CouchPotato API written in Go, created mainly for its role in [PotatoBot](https://github.com/noam09/potatobot), a Telegram bot for CouchPotato. 

## Usage

```go
// Initialize CP client
cp := couchpotato.NewClient(hostname, port, apiKey, urlBase, ssl, username, password)
// Get search results
val, err := cp.SearchMovie("big buck bunny")
```

## TODO

* Quality selection
* Check if title exists in library

## License

This is free software under the GPL v3 open source license. Feel free to do with it what you wish, but any modification must be open sourced. A copy of the license is included.
