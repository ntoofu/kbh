package lib

// define format of config.yml

import (
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type GlobalConfig struct {
	User string
	Endpoint []EndpointDef
	Board []BoardDef
	Command CommandOptions
}

type EndpointDef struct {
	Name string
	Type string
	Url string
	ApiVersion string `yaml:"api_version"`
	AuthToken string `yaml:"auth_token"`
}

type BoardDef struct {
	Name string
	Alias []string
	Endpoint string
	Mapping struct {
		State map[string]StateCondDef
	}
}

type CommandOptions struct {
	Show ShowOptions
}

type StateCondDef struct {
	Order uint
	Labels []string
	Asignee NullableString
	IsClosed NullableBool `yaml:"is_closed"`
	MaxDaysWOUpdate NullableUint `yaml:"max_days_without_update"`
}

type ShowOptions struct {
	Delimiter string
	Replacer string
	Field []string
}

func (x *StateCondDef) IsMatched(issue *Issue) bool {
	for _, label := range x.Labels {
		contained := false
		for _, l := range issue.Label {
			if label == l {
				contained = true
				break
			}
		}
		if !contained {
			return false
		}
	}
	if x.Asignee.Valid {
		if x.Asignee.Value != issue.Asignee {
			return false
		}
	}
	if x.IsClosed.Valid {
		if x.IsClosed.Value != issue.IsClosed {
			return false
		}
	}
	if x.MaxDaysWOUpdate.Valid {
		maxHours := time.Duration(x.MaxDaysWOUpdate.Value) * 24 * time.Hour
		if time.Now().After(issue.UpdateTime().Add(maxHours)) {
			return false
		}
	}
	return true
}

func ParseConfig(filename string) (GlobalConfig, error) {
	globalConf := GlobalConfig{}
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return globalConf, fmt.Errorf("Cannot open config file: %v\n%v", filename, err.Error())
	}

	err = yaml.Unmarshal(file, &globalConf)
	if err != nil {
		return globalConf, fmt.Errorf("Cannot parse config file: %v\n%v", filename, err.Error())
	}

	return globalConf, nil
}
