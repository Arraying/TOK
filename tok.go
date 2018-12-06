package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type (
	// Rule represents a pair between a regular expression and a response.
	Rule struct {
		Regex      string `json:"regex"`
		Response   string `json:"response"`
		expression *regexp.Regexp
	}
)

var (
	// Whether or not the process should verbosely log everything that happens.
	verbose bool
	// The Discord bot token used in the IDENTIFY packet.
	token string
	// The config path, relative or absolute. This should contain the file name.
	config string
	// A collection of all rules.
	rules []*Rule
)

// init is the initialization function. This will load all the flags required to run the bot.
func init() {
	flag.BoolVar(&verbose, "verbose", false, "Verbosely prints out data, useful for debugging")
	flag.StringVar(&token, "token", "", "The Discord bot token")
	flag.StringVar(&config, "config", "config.json", "The config file path")
	flag.Parse()
}

// main is the entry point of the program.
func main() {
	log.Println("Loading config...")
	b, err := ioutil.ReadFile(config)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &rules)
	if err != nil {
		panic(err)
	}
	log.Printf("Starting bot with token \"%s\".\n", token)
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
	dg.AddHandler(message)
	err = dg.Open()
	if err != nil {
		panic(err)
	}
	log.Println("Opened session. Use CTRL^C to close.")
	k := make(chan os.Signal, 1)
	signal.Notify(k, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-k
	dg.Close()
}

// Expression returns the text-based regular expression as a Go regular expression pointer. May be nil in the case of errors.
func (r *Rule) Expression() *regexp.Regexp {
	if r.expression == nil {
		log.Println("Expression nil, compiling.")
		r.expression = regexp.MustCompile(r.Regex)
		log.Printf("Done compiling expression.")
	}
	return r.expression
}

// message handles incoming MESSAGE_CREATE events from the web socket.
func message(s *discordgo.Session, m *discordgo.MessageCreate) {
	if s.State.User.ID == m.Author.ID {
		return
	}
	c := m.Message.Content
	if verbose {
		log.Printf("Received message %s.\n", c)
	}
	for i := range rules {
		r := rules[i]
		if verbose {
			log.Printf("Checking match for %s: ", r.Regex)
		}
		if r.Expression() == nil {
			log.Println("Expression nil.")
			continue
		}
		x := r.Expression().MatchString(c)
		if verbose {
			log.Println(x)
		}
		if x {
			if verbose {
				log.Printf("Dispatching message %s.\n", r.Response)
			}
			_, err := s.ChannelMessageSend(m.ChannelID, r.Response)
			if err != nil {
				log.Println(err)
			}
			log.Println("Completed with event, returning.")
			return
		}
	}
}
