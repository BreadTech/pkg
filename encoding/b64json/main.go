package b64json

import (
	"encoding/base64"
	"encoding/json"
)

var (
	enc = base64.StdEncoding
)

func SetEncoding(x *base64.Encoding) {
	enc = x
}

func Unmarshal(b []byte, obj interface{}) error {
	dst := make([]byte, enc.DecodedLen(len(b)))
	n, err := enc.Decode(dst, b)
	if err != nil {
		return err
	}
	return json.Unmarshal(dst[:n], obj)
}

func Marshal(obj interface{}) ([]byte, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	dst := make([]byte, enc.EncodedLen(len(b)))
	enc.Encode(dst, b)
	return dst, nil
}
