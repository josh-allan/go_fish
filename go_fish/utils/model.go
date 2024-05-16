package shared

type Entry struct {
	Title     string `bson:"title"`
	Published string `bson:"published"`
}

