package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var path = os.Getenv("PWD")

func main() {
	http.Handle("/", http.FileServer(http.Dir(path)))
	log.Fatal(http.ListenAndServe(":8000", getHandlerFunc()))
}

func getHandlerFunc() http.HandlerFunc {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if "/" == r.URL.Path {
			// f, err := os.Open("/home/d/Download/upload.html")
			// if nil != err {
			// 	w.WriteHeader(http.StatusNotFound)
			// 	return
			// }
			// io.Copy(w, f)
			http.Redirect(w, r, "http://"+r.Host+"/upload.html", http.StatusTemporaryRedirect)
			return
		}
		http.DefaultServeMux.ServeHTTP(w, r)
	})
	mux.HandleFunc("/upload.php", func(w http.ResponseWriter, r *http.Request) {
		f, fh, err := r.FormFile("fileToUpload")
		if nil != err {
			panic(err)
		}
		wf, err := os.OpenFile(path+"/"+fh.Filename, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if nil != err {
			panic(err)
		}
		io.Copy(wf, f)
	})

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		mux.ServeHTTP(w, r)
	}
}
