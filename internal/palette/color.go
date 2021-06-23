package palette

// All functions without any colors are transparent: C + transparent = C.
const SpecialColorTransparent uint64 = 0

// @color remover works so: C + remover = transparent.
const SpecialColorRemover uint64 = 1 << 0

// Number of colors with which we occupy the initial bits of the first mask.
// The color SpecialColorTransparent is zero, so it is not counted.
const CountSpecialColors int = 1

// MaxColorsInMask stores the maximum number of colors in the mask.
// Used in tests. It cannot be more than 64.
var MaxColorsInMask = 64

type ColorMask Color
type ColorMasks []ColorMask

// Color structure describes a color.
// The color is represented by two components.
// The first component is a mask, which contains
// single one bit that defines the color number.
// The second component is the index of the mask
// in which this one bit is.
//
// The length of the mask in any color cannot be
// more than MaxColorsInMask.
//
// So, for example, if the number of colors in the
// mask can be a maximum of 3, then if you need to
// keep the fourth color, then for this you need
// to use the following mask, so the mask Index will
// be 1, and the Val in the mask will be
//
//   1 << 4 % MaxColorsInMask,
//
// then there is 0b01.
type Color struct {
	// Val is the mask of the current color.
	// A one bit determines the number of the current color.
	Val uint64

	// Index is the ordinal number of the mask,
	// in which this color is indicated by bit one.
	Index int
}

func NewColor(val uint64, index int) Color {
	return Color{
		Val:   val,
		Index: index,
	}
}

func NewEmptyColorMasks() ColorMasks {
	return make([]ColorMask, 1)
}

func NewColorMasks(colors []Color) ColorMasks {
	var maxIndex int
	for _, color := range colors {
		if color.Index > maxIndex {
			maxIndex = color.Index
		}
	}

	masks := make(ColorMasks, maxIndex+1)
	for i := range masks {
		masks[i].Index = i
	}

	for _, color := range colors {
		masks[color.Index].Val |= color.Val
	}

	return masks
}

func (masks ColorMasks) Add(color Color) ColorMasks {
	curLen := len(masks)
	if color.Index >= curLen {
		for i := 0; i < color.Index-(curLen-1); i++ {
			masks = append(masks, ColorMask{Index: curLen + i})
		}
	}

	masks[color.Index].Val |= color.Val

	return masks
}

func (masks ColorMasks) Contains(color Color) bool {
	if len(masks) <= color.Index {
		return false
	}

	return (masks[color.Index].Val&color.Val) != 0 && masks[color.Index].Index == color.Index
}
