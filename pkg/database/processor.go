package database

import (
	"database/sql"
	"fmt"
)

type Connector struct {
	DB *sql.DB
}

func NewConnector(db *sql.DB) *Connector {
	return &Connector{DB:db}
}

func GetQuery(s string, args... interface{}) string {
	return fmt.Sprintf(s, args...)
}

func (c *Connector) SelectOne(what, from, where string, args... interface{}) *sql.Row {
	query := GetQuery("SELECT %s FROM %s WHERE %s LIMIT 1", what, from, where)
	return c.DB.QueryRow(query, args...)
}

func (c *Connector) SelectAll(what, from string) (*sql.Rows, error) {
	query := GetQuery("SELECT %s FROM %s", what, from)
	return c.DB.Query(query)
}

func (c *Connector) SelectWhere(what, from, where string, args... interface{}) (*sql.Rows, error) {
	query := GetQuery("SELECT %s FROM %s WHERE %s", what, from, where)
	return c.DB.Query(query, args...)
}

func (c *Connector) Insert(table, what, values string, args... interface{}) (sql.Result, error) {
	query := GetQuery("INSERT INTO %s (%s) VALUES (%s)", table, what, values)
	return c.DB.Exec(query, args...)
}

func (c *Connector) Update(table, set, where string, args... interface{}) (sql.Result, error) {
	query := GetQuery("UPDATE %s SET %s WHERE %s", table, set, where)
	return c.DB.Exec(query, args...)
}

func (c *Connector) Delete(from, where string, args... interface{}) (sql.Result, error) {
	query := GetQuery("DELETE FROM %s WHERE %s", from, where)
	return c.DB.Exec(query, args...)
}
