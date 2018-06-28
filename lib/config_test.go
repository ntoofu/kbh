package lib

import (
	"github.com/go-test/deep"
	"testing"
)

func TestParseConfig(t *testing.T) {
	expectedConf := GlobalConfig{
		"my.name",
		[]EndpointDef{
			EndpointDef{"repo1", "dummy", "https://foo.bar/baz", "1", "AAA"},
			EndpointDef{"repo2", "dummy", "https://hoge.hoge/fuga/fuga", "2", "BBB"},
		},
		[]BoardDef{
			BoardDef{
				"team1/proj1",
				[]string{"proj1"},
				"repo1",
				struct{ State map[string]StateCondDef }{
					map[string]StateCondDef{
						"todo": StateCondDef{
							3,
							nil,
							NullableString{"my.name", true},
							NullableBool{false, true},
							NullableUint{0, false},
						},
						"doing": StateCondDef{
							1,
							[]string{"doing-label"},
							NullableString{"my.name", true},
							NullableBool{false, true},
							NullableUint{0, false},
						},
						"waiting": StateCondDef{
							2,
							[]string{"waiting-label"},
							NullableString{"my.name", true},
							NullableBool{false, true},
							NullableUint{0, false},
						},
						"closed": StateCondDef{
							4,
							nil,
							NullableString{"my.name", true},
							NullableBool{true, true},
							NullableUint{14, true},
						},
					},
				},
			},
			BoardDef{
				"team1/proj2",
				[]string{"proj2", "proj02"},
				"repo1",
				struct{ State map[string]StateCondDef }{
					map[string]StateCondDef{
						"todo": StateCondDef{
							3,
							[]string{"my-name"},
							NullableString{"my.name", true},
							NullableBool{false, true},
							NullableUint{0, false},
						},
						"doing": StateCondDef{
							1,
							[]string{"doing-label", "my-name"},
							NullableString{"my.name", true},
							NullableBool{false, true},
							NullableUint{0, false},
						},
						"waiting": StateCondDef{
							2,
							[]string{"waiting-label", "my-name"},
							NullableString{"my.name", true},
							NullableBool{false, true},
							NullableUint{0, false},
						},
						"closed": StateCondDef{
							4,
							[]string{"my-name"},
							NullableString{"my.name", true},
							NullableBool{true, true},
							NullableUint{14, true},
						},
					},
				},
			},
			BoardDef{
				"individual",
				[]string{"mytask"},
				"repo2",
				struct{ State map[string]StateCondDef }{
					map[string]StateCondDef{
						"todo": StateCondDef{
							2,
							nil,
							NullableString{"", false},
							NullableBool{false, true},
							NullableUint{0, false},
						},
						"now": StateCondDef{
							1,
							[]string{"doing-label"},
							NullableString{"", false},
							NullableBool{false, true},
							NullableUint{3, true},
						},
						"fin": StateCondDef{
							3,
							nil,
							NullableString{"", false},
							NullableBool{true, true},
							NullableUint{7, true},
						},
					},
				},
			},
		},
		CommandOptions{
			ShowOptions{
				"	",
				" ",
				[]string{"uri", "title", "state"},
			},
		},
	}
	parsedConf, err := ParseConfig("../config_test.yml")
	if err != nil {
		t.Errorf("An error has occured during parsing config: %v", err)
	}
	if diff := deep.Equal(parsedConf, expectedConf); diff != nil {
		t.Errorf("config parsed from config_test.yml differs from expected one\n%v", diff)
	}
}
