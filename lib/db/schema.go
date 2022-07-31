package db

import (
	"Mashiron-V/lib/config"
	"errors"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/koron/go-dproxy"
	"github.com/patrickmn/go-cache"
)

type Schema struct {
	Trigger struct {
		Content []string `yaml:",omitempty"`
		Uid     []uint64 `yaml:",omitempty"`
	} `yaml:",omitempty"`
	ReturnStr []string `yaml:"return,omitempty"`
	React     []string `yaml:",omitempty"`
	Js        string   `yaml:",omitempty"`
}

var dataCache *cache.Cache

func init() {
	dataCache = cache.New(24*time.Hour, 1*time.Hour)
	err := os.MkdirAll(config.CurrentConfig.Data, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func Convert(v *interface{}) (*Schema, error) {
	var contentRes []string
	var uidRes []uint64
	var returnStrRes []string
	var reactRes []string
	var jsRes string

	proxy := dproxy.New(*v)

	switch proxy.(type) {
	case dproxy.Error:
		return nil, errors.New("there's nothing")
	}
	str, err := proxy.String()
	if err != nil {
		trigger := proxy.M("trigger")
		switch trigger.(type) {
		case dproxy.Error:
			// Uncomment for absolute "fly" compatibility
			// return nil, errors.New("trigger is required")
		}

		str, err = trigger.String()
		if err != nil {
			var noContent bool
			content := trigger.M("content")
			switch content.(type) {
			case dproxy.Error:
				noContent = true
			default:
				noContent = false
				str, err := content.String()
				if err != nil {
					arr, err := content.Array()
					if err != nil {
						return nil, err
					} else {
						for _, v := range arr {
							v, ok := v.(string)
							if !ok {
								return nil, err
							}
							contentRes = append(contentRes, v)
						}
					}
				} else {
					contentRes = append(contentRes, str)
				}
			}

			uid := trigger.M("uid")
			switch uid.(type) {
			case dproxy.Error:
				// Uncomment for absolute "fly" compatibility
				if noContent {
					// 	return nil, errors.New("content or uid is required")
				}
			default:
				val, err := uid.Value()
				if err != nil {
					return nil, err
				}
				uidVal, ok := val.(uint64)
				if !ok {
					arr, err := uid.Array()
					if err != nil {
						return nil, err
					} else {
						for _, v := range arr {
							v, ok := v.(uint64)
							if !ok {
								return nil, errors.New("uid must be uint64")
							}
							uidRes = append(uidRes, v)
						}
					}
				} else {
					uidRes = append(uidRes, uidVal)
				}
			}
		} else {
			contentRes = append(contentRes, str)
		}
		var noReturn bool
		var noReact bool
		returnStr := proxy.M("return")
		switch returnStr.(type) {
		case dproxy.Error:
			noReturn = true
		default:
			noReturn = false
			str, err = returnStr.String()
			if err != nil {
				arr, err := returnStr.Array()
				if err != nil {
					return nil, err
				} else {
					for _, v := range arr {
						v, ok := v.(string)
						if !ok {
							return nil, err
						}
						returnStrRes = append(returnStrRes, v)
					}
				}
			} else {
				returnStrRes = append(returnStrRes, str)
			}
		}

		for _, v := range reactRes {
			if v == "" {
				panic("what")
			}
		}

		react := proxy.M("react")
		switch react.(type) {
		case dproxy.Error:
			noReact = true
		default:
			noReact = false
			str, err = react.String()
			if err != nil {
				val, err := react.Value()
				if err != nil {
					return nil, err
				}
				reactVal, ok := val.(uint64)
				if !ok {
					arr, err := react.Array()
					if err != nil {
						return nil, errors.New("react must be uint64 or string")
					} else {
						for _, v := range arr {
							res, ok := v.(string)
							if !ok {
								v, ok := v.(uint64)
								if !ok {
									return nil, errors.New("react must be uint64 or string")
								} else {
									reactRes = append(reactRes, strconv.FormatUint(v, 10))
								}
							} else {
								reactRes = append(reactRes, res)
							}
						}
					}
				} else {
					reactRes = append(reactRes, strconv.FormatUint(reactVal, 10))
				}
			} else {
				reactRes = append(reactRes, str)
			}
		}

		js := proxy.M("js")
		switch js.(type) {
		case dproxy.Error:
			if noReact && noReturn {
				return nil, errors.New("return or react or js is required")
			}
		default:
			jsRes, err = js.String()
			if err != nil {
				return nil, err
			}
		}
	} else {
		returnStrRes = append(returnStrRes, str)
	}

	return &Schema{
		Trigger: struct {
			Content []string `yaml:",omitempty"`
			Uid     []uint64 `yaml:",omitempty"`
		}{
			Content: contentRes,
			Uid:     uidRes,
		},
		ReturnStr: returnStrRes,
		React:     reactRes,
		Js:        jsRes,
	}, nil
}

func VerifySchema(schema *Schema) error {
	for _, v := range schema.Trigger.Content {
		_, err := regexp.Compile(v)
		if err != nil {
			return err
		}
	}
	for _, v := range schema.React {
		if len([]rune(v)) != 1 {
			_, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return err
			}
		}
	}
	if len(schema.ReturnStr) == 0 && len(schema.React) == 0 && schema.Js == "" {
		return errors.New("nothing to do")
	}
	return nil
}

func CountCache() int {
	return dataCache.ItemCount()
}
