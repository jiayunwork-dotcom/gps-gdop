package main

import (
	"embed"
	"flag"
	"log"
	"net/http"

	"gps-gdop/internal/web"
)

//go:embed web
var webFS embed.FS

//go:embed example/four-good.json example/four-poor.json
var exampleFiles embed.FS

func main() {
	httpAddr := flag.String("http", ":8080", "serve the web console on this address (e.g. :8080)")
	flag.Parse()

	examples := map[string][]byte{}
	for _, name := range []string{"four-good", "four-poor"} {
		payload, err := exampleFiles.ReadFile("example/" + name + ".json")
		if err != nil {
			log.Fatalf("load example %s: %v", name, err)
		}
		examples[name] = payload
	}

	handler := web.NewServer(web.Assets{
		WebFS:    webFS,
		Examples: examples,
	})
	log.Printf("gps-gdop web console on http://localhost%s", *httpAddr)
	if err := http.ListenAndServe(*httpAddr, handler); err != nil {
		log.Fatal(err)
	}
}
