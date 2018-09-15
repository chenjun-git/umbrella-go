package json

import "testing"

type S struct {
	Int    int     `json:"int"`
	String string  `json:"string"`
	Bool   bool    `json:"bool"`
	Float  float64 `json:"float64"`
}

func TestJSON(t *testing.T) {

	a := S{
		Int:    2018,
		String: "meiqia",
		Bool:   true,
		Float:  3.14,
	}

	data, err := Marshal(a)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))

	var b S
	err = Unmarshal(data, &b)
	if err != nil {
		t.Fatal(err)
	}

	if a != b {
		t.Fatal("a != b")
	}
}
