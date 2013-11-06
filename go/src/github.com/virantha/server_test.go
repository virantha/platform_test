package main

import (
    "testing"
    "fmt"
    "sort"
)


func Test_divide(t *testing.T) {
    //Test if dividing the message works
    var keys []int
    for key,_ := range IPs {
        keys = append(keys, key)
    }
    sort.Sort(sort.Reverse(sort.IntSlice(keys)))

    fixtures := map[int] map[int]int {
        11: map[int]int{25:0, 10:1, 5:0, 1:1},
        100: map[int]int{25:4, 10:0, 5:0, 1:0},
        107: map[int]int{25:4, 10:0, 5:1, 1:2},
        4: map[int]int{25:0, 10:0, 5:0, 1:4},
    }
    for input,expected := range fixtures {
        result := divide(input, keys)
        failed := false
        for k,v := range expected {
            if result[k] != v { failed = true }
        }
        if failed {
            t.Error(fmt.Sprintf("Dividing %d expected %v, got %v", input, expected,result))
        }
    }

