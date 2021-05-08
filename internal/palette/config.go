package palette

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Palette [][]map[string]string `json:"palette"`
}

// OpenPaletteFromFile returns a ready-use palette from a file.
func OpenPaletteFromFile(path string) (*Palette, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if filepath.Ext(path) == ".json" {
		return ReadPaletteFileJSON(data)
	}

	return ReadPaletteFileYAML(data)
}

func ReadPaletteFileYAML(data []byte) (*Palette, error) {
	config := &Config{}

	err := yaml.Unmarshal(data, &config.Palette)
	if err != nil {
		return nil, err
	}

	return parsePaletteRaw(config), nil
}

func ReadPaletteFileJSON(data []byte) (*Palette, error) {
	var config *Config

	err := json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
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
