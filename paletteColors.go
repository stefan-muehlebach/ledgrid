//go:generate go run gen.go

package ledgrid

var (
	// Dies sind die Farblisten, welche fuer die einzelnen Paletten verwendet
	// werden. Diese Namen werden vom Programm 'gen' gelesen und zur
	// Erstellung des Files paletteNames.go verwendet. Die Namen der
	// Farblisten werden NICHT exportiert (d.h. beginnen mit Kleinbuchstaben),
	// muessen aber mit 'Gradient' enden, damit daraus Paletten erstellt
	// werden sollen. Endet der Name mit 'GradientNoCycle', dann wird eine
	// Palette ohne 'cycle'-Flag erstellt.

	blackAndWhiteGradient = []ColorStop{
		{0.0, NewLedColor(0x000000)},
		{0.5, NewLedColor(0xFFFFFF)},
		{1.0, NewLedColor(0x000000)},
	}

	// Paletten aus der Mandelbrot-Kueche
	defaultGradient = []ColorStop{
		{0.00, NewLedColor(0xFEBC08)},
		{0.25, NewLedColor(0xDAF9FE)},
		{0.50, NewLedColor(0x010662)},
		{0.75, NewLedColor(0xDAF9FE)},
		{1.00, NewLedColor(0xFEBC08)},
	}

	earthAndSkyGradient = []ColorStop{
		{0.00, NewLedColor(0xffffff)},
		{0.13, NewLedColor(0xffff00)},
		{0.46, NewLedColor(0xff3300)},
		{0.76, NewLedColor(0x000099)},
		{0.90, NewLedColor(0x0066ff)},
		{1.00, NewLedColor(0xffffff)},
	}

	seashoreGradient = []ColorStop{
		{0.00, NewLedColor(0xC9FEC2)},
		{0.16, NewLedColor(0xE4E4A7)},
		{0.33, NewLedColor(0xF15020)},
		{0.50, NewLedColor(0x841C17)},
		{0.66, NewLedColor(0x0574AE)},
		{0.83, NewLedColor(0x89D2D0)},
		{1.00, NewLedColor(0xC9FEC2)},
	}

	pastellGradient = []ColorStop{
		{0.00, NewLedColor(0xCDD0D1)},
		{0.18, NewLedColor(0x6F85FF)},
		{0.42, NewLedColor(0xFF5B94)},
		{0.63, NewLedColor(0xFFFF84)},
		{0.86, NewLedColor(0x8BEE91)},
		{1.00, NewLedColor(0xCDD0D1)},
	}

	darkGradient = []ColorStop{
		{0.00, NewLedColor(0xA80000)},
		{0.18, NewLedColor(0x004D95)},
		{0.39, NewLedColor(0xD06912)},
		{0.57, NewLedColor(0x007C2A)},
		{0.78, NewLedColor(0x4B23BF)},
		{1.00, NewLedColor(0xA80000)},
	}

	darkerGradient = []ColorStop{
		{0.00, NewLedColor(0x700000)},
		{0.18, NewLedColor(0x003363)},
		{0.39, NewLedColor(0x8a460c)},
		{0.57, NewLedColor(0x00521c)},
		{0.78, NewLedColor(0x32177f)},
		{1.00, NewLedColor(0x700000)},
	}

	hotAndColdGradient = []ColorStop{
		{0.00, NewLedColor(0xFFFFFF)},
		{0.16, NewLedColor(0x0066FF)},
		{0.50, NewLedColor(0x333333)},
		{0.84, NewLedColor(0xFF00CC)},
		{1.00, NewLedColor(0xFFFFFF)},
	}

	neonColorList = []LedColor{
		NewLedColor(0xffffff),
		NewLedColor(0x0099cc),
		NewLedColor(0xff0099),
	}

	fadeRedColorListNonCyc = []LedColor{
		NewLedColor(0x000000),
		NewLedColor(0xff0000),
	}

	fadeGreenColorListNonCyc = []LedColor{
		NewLedColor(0x000000),
		NewLedColor(0x00ff00),
	}

	fadeBlueColorListNonCyc = []LedColor{
		NewLedColor(0x000000),
		NewLedColor(0x0000ff),
	}

	fadeYellowColorListNonCyc = []LedColor{
		NewLedColor(0x000000),
		NewLedColor(0xffff00),
	}

	fadeCyanColorListNonCyc = []LedColor{
		NewLedColor(0x000000),
		NewLedColor(0x00ffff),
	}

	fadeMagentaColorListNonCyc = []LedColor{
		NewLedColor(0x000000),
		NewLedColor(0xff00ff),
	}

	hipsterColorList = []LedColor{
		NewLedColor(0x8D0E6B),
		NewLedColor(0x053434),
		NewLedColor(0x008080),
		NewLedColor(0x4D4208),
		NewLedColor(0xCDAD0A),
	}

	hotFireColorList = []LedColor{
		NewLedColor(0xbe1013),
		NewLedColor(0xd23008),
		NewLedColor(0xe45323),
		NewLedColor(0xee771c),
		NewLedColor(0xf6960e),
		NewLedColor(0xffcd06),
		NewLedColor(0xf6960e),
		NewLedColor(0xee771c),
		NewLedColor(0xe45323),
		NewLedColor(0xd23008),
	}

	fireGradient = []ColorStop{
		{0.00, NewLedColor(0x000000)},
		{0.14, NewLedColor(0x5f0809)},
		{0.29, NewLedColor(0xbe1013)},
		{0.43, NewLedColor(0xd23008)},
		{0.57, NewLedColor(0xe45323)},
		{0.71, NewLedColor(0xee771c)},
		{0.86, NewLedColor(0xf6960e)},
		{1.00, NewLedColor(0xffcd06)},
	}

	darkJungleColorList = []LedColor{
		NewLedColor(0x184918),
		NewLedColor(0x406a3a),
		NewLedColor(0xb1a658),
		NewLedColor(0xa28d33),
		NewLedColor(0x6a6232),
		NewLedColor(0xa28d33),
		NewLedColor(0xb1a658),
		NewLedColor(0x406a3a),
	}

	Pico08Colors = []LedColor{
		// 16 Farben der Default-Palette
		NewLedColorAlpha(0x00000000), // black / transparent
		NewLedColor(0x1d2b53), // dark blue
		NewLedColor(0x7e2553), // dark purple
		NewLedColor(0x008751), // dark green
		NewLedColor(0xab5236), // brown
		NewLedColor(0x5f574f), // dark gray
		NewLedColor(0xc2c3c7), // light gray
		NewLedColor(0xfff1e8), // white
		NewLedColor(0xff004d), // red
		NewLedColor(0xffa300), // orange
		NewLedColor(0xffff27), // yellow
		NewLedColor(0x00e756), // green
		NewLedColor(0x29adff), // blue
		NewLedColor(0x83769c), // indigo
		NewLedColor(0xff77a8), // pink
		NewLedColor(0xffccaa), // peach
		// 16 Farben der Hidden-Palette
		NewLedColor(0x291814), // onyx
		NewLedColor(0x111d35), // midnight
		NewLedColor(0x422136), // plum
		NewLedColor(0x125359), // forest
		NewLedColor(0x742f29), // chocolate
		NewLedColor(0x49333b), // eggplant
		NewLedColor(0xa28879), // beige
		NewLedColor(0xf3ef7d), // lemon
		NewLedColor(0xbe1250), // burgundy
		NewLedColor(0xff6c24), // pumpkin
		NewLedColor(0xa8e72e), // lime
		NewLedColor(0x00b543), // jade
		NewLedColor(0x065ab5), // royal
		NewLedColor(0x754665), // mauve
		NewLedColor(0xff6e59), // coral
		NewLedColor(0xff9d81), // salmon
	}

	// Kopierte Paletten aus PixelController
	// aAAColorList = []LedColor{
	//     NewLedColor(0x00ffff),
	//     NewLedColor(0x0000ff),
	//     NewLedColor(0x0080ff),
	//     NewLedColor(0x80ffff),
	//     NewLedColor(0x000080),
	// }
	// adribbleColorList = []LedColor{
	//     NewLedColor(0x3D4C53),
	//     NewLedColor(0x70B7BA),
	//     NewLedColor(0xF1433F),
	//     NewLedColor(0xE7E1D4),
	//     NewLedColor(0xFFFFFF),
	// }
	// brazilColorList = []LedColor{
	//     NewLedColor(0x008c53),
	//     NewLedColor(0x2e00e4),
	//     NewLedColor(0xdfea00),
	// }
	// bW01ColorList = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0xffffff),
	// }
	// bW02ColorList = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0xffffff),
	//     NewLedColor(0x000000),
	//     NewLedColor(0xffffff),
	// }
	// bW03ColorList = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0xffffff),
	//     NewLedColor(0x000000),
	//     NewLedColor(0xffffff),
	//     NewLedColor(0x000000),
	//     NewLedColor(0xffffff),
	//     NewLedColor(0x000000),
	//     NewLedColor(0xffffff),
	// }
	// cakeColorList = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0xBD4141),
	//     NewLedColor(0xD97570),
	//     NewLedColor(0xF2E8DF),
	//     NewLedColor(0xB2CFC0),
	//     NewLedColor(0x719D98),
	// }
	// castleColorList = []LedColor{
	//     NewLedColor(0x4B345C),
	//     NewLedColor(0x946282),
	//     NewLedColor(0xE5A19B),
	// }
	// cGAColorList = []LedColor{
	//     NewLedColor(0xd3517d),
	//     NewLedColor(0x15a0bf),
	//     NewLedColor(0xffc062),
	// }
	// civilColorList = []LedColor{
	//     NewLedColor(0x362F2D),
	//     NewLedColor(0x4C4C4C),
	//     NewLedColor(0x94B73E),
	//     NewLedColor(0xB5C0AF),
	//     NewLedColor(0xFAFDF2),
	// }
	// cold00ColorList = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0x00ffff),
	// }
	// cold01ColorList = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0x00ffff),
	//     NewLedColor(0x000080),
	// }
	// cold02ColorList = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0x00ffff),
	//     NewLedColor(0x000080),
	//     NewLedColor(0x008080),
	// }
	// cold03ColorList = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0x00ffff),
	//     NewLedColor(0x000080),
	//     NewLedColor(0x008080),
	//     NewLedColor(0x80ffff),
	// }
	// cold04ColorList = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0x00ffff),
	//     NewLedColor(0x000080),
	//     NewLedColor(0x008080),
	//     NewLedColor(0x80ffff),
	//     NewLedColor(0x8080ff),
	// }
	// cold05ColorList = []LedColor{
	//     NewLedColor(0x80ffff),
	//     NewLedColor(0x0080ff),
	//     NewLedColor(0xc0c0c0),
	//     NewLedColor(0x80a3ff),
	//     NewLedColor(0x585858),
	//     NewLedColor(0xb9f4f9),
	// }
	// colorBirdColorList = []LedColor{
	//     NewLedColor(0x1FA698),
	//     NewLedColor(0xC5D932),
	//     NewLedColor(0xF25922),
	//     NewLedColor(0x401E11),
	//     NewLedColor(0xD7195A),
	// }
	// cornColorList = []LedColor{
	//     NewLedColor(0x29231F),
	//     NewLedColor(0xEBE1CC),
	//     NewLedColor(0xDB9B1A),
	// }
	// fadeAllColorList = []LedColor{
	//     NewLedColor(0xff0000),
	//     NewLedColor(0xffff00),
	//     NewLedColor(0x00ff00),
	//     NewLedColor(0x00ffff),
	//     NewLedColor(0x0000ff),
	//     NewLedColor(0xff00ff),
	// }
	// fadeRedColorList = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0xff0000),
	// }
	// fadeGreenColorList = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0x00ff00),
	// }
	// fadeBlueColorList = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0x0000ff),
	// }
	// fadeYellowColorList = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0xffff00),
	// }
	// fadeCyanColorList = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0x00ffff),
	// }
	// fadeMagentaColorList = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0xff00ff),
	// }
	// cycleRedColorList = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0xff0000),
	//     NewLedColor(0x000000),
	// }
	// fame575ColorList = []LedColor{
	//     NewLedColor(0x540c0d),
	//     NewLedColor(0xfb7423),
	//     NewLedColor(0xf9f48e),
	//     NewLedColor(0x4176c4),
	//     NewLedColor(0x5aaf2e),
	// }
	// fireColorList = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0x000040),
	//     NewLedColor(0xff0000),
	//     NewLedColor(0xffff00),
	//     NewLedColor(0xffffff),
	//     NewLedColor(0xffffff),
	// }
	// fizzColorList = []LedColor{
	//     NewLedColor(0x588F27),
	//     NewLedColor(0x04BFBF),
	//     NewLedColor(0xF7E967),
	// }
	// foldColorList = []LedColor{
	//     NewLedColor(0x2A0308),
	//     NewLedColor(0x924F1B),
	//     NewLedColor(0xE2AC3F),
	//     NewLedColor(0xF8EDC6),
	//     NewLedColor(0x7BA58D),
	// }
	// kittyHCColorList = []LedColor{
	//     NewLedColor(0xc756a7),
	//     NewLedColor(0xe0dd00),
	//     NewLedColor(0xc9cdd0),
	// }
	// kittyColorList = []LedColor{
	//     NewLedColor(0x9f456b),
	//     NewLedColor(0x4f7a9a),
	//     NewLedColor(0xe6c84c),
	// }
	// lanternColorList = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0x0D9A0D),
	//     NewLedColor(0xFFFFFF),
	// }
	// leBronColorList = []LedColor{
	//     NewLedColor(0x3e3e3e),
	//     NewLedColor(0xd4b600),
	//     NewLedColor(0xffffff),
	// }
	// lemmingColorList = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0x0000ff),
	//     NewLedColor(0x00ff00),
	//     NewLedColor(0xffffff),
	// }
	// miamiViceColorList = []LedColor{
	//     NewLedColor(0x1be3ff),
	//     NewLedColor(0xff82dc),
	//     NewLedColor(0xffffff),
	// }
	// mIUSAColorList = []LedColor{
	//     NewLedColor(0x504b46),
	//     NewLedColor(0x1a3c53),
	//     NewLedColor(0xa00028),
	// }
	// mL581ATColorList = []LedColor{
	//     NewLedColor(0x699655),
	//     NewLedColor(0xf26a36),
	//     NewLedColor(0xffffff),
	// }
	// newSeasonColorList = []LedColor{
	//     NewLedColor(0x4C3933),
	//     NewLedColor(0x005B5B),
	//     NewLedColor(0xFE5502),
	//     NewLedColor(0xFEBF51),
	//     NewLedColor(0xF8EAB2),
	// }
	// nightspellColorList = []LedColor{
	//     NewLedColor(0xEDEEBA),
	//     NewLedColor(0xFEA81C),
	//     NewLedColor(0xB20152),
	//     NewLedColor(0x4B0B44),
	//     NewLedColor(0x240F37),
	// }
	// violetColorList = []LedColor{
	//     NewLedColor(0x4B0B44),
	//     NewLedColor(0x4B0B44),
	// }
	// rainbowsColorList = []LedColor{
	//     NewLedColor(0x492D61),
	//     NewLedColor(0x048091),
	//     NewLedColor(0x61C155),
	//     NewLedColor(0xF2D43F),
	//     NewLedColor(0xD1026C),
	// }
	// rastaColorList = []LedColor{
	//     NewLedColor(0xdc323c),
	//     NewLedColor(0xf0cb58),
	//     NewLedColor(0x3c825e),
	// }
	// rGBColorList = []LedColor{
	//     NewLedColor(0xff0000),
	//     NewLedColor(0x00ff00),
	//     NewLedColor(0x0000ff),
	// }
	// simpsonColorList = []LedColor{
	//     NewLedColor(0xd9c23e),
	//     NewLedColor(0xa96a95),
	//     NewLedColor(0x7d954b),
	//     NewLedColor(0x4b396b),
	// }
	// smurfColorList = []LedColor{
	//     NewLedColor(0x1d1628),
	//     NewLedColor(0x44bdf4),
	//     NewLedColor(0xe31e3a),
	//     NewLedColor(0xe8b118),
	//     NewLedColor(0xffffff),
	// }
	// springColorList = []LedColor{
	//     NewLedColor(0x1D1929),
	//     NewLedColor(0xFE5324),
	//     NewLedColor(0xA8143A),
	//     NewLedColor(0x66A595),
	//     NewLedColor(0xCFBD81),
	// }
	// sunsetColorList = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0x352438),
	//     NewLedColor(0x612949),
	//     NewLedColor(0x992839),
	//     NewLedColor(0xED6137),
	//     NewLedColor(0xD6AB5C),
	// }
	// warmthColorList = []LedColor{
	//     NewLedColor(0xFCEBB6),
	//     NewLedColor(0x5E412F),
	//     NewLedColor(0xF07818),
	//     NewLedColor(0x78C0A8),
	//     NewLedColor(0xF0A830),
	// }
	// wayyouColorList = []LedColor{
	//     NewLedColor(0x1C2130),
	//     NewLedColor(0x028F76),
	//     NewLedColor(0xB3E099),
	//     NewLedColor(0xFFEAAD),
	//     NewLedColor(0xD14334),
	// }
	// wepartedColorList = []LedColor{
	//     NewLedColor(0x027B7F),
	//     NewLedColor(0xFFA588),
	//     NewLedColor(0xD62957),
	//     NewLedColor(0xBF1E62),
	//     NewLedColor(0x572E4F),
	// }
)
