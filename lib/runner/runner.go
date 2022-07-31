package runner

import (
	"Mashiron-V/lib/db"
	"Mashiron-V/lib/embed"
	"Mashiron-V/lib/utils"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/patrickmn/go-cache"
)

var regExCache *cache.Cache

func init() {
	regExCache = cache.New(24*time.Hour, 1*time.Hour)
}

func Run(def *db.Schema, msgInfo *embed.MsgInfo) error {
	rand.Seed(msgInfo.OrgMsg.Timestamp.UnixNano())
	if def.ReturnStr != nil && len(def.ReturnStr) != 0 {
		v := def.ReturnStr[rand.Intn(len(def.ReturnStr))]
		if v != "" {
			_, err := msgInfo.Session.ChannelMessageSendReply(msgInfo.OrgMsg.ChannelID, v, msgInfo.OrgMsg.Reference())
			if err != nil {
				return err
			}
		}
	}
	if def.React != nil && len(def.React) != 0 {
		v := def.React[rand.Intn(len(def.React))]
		if v != "" {
			emoji := ""
			if len([]rune(v)) == 1 {
				emoji = v
			} else {
				guildEmoji, err := msgInfo.Session.GuildEmoji(msgInfo.OrgMsg.GuildID, v)
				if err != nil {
					return err
				}
				emoji = guildEmoji.Name + ":" + guildEmoji.ID
			}
			err := msgInfo.Session.MessageReactionAdd(msgInfo.OrgMsg.ChannelID, msgInfo.OrgMsg.ID, emoji)
			if err != nil {
				return err
			}
		}
	}
	if def.Js != "" {
		str, err := js(msgInfo, def.Js)
		if err != nil {
			err := embed.SendErrorEmbed(msgInfo, err.Error())
			if err != nil {
				embed.SendUnknownErrorEmbed(msgInfo, err)
			}
		} else {
			if str != nil && *str != "" {
				msgInfo.Session.ChannelMessageSendReply(msgInfo.OrgMsg.ChannelID, *str, msgInfo.OrgMsg.Reference())
			}
		}
	}
	return nil
}

func Trigger(def *db.Schema, msgInfo *embed.MsgInfo) error {
	hit := true
	if def.Trigger.Uid != nil && len(def.Trigger.Uid) != 0 {
		hit = false
		for _, v := range def.Trigger.Uid {
			if strconv.FormatUint(v, 10) == msgInfo.OrgMsg.Author.ID {
				hit = true
			}
		}
	}
	if !hit {
		return nil
	}
	hit = false
	if def.Trigger.Content != nil && len(def.Trigger.Content) != 0 {
		for _, v := range def.Trigger.Content {
			var regex *regexp.Regexp
			crc := utils.Crc(v)
			res, ok := regExCache.Get(crc)
			if ok {
				regex = res.(*regexp.Regexp)
			} else {
				res, err := regexp.Compile(v)
				if err != nil {
					return err
				}
				regex = res
				regExCache.Set(crc, res, cache.DefaultExpiration)
			}
			if regex.MatchString(msgInfo.OrgMsg.Content) {
				hit = true
			}
		}
	}
	if hit {
		return Run(def, msgInfo)
	}
	return nil
}

func CountRegExCache() int {
	return regExCache.ItemCount()
}
