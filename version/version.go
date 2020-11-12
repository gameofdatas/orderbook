package version

import (
	"fmt"
	"runtime"
)

// GitCommit that was compiled. This will be filled in by the compiler.
var GitCommit string

// Version number that is being run at the moment.
const Version = "0.1.0"

// BuildDate is the date when code was build
var BuildDate = ""

// GoVersion is the version of go used by user
var GoVersion = runtime.Version()

// OsArch is the OS information
var OsArch = fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)
