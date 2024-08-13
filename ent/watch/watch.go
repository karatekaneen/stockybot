// Code generated by ent, DO NOT EDIT.

package watch

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the watch type in the database.
	Label = "watch"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldWatchedSince holds the string denoting the watched_since field in the database.
	FieldWatchedSince = "watched_since"
	// FieldUserID holds the string denoting the user_id field in the database.
	FieldUserID = "user_id"
	// EdgeWatching holds the string denoting the watching edge name in mutations.
	EdgeWatching = "watching"
	// Table holds the table name of the watch in the database.
	Table = "watches"
	// WatchingTable is the table that holds the watching relation/edge.
	WatchingTable = "watches"
	// WatchingInverseTable is the table name for the Security entity.
	// It exists in this package in order to avoid circular dependency with the "security" package.
	WatchingInverseTable = "securities"
	// WatchingColumn is the table column denoting the watching relation/edge.
	WatchingColumn = "security_watchers"
)

// Columns holds all SQL columns for watch fields.
var Columns = []string{
	FieldID,
	FieldWatchedSince,
	FieldUserID,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "watches"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"security_watchers",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultWatchedSince holds the default value on creation for the "watched_since" field.
	DefaultWatchedSince func() time.Time
	// UserIDValidator is a validator for the "user_id" field. It is called by the builders before save.
	UserIDValidator func(string) error
)

// OrderOption defines the ordering options for the Watch queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByWatchedSince orders the results by the watched_since field.
func ByWatchedSince(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldWatchedSince, opts...).ToFunc()
}

// ByUserID orders the results by the user_id field.
func ByUserID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUserID, opts...).ToFunc()
}

// ByWatchingField orders the results by watching field.
func ByWatchingField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newWatchingStep(), sql.OrderByField(field, opts...))
	}
}
func newWatchingStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(WatchingInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, WatchingTable, WatchingColumn),
	)
}
