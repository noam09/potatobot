# PotatoBot

🥔🤖 A simple [Telegram](https://telegram.org) bot for controlling [CouchPotato](https://github.com/CouchPotato/CouchPotatoServer).

## Dependencies

* [go-couchpotato-api](https://github.com/noam09/go-couchpotato-api)
* [telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api)
* [docopt-go](https://github.com/docopt/docopt-go)

## Build

Clone this repository and `go build`:

```console
git clone https://github.com/noam09/potatobot
cd potatobot
go build main.go
```

## Install

Use `go install` to get and build PotatoBot, making it available in `$GOPATH/bin/potatobot`:

```console
go get -u github.com/noam09/potatobot
go install github.com/noam09/potatobot
```

## docker-compose

A sample `docker-compose.yml` is included.

```console
# Pull the latest code
git clone https://github.com/noam09/potatobot
cd potatobot
# Modify YAML according to your setup
nano docker-compose.yml
# Run the container and send to background
docker-compose up -d
```

The `docker-compose.yml` is based on the official `golang:1.12.7-alpine` image:

```yaml
---
version: "2"

services:
  potatobot:
    image: golang:1.12.7-alpine
    volumes:
      - .:/go/src/potatobot
    working_dir: /go/src/potatobot
    command: >
      sh -c 'go run main.go
      --token=<bot>
      --key=<apikey>
      -w <chatid>
      --host=<host>
      --port=<port>
      --base=<urlbase>
      --ssl'
```

Modify the `command` section's parameters based on the help-text below.

## Usage

Running the bot:

```console
PotatoBot

Usage:
  potatobot --token=<bot> --user=<username> --pass=<password> -w <chatid>... [--host=<host>] [--port=<port>] [--base=<urlbase>] [--ssl] [-d]
  potatobot --token=<bot> --key=<apikey> -w <chatid>... [--host=<host>] [--port=<port>] [--base=<base>] [--ssl] [-d]
  potatobot -h | --help

Options:
  -h, --help                Show this screen.
  -t, --token=<bot>         Telegram bot token.
  -k, --key=<apikey>        API key.
  -u, --user=<username>     Username for web interface.
  -p, --pass=<password>     Password for web interface.
  -w, --whitelist=<chatid>  Telegram chat ID(s) allowed to communicate with the bot (contact @myidbot).
  -o, --host=<host>         Hostname or address CouchPotato runs on [default: 127.0.0.1].
  -r, --port=<port>         Port CouchPotato runs on [default: 5050].
  -b, --base=<urlbase>      Path which should follow the base URL.
  -s, --ssl                 Use TLS/SSL (HTTPS) [default: false].
  -d, --debug               Debug logging [default: false].
```

Controlling the bot:

```
📺 /q - Movie search

🔍 /f - Run full search for all wanted movies

❎ /c - Cancel current operation
```

**💡 Protip!** Sending PotatoBot an IMDB title ID by itself (e.g. `tt123456`) will add the title to the snatchlist.

## Screenshots

Start the bot:

![Start Bot](https://i.imgur.com/4ni9fDm.png)

Send `/help` to show the list of available commands:

![List Commands](https://i.imgur.com/okomCfX.png)

Send `/q QUERY` and use the custom keyboard to select a movie result:

![Select Movie](https://i.imgur.com/zGkv7Pm.png)

## TODO

* Makefile
* systemd service file
* Add group command support (`/command@bot`)
* Check if exists in library
  * Prompt to re-add
* On-the-fly user whitelisting?
* Choose quality profile other than the default

## Development

Contributions are always welcome, just create a [pull request](https://github.com/noam09/potatobot/pulls).

## License

This is free software under the GPL v3 open source license. Feel free to do with it what you wish, but any modification must be open sourced. A copy of the license is included.
