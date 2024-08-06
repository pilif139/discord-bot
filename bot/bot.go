package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var BotToken string

type task struct {
	name     string
	finished bool
	author   string
}

func Run() {

	// create a session
	discord, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatal(err)
	}

	// add a event handler
	tasks := []task{}
	discord.AddHandler(func(d *discordgo.Session, m *discordgo.MessageCreate) {
		newMessage(d, m, &tasks)
	})

	discord.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// open session
	err = discord.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer discord.Close() // close session, after function termination

	// keep bot running untill there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func newMessage(d *discordgo.Session, m *discordgo.MessageCreate, tasks *[]task) {
	if m.Author.ID == d.State.User.ID {
		return
	}

	words := strings.SplitN(m.Content, " ", 2)

	switch words[0] {
	case "!help":
		d.ChannelMessageSend(m.ChannelID, "Lista komend:\n!dodaj-zadanie <nazwa zadania> - dodaje nowe zadanie")
	case "!dodaj-zadanie":
		if len(words) > 1 {
			*tasks = append(*tasks, task{name: words[1], finished: false, author: m.Author.Username})
			d.ChannelMessageSend(m.ChannelID, "Dodano zadanie: "+words[1])
		} else {
			d.ChannelMessageSend(m.ChannelID, "Podaj nazwe zadania")
		}
	case "!zadania":
		d.ChannelMessageSend(m.ChannelID, printTasks(*tasks, m.Author.Username))
	}
}

func printTasks(tasks []task, author string) string {
	var result string
	for i, t := range tasks {
		if t.author == author {
			result += fmt.Sprintf("%d. %s\n", i+1, t.name)
		}
	}
	return result
}
