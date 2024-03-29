
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
        "Darker",
        "Default",
        "EarthAndSky",
        "FadeBlue",
        "FadeGreen",
        "FadeRed",
        "HotAndCold",
        "Neon",
        "Pastell",
        "Seashore",
    }

    // PaletteMap ist die Verbindung zwischen Palettenname und einer Variable
    // vom Typ Palette.
    PaletteMap = map[string]Colorable{
        "Dark": Dark,
        "Darker": Darker,
        "Default": Default,
        "EarthAndSky": EarthAndSky,
        "FadeBlue": FadeBlue,
        "FadeGreen": FadeGreen,
        "FadeRed": FadeRed,
        "HotAndCold": HotAndCold,
        "Neon": Neon,
        "Pastell": Pastell,
        "Seashore": Seashore,
    }

    // In diesem Block werden die Paletten konkret erstellt. Im Moment
    // koennen so nur Paletten mit aequidistanten Farbstuetzstellen
    // erzeugt werden.
    Dark                 = NewGradientPalette(darkGradient...)
    Darker               = NewGradientPalette(darkerGradient...)
    Default              = NewGradientPalette(defaultGradient...)
    EarthAndSky          = NewGradientPalette(earthAndSkyGradient...)
    FadeBlue             = NewGradientPalette(fadeBlueGradient...)
    FadeGreen            = NewGradientPalette(fadeGreenGradient...)
    FadeRed              = NewGradientPalette(fadeRedGradient...)
    HotAndCold           = NewGradientPalette(hotAndColdGradient...)
    Neon                 = NewGradientPalette(neonGradient...)
    Pastell              = NewGradientPalette(pastellGradient...)
    Seashore             = NewGradientPalette(seashoreGradient...)
)
