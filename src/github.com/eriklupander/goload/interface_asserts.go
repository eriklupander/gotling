package main
import "github.com/eriklupander/goload/model"

var _ model.Action = (*HttpAction)(nil)
var _ model.Action = (*model.SleepAction)(nil)