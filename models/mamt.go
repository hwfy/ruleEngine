package models

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/hwfy/redis"
)

const RMDS = "rule"

type (
	Rules struct {
		Name  string      `json:"name"`
		Rules interface{} `json:"rules"`
	}
	Rule struct {
		Conditions  []Condition `json:"conditions"`
		Result      string      `json:"result"`
		Description string      `json:"description"`
	}
	Condition struct {
		Name  string
		Opera string
		Value interface{}
	}
	Result struct {
		id Rule
	}
)

func newRules(name string) ([]Rule, error) {
	client, err := redis.NewClient(RMDS)
	if err != nil {
		return nil, err
	}
	data, err := client.HValues(name)
	if err != nil {
		return nil, err
	}
	var rules []Rule

	err = json.Unmarshal(data, &rules)
	if err != nil {
		return nil, err
	}
	return rules, nil
}

func SetAll(list []Rules) error {
	client, err := redis.NewClient(RMDS)
	if err != nil {
		return err
	}
	for _, r := range list {
		switch rules := r.Rules.(type) {
		// add a rule set
		case []interface{}:
			for _, rule := range rules {
				id := strconv.Itoa(client.HLen(r.Name))

				err = client.HSet(r.Name, id, rule)
				if err != nil {
					return err
				}
			}
		// update the rule set
		case map[string]interface{}:
			for id, rule := range rules {
				err = client.HSet(r.Name, id, rule)
				if err != nil {
					return err
				}
			}
		default:
			return errors.New("The rule set type is not an array or map")
		}
		// generate parsing files
		err = genAnalyFile(r.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetNames() ([]string, error) {
	client, err := redis.NewClient(RMDS)
	if err != nil {
		return nil, err
	}
	return client.Keys()
}

func GetAll(name string) (map[string]Rule, error) {
	client, err := redis.NewClient(RMDS)
	if err != nil {
		return nil, err
	}
	bytes, err := client.HGetAll(name)
	if err != nil {
		return nil, err
	}
	var rules map[string]Rule

	err = json.Unmarshal(bytes, &rules)
	if err != nil {
		return nil, err
	}
	return rules, nil
}

func DeleteAll(data map[string][]string) error {
	client, err := redis.NewClient(RMDS)
	if err != nil {
		return err
	}
	for name, ids := range data {
		for _, id := range ids {
			// delete the rule set
			err = client.HDel(name, id)
			if err != nil {
				return err
			}
		}
		// if the rule exists, update the parsing file
		if client.Exist(name) {
			err = genAnalyFile(name)
			if err != nil {
				return err
			}
		} else {
			delAnalyFile(name)
		}
	}
	return nil
}

func GetOne(name, id string) (*Rule, error) {
	client, err := redis.NewClient(RMDS)
	if err != nil {
		return nil, err
	}
	bytes, err := client.HGet(name, id)
	if err != nil {
		return nil, err
	}
	rule := new(Rule)

	err = json.Unmarshal(bytes, rule)
	if err != nil {
		return nil, err
	}
	return rule, nil
}

func Update(name, id string, rule Rule) error {
	client, err := redis.NewClient(RMDS)
	if err != nil {
		return err
	}
	err = client.HSet(name, id, rule)
	if err != nil {
		return err
	}
	// update the parsing file
	return genAnalyFile(name)
}

func Delete(name, id string) error {
	client, err := redis.NewClient(RMDS)
	if err != nil {
		return err
	}
	err = client.HDel(name, id)
	if err != nil {
		return err
	}
	// if the rule exists, update the parsing file
	if client.Exist(name) {
		return genAnalyFile(name)
	}
	return delAnalyFile(name)
}
