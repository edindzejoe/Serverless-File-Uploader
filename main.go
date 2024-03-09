package p

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "os"

    "cloud.google.com/go/storage"
)

func UploadFileToGCS(w http.ResponseWriter, r *http.Request) {
    // Parse input data from the request
    file, _, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Error retrieving the file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Setup GCS client
    ctx := context.Background()
    client, err := storage.NewClient(ctx)
    if err != nil {
        http.Error(w, "Error creating storage client", http.StatusInternalServerError)
        return
    }
    defer client.Close()

    // Name of the GCS bucket
    bucketName := os.Getenv("BUCKET_NAME")

    // Specify the name for the stored file
    objectName := "your-object-name" // TODO: Generate or extract object name

    // Upload file to GCS
    wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
    if _, err = io.Copy(wc, file); err != nil {
        http.Error(w, "Error uploading file to GCS", http.StatusInternalServerError)
        return
    }
    if err := wc.Close(); err != nil {
        http.Error(w, "Error closing GCS writer", http.StatusInternalServerError)
        return
    }

    // Respond with the URL of the uploaded file
    fmt.Fprintf(w, "File uploaded successfully: gs://%s/%s", bucketName, objectName)
}

