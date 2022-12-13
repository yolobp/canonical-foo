package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Println("Listening on port 8080...")
	if err := http.ListenAndServe(":8080", http.HandlerFunc(handleFunc)); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	log.Println("Server exited.")
}

func handleFunc(w http.ResponseWriter, req *http.Request) {
	var b bytes.Buffer
	fmt.Fprintln(&b, "Your request:")
	fmt.Fprintf(&b, "Method: %s\n", req.Method)
	fmt.Fprintf(&b, "Host: %s\n", req.Host)
	fmt.Fprintf(&b, "URL: %s\n", req.URL.String())
	fmt.Fprintf(&b, "Headers:\n%v", req.Header)
	w.Write(b.Bytes())
}
