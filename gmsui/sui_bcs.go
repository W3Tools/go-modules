package gmsui

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"time"
)

func Unmarshal(data []byte, v any) error {
	return NewDecoder(bytes.NewReader(data)).Decode(v)
}

type Decoder struct {
	reader io.Reader
}

func NewDecoder(reader io.Reader) *Decoder {
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
		return decoder.decodePointer(v.Elem())
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
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, // ints
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return decoder.decodeUint(v)
	default:
		bigInt := reflect.ValueOf(big.NewInt(0)).Type()

		switch v.Type() {
		case bigInt:
			return decoder.decodeBigInt(v)
		case reflect.ValueOf(time.Now()).Type():
			return decoder.decodeTime(v)
		}

		return fmt.Errorf("invalid kind [%v], type [%v]", kind, v.Type())
	}
}

func (decoder *Decoder) decodeSlice(v reflect.Value) error {
	data, err := io.ReadAll(decoder.reader)
	if err != nil {
		return fmt.Errorf("io.ReadAll %v", err)
	}

	size := int(data[0])
	data = data[1:]
	length := len(data) / size
	if length*size != len(data) {
		return fmt.Errorf("invalid data length")
	}

	valueType := v.Type()
	tmp := reflect.MakeSlice(valueType, 0, size)
	for i := 0; i < size; i++ {
		ind := reflect.New(valueType.Elem())
		innerData := data[i*length : (i+1)*length]
		decoder.reader = bytes.NewReader(innerData)
		if err := decoder.decode(ind); err != nil {
			return err
		}

		tmp = reflect.Append(tmp, ind.Elem())
	}
	v.Set(tmp)
	return nil
}

func (decoder *Decoder) decodeUint(v reflect.Value) error {
	data, err := io.ReadAll(decoder.reader)
	if err != nil {
		return fmt.Errorf("io.ReadAll %v", err)
	}

	v.SetUint(uint64(binary.LittleEndian.Uint16(data)))
	return nil
}

func (decoder *Decoder) decodeBool(v reflect.Value) error {
	data, err := io.ReadAll(decoder.reader)
	if err != nil {
		return fmt.Errorf("io.ReadAll %v", err)
	}
	if len(data) != 1 {
		return fmt.Errorf("invalid bool type")
	}

	if data[0] == 0 {
		v.SetBool(false)
	} else {
		v.SetBool(true)
	}
	return nil
}

func (decoder *Decoder) decodeBigInt(v reflect.Value) error {
	data, err := io.ReadAll(decoder.reader)
	if err != nil {
		return fmt.Errorf("io.ReadAll %v", err)
	}

	switch len(data) {
	case 16:
		b := new(big.Int).SetUint64(binary.LittleEndian.Uint64(data))
		v.Set(reflect.ValueOf(b))
		return nil
	default:
		return fmt.Errorf("invalid data length: [%v]", len(data))
	}
}

func (decoder *Decoder) decodeTime(v reflect.Value) error {
	data, err := io.ReadAll(decoder.reader)
	if err != nil {
		return fmt.Errorf("io.ReadAll %v", err)
	}
	fmt.Printf("len: %v\n", len(data))
	t := time.UnixMilli(int64(binary.LittleEndian.Uint64(data)))
	v.Set(reflect.ValueOf(t))
	return nil
}
