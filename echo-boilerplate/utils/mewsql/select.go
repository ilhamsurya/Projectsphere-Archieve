package mewsql

import (
	"fmt"
	"strings"
)

const (
	selectOptionWhere = iota
	selectOptionOrderBy
	selectOptionLimit
	selectOptionOffset

	numSelectOptionKind
)

type SelectOption interface {
	getKind() int
	marshal(bindVarCount *int) (sql string, bindVars []interface{})
}

func Select(
	columns string,
	table string,
	options ...SelectOption,
) (sql string, bindVars []interface{}) {
	sqlBuffer := make([]strings.Builder, numSelectOptionKind)
	bindVarsBuffer := make([][]interface{}, numSelectOptionKind)

	numBindVars := 0
	for _, option := range options {
		kind := option.getKind()
		sql, bindVars := option.marshal(&numBindVars)

		sqlBuffer[kind].WriteString(sql)
		bindVarsBuffer = append(bindVarsBuffer, bindVars)
	}

	var sqlBufferString []string
	for _, buf := range sqlBuffer {
		sqlBufferString = append(sqlBufferString, buf.String())
	}

	sql = fmt.Sprintf(
		"SELECT %s FROM %s %s",
		columns,
		table,
		strings.Join(sqlBufferString, " "),
	)

	for _, curVars := range bindVarsBuffer {
		bindVars = append(bindVars, curVars...)
	}

	return
}

type genericSelectOptionImpl struct {
	kind      int
	statement string
}

func (opt *genericSelectOptionImpl) getKind() int {
	return opt.kind
}

func (opt *genericSelectOptionImpl) marshal(
	bindVarCount *int,
) (ret string, vars []interface{}) {
	ret = opt.statement
	return
}

func WithOrderBy(expression string, ascDesc string) SelectOption {
	ascDesc = strings.ToUpper(ascDesc)
	if ascDesc != "ASC" && ascDesc != "DESC" {
		ascDesc = ""
	}
	return &genericSelectOptionImpl{
		kind:      selectOptionOrderBy,
		statement: fmt.Sprintf("ORDER BY %s %s", expression, ascDesc),
	}
}

func WithLimit(count int) SelectOption {
	return &genericSelectOptionImpl{
		kind:      selectOptionLimit,
		statement: fmt.Sprintf("LIMIT %d", count),
	}
}

func WithOffset(count int) SelectOption {
	return &genericSelectOptionImpl{
		kind:      selectOptionOffset,
		statement: fmt.Sprintf("OFFSET %d", count),
	}
}

type Condition interface {
	marshal(bindVarCount *int) (sql string, bindVars []interface{})
}

type joiningCondition struct {
	conditions  []Condition
	joinKeyword string
}

func (c *joiningCondition) marshal(
	bindVarCount *int,
) (sql string, bindVars []interface{}) {
	var sqlBuf []string
	for _, cond := range c.conditions {
		curSql, curVars := cond.marshal(bindVarCount)

		sqlBuf = append(sqlBuf, curSql)
		bindVars = append(bindVars, curVars...)
	}
	if len(sqlBuf) > 0 {
		sql = "(" + strings.Join(sqlBuf, c.joinKeyword) + ")"
	}
	return
}

func And(conditions ...Condition) Condition {
	return &joiningCondition{conditions: conditions, joinKeyword: " AND "}
}

func Or(conditions ...Condition) Condition {
	return &joiningCondition{conditions: conditions, joinKeyword: " OR "}
}

type whereSelectOptionImpl struct {
	condition Condition
}

func (opt *whereSelectOptionImpl) getKind() int {
	return selectOptionWhere
}

func (opt *whereSelectOptionImpl) marshal(
	bindVarCount *int,
) (ret string, vars []interface{}) {
	retTemp, vars := opt.condition.marshal(bindVarCount)
	if len(retTemp) > 0 {
		ret = "WHERE " + retTemp
	}
	return
}

func WithWhere(
	conditions ...Condition,
) SelectOption {
	return &whereSelectOptionImpl{
		condition: And(conditions...),
	}
}

type basicCondition struct {
	sqlQuery string
	bindVar  interface{}
}

func (c *basicCondition) marshal(
	bindVarCount *int,
) (sql string, bindVars []interface{}) {
	sql = c.sqlQuery

	if c.bindVar != nil {
		(*bindVarCount)++
		sql = strings.ReplaceAll(
			sql,
			"?",
			fmt.Sprintf("$%d", *bindVarCount),
		)
		bindVars = append(bindVars, c.bindVar)
	}

	return
}

func WithCondition(sqlQuery string, bindVar interface{}) Condition {
	return &basicCondition{sqlQuery: sqlQuery, bindVar: bindVar}
}
