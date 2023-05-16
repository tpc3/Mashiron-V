package cmds

import (
	"Mashiron-V/lib/config"
	"Mashiron-V/lib/db"
	"Mashiron-V/lib/embed"
	"strings"
)

const Delete = "delete"

func DeleteCmd(msgInfo *embed.MsgInfo, data *map[string]*db.Schema) {
	split := strings.Split(msgInfo.OrgMsg.Content, " ")
	if len(split) < 2 {
		embed.SendErrorEmbed(msgInfo, config.Lang[msgInfo.Lang].Error.Invalid)
		return
	}
	for _, v := range split[1:] {
		_, exists := (*data)[v]
		if !exists {
			embed.SendErrorEmbed(msgInfo, config.Lang[msgInfo.Lang].Error.NoEntry+"\n`"+v+"`")
			return
		}
		delete(*data, v)
	}

	err := db.SaveData(&msgInfo.OrgMsg.GuildID, data)
	if err != nil {
		embed.SendUnknownErrorEmbed(msgInfo, err)
		return
	}
	err = embed.SendMessageEmbed(msgInfo, Delete, config.Lang[msgInfo.Lang].Done)
	if err != nil {
		embed.SendUnknownErrorEmbed(msgInfo, err)
	}
}
