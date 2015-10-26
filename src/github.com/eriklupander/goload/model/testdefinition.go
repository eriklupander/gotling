package model

const FIRST = "first"
const LAST = "last"
const RANDOM = "random"

type TestDef struct {
	Iterations int `yaml:"iterations"`
	Users int `yaml:"users"`
	Rampup int `yaml:"rampup"`
	Actions []map[string]interface{} `yaml:"actions"`
}