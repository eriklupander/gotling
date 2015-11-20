package main
import (
    "log"

)

// TODO refactor this so it runs before parsing the actions

func ValidateTestDefinition(t *TestDef) (bool) {
    var valid bool = true
    if t.Iterations == 0 {
        log.Println("Iterations not set, must be > 0")
        valid = false
    }
    if t.Rampup < 0 {
        log.Println("Rampup not defined. must be > -1")
        valid = false
    }
    if t.Users == 0 {
        log.Println("Users must be > 0")
        valid = false
    }
    return valid
}