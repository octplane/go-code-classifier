package main

import (
  "fmt"
  "os"
  "path/filepath"
  "github.com/jbrukh/bayesian"
  "stringslice/betterSlice"
  "code_classifier"
)

const (
    Ruby bayesian.Class = "ruby"
    Go bayesian.Class = "go"
)

func main() {
  
  allowed := betterSlice.StringSlice{"go", "ruby", "python", "perl", "obj-c", "auto"}

  if len(os.Args) == 3 {
      lang := os.Args[1]
      path := os.Args[2]
      fmt.Printf("Will read %s as language: %s.\n", path, lang)
      if allowed.Pos(lang) > -1 {
        my_scanner := scanner.InitFromFile("plop.data", Go, Ruby)

        my_scanner.Scan(path, Go)
        my_scanner.Snapshot()
      } else {
        fmt.Printf("%s is not an allowed language among %q.\n", lang, allowed)
      }
    } else {
      fmt.Printf("Missing arguments: %s lang path_or_file\n", filepath.Base(os.Args[0]))

    }
}

