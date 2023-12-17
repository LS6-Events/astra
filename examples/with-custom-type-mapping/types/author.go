package types

type Author struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	TestArray [5]string      `json:"test_array"`
	TestMap   map[string]int `json:"test_map"`
}
