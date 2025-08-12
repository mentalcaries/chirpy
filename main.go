package main

import (
	"log"
	"net/http"
)

func main(){

  const filePathRoot = "."
  const port = "8080"
  mux := http.NewServeMux()

  server := http.Server{
    Handler : mux,
    Addr : ":" + port,
  }


  mux.Handle("/app/", http.StripPrefix("/app",  http.FileServer(http.Dir(filePathRoot))))
  mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(200)
    w.Write([]byte("OK"))
  })

  log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(server.ListenAndServe())
}