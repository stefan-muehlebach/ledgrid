package colors

import (
	"errors"
)

// Die in diesem Package definierten Farben lassen sich in 10 Gruppen
// unterteilen
type ColorGroup int

const (
	Purples ColorGroup = iota
	Pinks
	Blues
	Reds
	Greens
	Oranges
	Yellows
	Browns
	Whites
	Grays
	numColorGroups
)

// Implementation des fmt.Stringer-Interfaces.
func (g ColorGroup) String() string {
	switch g {
	case Purples:
		return "Purples"
	case Pinks:
		return "Pinks"
	case Blues:
		return "Blues"
	case Reds:
		return "Reds"
	case Greens:
		return "Greens"
	case Oranges:
		return "Oranges"
	case Yellows:
		return "Yellows"
	case Browns:
		return "Browns"
	case Whites:
		return "Whites"
	case Grays:
		return "Grays"
	default:
		return "(unknown group)"
	}
}

// Mit String() zusammen implementiert ColorGroup das flag.Value Interface.
// Damit koennen Farbgruppen bspw. auf der Befehlszeile als Flags angegeben
// werden.
func (g *ColorGroup) Set(str string) error {
	switch str {
	case "Purples":
		*g = Purples
	case "Pinks":
		*g = Pinks
	case "Blues":
		*g = Blues
	case "Reds":
		*g = Reds
	case "Greens":
		*g = Greens
	case "Oranges":
		*g = Oranges
	case "Yellows":
		*g = Yellows
	case "Browns":
		*g = Browns
	case "Whites":
		*g = Whites
	case "Grays":
		*g = Grays
	default:
		return errors.New("Unknown color group: " + str)
	}
	return nil
}

// Groups ist ein Map, der fuer jede Farbgruppe einen Array von zugehoerigen
// Farbnamen enthaelt. Mit diesen Namen lassen sich ueber die Variable Map
// die eigentlichen Farben ermitteln.
var Groups = map[ColorGroup][]string{
	Browns: {
		"Cornsilk",
		"BlanchedAlmond",
		"Bisque",
		"NavajoWhite",
		"Wheat",
		"BurlyWood",
		"Tan",
		"RosyBrown",
		"SandyBrown",
		"Goldenrod",
		"DarkGoldenrod",
		"Peru",
		"Chocolate",
		"SaddleBrown",
		"Sienna",
		"Brown",
		"Maroon",
	},
	Reds: {
		"IndianRed",
		"LightCoral",
		"Salmon",
		"DarkSalmon",
		"LightSalmon",
		"Red",
		"Crimson",
		"FireBrick",
		"DarkRed",
	},
	Oranges: {
		"LightSalmon",
		"Coral",
		"Tomato",
		"OrangeRed",
		"DarkOrange",
		"Orange",
	},
	Yellows: {
		"Gold",
		"Yellow",
		"LightYellow",
		"LemonChiffon",
		"LightGoldenrodYellow",
		"PapayaWhip",
		"Moccasin",
		"PeachPuff",
		"PaleGoldenrod",
		"Khaki",
		"DarkKhaki",
	},
	Greens: {
		"GreenYellow",
		"Chartreuse",
		"LawnGreen",
		"Lime",
		"LimeGreen",
		"PaleGreen",
		"LightGreen",
		"MediumSpringGreen",
		"SpringGreen",
		"MediumSeaGreen",
		"SeaGreen",
		"ForestGreen",
		"Green",
		"DarkGreen",
		"YellowGreen",
		"OliveDrab",
		"Olive",
		"DarkOliveGreen",
		"MediumAquamarine",
		"DarkSeaGreen",
		"LightSeaGreen",
		"DarkCyan",
		"Teal",
	},
	Blues: {
		"Aqua",
		"Cyan",
		"LightCyan",
		"PaleTurquoise",
		"Aquamarine",
		"Turquoise",
		"MediumTurquoise",
		"DarkTurquoise",
		"CadetBlue",
		"SteelBlue",
		"LightSteelBlue",
		"PowderBlue",
		"LightBlue",
		"SkyBlue",
		"LightSkyBlue",
		"DeepSkyBlue",
		"DodgerBlue",
		"CornflowerBlue",
		"RoyalBlue",
		"Blue",
		"MediumBlue",
		"DarkBlue",
		"Navy",
		"MidnightBlue",
	},
	Purples: {
		"Lavender",
		"Thistle",
		"Plum",
		"Violet",
		"Orchid",
		"Fuchsia",
		"Magenta",
		"MediumOrchid",
		"MediumPurple",
		"BlueViolet",
		"DarkViolet",
		"DarkOrchid",
		"DarkMagenta",
		"Purple",
		"Indigo",
		"DarkSlateBlue",
		"SlateBlue",
		"MediumSlateBlue",
	},
	Pinks: {
		"Pink",
		"LightPink",
		"HotPink",
		"DeepPink",
		"MediumVioletRed",
		"PaleVioletRed",
	},
	Whites: {
		"White",
		"Snow",
		"Honeydew",
		"MintCream",
		"Azure",
		"AliceBlue",
		"GhostWhite",
		"WhiteSmoke",
		"Seashell",
		"Beige",
		"OldLace",
		"FloralWhite",
		"Ivory",
		"AntiqueWhite",
		"Linen",
		"LavenderBlush",
		"MistyRose",
	},
	Grays: {
		"Gainsboro",
		"LightGray",
		"Silver",
		"DarkGray",
		"Gray",
		"DimGray",
		"LightSlateGray",
		"SlateGray",
		"DarkSlateGray",
		"Black",
	},
}
