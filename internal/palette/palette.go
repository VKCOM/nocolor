package palette

import (
	"strconv"
)

type Color = uint64
type ColorMask = uint64

// All functions without any colors are transparent: C + transparent = C.
const SpecialColorTransparent Color = 0

// @color remover works so: C + remover = transparent.
const SpecialColorRemover Color = 1 << 63

// Ruleset is a group of rules, order is important.
// Typically, it looks like one error rule and some
// "exceptions" — more specific color chains with no error
type Ruleset []*Rule

func NewRuleset(rules ...*Rule) Ruleset {
	return append(Ruleset(nil), rules...)
}

// Rule are representation of human-written "api has-curl" => "error text"
// or "api allow-curl has-curl" => 1
// All colors are pre-converted to numeric while reading strings.
type Rule struct {
	Colors []Color
	Mask   ColorMask
	Error  string
}

func NewRule(colors []Color, error string) *Rule {
	var mask ColorMask
	for _, color := range colors {
		mask |= color
	}

	return &Rule{Colors: colors, Mask: mask, Error: error}
}

func (r *Rule) IsError() bool {
	return r.Error != ""
}

func (r *Rule) ContainsIn(colorMask ColorMask) bool {
	return (r.Mask & colorMask) == r.Mask
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
// All colors are stored as numbers, not as strings:
//   They are numbers 1 << n, as we want to use bitmasks
//   to quickly test whether to check a rule for callstack.
type Palette struct {
	Rulesets          []Ruleset
	ColorNamesMapping map[string]Color
}

func NewPalette() *Palette {
	return &Palette{
		ColorNamesMapping: map[string]Color{
			"transparent": SpecialColorTransparent,
			"remover":     SpecialColorRemover,
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

	bitShift := len(p.ColorNamesMapping) - 1 // User-defined colors are 1<<1, 1<<2, and so on.
	color = Color(1) << bitShift

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

	return strconv.FormatUint(needColor, 10)
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

func (c *ColorContainer) String(palette *Palette, withHighlights ColorMask) string {
	var desc string

	for _, color := range c.Colors {
		if withHighlights == 0 || withHighlights&color != 0 {
			desc += "@"
			desc += palette.GetNameByColor(color)
		}
	}

	return desc
}
