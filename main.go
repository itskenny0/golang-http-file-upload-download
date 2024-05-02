package main

import (
        "crypto/rand"
        "fmt"
        "html/template"
        "io/ioutil"
        "log"
        "net/http"
        "os"
        "path/filepath"
)

const maxUploadSize = 20 * 1024 * 1024 // 20 MB

var uploadPath = os.TempDir()

func main() {
        http.HandleFunc("/upload", uploadFileHandler())

        log.Print("Server started on 0.0.0.0:8979, use /upload for uploading files")
        log.Fatal(http.ListenAndServe(":8979", nil))
}

func uploadFileHandler() http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
                if r.Method == "GET" {
                        t, _ := template.ParseFiles("upload.gtpl")
                        t.Execute(w, nil)
                        return
                }

                if err := r.ParseMultipartForm(maxUploadSize); err != nil {
                        renderError(w, "CANT_PARSE_FORM", http.StatusInternalServerError)
                        return
                }

                file, fileHeader, err := r.FormFile("uploadFile")
                if err != nil {
                        renderError(w, "INVALID_FILE", http.StatusBadRequest)
                        return
                }
                defer file.Close()

                fileSize := fileHeader.Size
                if fileSize > maxUploadSize {
                        renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
                        return
                }

                fileBytes, err := ioutil.ReadAll(file)
                if err != nil {
                        renderError(w, "INVALID_FILE", http.StatusBadRequest)
                        return
                }

                fileName := randToken(12) + filepath.Ext(fileHeader.Filename)
                newPath := filepath.Join(uploadPath, "httpup-"+fileName)

                newFile, err := os.Create(newPath)
                if err != nil {
                        renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
                        return
                }
                defer newFile.Close()

                if _, err := newFile.Write(fileBytes); err != nil {
                        renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
                        return
                }

                w.Header().Set("Content-Type", "application/json")
                w.Write([]byte(fmt.Sprintf(`{"status": "SUCCESS", "path": "%v"}`, newPath)))
        }
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
        w.WriteHeader(statusCode)
        w.Write([]byte(fmt.Sprintf(`{"error": "%v"}`, message)))
}

func randToken(len int) string {
        b := make([]byte, len)
        rand.Read(b)
        return fmt.Sprintf("%x", b)
}
