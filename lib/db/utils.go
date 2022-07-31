package db

import (
	"errors"
	"log"

	"github.com/goccy/go-yaml"
)

func ToYaml(file *[]byte, flex bool) (*map[string]*Schema, error) {

	var tmpData map[string]*interface{}
	data := map[string]*Schema{}
	err := yaml.Unmarshal(*file, &data)
	if err != nil && flex {
		//fallback for old "fly" yaml
		err = yaml.Unmarshal(*file, &tmpData)
		if err != nil {
			return nil, errors.New(yaml.FormatError(err, false, true))
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
		return nil, errors.New(yaml.FormatError(err, false, true))
	}
	return &data, nil
}

func Marge(from *map[string]*Schema, to *map[string]*Schema) *map[string]*Schema {
	for k, v := range *from {
		(*to)[k] = v
	}
	return to
}
