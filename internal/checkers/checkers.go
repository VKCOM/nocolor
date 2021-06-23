package checkers

import (
	"github.com/VKCOM/noverify/src/linter"
)

// Contains checks for the presence of a checker in the list.
func Contains(name string) bool {
	for _, info := range List() {
		if info.Name == name {
			return true
		}
	}
	return false
}

// List returns a list of checkers.
func List() []linter.CheckerInfo {
	return []linter.CheckerInfo{
		{
			Name:     "errorColor",
			Default:  true,
			Quickfix: false,
			Comment:  `Report erroneous color in phpdoc`,
		},
	}
}
