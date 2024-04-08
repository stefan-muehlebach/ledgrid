
//----------------------------------------------------------------------------
//
//   ACHTUNG: dieses File wird automatisch durch das Tool 'gen' in diesem
//   Verzeichnis erzeugt! Manuelle Anpassungen koennen verloren gehen.
//
//----------------------------------------------------------------------------

package ledgrid

var (
    // PaletteNames ist ein Slice mit den Namen aller vorhandenen Paletten.
    PaletteList = []Colorable{
        BlackPalette,
        BlackWhitePalette,
        DarkPalette,
        DarkJunglePalette,
        DarkerPalette,
        DefaultPalette,
        EarthAndSkyPalette,
        FadeBluePalette,
        FadeGreenPalette,
        FadeRedPalette,
        FirePalette,
        HipsterPalette,
        HotAndColdPalette,
        HotFirePalette,
        NeonPalette,
        PastellPalette,
        RichDeepavaliPalette,
        SeashorePalette,
    }

    // PaletteMap ist die Verbindung zwischen Palettenname und einer Variable
    // vom Typ Palette.
    PaletteMap = map[string]Colorable{
        "Black": BlackPalette,
        "BlackWhite": BlackWhitePalette,
        "Dark": DarkPalette,
        "DarkJungle": DarkJunglePalette,
        "Darker": DarkerPalette,
        "Default": DefaultPalette,
        "EarthAndSky": EarthAndSkyPalette,
        "FadeBlue": FadeBluePalette,
        "FadeGreen": FadeGreenPalette,
        "FadeRed": FadeRedPalette,
        "Fire": FirePalette,
        "Hipster": HipsterPalette,
        "HotAndCold": HotAndColdPalette,
        "HotFire": HotFirePalette,
        "Neon": NeonPalette,
        "Pastell": PastellPalette,
        "RichDeepavali": RichDeepavaliPalette,
        "Seashore": SeashorePalette,
    }

    // In diesem Block werden die Paletten konkret erstellt. Im Moment
    // koennen so nur Paletten mit aequidistanten Farbstuetzstellen
    // erzeugt werden.
    BlackPalette = NewGradientPalette("Black", true, blackGradient...)
    BlackWhitePalette = NewGradientPalette("BlackWhite", false, blackWhiteGradientNoCycle...)
    DarkPalette = NewGradientPalette("Dark", true, darkGradient...)
    DarkJunglePalette = NewGradientPalette("DarkJungle", true, darkJungleGradient...)
    DarkerPalette = NewGradientPalette("Darker", true, darkerGradient...)
    DefaultPalette = NewGradientPalette("Default", true, defaultGradient...)
    EarthAndSkyPalette = NewGradientPalette("EarthAndSky", true, earthAndSkyGradient...)
    FadeBluePalette = NewGradientPalette("FadeBlue", true, fadeBlueGradient...)
    FadeGreenPalette = NewGradientPalette("FadeGreen", true, fadeGreenGradient...)
    FadeRedPalette = NewGradientPalette("FadeRed", true, fadeRedGradient...)
    FirePalette = NewGradientPalette("Fire", false, fireGradientNoCycle...)
    HipsterPalette = NewGradientPalette("Hipster", true, hipsterGradient...)
    HotAndColdPalette = NewGradientPalette("HotAndCold", true, hotAndColdGradient...)
    HotFirePalette = NewGradientPalette("HotFire", true, hotFireGradient...)
    NeonPalette = NewGradientPalette("Neon", true, neonGradient...)
    PastellPalette = NewGradientPalette("Pastell", true, pastellGradient...)
    RichDeepavaliPalette = NewGradientPalette("RichDeepavali", true, richDeepavaliGradient...)
    SeashorePalette = NewGradientPalette("Seashore", true, seashoreGradient...)
)
