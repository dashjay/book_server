package reply

import "github.com/golang/protobuf/proto"

func NewBaseMessage(t string, content string) []byte {
	var r BaseMessage
	r.Type = t
	r.Data = []byte(content)
	rb, _ := proto.Marshal(&r)
	return rb
}

type JsonReply struct {
	Status  uint8  `json:"status"`
	Content []byte `json:"content" bson:"content"`
}
