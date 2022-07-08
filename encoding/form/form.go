package form

import (
	"net/url"
	"reflect"

	"github.com/go-playground/form/v4"

	"skuld/encoding"
)

// Name is form codec name
const Name = "x-www-form-urlencoded"

func init() {
	decoder := form.NewDecoder()
	decoder.SetTagName("json")
	encoder := form.NewEncoder()
	encoder.SetTagName("json")
	encoding.RegisterCodec(codec{encoder: encoder, decoder: decoder})
}

type codec struct {
	encoder *form.Encoder
	decoder *form.Decoder
}

func (c codec) Marshal(v interface{}) ([]byte, error) {
	var (
		vs  url.Values
		err error
	)
	if vs, err = c.encoder.Encode(v); err != nil {
		return nil, err
	}
	for k, v := range vs {
		if len(v) == 0 {
			delete(vs, k)
		}
	}
	return []byte(vs.Encode()), nil
}

func (c codec) Unmarshal(data []byte, v interface{}) error {
	vs, err := url.ParseQuery(string(data))
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}
	return c.decoder.Decode(v, vs)
}

func (codec) Name() string {
	return Name
}
