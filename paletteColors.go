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

	blackGradient = []LedColor{
		NewLedColor(0x000000),
	}

	blackWhiteGradientNoCycle = []LedColor{
		NewLedColor(0x000000),
		NewLedColor(0xFFFFFF),
	}

	// Paletten aus der Mandelbrot-Kueche
	defaultGradient = []LedColor{
		NewLedColor(0xFEBC08),
		NewLedColor(0xDAF9FE),
		NewLedColor(0x010662),
		NewLedColor(0xDAF9FE),
	}

	earthAndSkyGradient = []LedColor{
		NewLedColor(0xffffff), // 0.00
		NewLedColor(0xffff00), // 0.13
		NewLedColor(0xff3300), // 0.46
		NewLedColor(0x000099), // 0.76
		NewLedColor(0x0066ff), // 0.90
	}

	seashoreGradient = []LedColor{
		NewLedColor(0xC9FEC2),
		NewLedColor(0xE4E4A7),
		NewLedColor(0xF15020),
		NewLedColor(0x841C17),
		NewLedColor(0x0574AE),
		NewLedColor(0x89D2D0),
	}

	pastellGradient = []LedColor{
		NewLedColor(0xCDD0D1), // 0.00
		NewLedColor(0x6F85FF), // 0.18
		NewLedColor(0xFF5B94), // 0.42
		NewLedColor(0xFFFF84), // 0.63
		NewLedColor(0x8BEE91), // 0.86
	}

	darkGradient = []LedColor{
		NewLedColor(0xA80000), // 0.00
		NewLedColor(0x004D95), // 0.18
		NewLedColor(0xD06912), // 0.39
		NewLedColor(0x007C2A), // 0.57
		NewLedColor(0x4B23BF), // 0.78
	}

	darkerGradient = []LedColor{
		NewLedColor(0x700000), // 0.00
		NewLedColor(0x003363), // 0.18
		NewLedColor(0x8a460c), // 0.39
		NewLedColor(0x00521c), // 0.57
		NewLedColor(0x32177f), // 0.78
	}

	hotAndColdGradient = []LedColor{
		NewLedColor(0xFFFFFF), // 0.00
		NewLedColor(0x0066FF), // 0.16
		NewLedColor(0x333333), // 0.50
		NewLedColor(0xFF00CC), // 0.84
	}

	neonGradient = []LedColor{
		NewLedColor(0xffffff),
		NewLedColor(0x0099cc),
		NewLedColor(0xff0099),
	}

	transp01Gradient = []LedColor{
		NewLedColor(0xffffff),
		NewLedColor(0x0099cc),
		{0x00, 0x00, 0x00, 0x00},
		NewLedColor(0xff0099),
	}

	fadeRedGradient = []LedColor{
		NewLedColor(0x000000),
		NewLedColor(0xff0000),
	}

	fadeGreenGradient = []LedColor{
		NewLedColor(0x000000),
		NewLedColor(0x00ff00),
	}

	fadeBlueGradient = []LedColor{
		NewLedColor(0x000000),
		NewLedColor(0x0000ff),
	}

	hipsterGradient = []LedColor{
		NewLedColor(0x8D0E6B),
		NewLedColor(0x053434),
		NewLedColor(0x008080),
		NewLedColor(0x4D4208),
		NewLedColor(0xCDAD0A),
	}

	hotFireGradient = []LedColor{
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

	fireGradientNoCycle = []LedColor{
		NewLedColor(0x000000),
		NewLedColor(0x5f0809),
		NewLedColor(0xbe1013),
		NewLedColor(0xd23008),
		NewLedColor(0xe45323),
		NewLedColor(0xee771c),
		NewLedColor(0xf6960e),
		NewLedColor(0xffcd06),
	}

	darkJungleGradient = []LedColor{
		NewLedColor(0x184918),
		NewLedColor(0x406a3a),
		NewLedColor(0xb1a658),
		NewLedColor(0xa28d33),
		NewLedColor(0x6a6232),
		NewLedColor(0xa28d33),
		NewLedColor(0xb1a658),
		NewLedColor(0x406a3a),
	}

	richDeepavaliGradient = []LedColor{
		NewLedColor(0xECAC3B),
		NewLedColor(0xF65A00),
		NewLedColor(0xD02626),
		NewLedColor(0x213C79),
		NewLedColor(0xD03979),
		NewLedColor(0x76C8BA),
		NewLedColor(0xD03979),
		NewLedColor(0x213C79),
		NewLedColor(0xD02626),
		NewLedColor(0xF65A00),
	}

    brownishGradient = []LedColor{
        NewLedColor(0x3c0000),
        NewLedColor(0x782d2d),
        NewLedColor(0xb44641),
        NewLedColor(0xc87d5a),
        NewLedColor(0xf5d7af),
    }

    dusterGradient = []LedColor{
        NewLedColor(0xd22323),
        NewLedColor(0x1e4b55),
        NewLedColor(0x733223),
        NewLedColor(0x3c5037),
        NewLedColor(0x4b3255),
    }

    icecreamGradient = []LedColor{
        NewLedColor(0xa5236e),
        NewLedColor(0xeb1e4b),
        NewLedColor(0xf06937),
        NewLedColor(0xf5dc50),
        NewLedColor(0x2d969b),
    }

	Pico08Colors = []LedColor{
		// 16 Farben der Default-Palette
		NewLedColor(0x000000), // black / transparent
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
	// aAAGradient = []LedColor{
	//     NewLedColor(0x00ffff),
	//     NewLedColor(0x0000ff),
	//     NewLedColor(0x0080ff),
	//     NewLedColor(0x80ffff),
	//     NewLedColor(0x000080),
	// }
	// adribbleGradient = []LedColor{
	//     NewLedColor(0x3D4C53),
	//     NewLedColor(0x70B7BA),
	//     NewLedColor(0xF1433F),
	//     NewLedColor(0xE7E1D4),
	//     NewLedColor(0xFFFFFF),
	// }
	// brazilGradient = []LedColor{
	//     NewLedColor(0x008c53),
	//     NewLedColor(0x2e00e4),
	//     NewLedColor(0xdfea00),
	// }
	// bW01Gradient = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0xffffff),
	// }
	// bW02Gradient = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0xffffff),
	//     NewLedColor(0x000000),
	//     NewLedColor(0xffffff),
	// }
	// bW03Gradient = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0xffffff),
	//     NewLedColor(0x000000),
	//     NewLedColor(0xffffff),
	//     NewLedColor(0x000000),
	//     NewLedColor(0xffffff),
	//     NewLedColor(0x000000),
	//     NewLedColor(0xffffff),
	// }
	// cakeGradient = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0xBD4141),
	//     NewLedColor(0xD97570),
	//     NewLedColor(0xF2E8DF),
	//     NewLedColor(0xB2CFC0),
	//     NewLedColor(0x719D98),
	// }
	// castleGradient = []LedColor{
	//     NewLedColor(0x4B345C),
	//     NewLedColor(0x946282),
	//     NewLedColor(0xE5A19B),
	// }
	// cGAGradient = []LedColor{
	//     NewLedColor(0xd3517d),
	//     NewLedColor(0x15a0bf),
	//     NewLedColor(0xffc062),
	// }
	// civilGradient = []LedColor{
	//     NewLedColor(0x362F2D),
	//     NewLedColor(0x4C4C4C),
	//     NewLedColor(0x94B73E),
	//     NewLedColor(0xB5C0AF),
	//     NewLedColor(0xFAFDF2),
	// }
	// cold00Gradient = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0x00ffff),
	// }
	// cold01Gradient = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0x00ffff),
	//     NewLedColor(0x000080),
	// }
	// cold02Gradient = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0x00ffff),
	//     NewLedColor(0x000080),
	//     NewLedColor(0x008080),
	// }
	// cold03Gradient = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0x00ffff),
	//     NewLedColor(0x000080),
	//     NewLedColor(0x008080),
	//     NewLedColor(0x80ffff),
	// }
	// cold04Gradient = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0x00ffff),
	//     NewLedColor(0x000080),
	//     NewLedColor(0x008080),
	//     NewLedColor(0x80ffff),
	//     NewLedColor(0x8080ff),
	// }
	// cold05Gradient = []LedColor{
	//     NewLedColor(0x80ffff),
	//     NewLedColor(0x0080ff),
	//     NewLedColor(0xc0c0c0),
	//     NewLedColor(0x80a3ff),
	//     NewLedColor(0x585858),
	//     NewLedColor(0xb9f4f9),
	// }
	// colorBirdGradient = []LedColor{
	//     NewLedColor(0x1FA698),
	//     NewLedColor(0xC5D932),
	//     NewLedColor(0xF25922),
	//     NewLedColor(0x401E11),
	//     NewLedColor(0xD7195A),
	// }
	// cornGradient = []LedColor{
	//     NewLedColor(0x29231F),
	//     NewLedColor(0xEBE1CC),
	//     NewLedColor(0xDB9B1A),
	// }
	// fadeAllGradient = []LedColor{
	//     NewLedColor(0xff0000),
	//     NewLedColor(0xffff00),
	//     NewLedColor(0x00ff00),
	//     NewLedColor(0x00ffff),
	//     NewLedColor(0x0000ff),
	//     NewLedColor(0xff00ff),
	// }
	// fadeRedGradient = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0xff0000),
	// }
	// fadeGreenGradient = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0x00ff00),
	// }
	// fadeBlueGradient = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0x0000ff),
	// }
	// fadeYellowGradient = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0xffff00),
	// }
	// fadeCyanGradient = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0x00ffff),
	// }
	// fadeMagentaGradient = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0xff00ff),
	// }
	// cycleRedGradient = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0xff0000),
	//     NewLedColor(0x000000),
	// }
	// fame575Gradient = []LedColor{
	//     NewLedColor(0x540c0d),
	//     NewLedColor(0xfb7423),
	//     NewLedColor(0xf9f48e),
	//     NewLedColor(0x4176c4),
	//     NewLedColor(0x5aaf2e),
	// }
	// fireGradient = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0x000040),
	//     NewLedColor(0xff0000),
	//     NewLedColor(0xffff00),
	//     NewLedColor(0xffffff),
	//     NewLedColor(0xffffff),
	// }
	// fizzGradient = []LedColor{
	//     NewLedColor(0x588F27),
	//     NewLedColor(0x04BFBF),
	//     NewLedColor(0xF7E967),
	// }
	// foldGradient = []LedColor{
	//     NewLedColor(0x2A0308),
	//     NewLedColor(0x924F1B),
	//     NewLedColor(0xE2AC3F),
	//     NewLedColor(0xF8EDC6),
	//     NewLedColor(0x7BA58D),
	// }
	// kittyHCGradient = []LedColor{
	//     NewLedColor(0xc756a7),
	//     NewLedColor(0xe0dd00),
	//     NewLedColor(0xc9cdd0),
	// }
	// kittyGradient = []LedColor{
	//     NewLedColor(0x9f456b),
	//     NewLedColor(0x4f7a9a),
	//     NewLedColor(0xe6c84c),
	// }
	// lanternGradient = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0x0D9A0D),
	//     NewLedColor(0xFFFFFF),
	// }
	// leBronGradient = []LedColor{
	//     NewLedColor(0x3e3e3e),
	//     NewLedColor(0xd4b600),
	//     NewLedColor(0xffffff),
	// }
	// lemmingGradient = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0x0000ff),
	//     NewLedColor(0x00ff00),
	//     NewLedColor(0xffffff),
	// }
	// miamiViceGradient = []LedColor{
	//     NewLedColor(0x1be3ff),
	//     NewLedColor(0xff82dc),
	//     NewLedColor(0xffffff),
	// }
	// mIUSAGradient = []LedColor{
	//     NewLedColor(0x504b46),
	//     NewLedColor(0x1a3c53),
	//     NewLedColor(0xa00028),
	// }
	// mL581ATGradient = []LedColor{
	//     NewLedColor(0x699655),
	//     NewLedColor(0xf26a36),
	//     NewLedColor(0xffffff),
	// }
	// newSeasonGradient = []LedColor{
	//     NewLedColor(0x4C3933),
	//     NewLedColor(0x005B5B),
	//     NewLedColor(0xFE5502),
	//     NewLedColor(0xFEBF51),
	//     NewLedColor(0xF8EAB2),
	// }
	// nightspellGradient = []LedColor{
	//     NewLedColor(0xEDEEBA),
	//     NewLedColor(0xFEA81C),
	//     NewLedColor(0xB20152),
	//     NewLedColor(0x4B0B44),
	//     NewLedColor(0x240F37),
	// }
	// violetGradient = []LedColor{
	//     NewLedColor(0x4B0B44),
	//     NewLedColor(0x4B0B44),
	// }
	// rainbowsGradient = []LedColor{
	//     NewLedColor(0x492D61),
	//     NewLedColor(0x048091),
	//     NewLedColor(0x61C155),
	//     NewLedColor(0xF2D43F),
	//     NewLedColor(0xD1026C),
	// }
	// rastaGradient = []LedColor{
	//     NewLedColor(0xdc323c),
	//     NewLedColor(0xf0cb58),
	//     NewLedColor(0x3c825e),
	// }
	// rGBGradient = []LedColor{
	//     NewLedColor(0xff0000),
	//     NewLedColor(0x00ff00),
	//     NewLedColor(0x0000ff),
	// }
	// simpsonGradient = []LedColor{
	//     NewLedColor(0xd9c23e),
	//     NewLedColor(0xa96a95),
	//     NewLedColor(0x7d954b),
	//     NewLedColor(0x4b396b),
	// }
	// smurfGradient = []LedColor{
	//     NewLedColor(0x1d1628),
	//     NewLedColor(0x44bdf4),
	//     NewLedColor(0xe31e3a),
	//     NewLedColor(0xe8b118),
	//     NewLedColor(0xffffff),
	// }
	// springGradient = []LedColor{
	//     NewLedColor(0x1D1929),
	//     NewLedColor(0xFE5324),
	//     NewLedColor(0xA8143A),
	//     NewLedColor(0x66A595),
	//     NewLedColor(0xCFBD81),
	// }
	// sunsetGradient = []LedColor{
	//     NewLedColor(0x000000),
	//     NewLedColor(0x352438),
	//     NewLedColor(0x612949),
	//     NewLedColor(0x992839),
	//     NewLedColor(0xED6137),
	//     NewLedColor(0xD6AB5C),
	// }
	// warmthGradient = []LedColor{
	//     NewLedColor(0xFCEBB6),
	//     NewLedColor(0x5E412F),
	//     NewLedColor(0xF07818),
	//     NewLedColor(0x78C0A8),
	//     NewLedColor(0xF0A830),
	// }
	// wayyouGradient = []LedColor{
	//     NewLedColor(0x1C2130),
	//     NewLedColor(0x028F76),
	//     NewLedColor(0xB3E099),
	//     NewLedColor(0xFFEAAD),
	//     NewLedColor(0xD14334),
	// }
	// wepartedGradient = []LedColor{
	//     NewLedColor(0x027B7F),
	//     NewLedColor(0xFFA588),
	//     NewLedColor(0xD62957),
	//     NewLedColor(0xBF1E62),
	//     NewLedColor(0x572E4F),
	// }
)
