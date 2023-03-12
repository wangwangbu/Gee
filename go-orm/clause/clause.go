package clause

import "strings"

// Clause connects the sub-sentences
type Clause struct {
	sql		map[Type]string
	sqlVars map[Type][]interface{}
}

// Type chooses the type of sql command
type Type int
const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
)

func (c *Clause) Set(name Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	sql, vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}