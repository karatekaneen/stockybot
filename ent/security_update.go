// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/karatekaneen/stockybot/ent/predicate"
	"github.com/karatekaneen/stockybot/ent/security"
	"github.com/karatekaneen/stockybot/ent/watch"
)

// SecurityUpdate is the builder for updating Security entities.
type SecurityUpdate struct {
	config
	hooks    []Hook
	mutation *SecurityMutation
}

// Where appends a list predicates to the SecurityUpdate builder.
func (su *SecurityUpdate) Where(ps ...predicate.Security) *SecurityUpdate {
	su.mutation.Where(ps...)
	return su
}

// SetName sets the "name" field.
func (su *SecurityUpdate) SetName(s string) *SecurityUpdate {
	su.mutation.SetName(s)
	return su
}

// SetNillableName sets the "name" field if the given value is not nil.
func (su *SecurityUpdate) SetNillableName(s *string) *SecurityUpdate {
	if s != nil {
		su.SetName(*s)
	}
	return su
}

// SetCountry sets the "country" field.
func (su *SecurityUpdate) SetCountry(s string) *SecurityUpdate {
	su.mutation.SetCountry(s)
	return su
}

// SetNillableCountry sets the "country" field if the given value is not nil.
func (su *SecurityUpdate) SetNillableCountry(s *string) *SecurityUpdate {
	if s != nil {
		su.SetCountry(*s)
	}
	return su
}

// SetLinkName sets the "link_name" field.
func (su *SecurityUpdate) SetLinkName(s string) *SecurityUpdate {
	su.mutation.SetLinkName(s)
	return su
}

// SetNillableLinkName sets the "link_name" field if the given value is not nil.
func (su *SecurityUpdate) SetNillableLinkName(s *string) *SecurityUpdate {
	if s != nil {
		su.SetLinkName(*s)
	}
	return su
}

// ClearLinkName clears the value of the "link_name" field.
func (su *SecurityUpdate) ClearLinkName() *SecurityUpdate {
	su.mutation.ClearLinkName()
	return su
}

// SetList sets the "list" field.
func (su *SecurityUpdate) SetList(s string) *SecurityUpdate {
	su.mutation.SetList(s)
	return su
}

// SetNillableList sets the "list" field if the given value is not nil.
func (su *SecurityUpdate) SetNillableList(s *string) *SecurityUpdate {
	if s != nil {
		su.SetList(*s)
	}
	return su
}

// SetType sets the "type" field.
func (su *SecurityUpdate) SetType(s security.Type) *SecurityUpdate {
	su.mutation.SetType(s)
	return su
}

// SetNillableType sets the "type" field if the given value is not nil.
func (su *SecurityUpdate) SetNillableType(s *security.Type) *SecurityUpdate {
	if s != nil {
		su.SetType(*s)
	}
	return su
}

// AddWatcherIDs adds the "watchers" edge to the Watch entity by IDs.
func (su *SecurityUpdate) AddWatcherIDs(ids ...int) *SecurityUpdate {
	su.mutation.AddWatcherIDs(ids...)
	return su
}

// AddWatchers adds the "watchers" edges to the Watch entity.
func (su *SecurityUpdate) AddWatchers(w ...*Watch) *SecurityUpdate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return su.AddWatcherIDs(ids...)
}

// Mutation returns the SecurityMutation object of the builder.
func (su *SecurityUpdate) Mutation() *SecurityMutation {
	return su.mutation
}

// ClearWatchers clears all "watchers" edges to the Watch entity.
func (su *SecurityUpdate) ClearWatchers() *SecurityUpdate {
	su.mutation.ClearWatchers()
	return su
}

// RemoveWatcherIDs removes the "watchers" edge to Watch entities by IDs.
func (su *SecurityUpdate) RemoveWatcherIDs(ids ...int) *SecurityUpdate {
	su.mutation.RemoveWatcherIDs(ids...)
	return su
}

