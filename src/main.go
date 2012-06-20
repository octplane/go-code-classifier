package main

import (
	"cclassifier"
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
)

func main() {

	var enabledServer = flag.Bool("server", false, "Enable the webserver for request")
	var serverPort = flag.Int("server-port", 9014, "Port for the webserver")
	var path = flag.String("path", "", "File or Path")

	flag.Parse()

	if *path != "" {
		log.Printf("Will read '%s'.\n", *path)
		my_scanner := scanner.InitFromFile(".", "plop")
		my_scanner.Scan(*path)
		my_scanner.Snapshot()
	} else if *enabledServer {
		log.Printf("Starting server at port %d\n", *serverPort)
		http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		})

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *serverPort), nil))

	} else {
		log.Printf("Please provide valid option: path or webserver \n")
		flag.PrintDefaults()
	}
}
