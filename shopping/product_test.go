package shopping

import (
	"reflect"
	"testing"

	"errors"
)

func TestSaveSuccess(t *testing.T) {

	p := &Product{Id: "123", Product: "product 123", Price: "£12"}
	ps := createProductServiceSUT(make(map[string]*Product))

	got, err := ps.Save(p)
	if err != nil {
		t.Fatalf("no error expected but got err: %v", err)
	}

	if p != got {
		t.Errorf("expected %v, but got %v", p, got)
	}
}

func TestSaveFailedWithCorruptDbError(t *testing.T) {

	p := new(Product)
	m := make(map[string]string)
	m["abc"] = "def"

	ps := createProductServiceSUT(m)

	_, err := ps.Save(p)
	if !errors.Is(err, ErrorCorruptDb) {
		t.Fatalf("corrupt db error expected but got no err %v", err)
	}

}

func TestFetchSuccess(t *testing.T) {

	p := &Product{Id: "123", Product: "product 123", Price: "£12"}
	ps := createProductServiceSUT(make(map[string]*Product))

	ps.Db.Type.(map[string]*Product)["123"] = p

	got, err := ps.Fetch("123")
	if err != nil {
		t.Fatalf("no error expected but got err: %v", err)
	}

	if p != got {
		t.Errorf("expected %v, but got %v", p, got)
	}
}

func TestFetchFailedNotFound(t *testing.T) {

	p := &Product{Id: "123", Product: "product 123", Price: "£12"}
	ps := createProductServiceSUT(make(map[string]*Product))

	ps.Db.Type.(map[string]*Product)["123"] = p

	_, err := ps.Fetch("456")
	if err == nil {
		t.Fatalf("fetch data error expected but got no err")
	}

}

func TestFetchAllSuccess(t *testing.T) {

	p := []*Product{
		&Product{Id: "123",
			Product: "product 123",
			Price:   "£12"},
	}
	ps := createProductServiceSUT(make(map[string]*Product))

	ps.Db.Type.(map[string]*Product)["123"] = &Product{Id: "123",
		Product: "product 123",
		Price:   "£12"}

	got, err := ps.FetchAll()
	if err != nil {
		t.Fatalf("no error expected but got err: %v", err)
	}

	if !reflect.DeepEqual(p, got) {
		t.Errorf("expected %v, but got %v", p, got)
	}
}

func TestFetchAllFailedCorruptDb(t *testing.T) {

	ps := createProductServiceSUT(make(map[string]string))

	_, err := ps.FetchAll()
	if !errors.Is(err, ErrorCorruptDb) {
		t.Fatalf("corrupt db error expected but got no err")
	}

}

func TestDeleteSuccess(t *testing.T) {

	p := &Product{Id: "123", Product: "product 123", Price: "£12"}
	ps := createProductServiceSUT(make(map[string]*Product))

	ps.Db.Type.(map[string]*Product)["123"] = p

	err := ps.Delete("123")
	if err != nil {
		t.Fatalf("no error expected but got err: %v", err)
	}

}

func TestDeleteFailedNotFound(t *testing.T) {

	p := &Product{Id: "123", Product: "product 123", Price: "£12"}
	ps := createProductServiceSUT(make(map[string]*Product))

	ps.Db.Type.(map[string]*Product)["123"] = p

	err := ps.Delete("456")
	if err == nil {
		t.Fatalf("fetch data error expected but got no err")
	}

}

func createProductServiceSUT(m interface{}) *ProductService {
	return NewProductService(NewDB(m))
}