// RemoveWatchers removes "watchers" edges to Watch entities.
func (su *SecurityUpdate) RemoveWatchers(w ...*Watch) *SecurityUpdate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return su.RemoveWatcherIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (su *SecurityUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, su.sqlSave, su.mutation, su.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (su *SecurityUpdate) SaveX(ctx context.Context) int {
	affected, err := su.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (su *SecurityUpdate) Exec(ctx context.Context) error {
	_, err := su.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (su *SecurityUpdate) ExecX(ctx context.Context) {
	if err := su.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (su *SecurityUpdate) check() error {
	if v, ok := su.mutation.Name(); ok {
		if err := security.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Security.name": %w`, err)}
		}
	}
	if v, ok := su.mutation.Country(); ok {
		if err := security.CountryValidator(v); err != nil {
			return &ValidationError{Name: "country", err: fmt.Errorf(`ent: validator failed for field "Security.country": %w`, err)}
		}
	}
	if v, ok := su.mutation.List(); ok {
		if err := security.ListValidator(v); err != nil {
			return &ValidationError{Name: "list", err: fmt.Errorf(`ent: validator failed for field "Security.list": %w`, err)}
		}
	}
	if v, ok := su.mutation.GetType(); ok {
		if err := security.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "Security.type": %w`, err)}
		}
	}
	return nil
}

func (su *SecurityUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := su.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(security.Table, security.Columns, sqlgraph.NewFieldSpec(security.FieldID, field.TypeInt64))
	if ps := su.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := su.mutation.Name(); ok {
		_spec.SetField(security.FieldName, field.TypeString, value)
	}
	if value, ok := su.mutation.Country(); ok {
		_spec.SetField(security.FieldCountry, field.TypeString, value)
	}
	if value, ok := su.mutation.LinkName(); ok {
		_spec.SetField(security.FieldLinkName, field.TypeString, value)
	}
	if su.mutation.LinkNameCleared() {
		_spec.ClearField(security.FieldLinkName, field.TypeString)
	}
	if value, ok := su.mutation.List(); ok {
		_spec.SetField(security.FieldList, field.TypeString, value)
	}
	if value, ok := su.mutation.GetType(); ok {
		_spec.SetField(security.FieldType, field.TypeEnum, value)
	}
	if su.mutation.WatchersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   security.WatchersTable,
			Columns: []string{security.WatchersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(watch.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.RemovedWatchersIDs(); len(nodes) > 0 && !su.mutation.WatchersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   security.WatchersTable,
			Columns: []string{security.WatchersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(watch.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.WatchersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   security.WatchersTable,
			Columns: []string{security.WatchersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(watch.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, su.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{security.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	su.mutation.done = true
	return n, nil
}

// SecurityUpdateOne is the builder for updating a single Security entity.
type SecurityUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *SecurityMutation
}

// SetName sets the "name" field.
func (suo *SecurityUpdateOne) SetName(s string) *SecurityUpdateOne {
	suo.mutation.SetName(s)
	return suo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (suo *SecurityUpdateOne) SetNillableName(s *string) *SecurityUpdateOne {
	if s != nil {
		suo.SetName(*s)
	}
	return suo
}

// SetCountry sets the "country" field.
func (suo *SecurityUpdateOne) SetCountry(s string) *SecurityUpdateOne {
	suo.mutation.SetCountry(s)
	return suo
}

// SetNillableCountry sets the "country" field if the given value is not nil.
func (suo *SecurityUpdateOne) SetNillableCountry(s *string) *SecurityUpdateOne {
	if s != nil {
		suo.SetCountry(*s)
	}
	return suo
}

// SetLinkName sets the "link_name" field.
func (suo *SecurityUpdateOne) SetLinkName(s string) *SecurityUpdateOne {
	suo.mutation.SetLinkName(s)
	return suo
}

// SetNillableLinkName sets the "link_name" field if the given value is not nil.
func (suo *SecurityUpdateOne) SetNillableLinkName(s *string) *SecurityUpdateOne {
	if s != nil {
		suo.SetLinkName(*s)
	}
	return suo
}

// ClearLinkName clears the value of the "link_name" field.
func (suo *SecurityUpdateOne) ClearLinkName() *SecurityUpdateOne {
	suo.mutation.ClearLinkName()
	return suo
}

// SetList sets the "list" field.
func (suo *SecurityUpdateOne) SetList(s string) *SecurityUpdateOne {
	suo.mutation.SetList(s)
	return suo
}

// SetNillableList sets the "list" field if the given value is not nil.
func (suo *SecurityUpdateOne) SetNillableList(s *string) *SecurityUpdateOne {
	if s != nil {
		suo.SetList(*s)
	}
	return suo
}

// SetType sets the "type" field.
func (suo *SecurityUpdateOne) SetType(s security.Type) *SecurityUpdateOne {
	suo.mutation.SetType(s)
	return suo
}

// SetNillableType sets the "type" field if the given value is not nil.
func (suo *SecurityUpdateOne) SetNillableType(s *security.Type) *SecurityUpdateOne {
	if s != nil {
		suo.SetType(*s)
	}
	return suo
}

// AddWatcherIDs adds the "watchers" edge to the Watch entity by IDs.
func (suo *SecurityUpdateOne) AddWatcherIDs(ids ...int) *SecurityUpdateOne {
	suo.mutation.AddWatcherIDs(ids...)
	return suo
}

// AddWatchers adds the "watchers" edges to the Watch entity.
func (suo *SecurityUpdateOne) AddWatchers(w ...*Watch) *SecurityUpdateOne {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return suo.AddWatcherIDs(ids...)
}

// Mutation returns the SecurityMutation object of the builder.
func (suo *SecurityUpdateOne) Mutation() *SecurityMutation {
	return suo.mutation
}

// ClearWatchers clears all "watchers" edges to the Watch entity.
func (suo *SecurityUpdateOne) ClearWatchers() *SecurityUpdateOne {
	suo.mutation.ClearWatchers()
	return suo
}

// RemoveWatcherIDs removes the "watchers" edge to Watch entities by IDs.
func (suo *SecurityUpdateOne) RemoveWatcherIDs(ids ...int) *SecurityUpdateOne {
	suo.mutation.RemoveWatcherIDs(ids...)
	return suo
}

// RemoveWatchers removes "watchers" edges to Watch entities.
func (suo *SecurityUpdateOne) RemoveWatchers(w ...*Watch) *SecurityUpdateOne {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return suo.RemoveWatcherIDs(ids...)
}

// Where appends a list predicates to the SecurityUpdate builder.
func (suo *SecurityUpdateOne) Where(ps ...predicate.Security) *SecurityUpdateOne {
	suo.mutation.Where(ps...)
	return suo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (suo *SecurityUpdateOne) Select(field string, fields ...string) *SecurityUpdateOne {
	suo.fields = append([]string{field}, fields...)
	return suo
}

// Save executes the query and returns the updated Security entity.
func (suo *SecurityUpdateOne) Save(ctx context.Context) (*Security, error) {
	return withHooks(ctx, suo.sqlSave, suo.mutation, suo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (suo *SecurityUpdateOne) SaveX(ctx context.Context) *Security {
	node, err := suo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (suo *SecurityUpdateOne) Exec(ctx context.Context) error {
	_, err := suo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (suo *SecurityUpdateOne) ExecX(ctx context.Context) {
	if err := suo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (suo *SecurityUpdateOne) check() error {
	if v, ok := suo.mutation.Name(); ok {
		if err := security.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Security.name": %w`, err)}
		}
	}
	if v, ok := suo.mutation.Country(); ok {
		if err := security.CountryValidator(v); err != nil {
			return &ValidationError{Name: "country", err: fmt.Errorf(`ent: validator failed for field "Security.country": %w`, err)}
		}
	}
	if v, ok := suo.mutation.List(); ok {
		if err := security.ListValidator(v); err != nil {
			return &ValidationError{Name: "list", err: fmt.Errorf(`ent: validator failed for field "Security.list": %w`, err)}
		}
	}
	if v, ok := suo.mutation.GetType(); ok {
		if err := security.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "Security.type": %w`, err)}
		}
	}
	return nil
}

func (suo *SecurityUpdateOne) sqlSave(ctx context.Context) (_node *Security, err error) {
	if err := suo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(security.Table, security.Columns, sqlgraph.NewFieldSpec(security.FieldID, field.TypeInt64))
	id, ok := suo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Security.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := suo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, security.FieldID)
		for _, f := range fields {
			if !security.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != security.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := suo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := suo.mutation.Name(); ok {
		_spec.SetField(security.FieldName, field.TypeString, value)
	}
	if value, ok := suo.mutation.Country(); ok {
		_spec.SetField(security.FieldCountry, field.TypeString, value)
	}
	if value, ok := suo.mutation.LinkName(); ok {
		_spec.SetField(security.FieldLinkName, field.TypeString, value)
	}
	if suo.mutation.LinkNameCleared() {
		_spec.ClearField(security.FieldLinkName, field.TypeString)
	}
	if value, ok := suo.mutation.List(); ok {
		_spec.SetField(security.FieldList, field.TypeString, value)
	}
	if value, ok := suo.mutation.GetType(); ok {
		_spec.SetField(security.FieldType, field.TypeEnum, value)
	}
	if suo.mutation.WatchersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   security.WatchersTable,
			Columns: []string{security.WatchersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(watch.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.RemovedWatchersIDs(); len(nodes) > 0 && !suo.mutation.WatchersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   security.WatchersTable,
			Columns: []string{security.WatchersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(watch.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.WatchersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   security.WatchersTable,
			Columns: []string{security.WatchersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(watch.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Security{config: suo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, suo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{security.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	suo.mutation.done = true
	return _node, nil
}