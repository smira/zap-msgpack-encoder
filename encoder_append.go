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

func (enc *encoder) AppendArray(arr zapcore.ArrayMarshaler) error {
	enc.sliceLen++
	return enc.encodeArray(arr)
}

func (enc *encoder) AppendObject(obj zapcore.ObjectMarshaler) error {
	enc.sliceLen++
	return enc.encodeObject(obj)
}

func (enc *encoder) AppendBool(val bool) {
	enc.sliceLen++
	_ = enc.enc.EncodeBool(val)
}

func (enc *encoder) AppendByteString(val []byte) { // for UTF-8 encoded bytes
	enc.sliceLen++
	_ = enc.enc.EncodeString(string(val))
}

func (enc *encoder) AppendComplex128(complex128) {
	panic("not implemented")
}

func (enc *encoder) AppendComplex64(complex64) {
	panic("not implemented")
}

func (enc *encoder) AppendDuration(val time.Duration) {
	enc.sliceLen++
	_ = enc.enc.EncodeFloat64(val.Seconds())
}

func (enc *encoder) AppendFloat64(val float64) {
	enc.sliceLen++
	_ = enc.enc.EncodeFloat64(val)
}

func (enc *encoder) AppendFloat32(val float32) {
	enc.sliceLen++
	_ = enc.enc.EncodeFloat32(val)
}

func (enc *encoder) AppendInt(val int) {
	enc.sliceLen++
	_ = enc.enc.EncodeInt(int64(val))
}

func (enc *encoder) AppendInt64(val int64) {
	enc.sliceLen++
	_ = enc.enc.EncodeInt64(val)
}

func (enc *encoder) AppendInt32(val int32) {
	enc.sliceLen++
	_ = enc.enc.EncodeInt32(val)
}

func (enc *encoder) AppendInt16(val int16) {
	enc.sliceLen++
	_ = enc.enc.EncodeInt16(val)
}

func (enc *encoder) AppendInt8(val int8) {
	enc.sliceLen++
	_ = enc.enc.EncodeInt8(val)
}
func (enc *encoder) AppendString(val string) {
	enc.sliceLen++
	_ = enc.enc.EncodeString(val)
}

func (enc *encoder) AppendTime(val time.Time) {
	enc.sliceLen++
	_ = enc.enc.EncodeTime(val)
}

func (enc *encoder) AppendUint(val uint) {
	enc.sliceLen++
	_ = enc.enc.EncodeUint(uint64(val))
}

func (enc *encoder) AppendUint64(val uint64) {
	enc.sliceLen++
	_ = enc.enc.EncodeUint64(val)
}

func (enc *encoder) AppendUint32(val uint32) {
	enc.sliceLen++
	_ = enc.enc.EncodeUint32(val)
}
func (enc *encoder) AppendUint16(val uint16) {
	enc.sliceLen++
	_ = enc.enc.EncodeUint16(val)
}

func (enc *encoder) AppendUint8(val uint8) {
	enc.sliceLen++
	_ = enc.enc.EncodeUint8(val)
}
func (enc *encoder) AppendUintptr(val uintptr) {
	enc.sliceLen++
	_ = enc.enc.EncodeUint64(uint64(val))
}

func (enc *encoder) AppendReflected(val interface{}) error {
	enc.sliceLen++
	return enc.enc.Encode(val)
}
