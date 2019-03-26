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

package zapmsgpack_test

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	zapmsgpack "github.com/smira/zap-msgpack-encoder"
)

func testEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "name",
		TimeKey:        "ts",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     "\n",
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func BenchmarkZapMsgpack(b *testing.B) {
	enc := zapmsgpack.NewEncoder(testEncoderConfig())
	enc.AddString("str", "foo")
	enc.AddInt64("int64-1", 1)
	enc.AddInt64("int64-2", 2)
	enc.AddFloat64("float64", 1.0)
	enc.AddString("string1", "\n")
	enc.AddString("string2", "💩")
	enc.AddString("string3", "🤔")
	enc.AddString("string4", "🙊")
	enc.AddBool("bool", true)

	entry := zapcore.Entry{
		Message: "fake",
		Level:   zap.DebugLevel,
	}

	fields := []zap.Field{
		zap.String("string5", "yeah"),
		zap.Int64("int64-3", 543210),
		zap.Bool("yes", false),
		zap.Int16s("setofints", []int16{0, 1, 2, 3, 4, 5}),
		zap.Object("obj2", zapcore.ObjectMarshalerFunc(func(obj zapcore.ObjectEncoder) error {
			obj.AddBinary("bits", []byte{0, 1, 2, 3})
			obj.AddBool("why", true)
			obj.AddByteString("bs", []byte{})
			obj.AddFloat32("f", 3.5)
			return nil
		})),
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf, _ := enc.EncodeEntry(entry, fields)
			buf.Free()
		}
	})
}
