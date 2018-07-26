package couchpotato

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

type client struct {
	hostname   string
	port       int
	apiKey     string
	urlBase    string
	ssl        bool
	username   string
	password   string
	auth       bool
	httpClient *http.Client
	// httpTimeout  time.Duration
	// debug        bool
}

type apikey struct {
	ApiKey  string `json:"api_key"`
	Success bool   `json:"success"`
}

type MovieResults struct {
	Movies []Movie `json:"movies,omitempty"`
}

type Movie struct {
	Tmdb int    `json:"tmdb_id,omitempty"`
	Imdb string `json:"imdb,omitempty"`
	Year int    `json:"year,omitempty"`
	// Titles      []string     `json:"titles,omitempty"`
	Title   string     `json:"original_title,omitempty"`
	Images  Poster     `json:"images,omitempty"`
	Runtime int        `json:"runtime,omitempty"`
	Rating  ImdbRating `json:"rating,omitempty"`
	// Library InLibrary `json:"in_library,string"`
}

type InLibrary struct {
	IsSet bool `json:"-"`
	InLibraryStatus
}

type InLibraryStatus struct {
	Status string `json:"status,omitempty"`
}

type ImdbRating struct {
	Score []float64 `json:"imdb,omitempty"`
	// Score  float64  `json:"score,omitempty"`
	// Votes  int      `json:"votes,omitempty"`
}

type Poster struct {
	Original []string `json:"poster_original,omitempty"`
	More     []string `json:"poster,omitempty"`
}

// NewClient returns a new CP HTTP client
func NewClient(hostname string, port int, apiKey, urlBase string, ssl bool, username, password string) (c *client) {
	auth := true
	return &client{hostname, port, apiKey, urlBase, ssl, username, password, auth, &http.Client{Timeout: 10 * time.Second}}
}

/*
func (w *InLibrary) UnmarshalJSON(data []byte) error {
  if id, err := strconv.Atoi(string(data)); err == nil {
    w.IsSet = id
    w.IsSet = true
    return nil
  }
  return json.Unmarshal(data, &w.InLibrary)
}
*/

func (c *client) GetApiKey(serverUrl string) (string, error) {
	var errorMsg string
	// Username MD5 hash
	h := md5.New()
	io.WriteString(h, c.username)
	u := fmt.Sprintf("%x", h.Sum(nil))
	// Password MD5 hash
	h = md5.New()
	io.WriteString(h, c.password)
	p := fmt.Sprintf("%x", h.Sum(nil))
	// Make URL
	var Url *url.URL
	Url, err := url.Parse(serverUrl)
	if err != nil {
		log.Println("Bad serverUrl")
		return "Bad serverUrl", err
	}
	// Set URL path
	Url.Path += "/getkey/"
	// Set URL parameters (username and password hashes)
	parameters := url.Values{}
	parameters.Add("p", p)
	parameters.Add("u", u)
	// Prepare URL
	Url.RawQuery = parameters.Encode()
	// GET URL
	log.Printf("Will try to GET: %q\n", Url.String())
	resp, err := c.httpClient.Get(Url.String())
	if err != nil {
		log.Println(err)
		return "Can't GET " + Url.String(), err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	apikey := apikey{}
	jerr := json.Unmarshal(body, &apikey)
	if err != nil {
		log.Println(jerr)
		return "Failed unmarshalling JSON", jerr

	}
	if apikey.Success {
		return apikey.ApiKey, err
	} else {
		errorMsg = "Failed getting API key"
		log.Println(errorMsg)
		log.Println(string(body))
		return "Failed getting API key", errors.New(errorMsg)
	}
}

func (c *client) GetAddr() string {
	// var errorMsg string
	var err error
	log.Println("Start building full address")
	serverUrl := "http"
	if c.ssl == true {
		serverUrl += "s"
	}
	serverUrl += "://" + c.hostname + ":" + strconv.Itoa(c.port)

	if c.urlBase != "" {
		serverUrl += c.urlBase
	}

	if c.apiKey == "" {
		log.Println("No API key specified")
		// Get API key
		if c.username == "" || c.password == "" {
			log.Panic("Username and/or password not specified, API key cannot be computed")
		} else {
			log.Println("Compute API key with c.GetApiKey()")
			c.apiKey, err = c.GetApiKey(serverUrl)
			if err != nil {
				log.Println(err)
			}
		}
	}

	if len(c.apiKey) != 32 {
		log.Panic("API key length not 32")
	}

	match, _ := regexp.MatchString("[a-z0-9]{32}", c.apiKey)
	if !match {
		log.Panic("API key contains invalid characters")
	}

	// Complete API URL
	serverApi := serverUrl + "/api/" + c.apiKey + "/"
	log.Println("Return full API addr")
	return serverApi
}

func (c *client) AddMovie(id, title string) (string, error) {
	var errorMsg, bodyString string
	//var bodyBytes []byte
	var Url *url.URL
	Url, err := url.Parse(c.GetAddr())
	if err != nil {
		panic("Arghh bad serverUrl")
	}
	// Set API method
	Url.Path += "movie.add"
	// Set method parameters
	// IMDB identifier, original title, and whether to re-add existing title
	parameters := url.Values{}
	parameters.Add("identifier", id)
	parameters.Add("title", title)
	parameters.Add("force_readd", "false")
	// Prepare URL
	Url.RawQuery = parameters.Encode()
	// GET URL
	log.Printf("Will try to GET: %q\n", Url.String())
	resp, err := c.httpClient.Get(Url.String())
	if err != nil {
		log.Println(err)
		return "Can't GET " + Url.String(), err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			// Return JSON result
			bodyString = string(bodyBytes)
		}
		return bodyString, err2
	}
	return bodyString, errors.New(errorMsg)
}

