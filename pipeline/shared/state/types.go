package state

import (
	"encoding/json"
	"path/filepath"
)

type InflightRef struct {
	Bucket string
	Path   string
	Object string
}

func (i InflightRef) mustConvertToBytes() []byte {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err.Error())
	}
	return b
}

func (i *InflightRef) mustConvertFromBytes(b []byte) {
	err := json.Unmarshal(b, i)
	if err != nil {
		panic(err.Error())
	}
}

func (i *InflightRef) S3URI() string {
	return filepath.Join(i.Bucket, i.Path, i.Object)
}
