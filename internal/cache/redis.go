package cache

import (
    "context"
    "fmt"
    "io"
    "time"

    "github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func CacheVideoSegment(videoID string, file io.Reader) error {
    // افتراضيا، يتم تخزين أول 10 ثوانٍ من الفيديو فقط
    buffer := make([]byte, 10*1024*1024)
    _, err := file.Read(buffer)
    if err != nil {
        return fmt.Errorf("failed to read video segment: %v", err)
    }

    client := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    err = client.Set(ctx, videoID, buffer, 10*time.Minute).Err()
    if err != nil {
        return fmt.Errorf("failed to cache video segment: %v", err)
    }

    return nil
}

func GetVideoSegment(videoID string) ([]byte, error) {
    client := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    segment, err := client.Get(ctx, videoID).Bytes()
    if err != nil {
        return nil, fmt.Errorf("failed to get video segment: %v", err)
    }

    return segment, nil
}
