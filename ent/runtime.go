// Code generated by ent, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/karatekaneen/stockybot/ent/schema"
	"github.com/karatekaneen/stockybot/ent/security"
	"github.com/karatekaneen/stockybot/ent/watch"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	securityFields := schema.Security{}.Fields()
	_ = securityFields
	// securityDescName is the schema descriptor for name field.
	securityDescName := securityFields[1].Descriptor()
	// security.NameValidator is a validator for the "name" field. It is called by the builders before save.
	security.NameValidator = securityDescName.Validators[0].(func(string) error)
	// securityDescCountry is the schema descriptor for country field.
	securityDescCountry := securityFields[2].Descriptor()
	// security.DefaultCountry holds the default value on creation for the country field.
	security.DefaultCountry = securityDescCountry.Default.(string)
	// securityDescList is the schema descriptor for list field.
	securityDescList := securityFields[4].Descriptor()
	// security.DefaultList holds the default value on creation for the list field.
	security.DefaultList = securityDescList.Default.(string)
	// securityDescID is the schema descriptor for id field.
	securityDescID := securityFields[0].Descriptor()
	// security.IDValidator is a validator for the "id" field. It is called by the builders before save.
	security.IDValidator = securityDescID.Validators[0].(func(int64) error)
	watchFields := schema.Watch{}.Fields()
	_ = watchFields
	// watchDescWatchedSince is the schema descriptor for watched_since field.
	watchDescWatchedSince := watchFields[0].Descriptor()
	// watch.DefaultWatchedSince holds the default value on creation for the watched_since field.
	watch.DefaultWatchedSince = watchDescWatchedSince.Default.(func() time.Time)
	// watchDescUserID is the schema descriptor for user_id field.
	watchDescUserID := watchFields[1].Descriptor()
	// watch.UserIDValidator is a validator for the "user_id" field. It is called by the builders before save.
	watch.UserIDValidator = watchDescUserID.Validators[0].(func(string) error)
}
