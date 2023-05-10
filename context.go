package intake

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

//var bufferPool = &sync.Pool{
//	New: func() interface{} {
//		return new(bytes.Buffer)
//	},
//}
//
//// AddToContext adds v to the request context.
//func AddToContext(r *http.Request, key string, v interface{}) error {
//	buf := bufferPool.Get().(*bytes.Buffer)
//
//	if err := gob.NewEncoder(buf).Encode(v); err != nil {
//		return err
//	}
//
//	*r = *r.WithContext(context.WithValue(r.Context(), key, buf.Bytes()))
//
//	buf.Reset()
//	bufferPool.Put(buf)
//	return nil
//}
//
//// FromContext adds v to the request context. v must be json decode-able
//func FromContext(r *http.Request, key string, v interface{}) error {
//	data, _ := r.Context().Value(key).([]byte)
//	//buf := bufferPool.Get().(*bytes.Buffer)
//
//	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(v); err != nil {
//		return err
//	}
//	return nil
//}
//
//func encodeGob(employee interface{}) (*bytes.Buffer, error) {
//	buf := new(bytes.Buffer)
//
//	// Create a new gob.Encoder and encode an Employee instance
//	encoder := gob.NewEncoder(buf)
//	err := encoder.Encode(employee)
//	if err != nil {
//		return nil, fmt.Errorf("error encoding data: %v", err)
//	}
//
//	return buf, nil
//}
//
//func decodeGob(buf *bytes.Buffer, v interface{})  error {
//	// Create a new gob.Decoder and decode the data into an Employee instance
//	decoder := gob.NewDecoder(buf)
//	err := decoder.Decode(v)
//	if err != nil {
//		return fmt.Errorf("error decoding data: %v", err)
//	}
//
//	return  nil
//}

func AddToContextV2(r *http.Request, key string, v interface{}) error{
	var buf bytes.Buffer
	if err:=json.NewEncoder(&buf).Encode(v);err != nil {
		return err
	}
	*r = *r.WithContext(context.WithValue(r.Context(), key, &buf))
	return nil
}

func FromContextV2(r *http.Request, key string, v interface{}) error {
	data, ok := r.Context().Value(key).(*bytes.Buffer)
	if !ok {
		return fmt.Errorf("error casting to []byte for key %s", key)
	}
	return json.NewDecoder(data).Decode(v)
}



func AddToContext(r *http.Request, key string, v interface{}) error {
	encoded, err := json.Marshal(v)
	if err != nil {
		return err
	}
	*r = *r.WithContext(context.WithValue(r.Context(), key, encoded))
	return nil
}

func FromContext(r *http.Request, key string, v interface{}) error {
	data, ok := r.Context().Value(key).([]byte)
	if !ok {
		return fmt.Errorf("error casting to []byte for key %s", key)
	}
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}



//
//func encodeJson(v interface{}) (*bytes.Buffer, error) {
//	var buf bytes.Buffer
//
//	// Create a new json.Encoder and encode the data
//	encoder := json.NewEncoder(&buf)
//	err := encoder.Encode(v)
//	if err != nil {
//		return nil, fmt.Errorf("error encoding data: %v", err)
//	}
//
//	return buf, nil
//}
//
//func decodeJson(buf *bytes.Buffer, v interface{}) error {
//	// Create a new json.Decoder and decode the data
//	decoder := json.NewDecoder(buf)
//	err := decoder.Decode(v)
//	if err != nil {
//		return fmt.Errorf("error decoding data: %v", err)
//	}
//
//	return nil
//}