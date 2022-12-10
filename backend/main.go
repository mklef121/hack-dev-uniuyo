package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Hello Welcome to golang")

	FileUploadHandler := http.HandlerFunc(FileUpload)
	http.Handle("/upload", FileUploadHandler)
	http.ListenAndServe(":2023", nil)

}

func FileUpload(w http.ResponseWriter, r *http.Request) {
	// Bypass CORS problems
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if r.Method == "OPTIONS" {
		// if it's options request just bypass
		return
	}

	stream, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	fileName := r.Header.Get("file-name")

	_, err = os.Stat("./" + fileName)

	if err != nil {
		// when file is a new file and has never been created before
		if errors.Is(err, os.ErrNotExist) {
			tmpfile, err := os.Create("./" + fileName)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			defer tmpfile.Close()

			tmpfile.Write(stream)
			w.Write([]byte("done writing"))
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// if file exists

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	// Append to file, if it exists already
	f.Write(stream)

	w.Write([]byte("done writing existing file"))

}
