// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/karatekaneen/stockybot/ent/security"
	"github.com/karatekaneen/stockybot/ent/watch"
)

// WatchCreate is the builder for creating a Watch entity.
type WatchCreate struct {
	config
	mutation *WatchMutation
	hooks    []Hook
}

// SetWatchedSince sets the "watched_since" field.
func (wc *WatchCreate) SetWatchedSince(t time.Time) *WatchCreate {
	wc.mutation.SetWatchedSince(t)
	return wc
}

// SetNillableWatchedSince sets the "watched_since" field if the given value is not nil.
func (wc *WatchCreate) SetNillableWatchedSince(t *time.Time) *WatchCreate {
	if t != nil {
		wc.SetWatchedSince(*t)
	}
	return wc
}

// SetUserID sets the "user_id" field.
func (wc *WatchCreate) SetUserID(s string) *WatchCreate {
	wc.mutation.SetUserID(s)
	return wc
}

// SetSecurityID sets the "security" edge to the Security entity by ID.
func (wc *WatchCreate) SetSecurityID(id int64) *WatchCreate {
	wc.mutation.SetSecurityID(id)
	return wc
}

// SetNillableSecurityID sets the "security" edge to the Security entity by ID if the given value is not nil.
func (wc *WatchCreate) SetNillableSecurityID(id *int64) *WatchCreate {
	if id != nil {
		wc = wc.SetSecurityID(*id)
	}
	return wc
}

// SetSecurity sets the "security" edge to the Security entity.
func (wc *WatchCreate) SetSecurity(s *Security) *WatchCreate {
	return wc.SetSecurityID(s.ID)
}

// Mutation returns the WatchMutation object of the builder.
func (wc *WatchCreate) Mutation() *WatchMutation {
	return wc.mutation
}

// Save creates the Watch in the database.
func (wc *WatchCreate) Save(ctx context.Context) (*Watch, error) {
	wc.defaults()
	return withHooks(ctx, wc.sqlSave, wc.mutation, wc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (wc *WatchCreate) SaveX(ctx context.Context) *Watch {
	v, err := wc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (wc *WatchCreate) Exec(ctx context.Context) error {
	_, err := wc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (wc *WatchCreate) ExecX(ctx context.Context) {
	if err := wc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (wc *WatchCreate) defaults() {
	if _, ok := wc.mutation.WatchedSince(); !ok {
		v := watch.DefaultWatchedSince()
		wc.mutation.SetWatchedSince(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (wc *WatchCreate) check() error {
	if _, ok := wc.mutation.WatchedSince(); !ok {
		return &ValidationError{Name: "watched_since", err: errors.New(`ent: missing required field "Watch.watched_since"`)}
	}
	if _, ok := wc.mutation.UserID(); !ok {
		return &ValidationError{Name: "user_id", err: errors.New(`ent: missing required field "Watch.user_id"`)}
	}
	if v, ok := wc.mutation.UserID(); ok {
		if err := watch.UserIDValidator(v); err != nil {
			return &ValidationError{Name: "user_id", err: fmt.Errorf(`ent: validator failed for field "Watch.user_id": %w`, err)}
		}
	}
	return nil
}

func (wc *WatchCreate) sqlSave(ctx context.Context) (*Watch, error) {
	if err := wc.check(); err != nil {
		return nil, err
	}
	_node, _spec := wc.createSpec()
	if err := sqlgraph.CreateNode(ctx, wc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	wc.mutation.id = &_node.ID
	wc.mutation.done = true
	return _node, nil
}

func (wc *WatchCreate) createSpec() (*Watch, *sqlgraph.CreateSpec) {
	var (
		_node = &Watch{config: wc.config}
		_spec = sqlgraph.NewCreateSpec(watch.Table, sqlgraph.NewFieldSpec(watch.FieldID, field.TypeInt))
	)
	if value, ok := wc.mutation.WatchedSince(); ok {
		_spec.SetField(watch.FieldWatchedSince, field.TypeTime, value)
		_node.WatchedSince = value
	}
	if value, ok := wc.mutation.UserID(); ok {
		_spec.SetField(watch.FieldUserID, field.TypeString, value)
		_node.UserID = value
	}
	if nodes := wc.mutation.SecurityIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   watch.SecurityTable,
			Columns: []string{watch.SecurityColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(security.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.security_watchers = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// WatchCreateBulk is the builder for creating many Watch entities in bulk.
type WatchCreateBulk struct {
	config
	err      error
	builders []*WatchCreate
}

// Save creates the Watch entities in the database.
func (wcb *WatchCreateBulk) Save(ctx context.Context) ([]*Watch, error) {
	if wcb.err != nil {
		return nil, wcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(wcb.builders))
	nodes := make([]*Watch, len(wcb.builders))
	mutators := make([]Mutator, len(wcb.builders))
	for i := range wcb.builders {
		func(i int, root context.Context) {
			builder := wcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*WatchMutation)
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
					_, err = mutators[i+1].Mutate(root, wcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, wcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
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
		if _, err := mutators[0].Mutate(ctx, wcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (wcb *WatchCreateBulk) SaveX(ctx context.Context) []*Watch {
	v, err := wcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (wcb *WatchCreateBulk) Exec(ctx context.Context) error {
	_, err := wcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (wcb *WatchCreateBulk) ExecX(ctx context.Context) {
	if err := wcb.Exec(ctx); err != nil {
		panic(err)
	}
}
