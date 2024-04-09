
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
        BrownishPalette,
        DarkPalette,
        DarkJunglePalette,
        DarkerPalette,
        DefaultPalette,
        DusterPalette,
        EarthAndSkyPalette,
        FadeBluePalette,
        FadeGreenPalette,
        FadeRedPalette,
        FirePalette,
        HipsterPalette,
        HotAndColdPalette,
        HotFirePalette,
        IcecreamPalette,
        NeonPalette,
        PastellPalette,
        RichDeepavaliPalette,
        SeashorePalette,
    }

    // In diesem Block werden die Paletten konkret erstellt. Im Moment
    // koennen so nur Paletten mit aequidistanten Farbstuetzstellen
    // erzeugt werden.
    BlackPalette = NewGradientPalette("Black", true, blackGradient...)
    BlackWhitePalette = NewGradientPalette("BlackWhite", false, blackWhiteGradientNoCycle...)
    BrownishPalette = NewGradientPalette("Brownish", true, brownishGradient...)
    DarkPalette = NewGradientPalette("Dark", true, darkGradient...)
    DarkJunglePalette = NewGradientPalette("DarkJungle", true, darkJungleGradient...)
    DarkerPalette = NewGradientPalette("Darker", true, darkerGradient...)
    DefaultPalette = NewGradientPalette("Default", true, defaultGradient...)
    DusterPalette = NewGradientPalette("Duster", true, dusterGradient...)
    EarthAndSkyPalette = NewGradientPalette("EarthAndSky", true, earthAndSkyGradient...)
    FadeBluePalette = NewGradientPalette("FadeBlue", true, fadeBlueGradient...)
    FadeGreenPalette = NewGradientPalette("FadeGreen", true, fadeGreenGradient...)
    FadeRedPalette = NewGradientPalette("FadeRed", true, fadeRedGradient...)
    FirePalette = NewGradientPalette("Fire", false, fireGradientNoCycle...)
    HipsterPalette = NewGradientPalette("Hipster", true, hipsterGradient...)
    HotAndColdPalette = NewGradientPalette("HotAndCold", true, hotAndColdGradient...)
    HotFirePalette = NewGradientPalette("HotFire", true, hotFireGradient...)
    IcecreamPalette = NewGradientPalette("Icecream", true, icecreamGradient...)
    NeonPalette = NewGradientPalette("Neon", true, neonGradient...)
    PastellPalette = NewGradientPalette("Pastell", true, pastellGradient...)
    RichDeepavaliPalette = NewGradientPalette("RichDeepavali", true, richDeepavaliGradient...)
    SeashorePalette = NewGradientPalette("Seashore", true, seashoreGradient...)
)
