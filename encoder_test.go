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
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack"

	zapmsgpack "github.com/smira/zap-msgpack-encoder"
)

func TestEncodeEntry(t *testing.T) {
	type bar struct {
		Key string  `msgpack:"key"`
		Val float64 `msgpack:"val"`
	}

	type foo struct {
		A string  `msgpack:"aee"`
		B int     `msgpack:"bee"`
		C float64 `msgpack:"cee"`
		D []bar   `msgpack:"dee"`
	}

	ts := time.Date(2018, 6, 19, 16, 33, 42, 99, time.Local)

	tests := []struct {
		desc     string
		expected interface{}
		ent      zapcore.Entry
		fields   []zapcore.Field
	}{
		{
			desc: "info entry with some fields",
			expected: []interface{}{
				&ts,
				map[string]interface{}{
					"L":          "info",
					"T":          &ts,
					"N":          "bob",
					"M":          "lob law",
					"so":         "passes",
					"answer":     int64(42),
					"common_pie": 3.14,
					"such": map[string]interface{}{
						"aee": "lol",
						"bee": int64(123),
						"cee": 0.9999,
						"dee": []interface{}{
							map[string]interface{}{"key": "pi", "val": 3.141592653589793},
							map[string]interface{}{"key": "tau", "val": 6.283185307179586},
						},
					},
				},
			},
			ent: zapcore.Entry{
				Level:      zapcore.InfoLevel,
				Time:       ts,
				LoggerName: "bob",
				Message:    "lob law",
			},
			fields: []zapcore.Field{
				zap.String("so", "passes"),
				zap.Int("answer", 42),
				zap.Float64("common_pie", 3.14),
				zap.Reflect("such", foo{
					A: "lol",
					B: 123,
					C: 0.9999,
					D: []bar{
						{"pi", 3.141592653589793},
						{"tau", 6.283185307179586},
					},
				}),
			},
		},
		{
			desc: "info entry with array fields",
			expected: []interface{}{
				&ts,
				map[string]interface{}{
					"L":    "debug",
					"T":    &ts,
					"M":    "lob law",
					"so":   "passes",
					"arr1": []interface{}{},
					"arr2": []interface{}{false, true, int8(34)},
					"arr3": []interface{}{"foo", 5.0, []interface{}{3.0, uint16(45678)}, map[string]interface{}{"x": uint16(3344)}},
					"ints": []interface{}{int16(1), int16(2), int16(3)},
					"bss":  []interface{}{"\x01\x02"},
				},
			},
			ent: zapcore.Entry{
				Level:   zapcore.DebugLevel,
				Time:    ts,
				Message: "lob law",
			},
			fields: []zapcore.Field{
				zap.Array("arr1", zapcore.ArrayMarshalerFunc(func(arr zapcore.ArrayEncoder) error {
					return nil
				})),
				zap.Array("arr2", zapcore.ArrayMarshalerFunc(func(arr zapcore.ArrayEncoder) error {
					arr.AppendBool(false)
					arr.AppendBool(true)
					arr.AppendInt(34)
					return nil
				})),
				zap.Array("arr3", zapcore.ArrayMarshalerFunc(func(arr zapcore.ArrayEncoder) error {
					arr.AppendString("foo")
					arr.AppendDuration(5 * time.Second)
					_ = arr.AppendArray(zapcore.ArrayMarshalerFunc(func(arrr zapcore.ArrayEncoder) error {
						arrr.AppendFloat64(3.0)
						arrr.AppendUint(45678)
						return nil
					}))
					_ = arr.AppendObject(zapcore.ObjectMarshalerFunc(func(obj zapcore.ObjectEncoder) error {
						obj.AddUint16("x", 3344)
						return nil
					}))
					return nil
				})),
				zap.String("so", "passes"),
				zap.Int16s("ints", []int16{1, 2, 3}),
				zap.ByteStrings("bss", [][]byte{[]byte{1, 2}}),
			},
		},
		{
			desc: "info entry with object fields",
			expected: []interface{}{
				&ts,
				map[string]interface{}{
					"L":    "debug",
					"T":    &ts,
					"M":    "",
					"d":    0.5,
					"obj1": map[string]interface{}{},
					"obj2": map[string]interface{}{
						"bits": []byte{0, 1, 2, 3},
						"why":  true,
						"bs":   "",
						"f":    float32(3.5),
						"oo": map[string]interface{}{
							"x":  int32(123456789),
							"y":  int16(-345),
							"z":  int16(-1),
							"zz": int8(0),
						},
					},
				},
			},
			ent: zapcore.Entry{
				Level: zapcore.DebugLevel,
				Time:  ts,
			},
			fields: []zapcore.Field{
				zap.Duration("d", 500*time.Millisecond),
				zap.Object("obj1", zapcore.ObjectMarshalerFunc(func(obj zapcore.ObjectEncoder) error {
					return nil
				})),
				zap.Object("obj2", zapcore.ObjectMarshalerFunc(func(obj zapcore.ObjectEncoder) error {
					obj.AddBinary("bits", []byte{0, 1, 2, 3})
					obj.AddBool("why", true)
					obj.AddByteString("bs", []byte{})
					obj.AddFloat32("f", 3.5)
					_ = obj.AddObject("oo", zapcore.ObjectMarshalerFunc(func(oobj zapcore.ObjectEncoder) error {
						oobj.AddInt("y", -345)
						oobj.AddInt32("x", 123456789)
						oobj.AddInt16("z", -1)
						oobj.AddInt8("zz", 0)
						return nil
					}))
					return nil
				})),
			},
		},
	}

	enc := zapmsgpack.NewEncoder(zapcore.EncoderConfig{
		MessageKey:     "M",
		LevelKey:       "L",
		TimeKey:        "T",
		NameKey:        "N",
		CallerKey:      "C",
		StacktraceKey:  "S",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			buf, err := enc.EncodeEntry(tt.ent, tt.fields)

			if assert.NoError(t, err, "Unexpected msgpack encoding error.") {
				var v interface{}

				err = msgpack.Unmarshal(buf.Bytes(), &v)
				if assert.NoErrorf(t, err, "Unexpected msgpack unmarshal error: %#v", buf.String()) {
					assert.EqualValues(t, tt.expected, v, "Incorrect encoded msgpack entry")
				}
			}
			buf.Free()
		})
	}
}
