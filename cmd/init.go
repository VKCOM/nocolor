package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/i582/cfmt/cmd/cfmt"
)

var configTemplate = `# This is a palette file where you can write your own rules.
# All rules are divided into groups, where one group represents an associated set of rules.
# Each group is a pair, where the group identifier is the key, and the value is a list of rules.
#
# For example:
test group:
- red green: forbidden to call green functions from red ones
- red blue green: ""
# There can be several groups.
test group 2:
- highload no-highload: forbidden to call no-highload functions from highload ones
`

// Init create template config for project.
func Init() (status int, err error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return 2, fmt.Errorf("some unexpected error (%v), try again", err)
	}

	filePath := filepath.Join(workingDir, "palette.yaml")

	// If the file exists.
	if _, err = os.Stat(filePath); !os.IsNotExist(err) {
		cfmt.Println("The palette file {{already exists}}::bold.")
		return 0, nil
	}

	err = createFile(filePath)
	if err != nil {
		return 2, err
	}
	return 0, nil
}

func createFile(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0677)
	if err != nil {
		return fmt.Errorf("error create file '%s': %v", filePath, err)
	}
	defer file.Close()

	fmt.Fprint(file, configTemplate)
	cfmt.Println("The palette file was created {{successfully}}::green. It is located in the root folder of the project.")
	return nil
}
