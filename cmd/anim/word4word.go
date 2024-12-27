//go:build ignore 

package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
)

func main() {
    // open file
    f, err := os.Open("Faust.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    // read the file word by word using scanner
    scanner := bufio.NewScanner(f)
    scanner.Split(bufio.ScanWords)

    for scanner.Scan() {
        fmt.Println(scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}

