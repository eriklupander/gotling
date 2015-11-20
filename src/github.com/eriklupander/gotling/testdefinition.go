package main

const FIRST = "first"
const LAST = "last"
const RANDOM = "random"


type TestDef struct {
	Iterations int `yaml:"iterations"`
	Users int `yaml:"users"`
	Rampup int `yaml:"rampup"`
	Feeder Feeder `yaml:"feeder"`
	Actions []map[string]interface{} `yaml:"actions"`
}

type Feeder struct {
	Type string `yaml:"type"`
	Filename string `yaml:"filename`
}