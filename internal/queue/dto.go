package queue

import "encoding/json"

type QueueDto struct {
	Filename string `json:"filename"`
	Path     string `json:"path"`
	ID       int    `json:"id"`
}

func (q *QueueDto) ToJson() ([]byte, error) {
	return json.Marshal(q)
}

func (q *QueueDto) FromJson(b []byte) error {
	return json.Unmarshal(b, q)
}
