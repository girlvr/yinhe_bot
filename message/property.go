package message

import (
	jsoniter "github.com/json-iterator/go"
	"os"
)

var property *Property

// Property ...
type Property struct {
	Welcome     string `json:"welcome"`
	GroupName   string `json:"group_name"`
	BotName     string `json:"bot_name"`
	Host        string `json:"host"`
	HookAddress string `json:"hook_address"`
	Token       string `json:"token"`
	Download    string `json:"download"`
	Rule        string `json:"rule"`
}

// LoadProperty ...
func LoadProperty(pathname string) error {
	property = &Property{}
	file, e := os.OpenFile(pathname, os.O_RDONLY, 0755)
	if e != nil {
		return e
	}
	dec := jsoniter.NewDecoder(file)
	e = dec.Decode(property)
	if e != nil {
		return e
	}
	log.Infof("property:%+v", *property)
	return nil
}

// GetProperty ...
func GetProperty() *Property {
	return property
}
