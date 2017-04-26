# ruleEngine
A rule engine includes administrative rules, binding forms, and call parsing services
# Installation
> go get github.com/hwfy/ruleEngine
# usage
Analysis services need to build etcd, the default address 127.0.0.1:2379. http listening port 8080, can enter the app.conf to amend.  api document path 127.0.0.1:8080/apidoc.
Create a rule set


----------


 post:  127.0.0.1:8080/v1/rule/mamt/
 data: 
 ```json
[
    {
        "name": "car",
        "rules": [
            {
                "conditions": [
                    {
                        "name": "age",
                        "opera": ">",
                        "value": 18
                    },
                    {
                        "name": "age",
                        "opera": "<=",
                        "value": 25
                    },
                    {
                        "name": "color",
                        "opera": "==",
                        "value": "red"
                    },
                    {
                        "name": "sex",
                        "opera": "==",
                        "value": "male"
                    }
                ],
                "result": "insurance+20%",
                "description": "rules of car insurance"
            },
            {
                "conditions": [
                    {
                        "name": "age",
                        "opera": ">",
                        "value": 25
                    },
                    {
                        "name": "age",
                        "opera": "<=",
                        "value": 35
                    },
                    {
                        "name": "color",
                        "opera": "in",
                        "value": [
                            "red",
                            "black"
                        ]
                    },
                    {
                        "name": "sex",
                        "opera": "==",
                        "value": "male"
                    }
                ],
                "result": "insurance+10%",
                "description": "rules of car insurance"
            }
        ]
    }
]
```
After the success of the two rules stored in redis at the same time generate the parse file car.go


----------


```go
package models

import "errors"

func (a *analyServer) Car(data map[string]interface{}) (string, error) {
	age, ok := data["age"].(float64)
	if !ok {
		err = "age does not exist or type is not float64"
	}
	color, ok := data["color"]
	if !ok {
		err = "color does not exist"
	}
	sex, ok := data["sex"].(string)
	if !ok {
		err = "sex does not exist or type is not string"
	}
	if age > 25 && age <= 35 && contains([]interface{}{"red", "black"}, color) && sex == "male" {
		return "insurance+10%", nil
	}
	if age > 18 && age <= 25 && color == "red" && sex == "male" {
		return "insurance+20%", nil
	}
	return "", errors.New(err)
}
```
Above to complete the rules into executable code for analysis, send business data to the resolution interface


----------


post:  127.0.0.1:8080/v1/rule/analy/
data:  
```json
{
  "name": "car",
  "data": {
     "age":20,
     "sex":"male",
     "color":"red"
  }
}
```
Finally return the expression


----------


```json
[
    {
        "type": "identity",
        "value": "insurance"
    },
    {
        "type": "operator",
        "value": "+"
    },
    {
        "type": "numeric",
        "value": "20"
    },
    {
        "type": "operator",
        "value": "%"
    }
]
```