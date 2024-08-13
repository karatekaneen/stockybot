// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// SecuritiesColumns holds the columns for the "securities" table.
	SecuritiesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "country", Type: field.TypeString},
		{Name: "link_name", Type: field.TypeString, Nullable: true},
		{Name: "list", Type: field.TypeString},
		{Name: "type", Type: field.TypeEnum, Enums: []string{"stock", "index"}},
	}
	// SecuritiesTable holds the schema information for the "securities" table.
	SecuritiesTable = &schema.Table{
		Name:       "securities",
		Columns:    SecuritiesColumns,
		PrimaryKey: []*schema.Column{SecuritiesColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "security_name",
				Unique:  true,
				Columns: []*schema.Column{SecuritiesColumns[1]},
			},
		},
	}
	// WatchesColumns holds the columns for the "watches" table.
	WatchesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "watched_since", Type: field.TypeTime},
		{Name: "user_id", Type: field.TypeString},
		{Name: "security_watchers", Type: field.TypeInt, Nullable: true},
	}
	// WatchesTable holds the schema information for the "watches" table.
	WatchesTable = &schema.Table{
		Name:       "watches",
		Columns:    WatchesColumns,
		PrimaryKey: []*schema.Column{WatchesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "watches_securities_watchers",
				Columns:    []*schema.Column{WatchesColumns[3]},
				RefColumns: []*schema.Column{SecuritiesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "watch_user_id",
				Unique:  false,
				Columns: []*schema.Column{WatchesColumns[2]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		SecuritiesTable,
		WatchesTable,
	}
)

func init() {
	WatchesTable.ForeignKeys[0].RefTable = SecuritiesTable
}
