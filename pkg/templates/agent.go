package templates

// GetAgentMountScript contains the mount script for an agent.
func GetAgentMountScript() string {
	return `#!/bin/sh

set -eux

mkdir -p $1 && cd $1
shift

exec $*
`
}
