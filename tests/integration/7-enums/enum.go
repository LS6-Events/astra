package integration

type TestStringEnum string

const (
	TestStringEnumAvailable TestStringEnum = "available"
	TestStringEnumSold      TestStringEnum = "sold"
)

type TestStructWithStringEnum struct {
	// Enum
	Enum TestStringEnum `json:"enum,omitempty"`
}

type TestIntEnum int

const (
	TestIntEnumAvailable TestIntEnum = 1
	TestIntEnumSold      TestIntEnum = 2
)

type TestStructWithIntEnum struct {
	// Enum
	Enum TestIntEnum `json:"enum,omitempty"`
}
