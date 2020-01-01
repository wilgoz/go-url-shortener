package shortener

type Redirect struct {
	CreatedAt int64  `json:"created_at" bson:"created_at"`
	Original  string `json:"original" bson:"original" validate:"empty=false & format=url"`
	Shortened string `json:"shortened" bson:"shortened"`
}
