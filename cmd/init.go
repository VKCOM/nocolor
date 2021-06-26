package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/i582/cfmt/cmd/cfmt"
)

var configTemplate = `# This is an example palette file; your rules should expect the same format
# There are multiple groups (rulesets); each ruleset is a description (key) and a list of rules
# For more info, consider https://github.com/vkcom/nocolor

demo of green red from the docs:
- green red: calling red from green is prohibited

analyze performance:
- fast slow: potential performance leak
- fast slow-ignore slow: ""
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
