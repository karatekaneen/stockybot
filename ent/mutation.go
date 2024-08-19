// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/karatekaneen/stockybot/ent/predicate"
	"github.com/karatekaneen/stockybot/ent/security"
	"github.com/karatekaneen/stockybot/ent/watch"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeSecurity = "Security"
	TypeWatch    = "Watch"
)

// SecurityMutation represents an operation that mutates the Security nodes in the graph.
type SecurityMutation struct {
	config
	op              Op
	typ             string
	id              *int64
	name            *string
	country         *string
	link_name       *string
	list            *string
	_type           *security.Type
	clearedFields   map[string]struct{}
	watchers        map[int]struct{}
	removedwatchers map[int]struct{}
	clearedwatchers bool
	done            bool
	oldValue        func(context.Context) (*Security, error)
	predicates      []predicate.Security
}

var _ ent.Mutation = (*SecurityMutation)(nil)

// securityOption allows management of the mutation configuration using functional options.
type securityOption func(*SecurityMutation)

// newSecurityMutation creates new mutation for the Security entity.
func newSecurityMutation(c config, op Op, opts ...securityOption) *SecurityMutation {
	m := &SecurityMutation{
		config:        c,
		op:            op,
		typ:           TypeSecurity,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withSecurityID sets the ID field of the mutation.
func withSecurityID(id int64) securityOption {
	return func(m *SecurityMutation) {
		var (
			err   error
			once  sync.Once
			value *Security
		)
		m.oldValue = func(ctx context.Context) (*Security, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Security.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withSecurity sets the old Security of the mutation.
func withSecurity(node *Security) securityOption {
	return func(m *SecurityMutation) {
		m.oldValue = func(context.Context) (*Security, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m SecurityMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m SecurityMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of Security entities.
func (m *SecurityMutation) SetID(id int64) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *SecurityMutation) ID() (id int64, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *SecurityMutation) IDs(ctx context.Context) ([]int64, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int64{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Security.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetName sets the "name" field.
func (m *SecurityMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *SecurityMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the Security entity.
// If the Security object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SecurityMutation) OldName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldName is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldName: %w", err)
	}
	return oldValue.Name, nil
}

// ResetName resets all changes to the "name" field.
func (m *SecurityMutation) ResetName() {
	m.name = nil
}

// SetCountry sets the "country" field.
func (m *SecurityMutation) SetCountry(s string) {
	m.country = &s
}

// Country returns the value of the "country" field in the mutation.
func (m *SecurityMutation) Country() (r string, exists bool) {
	v := m.country
	if v == nil {
		return
	}
	return *v, true
}

// OldCountry returns the old "country" field's value of the Security entity.
// If the Security object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SecurityMutation) OldCountry(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCountry is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCountry requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCountry: %w", err)
	}
	return oldValue.Country, nil
}

// ResetCountry resets all changes to the "country" field.
func (m *SecurityMutation) ResetCountry() {
	m.country = nil
}

// SetLinkName sets the "link_name" field.
func (m *SecurityMutation) SetLinkName(s string) {
	m.link_name = &s
}

// LinkName returns the value of the "link_name" field in the mutation.
func (m *SecurityMutation) LinkName() (r string, exists bool) {
	v := m.link_name
	if v == nil {
		return
	}
	return *v, true
}

// OldLinkName returns the old "link_name" field's value of the Security entity.
// If the Security object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SecurityMutation) OldLinkName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldLinkName is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldLinkName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldLinkName: %w", err)
	}
	return oldValue.LinkName, nil
}

// ClearLinkName clears the value of the "link_name" field.
func (m *SecurityMutation) ClearLinkName() {
	m.link_name = nil
	m.clearedFields[security.FieldLinkName] = struct{}{}
}

// LinkNameCleared returns if the "link_name" field was cleared in this mutation.
func (m *SecurityMutation) LinkNameCleared() bool {
	_, ok := m.clearedFields[security.FieldLinkName]
	return ok
}

// ResetLinkName resets all changes to the "link_name" field.
func (m *SecurityMutation) ResetLinkName() {
	m.link_name = nil
	delete(m.clearedFields, security.FieldLinkName)
}

// SetList sets the "list" field.
func (m *SecurityMutation) SetList(s string) {
	m.list = &s
}

// List returns the value of the "list" field in the mutation.
func (m *SecurityMutation) List() (r string, exists bool) {
	v := m.list
	if v == nil {
		return
	}
	return *v, true
}

// OldList returns the old "list" field's value of the Security entity.
// If the Security object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SecurityMutation) OldList(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldList is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldList requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldList: %w", err)
	}
	return oldValue.List, nil
}

// ResetList resets all changes to the "list" field.
func (m *SecurityMutation) ResetList() {
	m.list = nil
}

// SetType sets the "type" field.
func (m *SecurityMutation) SetType(s security.Type) {
	m._type = &s
}

// GetType returns the value of the "type" field in the mutation.
func (m *SecurityMutation) GetType() (r security.Type, exists bool) {
	v := m._type
	if v == nil {
		return
	}
	return *v, true
}

// OldType returns the old "type" field's value of the Security entity.
// If the Security object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SecurityMutation) OldType(ctx context.Context) (v security.Type, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldType is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldType requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldType: %w", err)
	}
	return oldValue.Type, nil
}

