package resources

import "fmt"

type Resources struct {
	ID         int    `json:"id" bson:"id"`
	FileName   string `json:"file_name" bson:"file_name"`
	Link       string `json:"link" bson:"link"`
	Size       string `json:"size" bson:"size"`
	UpdateTime string `json:"update_time" bson:"update_time"`
}

func (r *Resources) GetFileName() string {
	return fmt.Sprintf("#%d %s", r.ID, r.FileName)
}
