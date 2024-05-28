//go:generate go run gen.go

package ledgrid

var (
	// Dies sind die Farblisten, welche fuer die einzelnen Paletten verwendet
	// werden. Diese Namen werden vom Programm 'gen' gelesen und zur
	// Erstellung des Files paletteNames.go verwendet. Die Namen der
	// Farblisten werden NICHT exportiert.
    //
    // Farblisten, deren Name mit 'Gradient' enden, haben als Elemente sog.
    // 'ColorStop's. Diese verbinden einen Fliesskommawert in [0,1] mit einem
    // konkreten Farbwert. Bei der Verwendung werden die Farben zwischen den
    // einzelnen Stops interpoliert.
    //
    // Farblisten, deren Namen mit 'ColorList' enden, enthalten nur Farbwerte
    // und verteilen die spezifizierten Farben aequidistant ueber dem
    // Intervall [0,1], wobei die erste Farbe in der Liste automatisch auch
    // als letzte Farbe der Palette verwendet wird.
    // Endet der Name ausserdem mit 'NonCyc', so wird dies nicht gemacht.
    //
    // Neben diesen automatisch verarbeiteten Farblisten, gibt es natürlich
    // die Möglichkeit, eigene Listen und Paletten zu definieren.

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

	fireGradient = []ColorStop{
		{0.00, NewLedColorAlpha(0x00000000)},
		{0.10, NewLedColorAlpha(0x5f080900)},
		{0.14, NewLedColorAlpha(0x5f0809e5)},
		{0.29, NewLedColor(0xbe1013)},
		{0.43, NewLedColor(0xd23008)},
		{0.57, NewLedColor(0xe45323)},
		{0.71, NewLedColor(0xee771c)},
		{0.86, NewLedColor(0xf6960e)},
		{1.00, NewLedColor(0xffcd06)},
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

	kittyColorList = []LedColor{
	    NewLedColor(0x9f456b),
	    NewLedColor(0x4f7a9a),
	    NewLedColor(0xe6c84c),
	}

	lanternColorList = []LedColor{
	    NewLedColor(0x000000),
	    NewLedColor(0x0D9A0D),
	    NewLedColor(0xFFFFFF),
	}

	lemmingColorList = []LedColor{
	    NewLedColor(0x000000),
	    NewLedColor(0x0000ff),
	    NewLedColor(0x00ff00),
	    NewLedColor(0xffffff),
	}

	miamiViceColorList = []LedColor{
	    NewLedColor(0x1be3ff),
	    NewLedColor(0xff82dc),
	    NewLedColor(0xffffff),
	}

	nightspellColorList = []LedColor{
	    NewLedColor(0xEDEEBA),
	    NewLedColor(0xFEA81C),
	    NewLedColor(0xB20152),
	    NewLedColor(0x4B0B44),
	    NewLedColor(0x240F37),
	}

	wayyouColorList = []LedColor{
	    NewLedColor(0x1C2130),
	    NewLedColor(0x028F76),
	    NewLedColor(0xB3E099),
	    NewLedColor(0xFFEAAD),
	    NewLedColor(0xD14334),
	}

	wepartedColorList = []LedColor{
	    NewLedColor(0x027B7F),
	    NewLedColor(0xFFA588),
	    NewLedColor(0xD62957),
	    NewLedColor(0xBF1E62),
	    NewLedColor(0x572E4F),
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
)
