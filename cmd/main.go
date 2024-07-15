package main

import (
    "log"
    "os"
    "video-streaming/internal/database"
    "video-streaming/internal/server"

    "github.com/joho/godotenv"
)



func main() {
    // تحميل ملف البيئة
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    log.Printf("AWS_REGION: %s", os.Getenv("AWS_REGION"))
    log.Printf("S3_BUCKET_NAME: %s", os.Getenv("S3_BUCKET_NAME"))
    log.Printf("AWS_ACCESS_KEY_ID: %s", os.Getenv("AWS_ACCESS_KEY_ID"))
    log.Printf("AWS_SECRET_ACCESS_KEY: %s", os.Getenv("AWS_SECRET_ACCESS_KEY"))

    // تهيئة MongoDB
    database.InitMongoDB()

    // بدء تشغيل الخادم
    server.StartServer()
}

