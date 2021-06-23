package palette

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config is a structure for storing a palette of colors as a config.
type Config struct {
	Palette [][]map[string]string `json:"palette"`
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
In .yaml syntax, it's a map from string key (description) to list (rules)`, path)
	}

	return parsePaletteRaw(config), nil
}

func parsePaletteRaw(config *Config) *Palette {
	pal := NewPalette()

	for _, group := range config.Palette {
		var rules []*Rule

		for _, rule := range group {
			var colorsRaw, desc string
			for rColor, rDesc := range rule {
				colorsRaw = rColor
				desc = rDesc
			}

			colors := strings.Split(colorsRaw, " ")
			colorsNums := make([]Color, 0, len(colors))

			for _, color := range colors {
				colorsNums = append(colorsNums, pal.RegisterColorName(color))
			}

			rules = append(rules, NewRule(colorsNums, desc))
		}

		pal.AddRuleset(NewRuleset(rules...))
	}
	return pal
}
