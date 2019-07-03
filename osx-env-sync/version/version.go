/*
Package version contains the command's standard version info.
*/
package version

const (
	application = "osx-env-sync"
	release     = "0.4.0"
)

// Application is the "friendly" name for this code
func Application() string {
	return application
}

// Release is the current version of this code
func Release() string {
	return release
}
