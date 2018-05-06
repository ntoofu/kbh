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
}

type EndpointDef struct {
	Name string
	Type string
	Url string
	ApiVersion string `yaml:"api_version"`
}

type BoardDef struct {
	Name string
	Alias []string
	Endpoint string
	Mapping struct {
		State map[string]StateCondDef
	}
}

type StateCondDef struct {
	Order uint
	LabelName NullableString `yaml:"label_name"`
	Asignee NullableString
	IsClosed NullableBool `yaml:"is_closed"`
	MaxDaysWOUpdate NullableUint `yaml:"max_days_without_update"`
}

func (x *StateCondDef) IsMatched(issue Issue) bool {
	if x.LabelName.Valid {
		contained := false
		for _, l := range issue.Label {
			if x.LabelName.Value == l {
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
		if x.MaxDaysWOUpdate.Valid {
			maxHours := time.Duration(x.MaxDaysWOUpdate.Value) * 24 * time.Hour
			if time.Now().After(issue.UpdateTime().Add(maxHours)) {
				return false
			}
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
