package main

import (
	"Mashiron-V/lib"
	"Mashiron-V/lib/config"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/common-nighthawk/go-figure"
)

var discord *discordgo.Session

func init() {
	if config.CurrentConfig.Debug {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}
	figure.NewColorFigure("Mashiron-V", "rectangles", "green", true).Print()
	log.SetPrefix("[Init]")
	var err error
	discord, err = discordgo.New("Bot " + config.CurrentConfig.Discord.Token)
	if err != nil {
		log.Fatal("Discordgo late init failure:", err)
	}
	discord.AddHandler(lib.MessageCreate)
	discord.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates
}

func main() {
	err := discord.Open()
	if err != nil {
		log.Print("Discordgo connection failure:", err)
		return
	}
	log.SetPrefix("[Main]")
	log.Print("Mashiron-V started successfly!")
	discord.UpdateGameStatus(0, config.CurrentConfig.Discord.Status)
	defer discord.Close()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Print("Mashiron-V is gracefully shutdowning!")
}
