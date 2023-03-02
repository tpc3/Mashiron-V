package cmds

import (
	"Mashiron-V/lib/config"
	"Mashiron-V/lib/db"
	"Mashiron-V/lib/embed"
	"strings"

	"github.com/goccy/go-yaml"
)

const Add = "add"

func AddCmd(msgInfo *embed.MsgInfo, data *map[string]*db.Schema) {
	lines := strings.Split(msgInfo.OrgMsg.Content, "\n")
	if len(lines) < 4 || !strings.HasPrefix(lines[1], "```") || lines[len(lines)-1] != "```" {
		err := embed.SendErrorEmbed(msgInfo, config.Lang[msgInfo.Lang].Error.Invalid)
		if err != nil {
			embed.SendUnknownErrorEmbed(msgInfo, err)
		}
		return
	}

	var lang db.FileLang
	switch lines[1][3:] {
	case "yaml", "yml":
		lang = db.FileLangYaml
	case "toml":
		lang = db.FileLangToml
	case "json":
		lang = db.FileLangJson
	default:
		err := embed.SendErrorEmbed(msgInfo, config.Lang[msgInfo.Lang].Error.InvalidLang)
		if err != nil {
			embed.SendUnknownErrorEmbed(msgInfo, err)
		}
		return
	}

	opt := strings.Split(lines[0], " ")
	flex := false
	for _, v := range opt[1:] {
		switch v {
		case "--flex":
			flex = true
		}
	}

	input := []byte(strings.Join(lines[2:len(lines)-1], "\n"))
	res, err := db.ParseData(lang, &input, flex)
	if err != nil {
		embed.SendErrorEmbed(msgInfo, "```\n"+err.Error()+"```")
		return
	}
	for k := range *res {
		if strings.Contains(k, " ") {
			err = embed.SendErrorEmbed(msgInfo, config.Lang[msgInfo.Lang].Error.Invalid)
			if err != nil {
				embed.SendUnknownErrorEmbed(msgInfo, err)
			}
			return
		}
		val, exists := (*data)[k]
		if exists {
			res, err := yaml.Marshal(map[string]*db.Schema{k: val})
			if err != nil {
				embed.SendUnknownErrorEmbed(msgInfo, err)
				return
			}
			resStr := string(res)
			if len([]rune(resStr)) > 1900 {
				resStr = string([]rune(resStr)[:1897]) + "..."
			}
			err = embed.SendWarningEmbed(msgInfo, config.Lang[msgInfo.Lang].Warning.Overwrite+"\n```yaml\n"+resStr+"\n```")
			if err != nil {
				embed.SendUnknownErrorEmbed(msgInfo, err)
				return
			}
		}
	}
	for _, v := range *res {
		err = db.VerifySchema(v)
		if err != nil {
			err = embed.SendErrorEmbed(msgInfo, err.Error())
			if err != nil {
				embed.SendUnknownErrorEmbed(msgInfo, err)
			}
			return
		}
	}

	db.Merge(res, data)
	err = db.SaveData(&msgInfo.OrgMsg.GuildID, data)
	if err != nil {
		embed.SendUnknownErrorEmbed(msgInfo, err)
		return
	}
	file, err := yaml.Marshal(res)
	if err != nil {
		embed.SendUnknownErrorEmbed(msgInfo, err)
	}
	err = embed.SendMessageEmbed(msgInfo, Add, config.Lang[msgInfo.Lang].Done+"\n```yaml\n"+string(file)+"\n```")
	if err != nil {
		embed.SendUnknownErrorEmbed(msgInfo, err)
	}
}
