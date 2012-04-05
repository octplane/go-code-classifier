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
        my_scanner, err := scanner.InitFromFile("plop.data", Go, Ruby)
        if err != nil {
          my_scanner = &scanner.Scanner{bayesian.NewClassifier(Ruby, Go), "plop.data"}
        }

        my_scanner.Scan(path, lang)
        my_scanner.Snapshot()
      } else {
        fmt.Printf("%s is not an allowed language among %q.\n", lang, allowed)
      }
    } else {
      fmt.Printf("Missing arguments: %s lang path_or_file\n", filepath.Base(os.Args[0]))

    }
}

