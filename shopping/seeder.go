package shopping

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Load must fill the Db with the data from the file path provided
func Load(path string, db *MapDb) (*MapDb, error) {

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var products []*Product

	err = json.Unmarshal([]byte(file), &products)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling json: %v", err)
	}

	m, ok := db.mp.(map[string]*Product)
	if !ok {
		return nil, fmt.Errorf("error corrupt data type")
	}

	for _, prod := range products {
		m[prod.Id] = prod
	}

	return db, nil
}