// movie.search
func (c *client) SearchMovie(q string) (*MovieResults, error) {
	var errorMsg string
	var Url *url.URL
	Url, err := url.Parse(c.GetAddr())
	if err != nil {
		panic("Bad serverUrl")
	}
	// Set API method
	Url.Path += "movie.search"
	// Set method parameters
	// q - Search query
	parameters := url.Values{}
	parameters.Add("q", q)
	// Prepare URL
	Url.RawQuery = parameters.Encode()
	// Tally tries
	tries := 0
	// GET URL
	for tries < 3 {
		tries += 1
		log.Printf("Will try to GET: %q\n", Url.String())
		resp, err := c.httpClient.Get(Url.String())
		if err != nil {
			log.Println("Can't GET " + Url.String())
			return &MovieResults{}, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			// Decode JSON result
			decoder := json.NewDecoder(resp.Body)
			val := &MovieResults{}
			err := decoder.Decode(val)
			if err != nil {
				log.Fatal(err)
			}
			// Tally results
			resultCount := 0
			for _, _ = range val.Movies {
				// log.Printf("%s (%d) [%s]\n", a.Titles[0], a.Year, a.Imdb)
				// log.Printf("Original poster: %s", a.Images.Original[0])
				// log.Printf("Small poster: %s", a.Images.More[0])
				resultCount += 1
			}
			if resultCount == 0 {
				errorMsg := "No results found, either try again or change search terms"
				log.Println(errorMsg)
				// tries += 1
			} else {
				return val, nil
			}
		} else {
			log.Println(resp.StatusCode)
			// tries += 1
		}
	}
	return &MovieResults{}, errors.New(errorMsg)
}

// movie.searcher.full_search
func (c *client) FullSearch() (string, error) {
	var errorMsg, bodyString string
	var Url *url.URL
	Url, err := url.Parse(c.GetAddr())
	if err != nil {
		panic("Bad serverUrl")
	}
	// Set API method
	Url.Path += "movie.searcher.full_search"
	// GET URL
	log.Printf("Will try to GET: %q\n", Url.String())
	resp, err := c.httpClient.Get(Url.String())
	if err != nil {
		log.Println(err)
		return "Can't GET " + Url.String(), err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			bodyString = string(bodyBytes)
		}
		return bodyString, err2
	}
	return bodyString, errors.New(errorMsg)
}
