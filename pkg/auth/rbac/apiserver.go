package rbac

import (
	"regexp"
	"strings"
)

type regex struct {
	Error error
	Regex *regexp.Regexp
}

func getRegex(pattern string) regex {
	r, err := regexp.Compile(pattern)
	return regex{
		Error: err,
		Regex: r,
	}
}

// RBACMapT is the type corresponding to the RBACMap containing the entire
// details corresponding to role based accesss to API server api.
//nolint
type RBACMapT map[string](map[string][]regex)

// APIServerRBACMap is the map containing the Role based access
var APIServerRBACMap RBACMapT = RBACMapT{
	// Make sure that the admin roles has access to everything
	// and all the verbs are specified for admin.
	"admin": {
		"get":    []regex{getRegex(".*")},
		"post":   []regex{getRegex(".*")},
		"put":    []regex{getRegex(".*")},
		"delete": []regex{getRegex(".*")},
		"patch":  []regex{getRegex(".*")},
	},
	"manager": {
		"get": []regex{getRegex(".*")},
		"post": []regex{
			getRegex(`/api/v1/registry/.*`),
			getRegex(`/api/v1/status/.*`),
		},
		"put": []regex{
			getRegex(`/api/v1/registry/.*`),
			getRegex(`/api/v1/status/.*`),
		},
		"delete": []regex{
			getRegex(`/api/v1/registry/.*`),
			getRegex(`/api/v1/status/.*`),
		},
		"patch": []regex{
			getRegex(`/api/v1/registry/.*`),
			getRegex(`/api/v1/status/.*`),
		},
	},
}

// Roles returns the list of all the available roles in the RBAC systems.
func (r RBACMapT) Roles() []string {
	var roles []string
	for role := range r {
		roles = append(roles, role)
	}

	return roles
}

func (r RBACMapT) filterNonAvailableRoles(roles []string) []string {
	var availableRoles []string
	for _, role := range roles {
		if _, ok := r[role]; ok {
			availableRoles = append(availableRoles, role)
		}
	}

	return availableRoles
}

// Verbs returns the list of all the available verbs in the RBAC systems.
func (r RBACMapT) Verbs() []string {
	var verbs []string
	for verb := range r["admin"] {
		verbs = append(verbs, verb)
	}

	return verbs
}

// ValidateAccessI validates the access of the provided roles corresponding to the verb and
// apiEndpoint. The function returns a boolean value indicating if the access is valid or not.
// The I in the name indicate that the matching is case insensitive.
func (r RBACMapT) ValidateAccessI(roles []string, verb, apiEndpoint string) bool {
	avRoles := r.filterNonAvailableRoles(roles)
	for _, role := range avRoles {
		if apiEndpoints, ok := r[role][strings.ToLower(verb)]; ok {
			for i := range apiEndpoints {
				if apiEndpoints[i].Error != nil {
					continue
				}

				if apiEndpoints[i].Regex.Match([]byte(apiEndpoint)) {
					return true
				}
			}
		}
	}

	return false
}