// ResetType resets all changes to the "type" field.
func (m *SecurityMutation) ResetType() {
	m._type = nil
}

// AddWatcherIDs adds the "watchers" edge to the Watch entity by ids.
func (m *SecurityMutation) AddWatcherIDs(ids ...int) {
	if m.watchers == nil {
		m.watchers = make(map[int]struct{})
	}
	for i := range ids {
		m.watchers[ids[i]] = struct{}{}
	}
}

// ClearWatchers clears the "watchers" edge to the Watch entity.
func (m *SecurityMutation) ClearWatchers() {
	m.clearedwatchers = true
}

// WatchersCleared reports if the "watchers" edge to the Watch entity was cleared.
func (m *SecurityMutation) WatchersCleared() bool {
	return m.clearedwatchers
}

// RemoveWatcherIDs removes the "watchers" edge to the Watch entity by IDs.
func (m *SecurityMutation) RemoveWatcherIDs(ids ...int) {
	if m.removedwatchers == nil {
		m.removedwatchers = make(map[int]struct{})
	}
	for i := range ids {
		delete(m.watchers, ids[i])
		m.removedwatchers[ids[i]] = struct{}{}
	}
}

// RemovedWatchers returns the removed IDs of the "watchers" edge to the Watch entity.
func (m *SecurityMutation) RemovedWatchersIDs() (ids []int) {
	for id := range m.removedwatchers {
		ids = append(ids, id)
	}
	return
}

// WatchersIDs returns the "watchers" edge IDs in the mutation.
func (m *SecurityMutation) WatchersIDs() (ids []int) {
	for id := range m.watchers {
		ids = append(ids, id)
	}
	return
}

// ResetWatchers resets all changes to the "watchers" edge.
func (m *SecurityMutation) ResetWatchers() {
	m.watchers = nil
	m.clearedwatchers = false
	m.removedwatchers = nil
}

