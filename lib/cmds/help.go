package cmds

import (
	"Mashiron-V/lib/config"
	"Mashiron-V/lib/embed"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const Help = "help"

func HelpCmd(msgInfo *embed.MsgInfo) {
	msg := embed.NewEmbed(msgInfo.Session, msgInfo.OrgMsg)
	msg.Title = cases.Title(language.Und, cases.NoLower).String(Help)
	msg.Description = config.Lang[msgInfo.Lang].Help + "\n" + config.CurrentConfig.Help
	msgInfo.Session.ChannelMessageSendEmbed(msgInfo.OrgMsg.ChannelID, msg)
}
