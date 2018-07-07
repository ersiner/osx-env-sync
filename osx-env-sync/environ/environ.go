package environ

import (
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type Environ map[string]string

func New() Environ {
	return make(Environ)
}

func (s Environ) AddOsEnviron() Environ {
	return s.AddEnviron(os.Environ())
}

func (s Environ) AddEnviron(environ []string) Environ {
	for _, envLine := range environ {
		kv := strings.SplitN(envLine, "=", 2)
		s[kv[0]] = kv[1]
	}

	log.WithField("Env", s).Debug("Environ imported")
	return s
}

func (s Environ) NewFromFiltered(names []string) Environ {
	newEnv := make(Environ)

	for _, name := range names {
		newEnv[name] = s[name]
	}

	log.WithField("Filtered Env", newEnv).Debug("New filtered environ created")
	return newEnv
}

func (s Environ) ToOsEnviron() []string {
	newEnv := make([]string, 0)

	for name, val := range s {
		line := fmt.Sprintf("%s=%s", name, val)
		newEnv = append(newEnv, line)
	}

	log.WithField("os.Environ", newEnv).Debug("os.Environ created")
	return newEnv
}
