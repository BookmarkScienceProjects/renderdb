package helpers

import "github.com/jmoiron/sqlx"

type constructorFunc func() interface{}

// GetAll is an internal function for retrieving a set of rows and parsing each of them as a
// struct.
// - tx is a transaction used to execute the SQL
// - constructor is a function that returns a pointer to a prototype element (i.e. `new(MyStruct)`).
// - query is the actual SQL statement (with parameter placeholders).
// - params is the parameters user to replace the placeholders in the SQL statement
func GetAll(tx *sqlx.Tx, constructor constructorFunc, query string, params ...interface{}) ([]interface{}, error) {
	rows, err := tx.Queryx(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]interface{}, 0, 16)
	for rows.Next() {
		prototype := constructor()
		err = rows.StructScan(prototype)
		if err != nil {
			return nil, err
		}
		items = append(items, prototype)
	}
	return items, nil
}

// Get is an internal function for getting one element from SQL by ID and parsing the element as
// a struct. If the query returns an empty result set the function will return nil (and no error).
// - tx is a transaction used to execute the SQL
// - constructor is a function that returns a pointer to a prototype element (i.e. `new(MyStruct)`).
// - query is the actual SQL statement (with parameter placeholders).
// - ids is one or more ID used to identify the of the object to retrieve (one per placeholder in the
//   the query)
func Get(tx *sqlx.Tx, constructor constructorFunc, query string, ids ...interface{}) (interface{}, error) {
	rows, err := tx.Queryx(query, ids...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	if !rows.Next() {
		return nil, nil // No results
	}

	prototype := constructor()
	if err = rows.StructScan(prototype); err != nil {
		return nil, err
	}
	return prototype, nil
}
