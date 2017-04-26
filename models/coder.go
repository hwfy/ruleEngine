package models

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// genDir Parse the file storage directory： ../ruleAnaly/models/
const genDir = "../ruleAnaly/models/"

// genAnalyFile Generate analytic source code according to rule name
func genAnalyFile(name string) error {
	head := `
	package models
	import "errors"
	func (a *analyServer)` + strings.Title(name) + `(data map[string]interface{})(string, error){
	`
	rules, err := newRules(name)
	if err != nil {
		return err
	}
	judge, err := getExpression(rules)
	if err != nil {
		return err
	}
	code := fmt.Sprintf("%s%s%s return %q,errors.New(err)}", head, getValue(rules), judge, "")

	return saver(name, code)
}

// getValue Obtain actual business data according to the condition name
func getValue(rules []Rule) (value string) {
	for _, r := range rules {
		for _, c := range r.Conditions {
			if !strings.Contains(value, c.Name) {
				err := c.Name + " does not exist or type is not "
				switch c.Value.(type) {
				case int:
					value += fmt.Sprintf("%s,ok:=data[%q].(int)\n if !ok{err=%q}\n", c.Name, c.Name, err+"int")
				case bool:
					value += fmt.Sprintf("%s,ok:=data[%q].(bool)\n if !ok{err=%q}\n", c.Name, c.Name, err+"bool")
				case float64:
					if c.Opera == "has" || c.Opera == "nhas" {
						value += fmt.Sprintf("%s,ok:=data[%q].([]interface{})\n if !ok{err=%q}\n", c.Name, c.Name, err+"[]interface{}")
					} else {
						value += fmt.Sprintf("%s,ok:=data[%q].(float64)\n if !ok{err=%q}\n", c.Name, c.Name, err+"float64")
					}
				case string:
					if c.Opera == "has" || c.Opera == "nhas" {
						value += fmt.Sprintf("%s,ok:=data[%q].([]interface{})\n if !ok{err=%q}\n", c.Name, c.Name, err+"[]interface{}")
					} else {
						value += fmt.Sprintf("%s,ok:=data[%q].(string)\n if !ok{err=%q}\n", c.Name, c.Name, err+"string")
					}
				case []interface{}:
					if c.Opera == "in" || c.Opera == "nin" {
						value += fmt.Sprintf("%s,ok:=data[%q]\n if !ok{err=%q}\n", c.Name, c.Name, c.Name+" does not exist")
					}
				}
			}
		}
	}
	return
}

// getExpression Get the if expression
func getExpression(rules []Rule) (judge string, err error) {
	for _, r := range rules {
		var exp string
		for _, c := range r.Conditions {
			switch actual := c.Value.(type) {
			case int:
				exp += fmt.Sprintf("%s%s%d&&", c.Name, c.Opera, actual)
			case bool:
				switch c.Opera {
				case "==":
					exp += fmt.Sprintf("%s==%t&&", c.Name, actual)
				case "!=":
					exp += fmt.Sprintf("%s!=%t&&", c.Name, actual)
				default:
					err = fmt.Errorf("The value of %s is the bool type and does not support the relation %s", c.Name, c.Opera)
				}
			case float64:
				switch c.Opera {
				case "has":
					exp += fmt.Sprintf(" contains(%v,%f)&&", c.Name, actual)
				case "nhas":
					exp += fmt.Sprintf(" !contains(%v,%f)&&", c.Name, actual)
				default:
					exp += fmt.Sprintf("%s%s%v&&", c.Name, c.Opera, actual)
				}
			case string:
				switch c.Opera {
				case "has":
					exp += fmt.Sprintf(" contains(%v,%q)&&", c.Name, actual)
				case "nhas":
					exp += fmt.Sprintf(" !contains(%v,%q)&&", c.Name, actual)
				case "==":
					exp += fmt.Sprintf("%s==%q&&", c.Name, actual)
				case "!=":
					exp += fmt.Sprintf("%s!=%q&&", c.Name, actual)
				default:
					err = fmt.Errorf("The value of %s is the string type and does not support the relation %s", c.Name, c.Opera)
				}
			case []interface{}:
				switch c.Opera {
				case "in":
					exp += fmt.Sprintf(" contains(%#v, %s)&&", actual, c.Name)
				case "nin":
					exp += fmt.Sprintf(" !contains(%#v, %s)&&", actual, c.Name)
				default:
					err = fmt.Errorf("The value of %s is the []interface{} type and does not support the relation %s", c.Name, c.Opera)
				}
			}
		}
		exp = strings.TrimSuffix(exp, "&&")

		judge += fmt.Sprintf("if %s{return %q,nil}\n", exp, r.Result)
	}
	return
}

// saver Generate go files based on rule name and parse source
func saver(name, code string) error {
	// format code
	cmd := exec.Command("gofmt")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	io.WriteString(stdin, code)
	stdin.Close()

	bytes, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("Generate parsing file failed: " + name + " syntax error")
	}
	// write file
	shortPath := filepath.Clean(genDir + name + ".go")

	file, err := os.Create(shortPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(bytes)

	return err
}

// delAnalyFile Remove the parsing file according to the rule name
func delAnalyFile(name string) error {
	file := filepath.Clean(genDir + name + ".go")

	return os.Remove(file)
}
