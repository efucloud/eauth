package embeds

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed all:web
var WebFiles embed.FS

//go:embed templates
var Templates embed.FS

//go:embed locales
var I18nFiles embed.FS

//go:embed face-recognition
var faceRecognition embed.FS

func FaceRecognitionModels() http.FileSystem {
	fsys, err := fs.Sub(faceRecognition, "face-recognition")
	if err != nil {
		log.Fatal(err)
	}
	return http.FS(fsys)
}
