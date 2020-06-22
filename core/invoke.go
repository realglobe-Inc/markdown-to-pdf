package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

const appName = "markdown-to-pdf"

// リクエストボディを書き出すファイル名を返す
func getFileName(r *http.Request) string {
	name := path.Base(r.URL.Path)
	if name == "/" || name == "." {
		name = "file"
	}

	extensions, _ := mime.ExtensionsByType(r.Header.Get("Content-Type"))
	if extensions == nil {
		return name
	}

	hasExtension := false
	currentExtension := path.Ext(name)
	for _, extension := range extensions {
		if currentExtension == extension {
			hasExtension = true
		}
	}
	if hasExtension {
		return name
	}

	return name + extensions[0]
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: received a request", appName)

	// 作業ディレクトリを用意
	workDir, err := ioutil.TempDir("", "")
	if err != nil {
		w.WriteHeader(500)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	defer func() {
		if err := os.RemoveAll(workDir); err != nil {
			log.Printf("%s: removing the working directory failed: %v", appName, err)
		}
	}()

	// リクエストボディをファイルに書き出す
	filePath := filepath.Join(workDir, getFileName(r))
	dest, err := os.Create(filePath)
	if err != nil {
		log.Printf("%s: creating a temporary file failed: %v", appName, err)
		w.WriteHeader(500)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	if _, err = io.Copy(dest, r.Body); err != nil {
		log.Printf("%s: saving the file failed: %v", appName, err)
		w.WriteHeader(500)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	// 変換スクリプトを実行する
	cmd := exec.Command("./script.sh", filePath)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err = cmd.Run(); err != nil {
		log.Printf("%s: conversion failed: %v\n%s\n%s", appName, err, stdout.String(), stderr.String())
		w.WriteHeader(500)
		_, _ = w.Write(stdout.Bytes())
		_, _ = w.Write(stderr.Bytes())
		return
	}
	if stderr.Len() > 0 {
		log.Printf("%s: conversion error: %s", appName, stderr.String())
	}

	// レスポンスを作る
	resultFile := strings.Trim(stdout.String(), " \n")
	reader, err := os.Open(resultFile)
	if err != nil {
		log.Printf("%s: open the result file failed: %v", appName, err)
		w.WriteHeader(500)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	// ファイルタイプを特定
	if contentType := mime.TypeByExtension(path.Ext(resultFile)); contentType != "" {
		w.Header().Set("Content-Type", contentType)
	} else {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	// ファイルをレスポンスボディに書き出す
	if _, err = io.Copy(w, reader); err != nil {
		log.Printf("%s: open the result file failed: %v", appName, err)
		w.WriteHeader(500)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

}

func main() {
	log.Printf("%s: starting server...", appName)

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("%s: listening on %s", appName, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
