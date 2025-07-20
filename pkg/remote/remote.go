package remote

import (
	"fmt"
	"sync"

	"github.com/abc1763613206/nabili/pkg/dbif"
)

// RemoteManager manages remote API sources
var (
	remoteSources map[string]dbif.DB
	sourcesMutex  sync.RWMutex
)

func init() {
	remoteSources = make(map[string]dbif.DB)
	RegisterRemoteSource(NewBiliSource())
}

// RegisterRemoteSource registers a remote API source
func RegisterRemoteSource(source dbif.DB) {
	sourcesMutex.Lock()
	defer sourcesMutex.Unlock()
	remoteSources[source.Name()] = source
}

// GetRemoteSource returns a remote source by name
func GetRemoteSource(name string) (dbif.DB, bool) {
	sourcesMutex.RLock()
	defer sourcesMutex.RUnlock()
	source, exists := remoteSources[name]
	return source, exists
}

// ListRemoteSources returns all available remote sources
func ListRemoteSources() []string {
	sourcesMutex.RLock()
	defer sourcesMutex.RUnlock()
	sources := make([]string, 0, len(remoteSources))
	for name := range remoteSources {
		sources = append(sources, name)
	}
	return sources
}

// RemoteResult wraps remote API results

type RemoteResult struct {
	SourceName string
	Result     fmt.Stringer
}

func (r *RemoteResult) String() string {
	return r.Result.String()
}