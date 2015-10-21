package main

type Action struct {
	Type string `yaml:"type"`
	Properties map[string]string `yaml:"properties"`
}

type SleepAction struct {
	Duration int `yaml:"duration"`
}

type HttpReqAction struct {
	Method string `yaml:"method"`
	Url string `yaml:"url"`
	Accept string `yaml:"accept"`
}

type TestDef struct {
	Iterations int `yaml:"iterations"`
	Users int `yaml:"users"`
	Rampup int `yaml:"rampup"`
	Actions []map[string]interface{} `yaml:"actions"`
}
