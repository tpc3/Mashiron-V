package db

import (
	"Mashiron-V/lib/config"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/patrickmn/go-cache"
)

func LoadData(id *string) (*map[string]*Schema, error) {
	val, exists := dataCache.Get(*id)
	if exists {
		return val.(*map[string]*Schema), nil
	}

	file, err := os.ReadFile(config.CurrentConfig.Data + *id + ".yaml")
	if os.IsNotExist(err) {
		return &map[string]*Schema{}, nil
	} else if err != nil {
		return nil, err
	}

	data, err := ToYaml(&file, true)
	if err != nil {
		return nil, err
	}

	for _, v := range *data {
		err := VerifySchema(v)
		if err != nil {
			return nil, err
		}
	}

	dataCache.Set(*id, data, cache.DefaultExpiration)
	return data, nil
}

func SaveData(id *string, def *map[string]*Schema) error {
	data, err := yaml.Marshal(def)
	if err != nil {
		return err
	}
	err = os.WriteFile(config.CurrentConfig.Data+*id+".yaml", data, os.ModePerm)
	if err != nil {
		return err
	}
	dataCache.Delete(*id)
	return nil
}
