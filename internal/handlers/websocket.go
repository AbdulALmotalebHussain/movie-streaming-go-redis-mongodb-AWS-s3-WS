package handlers

import (
	"log"
	"net/http"
	"video-streaming/internal/cache"
	"video-streaming/internal/models"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	videoID := r.URL.Query().Get("videoID")
	if videoID == "" {
		log.Printf("videoID is required")
		conn.WriteMessage(websocket.TextMessage, []byte("videoID is required"))
		return
	}

	video, err := models.GetVideoByID(videoID)
	if err != nil {
		log.Printf("Error fetching video by ID: %v", err)
		conn.WriteMessage(websocket.TextMessage, []byte("Error fetching video"))
		return
	}

	// Send the video URL to the client
	conn.WriteMessage(websocket.TextMessage, []byte(video.URL))

	// Retrieve cached video segment if exists
	cachedSegment, err := cache.GetVideoSegment(videoID)
	if err == nil {
		conn.WriteMessage(websocket.BinaryMessage, cachedSegment)
	} else {
		log.Printf("No cached segment found for videoID: %v", videoID)
	}
}
