package cmds

import (
	"Mashiron-V/lib/config"
	"Mashiron-V/lib/embed"
	"strings"

	"github.com/goccy/go-yaml"
)

const Config = "config"

func ConfigCmd(msgInfo *embed.MsgInfo, guild config.Guild) {
	split := strings.Split(msgInfo.OrgMsg.Content, " ")
	if len(split) < 3 {
		file, err := yaml.Marshal(guild)
		if err != nil {
			embed.SendUnknownErrorEmbed(msgInfo, err)
		}
		embed.SendMessageEmbed(msgInfo, Config, "```yaml\n"+string(file)+"\n```")
		if err != nil {
			embed.SendUnknownErrorEmbed(msgInfo, err)
		}
		return
	}
	switch split[1] {
	case "prefix":
		guild.Prefix = split[2]
	case "lang":
		guild.Lang = split[2]
	}
	err := config.VerifyGuild(&guild)
	if err != nil {
		err := embed.SendErrorEmbed(msgInfo, config.Lang[msgInfo.Lang].Error.Invalid)
		if err != nil {
			embed.SendUnknownErrorEmbed(msgInfo, err)
		}
		return
	}
	err = config.SaveGuild(&msgInfo.OrgMsg.GuildID, &guild)
	if err != nil {
		embed.SendUnknownErrorEmbed(msgInfo, err)
		return
	}
	err = embed.SendMessageEmbed(msgInfo, Config, config.Lang[msgInfo.Lang].Done)
	if err != nil {
		embed.SendUnknownErrorEmbed(msgInfo, err)
	}
}
