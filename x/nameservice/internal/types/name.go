package types

import (
	"fmt"
	"golang.org/x/net/idna"
	"strings"
)

// Name - provides useful functions for name manipulation
type Name string

// Split returns root,parent,child names
// Contract - it ensures name has second or third level name
func (n Name) Split() (rootName string, parentName string, childName string) {
	names := strings.Split(string(n), ".")
	if len(names) == 2 {
		return names[1], names[0], ""
	} else if len(names) == 3 {
		return names[2], names[1], names[0]
	}

	panic(fmt.Sprintf("name must be formatted level2_name.root_name or level3_name.level2_name.root_name; %s", n))
}

// GetLevels returns # of level in the name
func (n Name) Levels() int {
	return len(strings.Split(string(n), "."))
}

// NameHash return first 20 bytes of sha256(name string)
func (n Name) NameHash() (nameHash, childNameHash NameHash) {
	rootName, parentName, childName := n.Split()
	registryName := parentName + "." + rootName

	return GetNameHash(registryName), GetNameHash(childName)
}

// ValidateName validate the name properly normalised
// as described in _UTS46_ (https://unicode.org/reports/tr46/)
// with options `transitional=false` and `useSTD3AsciiRules=true`.
func (n Name) Validate() error {
	name := n.String()

	if len(name) > 255 {
		return fmt.Errorf("full name (%s) cannot exceed 255 charaters", n)
	}

	// length check is not provided for `transitional=false` (unicode) option
	labels := strings.Split(name, ".")
	for _, label := range labels {
		if len(label) > 63 {
			return fmt.Errorf("label (%s) cannot exceed 63 charaters", label)
		}
	}

	if normalized, err := idna.New(idna.ValidateForRegistration()).ToUnicode(name); err != nil {
		return err
	} else if normalized != name {
		return fmt.Errorf("given name (%s) is not properly normalized; should be %s", name, normalized)
	}

	return nil
}

// String convert Name to string
func (n Name) String() string {
	return string(n)
}
