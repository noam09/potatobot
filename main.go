package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/noam09/go-couchpotato-api"
)

func main() {
	// Set up signal catching
	sigs := make(chan os.Signal, 1)
	// Catch all signals
	signal.Notify(sigs)
	// signal.Notify(sigs,syscall.SIGQUIT)
	// Method invoked on signal receive
	go func() {
		s := <-sigs
		log.Printf("RECEIVED SIGNAL: %s", s)
		AppCleanup()
		os.Exit(1)
	}()

	usage := `PotatoBot

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
  -d, --debug               Debug logging [default: false].`

	opts, _ := docopt.ParseDoc(usage)
	log.Println(opts)

	// Get whitelist
	var whitelist []int
	w := opts["--whitelist"].([]string)
	for _, v := range w {
		i, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		whitelist = append(whitelist, i)
	}

	// Get arguments
	token, _ := opts.String("--token")
	hostname, _ := opts.String("--host")
	port, _ := opts.Int("--port")
	apiKey, _ := opts.String("--key")
	urlBase, _ := opts.String("--base")
	ssl, _ := opts.Bool("--ssl")
	debug, _ := opts.Bool("--debug")
	username, _ := opts.String("--user")
	password, _ := opts.String("--pass")

	// Initialize CP client
	cp := couchpotato.NewClient(hostname, port, apiKey, urlBase, ssl, username, password)

	BOT_TOKEN := token

	bot, err := tgbotapi.NewBotAPI(BOT_TOKEN)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = debug

	log.Printf("Authorized on account %s", bot.Self.UserName)
	lowerBotName := strings.ToLower(bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	imdb, _ := regexp.Compile("tt[0-9]{5,10}")
	emojiStar := "\u2b50\ufe0f"
	emojiFilm := "\U0001f4fa"
	emojiSearch := "\U0001f50d"
	emojiCancel := "\u274e"

	for update := range updates {
		if update.Message == nil {
			continue
		}
		// Check if user ID in whitelist
		if !intInSlice(update.Message.From.ID, whitelist) {
			// log.Println("not me")
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		lowerMessageText := strings.ToLower(update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		if strings.HasPrefix(lowerMessageText, "/q ") || strings.HasPrefix(lowerMessageText, "/q@"+lowerBotName+" ") {
			// TEST add group support (/command@bot)
			split := strings.Split(lowerMessageText, "/q@"+lowerBotName+" ")
			if len(split) < 2 {
				split = strings.Split(lowerMessageText, "/q ")
			}
			if len(split) == 2 {
				_, q := split[0], split[1]
				log.Println(q)
				val, err := cp.SearchMovie(strings.TrimSpace(q))
				resultCount := 0
				if val != nil {
					msgBody := ""
					var rating float64
					var votes int
					// log.Println(len(val.Movies))
					// Make 2D slice with enough rows for the number of search results
					// Add one more row for the "/cancel" button
					rows := make([][]tgbotapi.KeyboardButton, len(val.Movies)+1)
					// Loop through results
					for _, a := range val.Movies {
						//log.Println("STATUS")
						//log.Println(a.Library.Status)
						library := ""
						/*
						if a.Library.Status != "" {
							if a.Library.Status == "done" {
								log.Println("IN LIBRARY")
								library = "(Already in library)"
							} else {
								library = "(Already in snatchlist)"
							}
						} else {
							library = ""
						}
						*/
						// Create button text for each search result
						// button := fmt.Sprintf("%s (%d) [%s]\n", a.Titles[0], a.Year, a.Imdb)
						button := fmt.Sprintf("%s %s [%d] [%s]", emojiFilm, a.Title, a.Year, a.Imdb)
						// In each row, append one column containing a KeyboardButton with the button text
						rows[resultCount] = append(rows[resultCount], tgbotapi.NewKeyboardButton(button))
						log.Printf("%s [%d] [%s]\n", a.Title, a.Year, a.Imdb)
						// log.Printf("Original poster: %s", a.Images.Original[0])
						// log.Printf("Small poster: %s", a.Images.More[0])
						// Tally the number of results
						resultCount += 1
						// Check if an IMDB rating and the number of votes were fetched
						if len(a.Rating.Score) == 2 {
							// log.Printf("Rating: %v", a.Rating.Score[0])
							rating = a.Rating.Score[0]
							// log.Printf("Votes: %v", a.Rating.Score[1])
							votes = int(a.Rating.Score[1])
							// Add the result to the message containing the list of results
							msgBody += fmt.Sprintf("*%d)* [%s (%d)](https://imdb.com/title/%s) %s _%v (%v) - %dm_ %s\n",
								resultCount, a.Title, a.Year, a.Imdb, emojiStar, rating, votes, a.Runtime, library)
						} else {
							rating = 0
							votes = 0
							// Add the result to the message containing the list of results (without rating)
							msgBody += fmt.Sprintf("*%d)* [%s (%d)](https://imdb.com/title/%s) - _%dm_\n",
								resultCount, a.Title, a.Year, a.Imdb, a.Runtime)
						}
					}
					// If there is at least one result ready, create the custom keyboard
					if resultCount > 0 {
						button := fmt.Sprintf("/cancel")
						rows[resultCount] = append(rows[resultCount], tgbotapi.NewKeyboardButton(button))
						// Init keyboard variable
						var kb tgbotapi.ReplyKeyboardMarkup
						kb = tgbotapi.ReplyKeyboardMarkup{
							ResizeKeyboard:  true,
							Keyboard:        rows,
							OneTimeKeyboard: true,
						}
						// kb.OneTimeKeyboard = true
						// Append the custom keyboard to the reply message
						msg.ReplyMarkup = kb
						// Append a hint to the results
						msgBody += "\n\nIf you were expecting more results, try running the same search again"
						// Append the list of results to the reply message
						msg.Text = msgBody
						msg.ParseMode = "markdown"
						// msg.ReplyToMessageID = update.Message.MessageID
						// Avoid the first IMDB link being resolved for preview
						msg.DisableWebPagePreview = true
					} else {
						msgBody += "Looks like no results were found. "
						msgBody += "You might try running the same search again, "
						msgBody += "sometimes they get lost on the way back!"
						msg.Text = msgBody
					}
				} else {
					log.Println(err)
				}
			} else {
				log.Println("/q or /q@ prefix but no message?")
				continue
			}
			// Detect if IMDB IDs are sent to the bot (e.g. tt12345678)
			// Check if the IMDB ID regex matches the received message
			// IMDB IDs sent to groups are automatically set as a reply (custom keyboard)
			// so the bot will have access to the message text
		} else if imdb.MatchString(lowerMessageText) {
			log.Println("Found IMDB ID")
			// Find the IMDB ID matched in the received message
			imdbId := imdb.FindString(lowerMessageText)
			// Find the IMDB ID matched in the received message
			title, _ := regexp.Compile(emojiFilm+" (.*) \\[.*\\] \\[tt[0-9]{5,10}\\]")
			movieTitle := title.FindStringSubmatch(update.Message.Text)
			if len(movieTitle) > 1 {
				t := movieTitle[1]
				log.Println(t)
				val, err := cp.AddMovie(imdbId, t)
				if err != nil {
					log.Println(val)
				}
				msg.Text = fmt.Sprintf("Adding movie _%s_ (%s) to snatchlist", t, imdbId)
				msg.ParseMode = "markdown"
			} else {
				msg.Text = "Failed adding movie to snatchlist"
			}
			// val, err := cp.AddMovie()
			/*
			   match, _ := regexp.MatchString("[a-z0-9]{32}", strings.ToLower(update.Message.Text))
			   if !match {
			       log.Panic("API key contains invalid characters")
			   }
			*/
			// Remove the custom keyboard if it's still there
			msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{
				RemoveKeyboard: true,
				Selective:      false,
			}
		} else {
			log.Println(strings.TrimSpace(lowerMessageText))
			switch strings.TrimSpace(lowerMessageText) {
			case "/start",  "/start@"+lowerBotName:
				msg.Text = "Hi! Use /q to start searching"
			case "/help", "/h", "/help@"+lowerBotName, "/h@"+lowerBotName:
				msg.Text = fmt.Sprintf("%s /q - Movie search", emojiFilm)
				msg.Text += fmt.Sprintf("\n\n%s /f - Run full search for all wanted movies", emojiSearch)
				msg.Text += fmt.Sprintf("\n\n%s /c - Cancel current operation", emojiCancel)
				msg.ParseMode = "markdown"
			case "/q", "/q@"+lowerBotName:
				msg.Text = fmt.Sprintf("/q should be followed by a search query. Example:\n`/q Inception`\n")
				msg.ParseMode = "markdown"
			case "/cancel", "/c", "/cancel@"+lowerBotName, "/c@"+lowerBotName:
				// Cancel
				msg.Text = "Cancelling"
				msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{
					RemoveKeyboard: true,
					Selective:      false,
				}
			case "/full", "/f", "/full@"+lowerBotName, "/f@"+lowerBotName:
				_, err := cp.FullSearch()
				if err != nil {
					log.Println(err)
				}
				msg.Text = "Running full search for all wanted movies"
			default:
				msg.Text = "`¯\\_(ツ)_/¯`"
				msg.ParseMode = "markdown"
			}
		}

		if msg.Text != "" {
			bot.Send(msg)
		}
		continue
	} // end updates
} // end main

func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func AppCleanup() {
	log.Println("CLEANUP APP BEFORE EXIT!!!")
}
