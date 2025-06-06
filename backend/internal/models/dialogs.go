package models

type QueuedDialog struct {
	DialogBase `bson:",inline"`
}

type ActiveDialog struct {
	DialogBase `bson:",inline"`
}

type ArchivedDialog struct {
	DialogBase `bson:",inline"`
}
