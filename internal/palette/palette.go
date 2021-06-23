package palette

import (
	"strconv"
)

// Ruleset is a group of rules, where order is important.
// Typically, it looks like one error rule and some
// "exceptions" — more specific color chains with no error
type Ruleset []*Rule

// NewRuleset creates a new Ruleset.
func NewRuleset(rules ...*Rule) Ruleset {
	return append(Ruleset(nil), rules...)
}

// Rule are representation of human-written rule:
//   "api has-curl" => "error text"
// or
//   "api allow-curl has-curl" => 1
type Rule struct {
	Colors []Color
	Masks  ColorMasks
	Error  string
}

// NewRule creates a new Rule.
func NewRule(colors []Color, error string) *Rule {
	return &Rule{
		Colors: colors,
		Masks:  NewColorMasks(colors),
		Error:  error,
	}
}

// IsError checks if the rule describes an error.
func (r *Rule) IsError() bool {
	return r.Error != ""
}

// ContainsIn checks if the rule's colors are contained in the passed mask.
func (r *Rule) ContainsIn(colorMasks ColorMasks) bool {
	// If the number of masks in the received list of masks is
	// less than in the current rule, then this rule cannot
	// automatically be a part of the received list of masks.
	if len(colorMasks) < len(r.Masks) {
		return false
	}

	for i := range r.Masks {
		matched := (r.Masks[i].Val & colorMasks[i].Val) == r.Masks[i].Val
		if !matched {
			return false
		}
	}

	return true
}

func (r *Rule) String(palette *Palette) string {
	var res string

	for i, color := range r.Colors {
		res += palette.GetNameByColor(color)
		if i != len(r.Colors)-1 {
			res += " "
		}
	}

	return res
}

// Palette is a group of rulesets.
// All colors are stored as Color struct, not as strings.
type Palette struct {
	Rulesets          []Ruleset
	ColorNamesMapping map[string]Color
}

// NewPalette creates a new Palette.
func NewPalette() *Palette {
	return &Palette{
		ColorNamesMapping: map[string]Color{
			"transparent": NewColor(SpecialColorTransparent, 0),
			"remover":     NewColor(SpecialColorRemover, 0),
		},
	}
}

func (p *Palette) AddRuleset(ruleset Ruleset) {
	p.Rulesets = append(p.Rulesets, ruleset)
}

func (p *Palette) ColorExists(colorName string) bool {
	_, ok := p.ColorNamesMapping[colorName]
	return ok
}

func (p *Palette) RegisterColorName(colorName string) Color {
	color, ok := p.ColorNamesMapping[colorName]
	if ok {
		return color
	}

	index := 0

	// User-defined colors are 1 << 2, 1 << 3, and so on.
	bitShift := len(p.ColorNamesMapping) - CountSpecialColors

	// If there are no more empty bits in the current mask.
	if bitShift >= MaxColorsInMask {
		index = bitShift / MaxColorsInMask
		bitShift %= MaxColorsInMask
	}

	color = NewColor(uint64(1)<<bitShift, index)
	p.ColorNamesMapping[colorName] = color

	return color
}

func (p *Palette) GetColorByName(colorName string) Color {
	return p.ColorNamesMapping[colorName]
}

func (p *Palette) GetNameByColor(needColor Color) string {
	for name, color := range p.ColorNamesMapping {
		if color == needColor {
			return name
		}
	}

	return strconv.FormatUint(needColor.Val, 10)
}

// ColorContainer is class containing colors after @color
// parsing above each function (order is important).
type ColorContainer struct {
	Colors []Color
}

func (c *ColorContainer) Add(color Color) {
	c.Colors = append(c.Colors, color)
}

func (c *ColorContainer) Contains(needColor Color) bool {
	for _, color := range c.Colors {
		if color == needColor {
			return true
		}
	}
	return false
}

func (c *ColorContainer) Empty() bool {
	return len(c.Colors) == 0
}

func (c *ColorContainer) String(palette *Palette, withHighlights ColorMasks) string {
	var desc string

	for _, color := range c.Colors {
		if !withHighlights.Contains(color) {
			continue
		}

		desc += "@"
		desc += palette.GetNameByColor(color)
	}

	return desc
}
