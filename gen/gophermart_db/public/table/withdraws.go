//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var Withdraws = newWithdrawsTable("public", "withdraws", "")

type withdrawsTable struct {
	postgres.Table

	// Columns
	ID        postgres.ColumnInteger
	UserID    postgres.ColumnInteger
	Amount    postgres.ColumnFloat
	OrderID   postgres.ColumnString
	CreatedAt postgres.ColumnTimestamp

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type WithdrawsTable struct {
	withdrawsTable

	EXCLUDED withdrawsTable
}

// AS creates new WithdrawsTable with assigned alias
func (a WithdrawsTable) AS(alias string) *WithdrawsTable {
	return newWithdrawsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new WithdrawsTable with assigned schema name
func (a WithdrawsTable) FromSchema(schemaName string) *WithdrawsTable {
	return newWithdrawsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new WithdrawsTable with assigned table prefix
func (a WithdrawsTable) WithPrefix(prefix string) *WithdrawsTable {
	return newWithdrawsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new WithdrawsTable with assigned table suffix
func (a WithdrawsTable) WithSuffix(suffix string) *WithdrawsTable {
	return newWithdrawsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newWithdrawsTable(schemaName, tableName, alias string) *WithdrawsTable {
	return &WithdrawsTable{
		withdrawsTable: newWithdrawsTableImpl(schemaName, tableName, alias),
		EXCLUDED:       newWithdrawsTableImpl("", "excluded", ""),
	}
}

func newWithdrawsTableImpl(schemaName, tableName, alias string) withdrawsTable {
	var (
		IDColumn        = postgres.IntegerColumn("id")
		UserIDColumn    = postgres.IntegerColumn("user_id")
		AmountColumn    = postgres.FloatColumn("amount")
		OrderIDColumn   = postgres.StringColumn("order_id")
		CreatedAtColumn = postgres.TimestampColumn("created_at")
		allColumns      = postgres.ColumnList{IDColumn, UserIDColumn, AmountColumn, OrderIDColumn, CreatedAtColumn}
		mutableColumns  = postgres.ColumnList{UserIDColumn, AmountColumn, OrderIDColumn, CreatedAtColumn}
	)

	return withdrawsTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:        IDColumn,
		UserID:    UserIDColumn,
		Amount:    AmountColumn,
		OrderID:   OrderIDColumn,
		CreatedAt: CreatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
