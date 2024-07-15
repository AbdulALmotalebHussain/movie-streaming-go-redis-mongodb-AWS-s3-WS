package storage

import (
    "bytes"
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "os"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)

var (
    s3Session *s3.S3
    bucketName string
)

func init() {
    bucketName = os.Getenv("S3_BUCKET_NAME")
    awsRegion := os.Getenv("AWS_REGION")
    awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
    awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

    sess, err := session.NewSession(&aws.Config{
        Region:      aws.String(awsRegion),
        Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
    })

    if err != nil {
        panic(err)
    }

    s3Session = s3.New(sess)
}

func UploadToS3(fileName string, file multipart.File) (string, error) {
    buffer := make([]byte, fileSize(file))
    file.Read(buffer)

    _, err := s3Session.PutObject(&s3.PutObjectInput{
        Bucket:        aws.String(bucketName),
        Key:           aws.String(fileName),
        Body:          bytes.NewReader(buffer),
        ContentLength: aws.Int64(int64(len(buffer))),
        ContentType:   aws.String(http.DetectContentType(buffer)),
    })

    if err != nil {
        return "", fmt.Errorf("failed to upload file to S3: %v", err)
    }

    return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, fileName), nil
}

func fileSize(file multipart.File) int {
    size, _ := file.Seek(0, io.SeekEnd)
    file.Seek(0, 0)
    return int(size)
}

