package file

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/tus/tusd/v2/pkg/filelocker"
	"github.com/tus/tusd/v2/pkg/filestore"
	tusd "github.com/tus/tusd/v2/pkg/handler"
)

func UploadHandler() http.Handler {
	var upload_dir = "upload"

	if err := os.MkdirAll(upload_dir, 0755); err != nil {
		panic(err)
	}

	var composer = tusd.NewStoreComposer()
	var store = filestore.New(upload_dir)
	var locker = filelocker.New(upload_dir)
	store.UseIn(composer)
	locker.UseIn(composer)
	handler, err := tusd.NewHandler(tusd.Config{
		BasePath:              "/files/",
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
	})
	if err != nil {
		panic(err)
	}

	go func() {
		for event := range handler.CompleteUploads {
			info_file := filepath.Join(upload_dir, event.Upload.ID+".info")
			if err := os.Remove(info_file); err != nil {
				log.Printf("Removing info file failed: %v\n", err)
			}
			tmp_file := filepath.Join(upload_dir, event.Upload.ID)
			filename, ok := event.Upload.MetaData["filename"]
			if !ok {
				log.Printf("No filename found in Metadata for %q\n", event.Upload.ID)
				continue
			}
			if err != nil {
				log.Printf("Decoding filename failed: %v\n", err)
				continue
			}
			var save_path = filepath.Join(upload_dir, "files")
			if err := os.MkdirAll(save_path, os.ModePerm); err != nil {
				log.Printf("Creating directory failed: %v\n", err)
			}
			dst_path := filepath.Join(save_path, filename)
			if err := os.Rename(tmp_file, dst_path); err != nil {
				log.Printf("Moving file failed: %v\n", err)
				continue
			}
			log.Printf("[FILE] Upload finished, file saved as: %s, upload id: %s\n", filename, event.Upload.ID)
		}
	}()

	return handler
}
