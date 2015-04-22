package main

// Main entry point for the app. Handles command-line options, starts the web
// listener and any import, etc

import (
	"log"
	"math/rand"
	"os"
	"time"
)

var (
	httpPort       string
	httpHost       string
	numWords       int
	prefixLen      int
	stateFile      string
	responseChance int
	botUsername    string

	markovChain *Chain
)

func init() {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator.
}

func main() {
	// // Parse command-line options
	// flag.Usage = func() {
	// 	fmt.Fprintf(os.Stderr, "usage: ./slack-markov\n")
	// 	flag.PrintDefaults()
	// }

	// flag.IntVar(&httpPort, "port", 8000, "The HTTP port on which to listen")
	// flag.IntVar(&numWords, "words", 100, "Maximum number of words in the output")
	// flag.IntVar(&prefixLen, "prefix", 2, "Prefix length in words")
	// flag.IntVar(&responseChance, "responseChance", 10, "Percent chance to generate a response on each request")
	// flag.StringVar(&stateFile, "stateFile", "state", "File to use for maintaining our markov chain state")
	// flag.StringVar(&botUsername, "botUsername", "markov-bot", "The name of the bot when it speaks")

	// flag.StringVar(&twitterConsumerKey, "twitterConsumerKey", "", "Twitter API key")
	// flag.StringVar(&twitterConsumerSecret, "twitterConsumerSecret", "", "Twitter API key secret")
	// flag.StringVar(&twitterAccessToken, "twitterAccessToken", "", "Twitter access token")
	// flag.StringVar(&twitterAccessTokenSecret, "twitterAccessTokenSecret", "", "Twitter access token secret")

	// var importDir = flag.String("importDir", "", "The directory of a Slack export")
	// var importChan = flag.String("importChan", "", "Optional channel to limit the import to")

	// flag.Parse()

	// if httpPort == 0 {
	// 	flag.Usage()
	// 	os.Exit(2)
	// }

	httpPort = os.Getenv("PORT")
	httpHost = os.Getenv("HOST")
	//stateFile = os.Getenv("BRAIN")
	stateFile = "state"
	//botUsername = os.Getenv("NAME")
	botUsername = "Brainard"

	markovChain = NewChain(prefixLen) // Initialize a new Chain.

	// Import into the chain
	err := markovChain.Load(stateFile)
	if err != nil {
		//log.Fatal(err)
		log.Printf("Could not load from '%s'. This may be expected.", stateFile)
	} else {
		log.Printf("Loaded previous state from '%s' (%d suffixes).", stateFile, len(markovChain.Chain))
	}

	// Start the webserver
	StartServer(httpHost+":"+httpPort)
}
