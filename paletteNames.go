
// ACHTUNG: dieses File wird automatisch durch das Tool 'gen' in diesem
// Verzeichnis erzeugt! Manuelle Anpassungen koennen verloren gehen.

package ledgrid

var (
    // PaletteNames ist ein Slice mit den Namen aller vorhandenen Paletten.
    PaletteList = []Colorable{
        BlackAndWhitePalette,
        DarkPalette,
        DarkJunglePalette,
        DarkerPalette,
        DefaultPalette,
        EarthAndSkyPalette,
        FadeBluePalette,
        FadeCyanPalette,
        FadeGreenPalette,
        FadeMagentaPalette,
        FadeRedPalette,
        FadeYellowPalette,
        FirePalette,
        HipsterPalette,
        HotAndColdPalette,
        HotFirePalette,
        NeonPalette,
        PastellPalette,
        SeashorePalette,
    }

    // In diesem Block werden die Paletten konkret erstellt. Im Moment
    // koennen so nur Paletten mit aequidistanten Farbstuetzstellen
    // erzeugt werden.
    BlackAndWhitePalette = NewGradientPalette("BlackAndWhite", blackAndWhiteGradient...)
    DarkPalette = NewGradientPalette("Dark", darkGradient...)
    DarkJunglePalette = NewGradientPaletteByList("DarkJungle", true, darkJungleColorList...)
    DarkerPalette = NewGradientPalette("Darker", darkerGradient...)
    DefaultPalette = NewGradientPalette("Default", defaultGradient...)
    EarthAndSkyPalette = NewGradientPalette("EarthAndSky", earthAndSkyGradient...)
    FadeBluePalette = NewGradientPaletteByList("FadeBlue", false, fadeBlueColorListNonCyc...)
    FadeCyanPalette = NewGradientPaletteByList("FadeCyan", false, fadeCyanColorListNonCyc...)
    FadeGreenPalette = NewGradientPaletteByList("FadeGreen", false, fadeGreenColorListNonCyc...)
    FadeMagentaPalette = NewGradientPaletteByList("FadeMagenta", false, fadeMagentaColorListNonCyc...)
    FadeRedPalette = NewGradientPaletteByList("FadeRed", false, fadeRedColorListNonCyc...)
    FadeYellowPalette = NewGradientPaletteByList("FadeYellow", false, fadeYellowColorListNonCyc...)
    FirePalette = NewGradientPalette("Fire", fireGradient...)
    HipsterPalette = NewGradientPaletteByList("Hipster", true, hipsterColorList...)
    HotAndColdPalette = NewGradientPalette("HotAndCold", hotAndColdGradient...)
    HotFirePalette = NewGradientPaletteByList("HotFire", true, hotFireColorList...)
    NeonPalette = NewGradientPaletteByList("Neon", true, neonColorList...)
    PastellPalette = NewGradientPalette("Pastell", pastellGradient...)
    SeashorePalette = NewGradientPalette("Seashore", seashoreGradient...)
)
