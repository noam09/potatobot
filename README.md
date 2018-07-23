# PotatoBot

🥔🤖 A simple Telegram bot for controlling CouchPotato.

## Dependencies

* https://github.com/noam09/go-couchpotato-api
* https://github.com/go-telegram-bot-api/telegram-bot-api
* https://github.com/docopt/docopt-go

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

**Protip:** Sending PotatoBot an IMDB title ID by itself (e.g. `tt123456`) will add the title to the snatchlist.

## Screenshots

Start the bot:

![Start Bot](https://i.imgur.com/4ni9fDm.png)

Send `/help` to show the list of available commands:

![List Commands](https://i.imgur.com/okomCfX.png)

Send `/q QUERY` and use the custom keyboard to select a movie result:

![Select Movie](https://i.imgur.com/zGkv7Pm.png)

## TODO

* Add group command support (`/command@bot`)
* Check if exists in library
  * Prompt to re-add
* On-the-fly user whitelisting?
* Choose quality profile other than the default

## License

This is free software under the GPL v3 open source license. Feel free to do with it what you wish, but any modification must be open sourced. A copy of the license is included.
