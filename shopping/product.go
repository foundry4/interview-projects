package shopping

import "fmt"

// ProductService enacpsulates the Db which could be the main inventory or the cart
// N.B. datatype of contents might differ for inventory and cart
// but for the sake of simplicity, they are assume dto be the same
type ProductService struct {
	Db *DB
}

//Product model for capturing data to be stored
type Product struct {
	Id      string `json:"id"`
	Product string `json:"product"`
	Price   string `json:"price"`
}

// NewProductService creates  anew ProductService used for Product CRUD on db
func NewProductService(inv *DB) *ProductService {
	return &ProductService{
		Db: inv,
	}
}

// error types
var (
	ErrorCorruptDb = fmt.Errorf("db error: corrupt or nil db found")
)

func ErrorNoDataFound(id, action string) error {
	return fmt.Errorf("no data found for id: %v to %s", id, action)
}

// Save adds a new Product into the db
func (ps *ProductService) Save(p *Product) (*Product, error) {
	ps.Db.Mu.RLock()
	defer ps.Db.Mu.RUnlock()

	m, ok := ps.Db.Type.(map[string]*Product)
	if !ok {
		return nil, ErrorCorruptDb
	}

	m[p.Id] = p

	return m[p.Id], nil
}

// Fetch can find an existing product by id present in the db
func (ps *ProductService) Fetch(id string) (*Product, error) {
	ps.Db.Mu.RLock()
	defer ps.Db.Mu.RUnlock()

	m, ok := ps.Db.Type.(map[string]*Product)
	if !ok {
		return nil, ErrorCorruptDb
	}

	v, ok := m[id]
	if !ok {
		return nil, ErrorNoDataFound(id, "fetch")
	}
	return v, nil

}

// FetchAll can find all existing products present in the db
func (ps *ProductService) FetchAll() ([]*Product, error) {
	ps.Db.Mu.RLock()
	defer ps.Db.Mu.RUnlock()

	m, ok := ps.Db.Type.(map[string]*Product)
	if !ok {
		return nil, ErrorCorruptDb
	}

	l := []*Product{}

	for _, v := range m {
		l = append(l, v)
	}

	return l, nil
}

// Delete can find all existing products present in the db
func (ps *ProductService) Delete(id string) error {
	ps.Db.Mu.RLock()
	defer ps.Db.Mu.RUnlock()

	m, ok := ps.Db.Type.(map[string]*Product)
	if !ok {
		return ErrorCorruptDb
	}

	_, ok = m[id]
	if !ok {
		return ErrorNoDataFound(id, "delete")
	}
	delete(m, id)

	return nil
}
