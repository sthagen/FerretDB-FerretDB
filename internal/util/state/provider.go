// Copyright 2021 FerretDB Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package state

import (
	"encoding/json"
	"expvar"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/FerretDB/FerretDB/v2/internal/util/iface"
	"github.com/FerretDB/FerretDB/v2/internal/util/must"
)

// Provider provides access to FerretDB process state.
type Provider struct {
	filename string

	rw   sync.RWMutex
	s    *State
	subs map[chan struct{}]struct{}
}

// NewProvider creates a new Provider that stores state in the given file
// (that will be created automatically if needed).
//
// If filename is empty, then the state is not persisted.
//
// All provider's methods are thread-safe.
func NewProvider(filename string) (*Provider, error) {
	p := &Provider{
		filename: filename,
		s:        new(State),
		subs:     make(map[chan struct{}]struct{}, 1),
	}

	if p.filename != "" {
		b, _ := os.ReadFile(p.filename)
		_ = json.Unmarshal(b, p.s)
	}

	p.s.fill()

	// Simply overwrite state to handle all errors and edge cases
	// like missing directory, corrupted file, invalid UUID, etc.,
	// and also to check permissions.
	if err := persistState(p.s, p.filename); err != nil {
		return p, fmt.Errorf("failed to persist state: %w", err)
	}

	return p, nil
}

// NewProviderDir creates a new Provider that stores state in the state.json file in the given directory
// (that will be created automatically if needed).
func NewProviderDir(dir string) (*Provider, error) {
	if dir == "" {
		return nil, fmt.Errorf("state directory is not set")
	}

	f, err := filepath.Abs(filepath.Join(dir, "state.json"))
	if err != nil {
		return nil, err
	}

	sp, err := NewProvider(f)
	if err != nil {
		return nil, newProviderDirErr(f, err)
	}

	return sp, nil
}

// Var returns an unpublished [expvar.Var] for the state.
func (p *Provider) Var() expvar.Var {
	return iface.Stringer(func() string {
		b := must.NotFail(json.Marshal(p.Get().asMap()))
		return string(b)
	})
}

// MetricsCollector returns Prometheus metrics collector for that provider.
//
// If addUUID is true, then the "uuid" label is added.
func (p *Provider) MetricsCollector(addUUID bool) prometheus.Collector {
	return newMetricsCollector(p, addUUID)
}

// Get returns a copy of the current process state.
//
// It is okay to call this function often.
// The caller should not cache result; Provider does everything needed itself.
func (p *Provider) Get() *State {
	p.rw.RLock()
	defer p.rw.RUnlock()

	return p.s.deepCopy()
}

// Subscribe returns a channel that would receive notifications on state changes.
// One notification would be scheduled immediately.
func (p *Provider) Subscribe() chan struct{} {
	p.rw.Lock()
	defer p.rw.Unlock()

	ch := make(chan struct{}, 1)
	ch <- struct{}{}

	p.subs[ch] = struct{}{}

	return ch
}

// Update gets the current state, calls the given function, updates state, and notifies all subscribers.
func (p *Provider) Update(update func(s *State)) error {
	p.rw.Lock()
	defer p.rw.Unlock()

	update(p.s)
	p.s = p.s.deepCopy()
	p.s.fill()

	err := persistState(p.s, p.filename)
	if err != nil {
		err = fmt.Errorf("failed to persist state: %w", err)
	}

	// skip subscribers that already have notification waiting for them
	for ch := range p.subs {
		select {
		case ch <- struct{}{}:
		default:
		}
	}

	return err
}

// persistState saves state to the given file without modifying (filling) it.
//
// It exist immediately if filename is empty.
func persistState(s *State, filename string) error {
	if filename == "" {
		return nil
	}

	b, err := json.Marshal(s)

	if err == nil {
		_ = os.MkdirAll(filepath.Dir(filename), 0o777)
		err = os.WriteFile(filename, b, 0o666)
	}

	return err
}
