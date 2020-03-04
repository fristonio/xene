package version

// Build information. Populated at build-time by the build script.
var (
	Version   string
	Revision  string
	Branch    string
	BuildUser string
	BuildDate string
	GoVersion string
)

// Info provides the iterable version information.
var Info = map[string]string{
	"version":   Version,
	"revision":  Revision,
	"branch":    Branch,
	"buildUser": BuildUser,
	"buildDate": BuildDate,
	"goVersion": GoVersion,
}

var VersionStr = `
$$\   $$\  $$$$$$\  $$$$$$$\   $$$$$$\
\$$\ $$  |$$  __$$\ $$  __$$\ $$  __$$\
 \$$$$  / $$$$$$$$ |$$ |  $$ |$$$$$$$$ |
 $$  $$<  $$   ____|$$ |  $$ |$$   ____|
$$  /\$$\ \$$$$$$$\ $$ |  $$ |\$$$$$$$\
\__/  \__| \_______|\__|  \__| \_______|


    Version    : %s
    Revision   : %s
    Branch     : %s
    Build-User : %s
    Build-Date : %s
    Go-Version : %s
`
