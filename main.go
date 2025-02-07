package main

import (
	"embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

//go:embed page/*
var fs embed.FS

var path = os.Getenv("PWD")

func main() {
	http.Handle("/", http.FileServer(http.Dir(path)))
	log.Fatal(http.ListenAndServe(":8000", getHandlerFunc()))
}

func getHandlerFunc() http.HandlerFunc {
	uploadHTML, err := fs.ReadFile("page/upload.html")
	if nil != err { panic(err) }
	jqjs, err := fs.ReadFile("page/jquery-3.1.1.min.js")
	if nil != err { panic(err) }
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if "/upload" == r.URL.Path {
			w.Write(uploadHTML)
			return
		}
		if "/jquery-3.1.1.min.js" == r.URL.Path {
			w.Write(jqjs)
			return
		}
		//if "/" == r.URL.Path {
		//	// f, err := os.Open("/home/d/Download/upload.html")
		//	// if nil != err {
		//	// 	w.WriteHeader(http.StatusNotFound)
		//	// 	return
		//	// }
		//	// io.Copy(w, f)
		//	http.Redirect(w, r, "http://"+r.Host+"/upload.html", http.StatusTemporaryRedirect)
		//	return
		//}
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
