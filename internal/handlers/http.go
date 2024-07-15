package handlers

import (
	"html/template"
	"log"
	"net/http"
	"video-streaming/internal/models"
)

var templates = template.Must(template.ParseGlob("web/templates/*.html"))

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	videos, err := models.GetAllVideos()
	if err != nil {
		log.Printf("Error fetching videos: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "home.html", videos)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func UploadPageHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "upload.html", nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func StreamHandler(w http.ResponseWriter, r *http.Request) {
	videoID := r.URL.Query().Get("videoID")
	if videoID == "" {
		log.Printf("videoID is required")
		http.Error(w, "videoID is required", http.StatusBadRequest)
		return
	}

	video, err := models.GetVideoByID(videoID)
	if err != nil {
		log.Printf("Error fetching video by ID: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "stream.html", video)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
