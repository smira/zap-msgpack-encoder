// Copyright (c) 2019 Andrey Smirnov
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package zapmsgpack

import (
	"time"

	"go.uber.org/zap/zapcore"
)

func (enc *encoder) AddArray(key string, arr zapcore.ArrayMarshaler) error {
	enc.mapSize++
	enc.encodeKey(key)

	return enc.encodeArray(arr)
}

func (enc *encoder) AddObject(key string, obj zapcore.ObjectMarshaler) error {
	enc.mapSize++
	enc.encodeKey(key)

	return enc.encodeObject(obj)
}

func (enc *encoder) AddBinary(key string, val []byte) {
	enc.mapSize++
	enc.encodeKey(key)
	_ = enc.enc.EncodeBytes(val)
}

func (enc *encoder) AddByteString(key string, val []byte) {
	enc.mapSize++
	enc.encodeKey(key)
	_ = enc.enc.EncodeString(string(val))
}

func (enc *encoder) AddBool(key string, val bool) {
	enc.mapSize++
	enc.encodeKey(key)
	_ = enc.enc.EncodeBool(val)
}

func (enc *encoder) AddComplex128(key string, val complex128) {
	panic("complex numbers not supported in msgpack")
}

func (enc *encoder) AddComplex64(key string, val complex64) {
	panic("complex numbers not supported in msgpack")
}

func (enc *encoder) AddDuration(key string, val time.Duration) {
	enc.mapSize++
	enc.encodeKey(key)
	_ = enc.enc.EncodeFloat64(val.Seconds())
}

func (enc *encoder) AddFloat64(key string, val float64) {
	enc.mapSize++
	enc.encodeKey(key)
	_ = enc.enc.EncodeFloat64(val)
}

func (enc *encoder) AddFloat32(key string, val float32) {
	enc.mapSize++
	enc.encodeKey(key)
	_ = enc.enc.EncodeFloat32(val)
}

func (enc *encoder) AddInt(key string, val int) {
	enc.mapSize++
	enc.encodeKey(key)
	_ = enc.enc.EncodeInt(int64(val))
}

func (enc *encoder) AddInt64(key string, val int64) {
	enc.mapSize++
	enc.encodeKey(key)
	_ = enc.enc.EncodeInt64(val)
}

func (enc *encoder) AddInt32(key string, val int32) {
	enc.mapSize++
	enc.encodeKey(key)
	_ = enc.enc.EncodeInt32(val)
}

func (enc *encoder) AddInt16(key string, val int16) {
	enc.mapSize++
	enc.encodeKey(key)
	_ = enc.enc.EncodeInt16(val)
}

func (enc *encoder) AddInt8(key string, val int8) {
	enc.mapSize++
	enc.encodeKey(key)
	_ = enc.enc.EncodeInt8(val)
}

func (enc *encoder) AddString(key string, val string) {
	enc.mapSize++
	enc.encodeKey(key)
	_ = enc.enc.EncodeString(val)
}

func (enc *encoder) AddTime(key string, val time.Time) {
	enc.mapSize++
	enc.encodeKey(key)
	_ = enc.enc.EncodeTime(val)
}

func (enc *encoder) AddUint(key string, val uint) {
	enc.mapSize++
	enc.encodeKey(key)
	_ = enc.enc.EncodeUint(uint64(val))
}

func (enc *encoder) AddUint64(key string, val uint64) {
	enc.mapSize++
	enc.encodeKey(key)
	_ = enc.enc.EncodeUint64(val)
}

func (enc *encoder) AddUint32(key string, val uint32) {
	enc.mapSize++
	enc.encodeKey(key)
	_ = enc.enc.EncodeUint32(val)
}

func (enc *encoder) AddUint16(key string, val uint16) {
	enc.mapSize++
	enc.encodeKey(key)
	_ = enc.enc.EncodeUint16(val)
}

func (enc *encoder) AddUint8(key string, val uint8) {
	enc.mapSize++
	enc.encodeKey(key)
	_ = enc.enc.EncodeUint8(val)
}

func (enc *encoder) AddUintptr(key string, val uintptr) {
	enc.mapSize++
	enc.encodeKey(key)
	_ = enc.enc.EncodeUint(uint64(val))
}

// AddReflected uses reflection to serialize arbitrary objects, so it's slow
// and allocation-heavy.
func (enc *encoder) AddReflected(key string, val interface{}) error {
	enc.mapSize++
	enc.encodeKey(key)
	return enc.enc.Encode(val)
}
