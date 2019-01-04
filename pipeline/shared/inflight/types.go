package inflight

import "encoding/json"

type inflightRef struct {
	Bucket string
	Path   string
	Object string
}

func (i inflightRef) mustConvertToBytes() []byte {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err.Error())
	}
	return b
}

func (i *inflightRef) mustConvertFromBytes(b []byte) {
	err := json.Unmarshal(b, i)
	if err != nil {
		panic(err.Error())
	}
}
