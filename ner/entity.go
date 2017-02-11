package ner

// Entity represents an extracted entity
type Entity struct {
	Tag 	Tag
	Value      string
	Score     float64
}

// Tag represents an type of entity
type Tag struct {
	ID int
	Name string
}