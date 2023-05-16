// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"fmt"
	"sync"

	"github.com/open-component-model/ocm/pkg/errors"
	"github.com/open-component-model/ocm/pkg/runtime"
	"github.com/open-component-model/ocm/pkg/utils"
)

const (
	KIND_ACTION     = "action"
	KIND_ACTIONTYPE = "action type"
)

type ActionTypeRegistry interface {
	RegisterAction(name string, description string, usage string, attrs []string) error
	RegisterActionType(typ ActionType) error

	DecodeActionSpec(data []byte, unmarshaler runtime.Unmarshaler) (ActionSpec, error)
	EncodeActionSpec(spec ActionSpec, marshaler runtime.Marshaler) ([]byte, error)

	DecodeActionResult(data []byte, unmarshaler runtime.Unmarshaler) (ActionResult, error)
	EncodeActionResult(spec ActionResult, marshaler runtime.Marshaler) ([]byte, error)

	GetAction(name string) Action
	SupportedActionVersions(name string) []string

	Copy() ActionTypeRegistry
}

type action struct {
	name       string
	shortdesc  string
	usage      string
	attributes []string
	types      map[string]ActionType
}

var _ Action = (*action)(nil)

func (a *action) Name() string {
	return a.name
}

func (a *action) Description() string {
	return a.shortdesc
}

func (a *action) Usage() string {
	return a.usage
}

func (a *action) ConsumerAttributes() []string {
	return a.attributes
}

type actionRegistry struct {
	lock        sync.Mutex
	actions     map[string]*action
	actionspecs runtime.TypeScheme[ActionSpec, ActionSpecType]
	resultspecs runtime.TypeScheme[ActionResult, ActionResultType]
}

func NewActionTypeRegistry() ActionTypeRegistry {
	return &actionRegistry{
		actions:     map[string]*action{},
		actionspecs: runtime.NewTypeScheme[ActionSpec, ActionSpecType](),
		resultspecs: runtime.NewTypeScheme[ActionResult, ActionResultType](),
	}
}

func (r *actionRegistry) Copy() ActionTypeRegistry {
	r.lock.Lock()
	defer r.lock.Unlock()

	actions := map[string]*action{}

	for k, v := range r.actions {
		a := *v
		a.types = map[string]ActionType{}
		for _, t := range v.types {
			a.types[t.GetType()] = t
		}
		actions[k] = &a
	}
	actionspecs := runtime.NewTypeScheme[ActionSpec, ActionSpecType]()
	actionspecs.AddKnownTypes(r.actionspecs)
	resultspecs := runtime.NewTypeScheme[ActionResult, ActionResultType]()
	resultspecs.AddKnownTypes(r.resultspecs)
	return &actionRegistry{
		actions:     actions,
		actionspecs: actionspecs,
		resultspecs: resultspecs,
	}
}

func (r *actionRegistry) RegisterAction(name string, description string, usage string, attrs []string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	ai := r.actions[name]
	if ai != nil {
		return errors.ErrAlreadyExists(KIND_ACTION, name)
	}

	ai = &action{
		name:       name,
		shortdesc:  description,
		usage:      usage,
		attributes: append(attrs[:0:0], attrs...),
		types:      map[string]ActionType{},
	}
	r.actions[name] = ai
	return nil
}

func (r *actionRegistry) RegisterActionType(typ ActionType) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	k := typ.GetKind()

	ai := r.actions[k]
	if ai == nil {
		return errors.ErrNotFound(KIND_ACTION, k)
	}

	if typ.SpecificationType().GetType() != typ.ResultType().GetType() {
		return errors.ErrInvalidWrap(fmt.Errorf("version mismatch: request[%s]!=result[%s]", typ.SpecificationType().GetType(), typ.ResultType().GetType()), KIND_ACTIONTYPE, k)
	}
	if typ.SpecificationType().GetKind() != k {
		return errors.ErrInvalidWrap(fmt.Errorf("kind mismatch in types: %s", typ.SpecificationType().GetType()), KIND_ACTIONTYPE, k)
	}
	ai.types[typ.GetVersion()] = typ
	r.actionspecs.Register(typ.SpecificationType())
	r.resultspecs.Register(typ.ResultType())
	return nil
}

func (r *actionRegistry) GetAction(name string) Action {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.actions[name]
}

func (r *actionRegistry) DecodeActionSpec(data []byte, unmarshaler runtime.Unmarshaler) (ActionSpec, error) {
	return r.actionspecs.Decode(data, unmarshaler)
}

func (r *actionRegistry) DecodeActionResult(data []byte, unmarshaler runtime.Unmarshaler) (ActionResult, error) {
	return r.resultspecs.Decode(data, unmarshaler)
}

func (r *actionRegistry) EncodeActionSpec(spec ActionSpec, marshaler runtime.Marshaler) ([]byte, error) {
	return r.actionspecs.Encode(spec, marshaler)
}

func (r *actionRegistry) EncodeActionResult(spec ActionResult, marshaler runtime.Marshaler) ([]byte, error) {
	return r.resultspecs.Encode(spec, marshaler)
}

func (r *actionRegistry) SupportedActionVersions(name string) []string {
	r.lock.Lock()
	defer r.lock.Unlock()
	a := r.actions[name]
	if a == nil {
		return nil
	}
	return utils.StringMapKeys(a.types)
}

////////////////////////////////////////////////////////////////////////////////

var registry = NewActionTypeRegistry()

func DefaultRegistry() ActionTypeRegistry {
	return registry
}

func RegisterAction(name string, description string, usage string, attrs []string) error {
	return registry.RegisterAction(name, description, usage, attrs)
}

func RegisterType(typ ActionType) error {
	return registry.RegisterActionType(typ)
}

func GetAction(name string) Action {
	return registry.GetAction(name)
}

func DecodeActionSpec(data []byte, unmarshaler runtime.Unmarshaler) (ActionSpec, error) {
	return registry.DecodeActionSpec(data, unmarshaler)
}

func EncodeActionSpec(spec ActionSpec, marshaler runtime.Marshaler) ([]byte, error) {
	return registry.EncodeActionSpec(spec, marshaler)
}

func DecodeActionResult(data []byte, unmarshaler runtime.Unmarshaler) (ActionResult, error) {
	return registry.DecodeActionResult(data, unmarshaler)
}

func EncodeActionResult(spec ActionResult, marshaler runtime.Marshaler) ([]byte, error) {
	return registry.EncodeActionResult(spec, marshaler)
}

func SupportedActionVersions(name string) []string {
	return registry.SupportedActionVersions(name)
}