// Where appends a list predicates to the SecurityMutation builder.
func (m *SecurityMutation) Where(ps ...predicate.Security) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the SecurityMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *SecurityMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Security, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *SecurityMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *SecurityMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Security).
func (m *SecurityMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *SecurityMutation) Fields() []string {
	fields := make([]string, 0, 5)
	if m.name != nil {
		fields = append(fields, security.FieldName)
	}
	if m.country != nil {
		fields = append(fields, security.FieldCountry)
	}
	if m.link_name != nil {
		fields = append(fields, security.FieldLinkName)
	}
	if m.list != nil {
		fields = append(fields, security.FieldList)
	}
	if m._type != nil {
		fields = append(fields, security.FieldType)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *SecurityMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case security.FieldName:
		return m.Name()
	case security.FieldCountry:
		return m.Country()
	case security.FieldLinkName:
		return m.LinkName()
	case security.FieldList:
		return m.List()
	case security.FieldType:
		return m.GetType()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *SecurityMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case security.FieldName:
		return m.OldName(ctx)
	case security.FieldCountry:
		return m.OldCountry(ctx)
	case security.FieldLinkName:
		return m.OldLinkName(ctx)
	case security.FieldList:
		return m.OldList(ctx)
	case security.FieldType:
		return m.OldType(ctx)
	}
	return nil, fmt.Errorf("unknown Security field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *SecurityMutation) SetField(name string, value ent.Value) error {
	switch name {
	case security.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case security.FieldCountry:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCountry(v)
		return nil
	case security.FieldLinkName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLinkName(v)
		return nil
	case security.FieldList:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetList(v)
		return nil
	case security.FieldType:
		v, ok := value.(security.Type)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetType(v)
		return nil
	}
	return fmt.Errorf("unknown Security field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *SecurityMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *SecurityMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *SecurityMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Security numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *SecurityMutation) ClearedFields() []string {
	var fields []string
	if m.FieldCleared(security.FieldLinkName) {
		fields = append(fields, security.FieldLinkName)
	}
	return fields
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *SecurityMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *SecurityMutation) ClearField(name string) error {
	switch name {
	case security.FieldLinkName:
		m.ClearLinkName()
		return nil
	}
	return fmt.Errorf("unknown Security nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *SecurityMutation) ResetField(name string) error {
	switch name {
	case security.FieldName:
		m.ResetName()
		return nil
	case security.FieldCountry:
		m.ResetCountry()
		return nil
	case security.FieldLinkName:
		m.ResetLinkName()
		return nil
	case security.FieldList:
		m.ResetList()
		return nil
	case security.FieldType:
		m.ResetType()
		return nil
	}
	return fmt.Errorf("unknown Security field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *SecurityMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.watchers != nil {
		edges = append(edges, security.EdgeWatchers)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *SecurityMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case security.EdgeWatchers:
		ids := make([]ent.Value, 0, len(m.watchers))
		for id := range m.watchers {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *SecurityMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	if m.removedwatchers != nil {
		edges = append(edges, security.EdgeWatchers)
	}
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *SecurityMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case security.EdgeWatchers:
		ids := make([]ent.Value, 0, len(m.removedwatchers))
		for id := range m.removedwatchers {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *SecurityMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedwatchers {
		edges = append(edges, security.EdgeWatchers)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *SecurityMutation) EdgeCleared(name string) bool {
	switch name {
	case security.EdgeWatchers:
		return m.clearedwatchers
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *SecurityMutation) ClearEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown Security unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *SecurityMutation) ResetEdge(name string) error {
	switch name {
	case security.EdgeWatchers:
		m.ResetWatchers()
		return nil
	}
	return fmt.Errorf("unknown Security edge %s", name)
}

// WatchMutation represents an operation that mutates the Watch nodes in the graph.
type WatchMutation struct {
	config
	op              Op
	typ             string
	id              *int
	watched_since   *time.Time
	user_id         *string
	clearedFields   map[string]struct{}
	security        *int64
	clearedsecurity bool
	done            bool
	oldValue        func(context.Context) (*Watch, error)
	predicates      []predicate.Watch
}

var _ ent.Mutation = (*WatchMutation)(nil)

// watchOption allows management of the mutation configuration using functional options.
type watchOption func(*WatchMutation)

// newWatchMutation creates new mutation for the Watch entity.
func newWatchMutation(c config, op Op, opts ...watchOption) *WatchMutation {
	m := &WatchMutation{
		config:        c,
		op:            op,
		typ:           TypeWatch,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withWatchID sets the ID field of the mutation.
func withWatchID(id int) watchOption {
	return func(m *WatchMutation) {
		var (
			err   error
			once  sync.Once
			value *Watch
		)
		m.oldValue = func(ctx context.Context) (*Watch, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Watch.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withWatch sets the old Watch of the mutation.
func withWatch(node *Watch) watchOption {
	return func(m *WatchMutation) {
		m.oldValue = func(context.Context) (*Watch, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m WatchMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m WatchMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *WatchMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *WatchMutation) IDs(ctx context.Context) ([]int, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Watch.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetWatchedSince sets the "watched_since" field.
func (m *WatchMutation) SetWatchedSince(t time.Time) {
	m.watched_since = &t
}

// WatchedSince returns the value of the "watched_since" field in the mutation.
func (m *WatchMutation) WatchedSince() (r time.Time, exists bool) {
	v := m.watched_since
	if v == nil {
		return
	}
	return *v, true
}

// OldWatchedSince returns the old "watched_since" field's value of the Watch entity.
// If the Watch object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *WatchMutation) OldWatchedSince(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldWatchedSince is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldWatchedSince requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldWatchedSince: %w", err)
	}
	return oldValue.WatchedSince, nil
}

// ResetWatchedSince resets all changes to the "watched_since" field.
func (m *WatchMutation) ResetWatchedSince() {
	m.watched_since = nil
}

// SetUserID sets the "user_id" field.
func (m *WatchMutation) SetUserID(s string) {
	m.user_id = &s
}

// UserID returns the value of the "user_id" field in the mutation.
func (m *WatchMutation) UserID() (r string, exists bool) {
	v := m.user_id
	if v == nil {
		return
	}
	return *v, true
}

// OldUserID returns the old "user_id" field's value of the Watch entity.
// If the Watch object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *WatchMutation) OldUserID(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldUserID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldUserID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUserID: %w", err)
	}
	return oldValue.UserID, nil
}

// ResetUserID resets all changes to the "user_id" field.
func (m *WatchMutation) ResetUserID() {
	m.user_id = nil
}

// SetSecurityID sets the "security" edge to the Security entity by id.
func (m *WatchMutation) SetSecurityID(id int64) {
	m.security = &id
}

// ClearSecurity clears the "security" edge to the Security entity.
func (m *WatchMutation) ClearSecurity() {
	m.clearedsecurity = true
}

// SecurityCleared reports if the "security" edge to the Security entity was cleared.
func (m *WatchMutation) SecurityCleared() bool {
	return m.clearedsecurity
}

// SecurityID returns the "security" edge ID in the mutation.
func (m *WatchMutation) SecurityID() (id int64, exists bool) {
	if m.security != nil {
		return *m.security, true
	}
	return
}

// SecurityIDs returns the "security" edge IDs in the mutation.
// Note that IDs always returns len(IDs) <= 1 for unique edges, and you should use
// SecurityID instead. It exists only for internal usage by the builders.
func (m *WatchMutation) SecurityIDs() (ids []int64) {
	if id := m.security; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetSecurity resets all changes to the "security" edge.
func (m *WatchMutation) ResetSecurity() {
	m.security = nil
	m.clearedsecurity = false
}

// Where appends a list predicates to the WatchMutation builder.
func (m *WatchMutation) Where(ps ...predicate.Watch) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the WatchMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *WatchMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Watch, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *WatchMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *WatchMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Watch).
func (m *WatchMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *WatchMutation) Fields() []string {
	fields := make([]string, 0, 2)
	if m.watched_since != nil {
		fields = append(fields, watch.FieldWatchedSince)
	}
	if m.user_id != nil {
		fields = append(fields, watch.FieldUserID)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *WatchMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case watch.FieldWatchedSince:
		return m.WatchedSince()
	case watch.FieldUserID:
		return m.UserID()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *WatchMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case watch.FieldWatchedSince:
		return m.OldWatchedSince(ctx)
	case watch.FieldUserID:
		return m.OldUserID(ctx)
	}
	return nil, fmt.Errorf("unknown Watch field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *WatchMutation) SetField(name string, value ent.Value) error {
	switch name {
	case watch.FieldWatchedSince:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetWatchedSince(v)
		return nil
	case watch.FieldUserID:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUserID(v)
		return nil
	}
	return fmt.Errorf("unknown Watch field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *WatchMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *WatchMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *WatchMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Watch numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *WatchMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *WatchMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *WatchMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Watch nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *WatchMutation) ResetField(name string) error {
	switch name {
	case watch.FieldWatchedSince:
		m.ResetWatchedSince()
		return nil
	case watch.FieldUserID:
		m.ResetUserID()
		return nil
	}
	return fmt.Errorf("unknown Watch field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *WatchMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.security != nil {
		edges = append(edges, watch.EdgeSecurity)
	}
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *WatchMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case watch.EdgeSecurity:
		if id := m.security; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *WatchMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *WatchMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *WatchMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedsecurity {
		edges = append(edges, watch.EdgeSecurity)
	}
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *WatchMutation) EdgeCleared(name string) bool {
	switch name {
	case watch.EdgeSecurity:
		return m.clearedsecurity
	}
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *WatchMutation) ClearEdge(name string) error {
	switch name {
	case watch.EdgeSecurity:
		m.ClearSecurity()
		return nil
	}
	return fmt.Errorf("unknown Watch unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *WatchMutation) ResetEdge(name string) error {
	switch name {
	case watch.EdgeSecurity:
		m.ResetSecurity()
		return nil
	}
	return fmt.Errorf("unknown Watch edge %s", name)
}
