package ledgrid

var (
	// Dies sind die Farblisten, welche fuer die einzelnen Paletten verwendet
	// werden. Diese Namen werden vom Programm 'gen' gelesen und zur
	// Erstellung des Files paletteNames.go verwendet. Die Namen der
	// Farblisten werden NICHT exportiert (d.h. beginnen mit Kleinbuchstaben),
	// muessen aber mit 'Colors' enden.

	// Paletten aus der Mandelbrot-Kueche
	defaultGradient = []LedColor{
		{0xFE, 0xBC, 0x08, 0xFF},
		{0xDA, 0xF9, 0xFE, 0xFF},
		{0x01, 0x06, 0x62, 0xFF},
		{0xDA, 0xF9, 0xFE, 0xFF},
		{0xFE, 0xBC, 0x08, 0xFF},
	}

	earthAndSkyGradient = []LedColor{
		{0xff, 0xff, 0xff, 0xFF}, // 0.00
		{0xff, 0xff, 0x00, 0xFF}, // 0.13
		{0xff, 0x33, 0x00, 0xFF}, // 0.46
		{0x00, 0x00, 0x99, 0xFF}, // 0.76
		{0x00, 0x66, 0xff, 0xFF}, // 0.90
		{0xff, 0xff, 0xff, 0xFF}, // 1.00
	}

	seashoreGradient = []LedColor{
		{0xC9, 0xFE, 0xC2, 0xFF},
		{0xE4, 0xE4, 0xA7, 0xFF},
		{0xF1, 0x50, 0x20, 0xFF},
		{0x84, 0x1C, 0x17, 0xFF},
		{0x05, 0x74, 0xAE, 0xFF},
		{0x89, 0xD2, 0xD0, 0xFF},
		{0xC9, 0xFE, 0xC2, 0xFF},
	}

	pastellGradient = []LedColor{
		{0xCD, 0xD0, 0xD1, 0xFF}, // 0.00
		{0x6F, 0x85, 0xFF, 0xFF}, // 0.18
		{0xFF, 0x5B, 0x94, 0xFF}, // 0.42
		{0xFF, 0xFF, 0x84, 0xFF}, // 0.63
		{0x8B, 0xEE, 0x91, 0xFF}, // 0.86
		{0xCD, 0xD0, 0xD1, 0xFF}, // 1.00
	}

	darkGradient = []LedColor{
		{0xA8, 0x00, 0x00, 0xFF}, // 0.00
		{0x00, 0x4D, 0x95, 0xFF}, // 0.18
		{0xD0, 0x69, 0x12, 0xFF}, // 0.39
		{0x00, 0x7C, 0x2A, 0xFF}, // 0.57
		{0x4B, 0x23, 0xBF, 0xFF}, // 0.78
		{0xA8, 0x00, 0x00, 0xFF}, // 1.00
	}

	darkerGradient = []LedColor{
		{0x70, 0x00, 0x00, 0xFF}, // 0.00
		{0x00, 0x33, 0x63, 0xFF}, // 0.18
		{0x8a, 0x46, 0x0c, 0xFF}, // 0.39
		{0x00, 0x52, 0x1c, 0xFF}, // 0.57
		{0x32, 0x17, 0x7f, 0xFF}, // 0.78
		{0x70, 0x00, 0x00, 0xFF}, // 0.00
	}

	hotAndColdGradient = []LedColor{
		{0xFF, 0xFF, 0xFF, 0xFF}, // 0.00
		{0x00, 0x66, 0xFF, 0xFF}, // 0.16
		{0x33, 0x33, 0x33, 0xFF}, // 0.50
		{0xFF, 0x00, 0xCC, 0xFF}, // 0.84
		{0xFF, 0xFF, 0xFF, 0xFF}, // 1.00
	}

	neonGradient = []LedColor{
		{0xff, 0xff, 0xff, 0xFF},
		{0x00, 0x99, 0xcc, 0xFF},
		{0x00, 0x00, 0x00, 0xFF},
		{0xff, 0x00, 0x99, 0xFF},
		{0xff, 0xff, 0xff, 0xFF},
	}

    fadeRedGradient = []LedColor{
		NewLedColor(0x000000),
		NewLedColor(0xff0000),
		NewLedColor(0x000000),
    }

    fadeGreenGradient = []LedColor{
		NewLedColor(0x00ff00),
		NewLedColor(0x000000),
		NewLedColor(0x00ff00),
    }

    fadeBlueGradient = []LedColor{
		NewLedColor(0x0000ff),
		NewLedColor(0x000000),
		NewLedColor(0x0000ff),
    }

	Pico08Colors = []LedColor{
		NewLedColor(0x000000),
		NewLedColor(0x1d2b53),
		NewLedColor(0x7e2553),
		NewLedColor(0x008751),
		NewLedColor(0xab5236),
		NewLedColor(0x5f574f),
		NewLedColor(0xc2c3c7),
		NewLedColor(0xfff1e8),
		NewLedColor(0xff004d),
		NewLedColor(0xffa300),
		NewLedColor(0xffff27),
		NewLedColor(0x00e756),
		NewLedColor(0x29adff),
		NewLedColor(0x83769c),
		NewLedColor(0xff77a8),
		NewLedColor(0xffccaa),
	}

	// Kopierte Paletten aus PixelController
	// aAAGradient = []LedColor{
	//     {0x00, 0xff, 0xff, 0xFF},
	//     {0x00, 0x00, 0xff, 0xFF},
	//     {0x00, 0x80, 0xff, 0xFF},
	//     {0x80, 0xff, 0xff, 0xFF},
	//     {0x00, 0x00, 0x80, 0xFF},
	// }
	// adribbleGradient = []LedColor{
	//     {0x3D, 0x4C, 0x53, 0xFF},
	//     {0x70, 0xB7, 0xBA, 0xFF},
	//     {0xF1, 0x43, 0x3F, 0xFF},
	//     {0xE7, 0xE1, 0xD4, 0xFF},
	//     {0xFF, 0xFF, 0xFF, 0xFF},
	// }
	// brazilGradient = []LedColor{
	//     {0x00, 0x8c, 0x53, 0xFF},
	//     {0x2e, 0x00, 0xe4, 0xFF},
	//     {0xdf, 0xea, 0x00, 0xFF},
	// }
	// bW01Gradient = []LedColor{
	//     {0x00, 0x00, 0x00, 0xFF},
	//     {0xff, 0xff, 0xff, 0xFF},
	// }
	// bW02Gradient = []LedColor{
	//     {0x00, 0x00, 0x00, 0xFF},
	//     {0xff, 0xff, 0xff, 0xFF},
	//     {0x00, 0x00, 0x00, 0xFF},
	//     {0xff, 0xff, 0xff, 0xFF},
	// }
	// bW03Gradient = []LedColor{
	//     {0x00, 0x00, 0x00, 0xFF},
	//     {0xff, 0xff, 0xff, 0xFF},
	//     {0x00, 0x00, 0x00, 0xFF},
	//     {0xff, 0xff, 0xff, 0xFF},
	//     {0x00, 0x00, 0x00, 0xFF},
	//     {0xff, 0xff, 0xff, 0xFF},
	//     {0x00, 0x00, 0x00, 0xFF},
	//     {0xff, 0xff, 0xff, 0xFF},
	// }
	// cakeGradient = []LedColor{
	//     {0x00, 0x00, 0x00, 0xFF},
	//     {0xBD, 0x41, 0x41, 0xFF},
	//     {0xD9, 0x75, 0x70, 0xFF},
	//     {0xF2, 0xE8, 0xDF, 0xFF},
	//     {0xB2, 0xCF, 0xC0, 0xFF},
	//     {0x71, 0x9D, 0x98, 0xFF},
	// }
	// castleGradient = []LedColor{
	//     {0x4B, 0x34, 0x5C, 0xFF},
	//     {0x94, 0x62, 0x82, 0xFF},
	//     {0xE5, 0xA1, 0x9B, 0xFF},
	// }
	// cGAGradient = []LedColor{
	//     {0xd3, 0x51, 0x7d, 0xFF},
	//     {0x15, 0xa0, 0xbf, 0xFF},
	//     {0xff, 0xc0, 0x62, 0xFF},
	// }
	// civilGradient = []LedColor{
	//     {0x36, 0x2F, 0x2D, 0xFF},
	//     {0x4C, 0x4C, 0x4C, 0xFF},
	//     {0x94, 0xB7, 0x3E, 0xFF},
	//     {0xB5, 0xC0, 0xAF, 0xFF},
	//     {0xFA, 0xFD, 0xF2, 0xFF},
	// }
	// cold00Gradient = []LedColor{
	//     {0x00, 0x00, 0x00, 0xFF},
	//     {0x00, 0xff, 0xff, 0xFF},
	// }
	// cold01Gradient = []LedColor{
	//     {0x00, 0x00, 0x00, 0xFF},
	//     {0x00, 0xff, 0xff, 0xFF},
	//     {0x00, 0x00, 0x80, 0xFF},
	// }
	// cold02Gradient = []LedColor{
	//     {0x00, 0x00, 0x00, 0xFF},
	//     {0x00, 0xff, 0xff, 0xFF},
	//     {0x00, 0x00, 0x80, 0xFF},
	//     {0x00, 0x80, 0x80, 0xFF},
	// }
	// cold03Gradient = []LedColor{
	//     {0x00, 0x00, 0x00, 0xFF},
	//     {0x00, 0xff, 0xff, 0xFF},
	//     {0x00, 0x00, 0x80, 0xFF},
	//     {0x00, 0x80, 0x80, 0xFF},
	//     {0x80, 0xff, 0xff, 0xFF},
	// }
	// cold04Gradient = []LedColor{
	//     {0x00, 0x00, 0x00, 0xFF},
	//     {0x00, 0xff, 0xff, 0xFF},
	//     {0x00, 0x00, 0x80, 0xFF},
	//     {0x00, 0x80, 0x80, 0xFF},
	//     {0x80, 0xff, 0xff, 0xFF},
	//     {0x80, 0x80, 0xff, 0xFF},
	// }
	// cold05Gradient = []LedColor{
	//     {0x80, 0xff, 0xff, 0xFF},
	//     {0x00, 0x80, 0xff, 0xFF},
	//     {0xc0, 0xc0, 0xc0, 0xFF},
	//     {0x80, 0xa3, 0xff, 0xFF},
	//     {0x58, 0x58, 0x58, 0xFF},
	//     {0xb9, 0xf4, 0xf9, 0xFF},
	// }
	// colorBirdGradient = []LedColor{
	//     {0x1F, 0xA6, 0x98, 0xFF},
	//     {0xC5, 0xD9, 0x32, 0xFF},
	//     {0xF2, 0x59, 0x22, 0xFF},
	//     {0x40, 0x1E, 0x11, 0xFF},
	//     {0xD7, 0x19, 0x5A, 0xFF},
	// }
	// cornGradient = []LedColor{
	//     {0x29, 0x23, 0x1F, 0xFF},
	//     {0xEB, 0xE1, 0xCC, 0xFF},
	//     {0xDB, 0x9B, 0x1A, 0xFF},
	// }
	// fadeAllGradient = []LedColor{
	//     {0xff, 0x00, 0x00, 0xFF},
	//     {0xff, 0xff, 0x00, 0xFF},
	//     {0x00, 0xff, 0x00, 0xFF},
	//     {0x00, 0xff, 0xff, 0xFF},
	//     {0x00, 0x00, 0xff, 0xFF},
	//     {0xff, 0x00, 0xff, 0xFF},
	// }
	// fadeRedGradient = []LedColor{
	//     {0x00, 0x00, 0x00, 0xFF},
	//     {0xff, 0x00, 0x00, 0xFF},
	// }
	// fadeGreenGradient = []LedColor{
	//     {0x00, 0x00, 0x00, 0xFF},
	//     {0x00, 0xff, 0x00, 0xFF},
	// }
	// fadeBlueGradient = []LedColor{
	//     {0x00, 0x00, 0x00, 0xFF},
	//     {0x00, 0x00, 0xff, 0xFF},
	// }
	// fadeYellowGradient = []LedColor{
	//     {0x00, 0x00, 0x00, 0xFF},
	//     {0xff, 0xff, 0x00, 0xFF},
	// }
	// fadeCyanGradient = []LedColor{
	//     {0x00, 0x00, 0x00, 0xFF},
	//     {0x00, 0xff, 0xff, 0xFF},
	// }
	// fadeMagentaGradient = []LedColor{
	//     {0x00, 0x00, 0x00, 0xFF},
	//     {0xff, 0x00, 0xff, 0xFF},
	// }
	// cycleRedGradient = []LedColor{
	//     {0x00, 0x00, 0x00, 0xFF},
	//     {0xff, 0x00, 0x00, 0xFF},
	//     {0x00, 0x00, 0x00, 0xFF},
	// }
	// fame575Gradient = []LedColor{
	//     {0x54, 0x0c, 0x0d, 0xFF},
	//     {0xfb, 0x74, 0x23, 0xFF},
	//     {0xf9, 0xf4, 0x8e, 0xFF},
	//     {0x41, 0x76, 0xc4, 0xFF},
	//     {0x5a, 0xaf, 0x2e, 0xFF},
	// }
	// fireGradient = []LedColor{
	//     {0x00, 0x00, 0x00, 0xFF},
	//     {0x00, 0x00, 0x40, 0xFF},
	//     {0xff, 0x00, 0x00, 0xFF},
	//     {0xff, 0xff, 0x00, 0xFF},
	//     {0xff, 0xff, 0xff, 0xFF},
	//     {0xff, 0xff, 0xff, 0xFF},
	// }
	// fizzGradient = []LedColor{
	//     {0x58, 0x8F, 0x27, 0xFF},
	//     {0x04, 0xBF, 0xBF, 0xFF},
	//     {0xF7, 0xE9, 0x67, 0xFF},
	// }
	// foldGradient = []LedColor{
	//     {0x2A, 0x03, 0x08, 0xFF},
	//     {0x92, 0x4F, 0x1B, 0xFF},
	//     {0xE2, 0xAC, 0x3F, 0xFF},
	//     {0xF8, 0xED, 0xC6, 0xFF},
	//     {0x7B, 0xA5, 0x8D, 0xFF},
	// }
	// hipsterGradient = []LedColor{
	//     {0x8D, 0x0E, 0x6B, 0xFF},
	//     {0x05, 0x34, 0x34, 0xFF},
	//     {0x00, 0x80, 0x80, 0xFF},
	//     {0x4D, 0x42, 0x08, 0xFF},
	//     {0xCD, 0xAD, 0x0A, 0xFF},
	// }
	// kittyHCGradient = []LedColor{
	//     {0xc7, 0x56, 0xa7, 0xFF},
	//     {0xe0, 0xdd, 0x00, 0xFF},
	//     {0xc9, 0xcd, 0xd0, 0xFF},
	// }
	// kittyGradient = []LedColor{
	//     {0x9f, 0x45, 0x6b, 0xFF},
	//     {0x4f, 0x7a, 0x9a, 0xFF},
	//     {0xe6, 0xc8, 0x4c, 0xFF},
	// }
	// lanternGradient = []LedColor{
	//     {0x00, 0x00, 0x00, 0xFF},
	//     {0x0D, 0x9A, 0x0D, 0xFF},
	//     {0xFF, 0xFF, 0xFF, 0xFF},
	// }
	// leBronGradient = []LedColor{
	//     {0x3e, 0x3e, 0x3e, 0xFF},
	//     {0xd4, 0xb6, 0x00, 0xFF},
	//     {0xff, 0xff, 0xff, 0xFF},
	// }
	// lemmingGradient = []LedColor{
	//     {0x00, 0x00, 0x00, 0xFF},
	//     {0x00, 0x00, 0xff, 0xFF},
	//     {0x00, 0xff, 0x00, 0xFF},
	//     {0xff, 0xff, 0xff, 0xFF},
	// }
	// miamiViceGradient = []LedColor{
	//     {0x1b, 0xe3, 0xff, 0xFF},
	//     {0xff, 0x82, 0xdc, 0xFF},
	//     {0xff, 0xff, 0xff, 0xFF},
	// }
	// mIUSAGradient = []LedColor{
	//     {0x50, 0x4b, 0x46, 0xFF},
	//     {0x1a, 0x3c, 0x53, 0xFF},
	//     {0xa0, 0x00, 0x28, 0xFF},
	// }
	// mL581ATGradient = []LedColor{
	//     {0x69, 0x96, 0x55, 0xFF},
	//     {0xf2, 0x6a, 0x36, 0xFF},
	//     {0xff, 0xff, 0xff, 0xFF},
	// }
	// newSeasonGradient = []LedColor{
	//     {0x4C, 0x39, 0x33, 0xFF},
	//     {0x00, 0x5B, 0x5B, 0xFF},
	//     {0xFE, 0x55, 0x02, 0xFF},
	//     {0xFE, 0xBF, 0x51, 0xFF},
	//     {0xF8, 0xEA, 0xB2, 0xFF},
	// }
	// nightspellGradient = []LedColor{
	//     {0xED, 0xEE, 0xBA, 0xFF},
	//     {0xFE, 0xA8, 0x1C, 0xFF},
	//     {0xB2, 0x01, 0x52, 0xFF},
	//     {0x4B, 0x0B, 0x44, 0xFF},
	//     {0x24, 0x0F, 0x37, 0xFF},
	// }
	// violetGradient = []LedColor{
	//     {0x4B, 0x0B, 0x44, 0xFF},
	//     {0x4B, 0x0B, 0x44, 0xFF},
	// }
	// rainbowsGradient = []LedColor{
	//     {0x49, 0x2D, 0x61, 0xFF},
	//     {0x04, 0x80, 0x91, 0xFF},
	//     {0x61, 0xC1, 0x55, 0xFF},
	//     {0xF2, 0xD4, 0x3F, 0xFF},
	//     {0xD1, 0x02, 0x6C, 0xFF},
	// }
	// rastaGradient = []LedColor{
	//     {0xdc, 0x32, 0x3c, 0xFF},
	//     {0xf0, 0xcb, 0x58, 0xFF},
	//     {0x3c, 0x82, 0x5e, 0xFF},
	// }
	// rGBGradient = []LedColor{
	//     {0xff, 0x00, 0x00, 0xFF},
	//     {0x00, 0xff, 0x00, 0xFF},
	//     {0x00, 0x00, 0xff, 0xFF},
	// }
	// simpsonGradient = []LedColor{
	//     {0xd9, 0xc2, 0x3e, 0xFF},
	//     {0xa9, 0x6a, 0x95, 0xFF},
	//     {0x7d, 0x95, 0x4b, 0xFF},
	//     {0x4b, 0x39, 0x6b, 0xFF},
	// }
	// smurfGradient = []LedColor{
	//     {0x1d, 0x16, 0x28, 0xFF},
	//     {0x44, 0xbd, 0xf4, 0xFF},
	//     {0xe3, 0x1e, 0x3a, 0xFF},
	//     {0xe8, 0xb1, 0x18, 0xFF},
	//     {0xff, 0xff, 0xff, 0xFF},
	// }
	// springGradient = []LedColor{
	//     {0x1D, 0x19, 0x29, 0xFF},
	//     {0xFE, 0x53, 0x24, 0xFF},
	//     {0xA8, 0x14, 0x3A, 0xFF},
	//     {0x66, 0xA5, 0x95, 0xFF},
	//     {0xCF, 0xBD, 0x81, 0xFF},
	// }
	// sunsetGradient = []LedColor{
	//     {0x00, 0x00, 0x00, 0xFF},
	//     {0x35, 0x24, 0x38, 0xFF},
	//     {0x61, 0x29, 0x49, 0xFF},
	//     {0x99, 0x28, 0x39, 0xFF},
	//     {0xED, 0x61, 0x37, 0xFF},
	//     {0xD6, 0xAB, 0x5C, 0xFF},
	// }
	// warmthGradient = []LedColor{
	//     {0xFC, 0xEB, 0xB6, 0xFF},
	//     {0x5E, 0x41, 0x2F, 0xFF},
	//     {0xF0, 0x78, 0x18, 0xFF},
	//     {0x78, 0xC0, 0xA8, 0xFF},
	//     {0xF0, 0xA8, 0x30, 0xFF},
	// }
	// wayyouGradient = []LedColor{
	//     {0x1C, 0x21, 0x30, 0xFF},
	//     {0x02, 0x8F, 0x76, 0xFF},
	//     {0xB3, 0xE0, 0x99, 0xFF},
	//     {0xFF, 0xEA, 0xAD, 0xFF},
	//     {0xD1, 0x43, 0x34, 0xFF},
	// }
	// wepartedGradient = []LedColor{
	//     {0x02, 0x7B, 0x7F, 0xFF},
	//     {0xFF, 0xA5, 0x88, 0xFF},
	//     {0xD6, 0x29, 0x57, 0xFF},
	//     {0xBF, 0x1E, 0x62, 0xFF},
	//     {0x57, 0x2E, 0x4F, 0xFF},
	// }
)