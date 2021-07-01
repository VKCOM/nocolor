package palette

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

type ConfigRule map[string]string

// Config is a structure for storing a palette of colors as a config.
type Config struct {
	Palette map[string][]ConfigRule
}

type groupRules struct {
	name  string
	rules []ConfigRule
}

// OpenPaletteFromFile returns a ready-use palette from a file.
func OpenPaletteFromFile(path string) (*Palette, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		var perr *fs.PathError
		if errors.As(err, &perr) {
			absPath, err := filepath.Abs(path)
			if err != nil {
				absPath = path
			}

			return nil, fmt.Errorf(`cannot open palette file '%s', file not found. Full path: %s`, path, absPath)
		}

		return nil, fmt.Errorf(`cannot open palette file '%s': %v`, path, err)
	}

	return ReadPaletteFileYAML(path, data)
}

// The ReadPaletteFileYAML function interprets the passed text as a
// config in YAML format and returns a ready-made palette.
func ReadPaletteFileYAML(path string, data []byte) (*Palette, error) {
	config := &Config{}

	err := yaml.Unmarshal(data, &config.Palette)
	if err != nil {
		return nil, fmt.Errorf(`could not parse palette file '%s'. 
The correct format is:

ruleset description:
- rule
- rule

(optionally with many rulesets)
In .yaml syntax, it's a map from a string key (description) to a list (rules)`, path)
	}

	return parsePaletteRaw(path, config)
}

func parsePaletteRaw(path string, config *Config) (*Palette, error) {
	pal := NewPalette()

	groups := make([]groupRules, 0, len(config.Palette))
	for name, group := range config.Palette {
		var singleGroup groupRules
		singleGroup.name = name
		singleGroup.rules = group
		groups = append(groups, singleGroup)
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].name < groups[j].name
	})

	for _, group := range groups {
		var rules []*Rule

		for _, rule := range group.rules {
			var colorsRaw, desc string
			for rColor, rDesc := range rule {
				colorsRaw = rColor
				desc = rDesc
			}

			colors := strings.Split(colorsRaw, " ")
			colorsNums := make([]Color, 0, len(colors))

			for _, color := range colors {
				if color == "transparent" {
					return nil, fmt.Errorf("error in palette file '%s': use of 'transparent' color is prohibited in the rules", path)
				}
				if color == "*" {
					return nil, fmt.Errorf("error in palette file '%s': use of 'wildcard' color is prohibited in the rules", path)
				}

				colorsNums = append(colorsNums, pal.RegisterColorName(color))
			}

			rules = append(rules, NewRule(colorsNums, desc))
		}

		pal.AddRuleset(NewRuleset(rules...))
	}

	return pal, nil
}
