package handlers

import (
    "bytes"
    "fmt"
    "log"
    "mime/multipart"
    "net/http"
    "os"
    "path/filepath"
    "video-streaming/internal/models"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)

// S3Uploader handles the S3 upload process
type S3Uploader struct {
    s3Client *s3.S3
    bucket   string
}

func NewS3Uploader(bucket, region string) *S3Uploader {
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(region),
        Credentials: credentials.NewStaticCredentials(
            os.Getenv("AWS_ACCESS_KEY_ID"),
            os.Getenv("AWS_SECRET_ACCESS_KEY"),
            "",
        ),
    })
    if err != nil {
        log.Fatalf("failed to create session: %v", err)
    }

    return &S3Uploader{
        s3Client: s3.New(sess),
        bucket:   bucket,
    }
}

func (u *S3Uploader) UploadFileToS3(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
    defer file.Close()

    size := fileHeader.Size
    buffer := make([]byte, size)
    file.Read(buffer)
    fileBytes := bytes.NewReader(buffer)
    fileName := filepath.Base(fileHeader.Filename)
    path := "videos/" + fileName

    params := &s3.PutObjectInput{
        Bucket:   aws.String(u.bucket),
        Key:      aws.String(path),
        Body:     fileBytes,
        ACL:      aws.String("public-read"),
    }

    _, err := u.s3Client.PutObject(params)
    if err != nil {
        log.Printf("failed to upload data to %s/%s, %s", u.bucket, path, err.Error())
        return "", fmt.Errorf("failed to upload data to %s/%s, %s", u.bucket, path, err.Error())
    }

    return path, nil
}

var s3Uploader = NewS3Uploader(os.Getenv("S3_BUCKET_NAME"), os.Getenv("AWS_REGION"))

func UploadHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseMultipartForm(10 << 20) // maximum upload size 10 MB

    file, handler, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Error retrieving the file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    filePath, err := s3Uploader.UploadFileToS3(file, handler)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    videoID := handler.Filename
    video := models.Video{
        ID:   videoID,
        Name: r.FormValue("name"),
        URL:  filePath,
    }

    err = models.SaveVideo(video)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    err = models.UpdateVideoURL(videoID, filePath)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Write([]byte("Successfully Uploaded Video"))
}
