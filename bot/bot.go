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
	id       uint
	name     string
	finished bool
	author   string
}

func printTasks(tasks []task, author string) string {
	var result string
	i := 0
	for _, t := range tasks {
		if t.author == author {
			i++
			result += fmt.Sprintf("%d. %s\n", i, t.name)
		}
	}
	return result
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

	discord.Identify.Intents = discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsGuildVoiceStates

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
		d.ChannelMessageSend(m.ChannelID, "Lista komend:\n!dodaj-zadanie nazwa zadania - dodaje nowe zadanie\n!zadania - wyswietla liste zadan")
	case "!dodaj-zadanie":
		addTask(words, tasks, d, m)
	case "!zadania":
		d.ChannelMessageSend(m.ChannelID, "Zadania użytkownika "+m.Author.Username+":\n"+printTasks(*tasks, m.Author.Username))
	case "!dolacz":
		joinVoiceChannel(d, m)
	}
}

func joinVoiceChannel(d *discordgo.Session, m *discordgo.MessageCreate) {

	guildID := m.GuildID

	var voiceChannelID string
	guild, err := d.State.Guild(guildID)
	if err == nil {
		for _, vs := range guild.VoiceStates {
			if vs.UserID == m.Author.ID {
				voiceChannelID = vs.ChannelID
				break
			}
		}
	}

	if voiceChannelID != "" {
		_, err := d.ChannelVoiceJoin(guildID, voiceChannelID, false, false)
		if err != nil {
			d.ChannelMessageSend(m.ChannelID, "Nie udało się dołączyć do kanału głosowego")
		}
	} else {
		d.ChannelMessageSend(m.ChannelID, "Muszisz być na kanale głosowym!")
	}
}

func addTask(words []string, tasks *[]task, d *discordgo.Session, m *discordgo.MessageCreate) {
	if len(words) > 1 {
		*tasks = append(*tasks, task{
			name:     words[1],
			finished: false,
			author:   m.Author.Username,
			id:       uint(len(*tasks) + 1),
		})
		d.ChannelMessageSend(m.ChannelID, "Dodano zadanie: "+words[1])
	} else {
		d.ChannelMessageSend(m.ChannelID, "Podaj nazwe zadania")
	}
}
