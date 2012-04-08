package main

import (
  "fmt"
  "os"
  "path/filepath"
  "code_classifier"
)

func main() {

  if len(os.Args) == 2 {
      path := os.Args[1]
      fmt.Printf("Will read %s.\n", path)
      my_scanner := scanner.InitFromFile("plop.data")
      my_scanner.Scan(path)
      my_scanner.Snapshot()
    } else {
      fmt.Printf("Missing arguments: %s path_or_file\n", filepath.Base(os.Args[0]))

    }
}

