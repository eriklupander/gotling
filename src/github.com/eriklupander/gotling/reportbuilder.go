package main

import (
    "github.com/tobyhede/go-underscore"
    "fmt"
)

func SumZeroes(numbers []int) (int) {
    var sum int

    fn := func(v, i int) {
        sum += v
    }
    un.EachInt(fn, numbers)
    fmt.Printf("%#v\n", sum) //15
    return sum
}

func BuildReport() {

}