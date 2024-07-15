package models

import (
	"context"
	"video-streaming/internal/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Video struct {
	ID   string `bson:"_id,omitempty"`
	Name string `bson:"name"`
	URL  string `bson:"url"`
}

func SaveVideo(video Video) error {
	client := database.GetClient()
	collection := client.Database("videostreaming").Collection("videos")
	_, err := collection.InsertOne(context.TODO(), video)
	return err
}

func GetVideoByID(id string) (*Video, error) {
	client := database.GetClient()
	collection := client.Database("videostreaming").Collection("videos")
	var video Video
	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&video)
	return &video, err
}

func GetAllVideos() ([]Video, error) {
	client := database.GetClient()
	collection := client.Database("videostreaming").Collection("videos")
	cur, err := collection.Find(context.TODO(), bson.D{}, options.Find())
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var videos []Video
	for cur.Next(context.Background()) {
		var video Video
		err := cur.Decode(&video)
		if err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func UpdateVideoURL(id, url string) error {
    client := database.GetClient()
    collection := client.Database("videostreaming").Collection("videos")

    filter := bson.M{"_id": id}
    update := bson.M{"$set": bson.M{"url": url}}

    _, err := collection.UpdateOne(context.TODO(), filter, update)
    return err
}

