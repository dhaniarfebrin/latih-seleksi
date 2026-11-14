package spec

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
)

//go:embed data/*.json
var dataFS embed.FS

// Registry holds every embedded spec, already parsed and validated.
type Registry struct {
	byID map[string]Spec
	ids  []string
}

// Load parses and validates every embedded spec. Any invalid file or
// duplicate id is an error.
func Load() (*Registry, error) {
	return loadFS(dataFS, "data")
}

func loadFS(fsys fs.FS, dir string) (*Registry, error) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, fmt.Errorf("spec: read %s: %w", dir, err)
	}
	r := &Registry{byID: map[string]Spec{}}
	for _, e := range entries {
		name := dir + "/" + e.Name()
		raw, err := fs.ReadFile(fsys, name)
		if err != nil {
			return nil, fmt.Errorf("spec: read %s: %w", name, err)
		}
		s, err := Parse(raw)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		if err := s.Validate(); err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		if _, dup := r.byID[s.ID]; dup {
			return nil, fmt.Errorf("%s: duplicate spec id %q", name, s.ID)
		}
		r.byID[s.ID] = s
		r.ids = append(r.ids, s.ID)
	}
	sort.Strings(r.ids)
	return r, nil
}

// Get returns the spec with the given id.
func (r *Registry) Get(id string) (Spec, bool) {
	s, ok := r.byID[id]
	return s, ok
}

// All returns every spec ordered by id.
func (r *Registry) All() []Spec {
	out := make([]Spec, 0, len(r.ids))
	for _, id := range r.ids {
		out = append(out, r.byID[id])
	}
	return out
}
