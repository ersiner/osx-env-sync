/*
Package environ contains methods for parsing and manipulating various environ settings.
*/
package environ

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Environ type holds a Map of EnvVar=Val KV-pairs.
type Environ map[string]string

// New creates a new Environ struct
func New() Environ {
	return make(Environ)
}

// NewFromFiltered creates a new Environ struct, and then filters-out all except the Env-Vars provided in the "names" argument.
func (s Environ) NewFromFiltered(names []string) Environ {
	newEnv := New()

	for _, name := range names {
		newEnv[name] = s[name]
	}

	log.WithField("Filtered Env", newEnv).Debug("New filtered environ created")
	return newEnv
}

// AddEnviron parses an array of "EnvVar=Val" strings into KV-pairs, and adds them to the Environ struct.
func (s Environ) AddEnviron(environ []string) Environ {
	for _, envLine := range environ {
		kv := strings.SplitN(envLine, "=", 2)
		s[kv[0]] = kv[1]
	}

	log.WithField("Env", s).Debug("Environ imported")
	return s
}

// AddOsEnviron parses then adds the contents of "os.Environ()" to the Environ struct.
func (s Environ) AddOsEnviron() Environ {
	return s.AddEnviron(os.Environ())
}

// ToOsEnviron converts the Environ struct into and array of "EnvVar=Val" strings, suitable for passing to "exec.Command()".
func (s Environ) ToOsEnviron() []string {
	newEnv := make([]string, 0)

	for name, val := range s {
		line := fmt.Sprintf("%s=%s", name, val)
		newEnv = append(newEnv, line)
	}

	log.WithField("os.Environ", newEnv).Debug("os.Environ created")
	return newEnv
}
