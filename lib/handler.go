package lib

import (
	"Mashiron-V/lib/cmds"
	"Mashiron-V/lib/config"
	"Mashiron-V/lib/db"
	"Mashiron-V/lib/embed"
	"Mashiron-V/lib/runner"
	"log"
	"runtime/debug"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func MessageCreate(session *discordgo.Session, orgMsg *discordgo.MessageCreate) {
	var start time.Time
	if config.CurrentConfig.Debug {
		start = time.Now()
	}

	defer func() {
		if err := recover(); err != nil {
			log.Print("Trying to recover from fatal error: ", err)
			debug.PrintStack()
		}
	}()

	msgInfo := embed.MsgInfo{
		Session: session,
		OrgMsg:  orgMsg,
		Lang:    config.CurrentConfig.Guild.Lang,
	}

	guild, err := config.LoadGuild(&orgMsg.GuildID)
	if err != nil {
		embed.SendUnknownErrorEmbed(&msgInfo, err)
	}
	msgInfo.Lang = guild.Lang
	data, err := db.LoadData(&orgMsg.GuildID)
	if err != nil {
		embed.SendUnknownErrorEmbed(&msgInfo, err)
	}

	if orgMsg.Author.ID == session.State.User.ID || orgMsg.Content == "" {
		return
	}

	if strings.HasPrefix(orgMsg.Content, guild.Prefix) {
		cmd := strings.SplitN(strings.SplitN(orgMsg.Content, "\n", 2)[0], " ", 2)[0][len(guild.Prefix):]
		switch cmd {
		case cmds.Ping:
			cmds.PingCmd(&msgInfo)
		case cmds.Help:
			cmds.HelpCmd(&msgInfo)
		case cmds.List:
			cmds.ListCmd(&msgInfo, data)
		case cmds.Delete:
			cmds.DeleteCmd(&msgInfo, data)
		case cmds.Add:
			cmds.AddCmd(&msgInfo, data)
		case cmds.Config:
			cmds.ConfigCmd(&msgInfo, *guild)
		default:
			val, exists := (*data)[cmd]
			if exists {
				err := runner.Run(val, &msgInfo)
				if err != nil {
					embed.SendUnknownErrorEmbed(&msgInfo, err)
				}
			}
		}

		if config.CurrentConfig.Debug {
			log.Print("Processed in ", time.Since(start).Milliseconds(), "ms.")
		}
		return
	}

	for k, v := range *data {
		if v.Trigger.Content == nil && v.Trigger.Uid == nil && k == orgMsg.Content {
			err = runner.Run(v, &msgInfo)
			if err != nil {
				embed.SendUnknownErrorEmbed(&msgInfo, err)
			}
		} else {
			err = runner.Trigger(v, &msgInfo)
			if err != nil {
				embed.SendUnknownErrorEmbed(&msgInfo, err)
			}
		}
	}

	if config.CurrentConfig.Debug {
		log.Print("Processed in ", time.Since(start).Nanoseconds(), "ns.")
	}
}
