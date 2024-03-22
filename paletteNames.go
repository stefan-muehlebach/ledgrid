// ACHTUNG: dieses File wird automatisch erzeugt

package ledgrid

var (
	// PaletteNames ist ein Slice mit den Namen aller vorhandenen Paletten.
	PaletteNames = []string{
		"FractalDefault",
		"EarthAndSky",
		"Neon",
		"AAA",
		"Adribble",
		"Brazil",
		"Cake",
		"Castle",
		"CGA",
		"Civil",
		"ColorBird",
		"Corn",
		"FadeAll",
		"FadeRed",
		"FadeGreen",
		"FadeBlue",
		"FadeYellow",
		"FadeCyan",
		"FadeMagenta",
		"CycleRed",
		"Fire",
		"Fizz",
		"Fold",
		"Hipster",
		"KittyHC",
		"Kitty",
		"Lantern",
		"LeBron",
		"Lemming",
		"MiamiVice",
		"MIUSA",
		"NewSeason",
		"Nightspell",
		"Violet",
		"Rainbows",
		"Rasta",
		"RGB",
		"Simpson",
		"Smurf",
		"Spring",
		"Sunset",
		"Warmth",
		"Wayyou",
		"Weparted",
	}

	// PaletteMap ist die Verbindung zwischen Palettenname und einer Variable
	// vom Typ Palette.
	PaletteMap = map[string]*Palette{
		"FractalDefault": FractalDefault,
		"EarthAndSky":    EarthAndSky,
		"Neon":           Neon,
		"AAA":            AAA,
		"Adribble":       Adribble,
		"Brazil":         Brazil,
		"Cake":           Cake,
		"Castle":         Castle,
		"CGA":            CGA,
		"Civil":          Civil,
		"ColorBird":      ColorBird,
		"Corn":           Corn,
		"FadeAll":        FadeAll,
		"FadeRed":        FadeRed,
		"FadeGreen":      FadeGreen,
		"FadeBlue":       FadeBlue,
		"FadeYellow":     FadeYellow,
		"FadeCyan":       FadeCyan,
		"FadeMagenta":    FadeMagenta,
		"CycleRed":       CycleRed,
		"Fire":           Fire,
		"Fizz":           Fizz,
		"Fold":           Fold,
		"Hipster":        Hipster,
		"KittyHC":        KittyHC,
		"Kitty":          Kitty,
		"Lantern":        Lantern,
		"LeBron":         LeBron,
		"Lemming":        Lemming,
		"MiamiVice":      MiamiVice,
		"MIUSA":          MIUSA,
		"NewSeason":      NewSeason,
		"Nightspell":     Nightspell,
		"Violet":         Violet,
		"Rainbows":       Rainbows,
		"Rasta":          Rasta,
		"RGB":            RGB,
		"Simpson":        Simpson,
		"Smurf":          Smurf,
		"Spring":         Spring,
		"Sunset":         Sunset,
		"Warmth":         Warmth,
		"Wayyou":         Wayyou,
		"Weparted":       Weparted,
	}

	// In diesem Block werden die Paletten konkret erstellt. Im Moment
	// koennen so nur Paletten mit aequidistanten Farbstuetzstellen
	// erzeugt werden.
	FractalDefault = NewPaletteWithColors(fractalDefaultColors)
	EarthAndSky    = NewPaletteWithColors(earthAndSkyColors)
	Neon           = NewPaletteWithColors(neonColors)
	AAA            = NewPaletteWithColors(aAAColors)
	Adribble       = NewPaletteWithColors(adribbleColors)
	Brazil         = NewPaletteWithColors(brazilColors)
	Cake           = NewPaletteWithColors(cakeColors)
	Castle         = NewPaletteWithColors(castleColors)
	CGA            = NewPaletteWithColors(cGAColors)
	Civil          = NewPaletteWithColors(civilColors)
	ColorBird      = NewPaletteWithColors(colorBirdColors)
	Corn           = NewPaletteWithColors(cornColors)
	FadeAll        = NewPaletteWithColors(fadeAllColors)
	FadeRed        = NewPaletteWithColors(fadeRedColors)
	FadeGreen      = NewPaletteWithColors(fadeGreenColors)
	FadeBlue       = NewPaletteWithColors(fadeBlueColors)
	FadeYellow     = NewPaletteWithColors(fadeYellowColors)
	FadeCyan       = NewPaletteWithColors(fadeCyanColors)
	FadeMagenta    = NewPaletteWithColors(fadeMagentaColors)
	CycleRed       = NewPaletteWithColors(cycleRedColors)
	Fire           = NewPaletteWithColors(fireColors)
	Fizz           = NewPaletteWithColors(fizzColors)
	Fold           = NewPaletteWithColors(foldColors)
	Hipster        = NewPaletteWithColors(hipsterColors)
	KittyHC        = NewPaletteWithColors(kittyHCColors)
	Kitty          = NewPaletteWithColors(kittyColors)
	Lantern        = NewPaletteWithColors(lanternColors)
	LeBron         = NewPaletteWithColors(leBronColors)
	Lemming        = NewPaletteWithColors(lemmingColors)
	MiamiVice      = NewPaletteWithColors(miamiViceColors)
	MIUSA          = NewPaletteWithColors(mIUSAColors)
	NewSeason      = NewPaletteWithColors(newSeasonColors)
	Nightspell     = NewPaletteWithColors(nightspellColors)
	Violet         = NewPaletteWithColors(violetColors)
	Rainbows       = NewPaletteWithColors(rainbowsColors)
	Rasta          = NewPaletteWithColors(rastaColors)
	RGB            = NewPaletteWithColors(rGBColors)
	Simpson        = NewPaletteWithColors(simpsonColors)
	Smurf          = NewPaletteWithColors(smurfColors)
	Spring         = NewPaletteWithColors(springColors)
	Sunset         = NewPaletteWithColors(sunsetColors)
	Warmth         = NewPaletteWithColors(warmthColors)
	Wayyou         = NewPaletteWithColors(wayyouColors)
	Weparted       = NewPaletteWithColors(wepartedColors)
)
