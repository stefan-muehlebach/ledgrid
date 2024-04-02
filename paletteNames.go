
//----------------------------------------------------------------------------
//
//   ACHTUNG: dieses File wird automatisch durch das Tool 'gen' in diesem
//   Verzeichnis erzeugt! Manuelle Anpassungen koennen verloren gehen.
//
//----------------------------------------------------------------------------

package ledgrid

var (
    // PaletteNames ist ein Slice mit den Namen aller vorhandenen Paletten.
    PaletteNames = []string{
        "Dark",
        "DarkJungle",
        "Darker",
        "Default",
        "EarthAndSky",
        "FadeBlue",
        "FadeGreen",
        "FadeRed",
        "Hipster",
        "HotAndCold",
        "HotFire",
        "Neon",
        "Pastell",
        "RichDeepavali",
        "Seashore",
    }

    // PaletteMap ist die Verbindung zwischen Palettenname und einer Variable
    // vom Typ Palette.
    PaletteMap = map[string]Colorable{
        "Dark": Dark,
        "DarkJungle": DarkJungle,
        "Darker": Darker,
        "Default": Default,
        "EarthAndSky": EarthAndSky,
        "FadeBlue": FadeBlue,
        "FadeGreen": FadeGreen,
        "FadeRed": FadeRed,
        "Hipster": Hipster,
        "HotAndCold": HotAndCold,
        "HotFire": HotFire,
        "Neon": Neon,
        "Pastell": Pastell,
        "RichDeepavali": RichDeepavali,
        "Seashore": Seashore,
    }

    // In diesem Block werden die Paletten konkret erstellt. Im Moment
    // koennen so nur Paletten mit aequidistanten Farbstuetzstellen
    // erzeugt werden.
    Dark                 = NewGradientPalette(true, darkGradient...)
    DarkJungle           = NewGradientPalette(true, darkJungleGradient...)
    Darker               = NewGradientPalette(true, darkerGradient...)
    Default              = NewGradientPalette(true, defaultGradient...)
    EarthAndSky          = NewGradientPalette(true, earthAndSkyGradient...)
    FadeBlue             = NewGradientPalette(true, fadeBlueGradient...)
    FadeGreen            = NewGradientPalette(true, fadeGreenGradient...)
    FadeRed              = NewGradientPalette(true, fadeRedGradient...)
    Hipster              = NewGradientPalette(true, hipsterGradient...)
    HotAndCold           = NewGradientPalette(true, hotAndColdGradient...)
    HotFire              = NewGradientPalette(true, hotFireGradient...)
    Neon                 = NewGradientPalette(true, neonGradient...)
    Pastell              = NewGradientPalette(true, pastellGradient...)
    RichDeepavali        = NewGradientPalette(true, richDeepavaliGradient...)
    Seashore             = NewGradientPalette(true, seashoreGradient...)
)
