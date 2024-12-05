package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os/exec"
	"os"
	"path/filepath"
)

const port string = ":8080"
var output string

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
		t.Execute(w, map[string]string{"Output": output,})
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

		fileContent, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}

		tempFile, err := os.CreateTemp("./", "temp_*.exe") 
		if err != nil {
			http.Error(w, "Error creating temporary file", http.StatusInternalServerError)
			return
		}
		defer os.Remove(tempFile.Name())

		_, err = tempFile.Write(fileContent)
		if err != nil {
			http.Error(w, "Error writing to temporary file", http.StatusInternalServerError)
			return
		}

		tempFile.Close()

		var outBuffer, errBuffer bytes.Buffer
		cmd := exec.Command("wine", tempFile.Name())
		cmd.Stdout = &outBuffer
		cmd.Stderr = &errBuffer

		err = cmd.Run()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error executing .exe: %v\n%s", err, errBuffer.String()), http.StatusInternalServerError)
			return
		}

		output = fmt.Sprintf("%s", outBuffer.String())

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(output))
	})

	fmt.Printf("Server started at http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
