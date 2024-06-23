package utils

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestHash(t *testing.T) {
	hash := "a4fbea6e8fec7ba1429275042bedbe3f43c708e448031efe4292be04d1f58325"
	s := struct {
		test string
	}{"Test"}
	result := Hash(s)

	t.Run("Hash must be always same.", func(t *testing.T) {
		if result != hash {
			t.Errorf("Hash(%q) = %q, want %q", s, result, hash)
		}
	})
	t.Run("Hash is hex encoded", func(t *testing.T) {
		_, err := hex.DecodeString(result)
		if err != nil {
			t.Errorf("Hash(%q) = %q, want nil", s, result)
		}
	})
}

func ExampleHash() {
	s := struct {
		test string
	}{"Test"}
	x := Hash(s)
	fmt.Println(x)
	// Output: a4fbea6e8fec7ba1429275042bedbe3f43c708e448031efe4292be04d1f58325
}

func TestToBytes(t *testing.T) {
	s := "test"
	result := ToBytes(s)
	if reflect.Slice != reflect.TypeOf(result).Kind() {
		t.Errorf("ToBytes(%q) = %T, want %s", s, result, reflect.Slice)
	}
}

func TestSplitter(t *testing.T) {
	type test struct {
		input  string
		sep    string
		idx    int
		output string
	}

	tests := []test{
		{"0:6:0", ":", 1, "6"},
		{"0:6:0", ":", 5, ""},
		{"0:6:0", "/", 0, "0:6:0"},
	}

	for _, tc := range tests {
		result := Splitter(tc.input, tc.sep, tc.idx)
		if result != tc.output {
			t.Errorf("Splitter(%q, %q, %d) = %q, want %q", tc.input, tc.sep, tc.idx, result, tc.output)
		}
	}

}

func TestHandleError(t *testing.T) {
	oldLogger := logger
	defer func() { logger = oldLogger }()
	called := false
	logger = func(v ...any) {
		called = true
	}

	err := errors.New("test")
	HandleError(err)
	if !called {
		t.Errorf("HandleError(%q) = false, want true", "test")
	}
}

func TestFromBytes(t *testing.T) {
	type dummy struct {
		S string
	}

	d := dummy{"test"}
	bytes := ToBytes(d)

	var restored dummy
	FromBytes(&restored, bytes)

	if !reflect.DeepEqual(d, restored) {
		t.Errorf("FromBytes(%v) = %v, want %v", d, restored, d)
	}
}

func TestToJSON(t *testing.T) {
	type dummy struct {
		S string
	}

	d := dummy{"test"}
	bytes := ToJSON(d)
	var restored dummy
	json.Unmarshal(bytes, &restored)

	if reflect.TypeOf(bytes).Kind() != reflect.Slice {
		t.Errorf("ToJSON(%v) = %T, want slice", d, bytes)
	}
	if !reflect.DeepEqual(d, restored) {
		t.Errorf("ToJSON(%v) = %v, want %v", d, restored, d)
	}
}
