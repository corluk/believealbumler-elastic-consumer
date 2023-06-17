package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AudioItem struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Id           int                `json:"id"`
	Url          string             `json:"url"`
	Type         string             `json:"type"`
	Status       string             `json:"status"`
	LastModified time.Time          `json:"lastModified"`
	Title        string             `json:"title"`

	Fetched      bool      `json:"fetched"`
	ArtistName   string    `json:"artistName"`
	PlakSirketi  string    `json:"plakSirketi"`
	UPC          string    `json:"upc"`
	CLine        string    `json:"cline"`
	PLine        string    `json:"pline"`
	YapimYili    int       `json:"yapımyılı"`
	Genre        string    `json:"genre"`
	Parental     bool      `json:"parental"`
	Duration     string    `json:"duration"`
	CreationDate time.Time `json:"creation_date"`
	ReleaseDate  time.Time `json:"release_date"`
	ReleaseType  string    `json:"release_type"`
	PriceTier    string    `json:"priceTier"`
	ImageLink    string    `json:"image_link"`
	Tracks       []Track   `json:"tracks"`
	Source       string    `json:"source" default:"believe_album"`
}

type Track struct {
	Order     int      `json:"order"`
	Title     string   `json:"title"`
	Version   string   `json:"version"`
	Artists   []string `json:"artists"`
	Authors   []string `json:"authors"`
	Composers []string `json:"composers"`
	Genres    []string `json:"genres"`
	Isrc      string   `json:"isrc"`
	PriceTier string   `json:"pricetier"`
	Volume    int      `json:"volume"`
	TrackNo   int      `json:"trackNo"`
	Duration  string   `json:"duration"`
	UPC       string   `json:"upc"`
}
