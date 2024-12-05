package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
)

const port string = ":8080"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
	})

	http.HandleFunc("/gimmeEXE", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Just POST method dude!!", http.StatusMethodNotAllowed)
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Are you forget some file BURH~~", http.StatusBadRequest)
			return
		}
		defer file.Close()

		if filepath.Ext(handler.Filename) != ".exe" {
			http.Error(w, "Do you see name of URL on top -_-", http.StatusBadRequest)
			return
		}

		cmd := exec.Command("cmd.exe", "/C") // "/bin/sh" for Linux
		stdin, _ := cmd.StdinPipe()

		go func() {
			defer stdin.Close()
			io.Copy(stdin, file)
		}()

		output, err := cmd.CombinedOutput()
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		fmt.Println(string(output))
		w.Header().Set("Content-Type", "text/plain")
		w.Write(output)
	})

	fmt.Printf("Server started at http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
