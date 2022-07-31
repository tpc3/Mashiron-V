package config

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/patrickmn/go-cache"
)

type Config struct {
	Debug   bool
	Help    string
	Data    string
	Config  string
	Discord struct {
		Token  string
		Status string
	}
	Db struct {
		Kind string
		Path string
	}
	Js struct {
		Enabled bool
		Timeout uint
	}
	Guild Guild
}

type Guild struct {
	Prefix string `yaml:",omitempty"`
	Lang   string `yaml:",omitempty"`
}

const configFile = "./config.yaml"

var (
	CurrentConfig Config
	cachedGuild   *cache.Cache
)

func init() {
	loadLang()
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal("Config load failed: ", err)
	}
	err = yaml.Unmarshal(file, &CurrentConfig)
	if err != nil {
		log.Fatal("Config parse failed: ", err)
	}

	//verify
	if CurrentConfig.Debug {
		log.Print("Debug is enabled")
	}
	if CurrentConfig.Discord.Token == "" {
		log.Fatal("Token is empty")
	}
	err = VerifyGuild(&CurrentConfig.Guild)
	if err != nil {
		log.Fatal("Config verify failed: ", err)
	}

	cachedGuild = cache.New(24*time.Hour, 1*time.Hour)
}

func VerifyGuild(guild *Guild) error {
	if len(guild.Prefix) == 0 || len(guild.Prefix) >= 10 {
		return errors.New("prefix is too short or long")
	}
	_, exists := Lang[guild.Lang]
	if !exists {
		return errors.New("language does not exists")
	}
	return nil
}

func LoadGuild(id *string) (*Guild, error) {
	val, exists := cachedGuild.Get(*id)
	if exists {
		return val.(*Guild), nil
	}

	err := os.MkdirAll(CurrentConfig.Config, os.ModePerm)
	if err != nil {
		return nil, err
	}
	file, err := ioutil.ReadFile(CurrentConfig.Config + *id + ".yaml")
	if os.IsNotExist(err) {
		return &Guild{
			Prefix: CurrentConfig.Guild.Prefix,
			Lang:   CurrentConfig.Guild.Lang,
		}, nil
	} else if err != nil {
		return nil, err
	}

	var guild Guild
	err = yaml.Unmarshal(file, &guild)
	if err != nil {
		return nil, err
	}

	cachedGuild.Set(*id, guild, cache.DefaultExpiration)
	return &guild, nil
}

func SaveGuild(id *string, guild *Guild) error {
	if guild.Lang == CurrentConfig.Guild.Lang && guild.Prefix == CurrentConfig.Guild.Prefix {
		err := os.Remove(CurrentConfig.Config + *id + ".yaml")
		if err != nil {
			return err
		}
		cachedGuild.Delete(*id)
		return nil
	}
	data, err := yaml.Marshal(guild)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(CurrentConfig.Config+*id+".yaml", data, 0666)
	if err != nil {
		return err
	}
	cachedGuild.Set(*id, guild, cache.DefaultExpiration)
	return nil
}

func CountCache() int {
	return cachedGuild.ItemCount()
}
