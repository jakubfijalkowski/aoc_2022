package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "strings"
    "strconv"
    "sort"
)

func main() {
    rawContent, err := ioutil.ReadFile("data.txt")
    if err != nil {
        log.Fatal(err)
    }
    content := string(rawContent)
    groups := strings.Split(content, "\n\n")
    summedGroups := make([]int, 0, len(groups))

    for _, group := range groups {
        elements := strings.Split(strings.TrimSpace(group), "\n")
        sum := int(0)
        for _, e := range elements {
            p, err := strconv.Atoi(e)
            if err != nil {
                log.Fatal(err)
            }
            sum = sum + p
        }
        summedGroups = append(summedGroups, sum)
    }

    sort.Ints(summedGroups)

    top3 := 0
    for i := 1; i <= 3; i++ {
        top3 = top3 + summedGroups[len(summedGroups) - i]
    }

    fmt.Println("Part 1: ", summedGroups[len(summedGroups) - 1])
    fmt.Println("Part 2: ", top3)
}
