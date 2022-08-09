package runner

import (
	"Mashiron-V/lib/config"
	"Mashiron-V/lib/embed"
	"Mashiron-V/lib/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dop251/goja"
	"github.com/patrickmn/go-cache"
)

var codeCache *cache.Cache

func init() {
	codeCache = cache.New(24*time.Hour, 1*time.Hour)
}

func js(msgInfo *embed.MsgInfo, js string) (*string, error) {
	if !config.CurrentConfig.Js.Enabled {
		return nil, errors.New("js functions are disabled")
	}

	vm := goja.New()
	crc := utils.Crc(js)
	var program *goja.Program
	cached, exists := codeCache.Get(crc)
	if !exists {
		res, err := goja.Compile("", js, false)
		if err != nil {
			return nil, err
		}
		program = res
		codeCache.Set(crc, program, cache.DefaultExpiration)
	} else {
		program = cached.(*goja.Program)
	}

	vm.Set("rawMsg", msgInfo.OrgMsg.Content)
	vm.Set("scriptArgs", strings.Split(msgInfo.OrgMsg.Content, " "))
	vm.Set("name", msgInfo.OrgMsg.Author.Username)
	vm.Set("nickname", msgInfo.OrgMsg.Member.Nick)
	vm.Set("channelId", msgInfo.OrgMsg.ChannelID)
	vm.Set("id", msgInfo.OrgMsg.ID)
	if msgInfo.OrgMsg.ReferencedMessage != nil {
		vm.Set("referencedMessage", msgInfo.OrgMsg.ReferencedMessage.Content)
	}
	vm.Set("author_id", msgInfo.OrgMsg.Author.ID)
	vm.Set("author_avatar", msgInfo.OrgMsg.Author.Avatar)
	channel, err := msgInfo.Session.State.Channel(msgInfo.OrgMsg.ChannelID)
	if err != nil {
		channel, err = msgInfo.Session.Channel(msgInfo.OrgMsg.ChannelID)
		if err != nil {
			return nil, err
		}
	}
	vm.Set("channel_isNsfw", channel.NSFW)
	category, err := getCategory(msgInfo.Session, channel)
	if err != nil {
		return nil, err
	}
	vm.Set("channel_category_id", category)
	vm.Set("sendEmbed", func(data string) string {
		msg := &discordgo.MessageEmbed{}
		err := json.Unmarshal([]byte(data), msg)
		if err != nil {
			return err.Error()
		}
		_, err = msgInfo.Session.ChannelMessageSendEmbed(msgInfo.OrgMsg.ChannelID, msg)
		if err != nil {
			return err.Error()
		}
		return ""
	})

	timer := time.AfterFunc(time.Duration(config.CurrentConfig.Js.Timeout)*time.Millisecond, func() {
		vm.Interrupt("timeout")
	})
	val, err := vm.RunProgram(program)
	timer.Stop()
	if err != nil {
		return nil, err
	}

	v := fmt.Sprintf("%v", val)

	if config.CurrentConfig.Debug {
		log.Print(v)
	}
	return &v, nil
}

func getCategory(session *discordgo.Session, channel *discordgo.Channel) (*string, error) {
	parentCh, err := session.State.Channel(channel.ParentID)
	if err != nil {
		parentCh, err = session.Channel(channel.ParentID)
		if err != nil {
			return nil, err
		}
	}
	if parentCh.Type == discordgo.ChannelTypeGuildCategory {
		return &parentCh.ID, nil
	} else {
		return getCategory(session, channel)
	}
}

func CountJsCache() int {
	return codeCache.ItemCount()
}
