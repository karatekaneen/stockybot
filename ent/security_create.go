// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/karatekaneen/stockybot/ent/security"
	"github.com/karatekaneen/stockybot/ent/watch"
)

// SecurityCreate is the builder for creating a Security entity.
type SecurityCreate struct {
	config
	mutation *SecurityMutation
	hooks    []Hook
}

// SetName sets the "name" field.
func (sc *SecurityCreate) SetName(s string) *SecurityCreate {
	sc.mutation.SetName(s)
	return sc
}

// SetCountry sets the "country" field.
func (sc *SecurityCreate) SetCountry(s string) *SecurityCreate {
	sc.mutation.SetCountry(s)
	return sc
}

// SetNillableCountry sets the "country" field if the given value is not nil.
func (sc *SecurityCreate) SetNillableCountry(s *string) *SecurityCreate {
	if s != nil {
		sc.SetCountry(*s)
	}
	return sc
}

// SetLinkName sets the "link_name" field.
func (sc *SecurityCreate) SetLinkName(s string) *SecurityCreate {
	sc.mutation.SetLinkName(s)
	return sc
}

// SetNillableLinkName sets the "link_name" field if the given value is not nil.
func (sc *SecurityCreate) SetNillableLinkName(s *string) *SecurityCreate {
	if s != nil {
		sc.SetLinkName(*s)
	}
	return sc
}

// SetList sets the "list" field.
func (sc *SecurityCreate) SetList(s string) *SecurityCreate {
	sc.mutation.SetList(s)
	return sc
}

// SetNillableList sets the "list" field if the given value is not nil.
func (sc *SecurityCreate) SetNillableList(s *string) *SecurityCreate {
	if s != nil {
		sc.SetList(*s)
	}
	return sc
}

// SetType sets the "type" field.
func (sc *SecurityCreate) SetType(s security.Type) *SecurityCreate {
	sc.mutation.SetType(s)
	return sc
}

// SetID sets the "id" field.
func (sc *SecurityCreate) SetID(i int64) *SecurityCreate {
	sc.mutation.SetID(i)
	return sc
}

// AddWatcherIDs adds the "watchers" edge to the Watch entity by IDs.
func (sc *SecurityCreate) AddWatcherIDs(ids ...int) *SecurityCreate {
	sc.mutation.AddWatcherIDs(ids...)
	return sc
}

// AddWatchers adds the "watchers" edges to the Watch entity.
func (sc *SecurityCreate) AddWatchers(w ...*Watch) *SecurityCreate {
	ids := make([]int, len(w))
	for i := range w {
		ids[i] = w[i].ID
	}
	return sc.AddWatcherIDs(ids...)
}

// Mutation returns the SecurityMutation object of the builder.
func (sc *SecurityCreate) Mutation() *SecurityMutation {
	return sc.mutation
}

// Save creates the Security in the database.
func (sc *SecurityCreate) Save(ctx context.Context) (*Security, error) {
	sc.defaults()
	return withHooks(ctx, sc.sqlSave, sc.mutation, sc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (sc *SecurityCreate) SaveX(ctx context.Context) *Security {
	v, err := sc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sc *SecurityCreate) Exec(ctx context.Context) error {
	_, err := sc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sc *SecurityCreate) ExecX(ctx context.Context) {
	if err := sc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (sc *SecurityCreate) defaults() {
	if _, ok := sc.mutation.Country(); !ok {
		v := security.DefaultCountry
		sc.mutation.SetCountry(v)
	}
	if _, ok := sc.mutation.List(); !ok {
		v := security.DefaultList
		sc.mutation.SetList(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (sc *SecurityCreate) check() error {
	if _, ok := sc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Security.name"`)}
	}
	if v, ok := sc.mutation.Name(); ok {
		if err := security.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Security.name": %w`, err)}
		}
	}
	if _, ok := sc.mutation.Country(); !ok {
		return &ValidationError{Name: "country", err: errors.New(`ent: missing required field "Security.country"`)}
	}
	if _, ok := sc.mutation.List(); !ok {
		return &ValidationError{Name: "list", err: errors.New(`ent: missing required field "Security.list"`)}
	}
	if _, ok := sc.mutation.GetType(); !ok {
		return &ValidationError{Name: "type", err: errors.New(`ent: missing required field "Security.type"`)}
	}
	if v, ok := sc.mutation.GetType(); ok {
		if err := security.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "Security.type": %w`, err)}
		}
	}
	if v, ok := sc.mutation.ID(); ok {
		if err := security.IDValidator(v); err != nil {
			return &ValidationError{Name: "id", err: fmt.Errorf(`ent: validator failed for field "Security.id": %w`, err)}
		}
	}
	return nil
}

func (sc *SecurityCreate) sqlSave(ctx context.Context) (*Security, error) {
	if err := sc.check(); err != nil {
		return nil, err
	}
	_node, _spec := sc.createSpec()
	if err := sqlgraph.CreateNode(ctx, sc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != _node.ID {
		id := _spec.ID.Value.(int64)
		_node.ID = int64(id)
	}
	sc.mutation.id = &_node.ID
	sc.mutation.done = true
	return _node, nil
}

func (sc *SecurityCreate) createSpec() (*Security, *sqlgraph.CreateSpec) {
	var (
		_node = &Security{config: sc.config}
		_spec = sqlgraph.NewCreateSpec(security.Table, sqlgraph.NewFieldSpec(security.FieldID, field.TypeInt64))
	)
	if id, ok := sc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := sc.mutation.Name(); ok {
		_spec.SetField(security.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := sc.mutation.Country(); ok {
		_spec.SetField(security.FieldCountry, field.TypeString, value)
		_node.Country = value
	}
	if value, ok := sc.mutation.LinkName(); ok {
		_spec.SetField(security.FieldLinkName, field.TypeString, value)
		_node.LinkName = value
	}
	if value, ok := sc.mutation.List(); ok {
		_spec.SetField(security.FieldList, field.TypeString, value)
		_node.List = value
	}
	if value, ok := sc.mutation.GetType(); ok {
		_spec.SetField(security.FieldType, field.TypeEnum, value)
		_node.Type = value
	}
	if nodes := sc.mutation.WatchersIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// SecurityCreateBulk is the builder for creating many Security entities in bulk.
type SecurityCreateBulk struct {
	config
	err      error
	builders []*SecurityCreate
}

// Save creates the Security entities in the database.
func (scb *SecurityCreateBulk) Save(ctx context.Context) ([]*Security, error) {
	if scb.err != nil {
		return nil, scb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(scb.builders))
	nodes := make([]*Security, len(scb.builders))
	mutators := make([]Mutator, len(scb.builders))
	for i := range scb.builders {
		func(i int, root context.Context) {
			builder := scb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*SecurityMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, scb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, scb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil && nodes[i].ID == 0 {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int64(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, scb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (scb *SecurityCreateBulk) SaveX(ctx context.Context) []*Security {
	v, err := scb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (scb *SecurityCreateBulk) Exec(ctx context.Context) error {
	_, err := scb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (scb *SecurityCreateBulk) ExecX(ctx context.Context) {
	if err := scb.Exec(ctx); err != nil {
		panic(err)
	}
}
