package main_test

import (
    "github.com/eriklupander/gotling"
    "testing"
)

func init() {
}

func TestEach(t *testing.T) {
    var numbers = []int{1,2,3,4,5}
    result := main.SumZeroes(numbers)
    if result != 15 {
        t.Errorf("Expected sum to be 15")
    }
}