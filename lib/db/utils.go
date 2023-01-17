package db

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/goccy/go-yaml"
)

type FileLang string

const (
	FileLangYaml FileLang = "yaml"
	FileLangToml FileLang = "toml"
	FileLangJson FileLang = "json"
)

func ParseData(lang FileLang, file *[]byte, flex bool) (*map[string]*Schema, error) {
	data := map[string]*Schema{}
	err := unmarshal(lang, *file, &data)
	if err != nil && flex {
		//fallback for old "fly" yaml
		tmpData := map[string]*interface{}{}
		err = unmarshal(lang, *file, &tmpData)
		if err != nil {
			return nil, err
		}

		for k, v := range tmpData {
			schema, err := Convert(v)
			if err != nil {
				return nil, err
			}
			if schema == nil {
				log.Fatal("schema is nil")
			}
			data[k] = schema
		}
	} else if err != nil {
		return nil, err
	}
	return &data, nil
}

func unmarshal(lang FileLang, file []byte, result any) error {
	switch lang {
	case FileLangYaml:
		err := yaml.Unmarshal(file, result)
		if err != nil {
			return errors.New(yaml.FormatError(err, false, true))
		}
		return nil
	case FileLangToml:
		err := toml.Unmarshal(file, result)
		if parseErr, ok := err.(toml.ParseError); ok {
			err = errors.New(parseErr.ErrorWithUsage())
		}
		return err
	case FileLangJson:
		return json.Unmarshal(file, result)
	default:
		log.Panic("Invalid lang of data")
		return nil
	}
}

func Marge(from *map[string]*Schema, to *map[string]*Schema) *map[string]*Schema {
	for k, v := range *from {
		(*to)[k] = v
	}
	return to
}
