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
	"bytes"
	"sync"

	"github.com/vmihailenco/msgpack"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

const initialSize = 1024

type encoder struct {
	*zapcore.EncoderConfig

	buf      *bytes.Buffer
	enc      *msgpack.Encoder
	mapSize  int
	sliceLen int
	nsPrefix string
}

var bufPool = buffer.NewPool()

var encoderPool = sync.Pool{
	New: func() interface{} {
		enc := &encoder{
			buf: bytes.NewBuffer(make([]byte, 0, initialSize)),
		}

		enc.enc = msgpack.NewEncoder(enc.buf)

		return enc
	},
}

func getEncoder() *encoder {
	return encoderPool.Get().(*encoder)
}

func putEncoder(enc *encoder) {
	enc.buf.Reset()
	enc.EncoderConfig = nil
	enc.mapSize = 0
	enc.sliceLen = 0
	enc.nsPrefix = ""

	encoderPool.Put(enc)
}

// NewEncoder creates fast, low-allocation msgpack encoder
//
// Msgpack encoder could be used e.g. while delivering go.uber.org/zap logs
// to fluentd destination.
func NewEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	enc := getEncoder()
	enc.EncoderConfig = &cfg

	return enc
}

func (enc *encoder) encodeKey(key string) {
	_ = enc.enc.EncodeString(enc.nsPrefix + key)
}

func (enc *encoder) encodeArray(arr zapcore.ArrayMarshaler) error {
	sliceEnc := enc.clone()
	if err := arr.MarshalLogArray(sliceEnc); err != nil {
		return err
	}

	if err := enc.enc.EncodeArrayLen(sliceEnc.sliceLen); err != nil {
		return err
	}

	if _, err := enc.buf.Write(sliceEnc.buf.Bytes()); err != nil {
		return err
	}

	putEncoder(sliceEnc)

	return nil
}

func (enc *encoder) encodeObject(obj zapcore.ObjectMarshaler) error {
	mapEnc := enc.clone()
	if err := obj.MarshalLogObject(mapEnc); err != nil {
		return err
	}

	if err := enc.enc.EncodeMapLen(mapEnc.mapSize); err != nil {
		return err
	}

	if _, err := enc.buf.Write(mapEnc.buf.Bytes()); err != nil {
		return err
	}

	putEncoder(mapEnc)

	return nil
}

// OpenNamespace opens an isolated namespace where all subsequent fields will
// be added. Applications can use namespaces to prevent key collisions when
// injecting loggers into sub-components or third-party libraries.
func (enc *encoder) OpenNamespace(key string) {
	enc.nsPrefix += key + "."
}

func (enc *encoder) clone() *encoder {
	clone := getEncoder()
	clone.EncoderConfig = enc.EncoderConfig

	return clone
}

// Clone copies the encoder, ensuring that adding fields to the copy doesn't
// affect the original.
func (enc *encoder) Clone() zapcore.Encoder {
	clone := enc.clone()
	_, _ = clone.buf.Write(enc.buf.Bytes())
	clone.mapSize = enc.mapSize
	clone.sliceLen = enc.sliceLen
	clone.nsPrefix = enc.nsPrefix
	return clone
}

// EncodeEntry encodes an entry and fields, along with any accumulated
// context, into a byte buffer and returns it.
//
// This method serializes message as "Entry" from fluentd forward protocol specification:
// https://github.com/fluent/fluentd/wiki/Forward-Protocol-Specification-v1
//
// [ timestamp, {key : value, ... } ]
func (enc *encoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	finenc := enc.clone()

	_ = finenc.enc.EncodeArrayLen(2)
	_ = finenc.enc.EncodeTime(ent.Time)

	final := enc.clone()

	if final.LevelKey != "" {
		final.mapSize++
		_ = final.enc.EncodeString(final.LevelKey)
		cur := final.buf.Len()
		final.EncodeLevel(ent.Level, final)
		if cur == final.buf.Len() {
			// User-supplied EncodeLevel was a no-op. Fall back to strings to keep
			// output JSON valid.
			_ = final.enc.EncodeString(ent.Level.String())
		}
	}
	if final.TimeKey != "" {
		final.AddTime(final.TimeKey, ent.Time)
	}
	if ent.LoggerName != "" && final.NameKey != "" {
		final.mapSize++
		_ = final.enc.EncodeString(final.NameKey)
		cur := final.buf.Len()
		nameEncoder := final.EncodeName

		// if no name encoder provided, fall back to FullNameEncoder for backwards
		// compatibility
		if nameEncoder == nil {
			nameEncoder = zapcore.FullNameEncoder
		}

		nameEncoder(ent.LoggerName, final)
		if cur == final.buf.Len() {
			// User-supplied EncodeName was a no-op. Fall back to strings to
			// keep output valid.
			_ = final.enc.EncodeString(ent.LoggerName)
		}
	}
	if ent.Caller.Defined && final.CallerKey != "" {
		final.mapSize++
		_ = final.enc.EncodeString(final.CallerKey)
		cur := final.buf.Len()
		final.EncodeCaller(ent.Caller, final)
		if cur == final.buf.Len() {
			// User-supplied EncodeCaller was a no-op. Fall back to strings to
			// keep output valid.
			_ = final.enc.EncodeString(ent.Caller.String())
		}
	}
	if final.MessageKey != "" {
		final.mapSize++
		_ = final.enc.EncodeString(enc.MessageKey)
		_ = final.enc.EncodeString(ent.Message)
	}

	for i := range fields {
		fields[i].AddTo(final)
	}

	if ent.Stack != "" && final.StacktraceKey != "" {
		final.AddString(final.StacktraceKey, ent.Stack)
	}

	_ = finenc.enc.EncodeMapLen(final.mapSize)
	_, _ = finenc.buf.Write(final.buf.Bytes())

	buf := bufPool.Get()
	_, _ = buf.Write(finenc.buf.Bytes())

	putEncoder(final)
	putEncoder(finenc)

	return buf, nil
}
