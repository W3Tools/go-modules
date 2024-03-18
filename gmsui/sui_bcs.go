package gmsui

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

func UnmarshalSuiBCS(data []byte, v any) error {
	return NewSuiBCSDecoder(bytes.NewReader(data)).Decode(v)
}

type Decoder struct {
	reader io.Reader
}

func NewSuiBCSDecoder(reader io.Reader) *Decoder {
	return &Decoder{reader: reader}
}

func (decoder *Decoder) Decode(v any) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Pointer || value.IsNil() {
		return fmt.Errorf("not a pointer or nil pointer")
	}

	return decoder.decode(value)
}

func (decoder *Decoder) decode(v reflect.Value) error {
	if !v.CanInterface() {
		return nil
	}

	switch v.Kind() {
	case reflect.Pointer:
		return decoder.decodePointer(v.Elem())
	default:
		return decoder.decodePointer(v)
	}
}

func (decoder *Decoder) decodePointer(v reflect.Value) error {
	kind := v.Kind()

	if !v.CanSet() {
		return fmt.Errorf("cannot change value of kind %s", kind.String())
	}

	switch kind {
	case reflect.Bool:
		return decoder.decodeBool(v)
	case reflect.Slice:
		return decoder.decodeSlice(v)
	case reflect.Uint8:
		return decoder.decodeUint8(v)
	case reflect.Uint16:
		return decoder.decodeUint16(v)
	case reflect.Uint32:
		return decoder.decodeUint32(v)
	case reflect.Uint64:
		return decoder.decodeUint64(v)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return decoder.decodeUint32(v)
	default:
		return fmt.Errorf("invalid kind [%v], type [%v]", kind, v.Type())
	}
}

func (decoder *Decoder) decodeSlice(v reflect.Value) error {
	markSize, err := decoder.ReadBytes(1)
	if err != nil {
		return err
	}
	size := int(markSize[0])
	t := v.Type()
	tmp := reflect.MakeSlice(t, 0, size)
	for i := 0; i < size; i++ {
		innerData := reflect.New(t.Elem())
		if err := decoder.decode(innerData); err != nil {
			return err
		}
		tmp = reflect.Append(tmp, innerData.Elem())
	}

	v.Set(tmp)

	return nil
}

func (decoder *Decoder) decodeUint8(v reflect.Value) error {
	data, err := decoder.ReadBytes(1)
	if err != nil {
		return err
	}
	v.SetUint(uint64(data[0]))

	return nil
}

func (decoder *Decoder) decodeUint16(v reflect.Value) error {
	data, err := decoder.ReadBytes(2)
	if err != nil {
		return err
	}
	v.SetUint(uint64(binary.LittleEndian.Uint16(data)))

	return nil
}

func (decoder *Decoder) decodeUint32(v reflect.Value) error {
	data, err := decoder.ReadBytes(4)
	if err != nil {
		return err
	}
	v.SetUint(uint64(binary.LittleEndian.Uint32(data)))

	return nil
}

func (decoder *Decoder) decodeUint64(v reflect.Value) error {
	data, err := decoder.ReadBytes(8)
	if err != nil {
		return err
	}
	v.SetUint(binary.LittleEndian.Uint64(data))

	return nil
}

func (decoder *Decoder) decodeBool(v reflect.Value) error {
	data, err := decoder.ReadBytes(1)
	if err != nil {
		return err
	}

	if data[0] == 0 {
		v.SetBool(false)
	} else {
		v.SetBool(true)
	}
	return nil
}

func (decoder *Decoder) ReadBytes(len int) ([]byte, error) {
	b := make([]byte, len)
	n, err := decoder.reader.Read(b)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, fmt.Errorf("EOF")
	}
	return b, nil
}

func (decoder *Decoder) ReadByte() (byte, error) {
	b, err := decoder.ReadBytes(1)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}