package models

import (
	"encoding/json"

	"github.com/hwfy/redis"
)

const RBDS = "ruleBind"

type (
	Form struct {
		name Bind
	}
	Bind struct {
		key Field
	}
	Field struct {
		Name string `json:"name"`
		Flag string `json:"flag"`
	}
)

func Binding(name string, rules map[string]map[string]Field) error {
	client, err := redis.NewClient(RBDS)
	if err != nil {
		return err
	}
	for rule, field := range rules {
		err := client.HSet(name, rule, field)
		if err != nil {
			return err
		}
	}
	return nil
}

func LookUp(name string) (map[string]map[string]Field, error) {
	client, err := redis.NewClient(RBDS)
	if err != nil {
		return nil, err
	}
	bytes, err := client.HGetAll(name)
	if err != nil {
		return nil, err
	}
	var rules map[string]map[string]Field

	err = json.Unmarshal(bytes, &rules)
	if err != nil {
		return nil, err
	}
	return rules, nil
}

func UnBind(name, rule string) error {
	client, err := redis.NewClient(RBDS)
	if err != nil {
		return err
	}
	return client.HDel(name, rule)
}
