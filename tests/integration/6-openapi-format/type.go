package petstore

import (
	"github.com/google/uuid"
	"time"
)

type TestStructFormatter struct {
	// String
	String string `json:"string,omitempty"`
	// Int (int32)
	Int int `json:"int,omitempty"`
	// Int8
	Int8 int8 `json:"int8,omitempty"`
	// Int16
	Int16 int16 `json:"int16,omitempty"`
	// Int32
	Int32 int32 `json:"int32,omitempty"`
	// Int64
	Int64 int64 `json:"int64,omitempty"`
	// Uint
	Uint uint `json:"uint,omitempty"`
	// Uint8
	Uint8 uint8 `json:"uint8,omitempty"`
	// Uint16
	Uint16 uint16 `json:"uint16,omitempty"`
	// Uint32
	Uint32 uint32 `json:"uint32,omitempty"`
	// Uint64
	Uint64 uint64 `json:"uint64,omitempty"`
	// Float32
	Float32 float32 `json:"float32,omitempty"`
	// Float64
	Float64 float64 `json:"float64,omitempty"`
	// Bool
	Bool bool `json:"bool,omitempty"`
	// Byte
	Byte byte `json:"byte,omitempty"`
	// Rune
	Rune rune `json:"rune,omitempty"`
	// Struct
	Struct struct{} `json:"struct,omitempty"`
	// Map
	Map map[string]string `json:"map,omitempty"`
	// Slice
	Slice []string `json:"slice,omitempty"`
	// Any
	Any any `json:"any,omitempty"`

	// time.Time
	Time time.Time `json:"time,omitempty"`
	// uuid.UUID
	UUID uuid.UUID `json:"uuid,omitempty"`
}
