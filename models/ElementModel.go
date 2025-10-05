package models

import "cloud.google.com/go/datastore"

type Choice struct {
	Text    string `json:"text,omitempty"`
	Correct bool   `json:"correct,omitempty"`
}

// Element is a question or other content with a unique
// ID in case the SortKey changes (important for grading,
// comparing with UserAnswer, etc.)
type Element struct {
	KeyID         int64    `json:"id"` //gorm:"primary_key,autoIncrement"
	Text          string   `json:"text,omitempty"`
	Type          string   `json:"type,omitempty"` //default 'single'
	ImageLocation string   `json:"image_location,omitempty"`
	ImageCaption  string   `json:"image_caption,omitempty"`
	ImageCredit   string   `json:"image_credit,omitempty"`
	VideoLocation string   `json:"video_location,omitempty"`
	VideoCaption  string   `json:"video_caption,omitempty"`
	VideoCredit   string   `json:"video_credit,omitempty"`
	Choices       []Choice `json:"choices,omitempty"`
	TextRegex     string   `json:"text_regex,omitempty"`
	EssayRegex    string   `json:"essay_regex,omitempty"`
	ProjectID     int64    `json:"project_id,omitempty"`
	OwnerID       int64    `json:"owner_id,omitempty"`
}

type ElementRepository interface {
	CreateElement(Element *Element) (*datastore.Key, error)
	GetAllElements() ([]*Element, error)
	DeleteElement(id int64) error
	GetElementByID(id int64) (*Element, error)
	UpdateElement(id int64, Element *Element) (*datastore.Key, error)
}
