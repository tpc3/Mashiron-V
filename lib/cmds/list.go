package cmds

import (
	"Mashiron-V/lib/config"
	"Mashiron-V/lib/db"
	"Mashiron-V/lib/embed"
	"bytes"
	"io/ioutil"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

const List = "list"

func ListCmd(msgInfo *embed.MsgInfo, data *map[string]*db.Schema) {
	var file []byte
	var err error
	split := strings.Split(msgInfo.OrgMsg.Content, " ")
	if len(split) > 1 {
		tmpData := map[string]*db.Schema{}
		for _, v := range split[1:] {
			val, ok := (*data)[v]
			if !ok {
				embed.SendErrorEmbed(msgInfo, config.Lang[msgInfo.Lang].Error.NoEntry+"\n`"+v+"`")
				return
			}
			tmpData[v] = val
		}
		file, err = yaml.Marshal(tmpData)
		if err != nil {
			embed.SendUnknownErrorEmbed(msgInfo, err)
		}
	} else {
		file, err = ioutil.ReadFile(config.CurrentConfig.Data + msgInfo.OrgMsg.GuildID + ".yaml")
		if os.IsNotExist(err) {
			err = embed.SendErrorEmbed(msgInfo, config.Lang[msgInfo.Lang].Error.ZeroEntry)
			if err != nil {
				embed.SendUnknownErrorEmbed(msgInfo, err)
			}
			return
		} else if err != nil {
			embed.SendUnknownErrorEmbed(msgInfo, err)
			return
		}
	}

	str := string(file)
	if len(str) < 1980 {
		err = embed.SendMessageEmbed(msgInfo, List, "```yaml\n"+str+"\n```")
		if err != nil {
			embed.SendUnknownErrorEmbed(msgInfo, err)
		}
	} else {
		_, err = msgInfo.Session.ChannelFileSend(msgInfo.OrgMsg.ChannelID, "list.yaml", bytes.NewReader(file))
		if err != nil {
			embed.SendUnknownErrorEmbed(msgInfo, err)
		}
	}
}
