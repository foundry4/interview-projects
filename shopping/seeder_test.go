package shopping

import (
	"reflect"
	"testing"
)

var db = make(map[string]*Product)
var mapDb = NewDB(db)

var filePath = "data/dump/test/"

func Test_Populate_Valid_File(t *testing.T) {

	file := filePath + "valid.json"

	m := make(map[string]*Product)
	m["1"] = &Product{Id: "1", Product: "abc", Price: "£1.00"}
	m["2"] = &Product{Id: "2", Product: "def", Price: "£2.00"}

	want := NewDB(m)

	got, err := Load(file, mapDb)
	if err != nil {
		t.Fatalf("error loading data to db: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected: %+v,  got: %+v", want, got)
	}

}

func Test_Populate_File_Not_Found(t *testing.T) {

	file := filePath + "not_found.json"

	_, err := Load(file, mapDb)
	if err == nil {
		t.Fatalf("expected err but not got")
	}
}

func Test_Populate_Invalid_File(t *testing.T) {

	file := filePath + "invalid.json"

	_, err := Load(file, mapDb)
	if err == nil {
		t.Fatalf("expected err but not got")
	}
}

func Test_Populate_Mismatch_File(t *testing.T) {

	file := filePath + "mismatch.json"

	want := map[string]*Product{
		"1": &Product{Id: "1", Product: "abc", Price: "£1.00"},
		"2": &Product{Id: "2", Product: "def", Price: "£2.00"},
	}

	got, err := Load(file, mapDb)
	if err != nil {
		t.Fatalf("error loading data to db: %v", err)
	}

	if reflect.DeepEqual(got, want) {
		t.Errorf("expected mismatch, but no mismatch found: %+v", got)
	}
}
