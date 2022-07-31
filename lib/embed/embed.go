package embed

import (
	"Mashiron-V/lib/config"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// https://material.io/archive/guidelines/style/color.html#color-color-palette
const (
	ColorBlue   = 0xB3E5FC
	ColorPink   = 0xf50057
	ColorYellow = 0xFFF176
)

var UnknownErrorNum int

type MsgInfo struct {
	Session *discordgo.Session
	OrgMsg  *discordgo.MessageCreate
	Lang    string
}

func init() {
	UnknownErrorNum = 0
}

func NewEmbed(session *discordgo.Session, orgMsg *discordgo.MessageCreate) *discordgo.MessageEmbed {
	now := time.Now()
	msg := &discordgo.MessageEmbed{}
	msg.Author = &discordgo.MessageEmbedAuthor{}
	msg.Footer = &discordgo.MessageEmbedFooter{}
	msg.Author.IconURL = session.State.User.AvatarURL("256")
	msg.Author.Name = session.State.User.Username
	msg.Footer.IconURL = orgMsg.Author.AvatarURL("256")
	msg.Footer.Text = "Request from " + orgMsg.Author.Username + " @ " + now.String()
	msg.Color = ColorBlue
	return msg
}

func NewErrorEmbed(msgInfo *MsgInfo, description string) *discordgo.MessageEmbed {
	msg := NewEmbed(msgInfo.Session, msgInfo.OrgMsg)
	msg.Color = ColorPink
	msg.Title = config.Lang[msgInfo.Lang].Error.Title
	msg.Description = description
	return msg
}

func NewUnknownErrorEmbed(msgInfo *MsgInfo, err error) *discordgo.MessageEmbed {
	log.Print("WARN: UnknownError called:", err)
	UnknownErrorNum++
	return NewErrorEmbed(msgInfo, config.Lang[msgInfo.Lang].Error.Unknown)
}

func SendMessageEmbed(msgInfo *MsgInfo, title string, description string) error {
	msg := NewEmbed(msgInfo.Session, msgInfo.OrgMsg)
	msg.Title = cases.Title(language.Und, cases.NoLower).String(title)
	msg.Description = description
	_, err := msgInfo.Session.ChannelMessageSendEmbed(msgInfo.OrgMsg.ChannelID, msg)
	return err
}

func SendWarningEmbed(msgInfo *MsgInfo, description string) error {
	msg := NewEmbed(msgInfo.Session, msgInfo.OrgMsg)
	msg.Title = config.Lang[msgInfo.Lang].Warning.Title
	msg.Description = description
	msg.Color = ColorYellow
	_, err := msgInfo.Session.ChannelMessageSendEmbed(msgInfo.OrgMsg.ChannelID, msg)
	return err
}

func SendErrorEmbed(msgInfo *MsgInfo, description string) error {
	_, err := msgInfo.Session.ChannelMessageSendEmbed(msgInfo.OrgMsg.ChannelID, NewErrorEmbed(msgInfo, description))
	return err
}

func SendUnknownErrorEmbed(msgInfo *MsgInfo, err error) error {
	_, err = msgInfo.Session.ChannelMessageSendEmbed(msgInfo.OrgMsg.ChannelID, NewUnknownErrorEmbed(msgInfo, err))
	return err
}
