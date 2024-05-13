
// ACHTUNG: dieses File wird automatisch durch das Tool 'gen' in diesem
// Verzeichnis erzeugt! Manuelle Anpassungen koennen verloren gehen.

package ledgrid

import (
    "github.com/stefan-muehlebach/gg/colornames"
)

var (
    // PaletteList ist ein Slice mit den Namen aller vorhandenen Paletten.
    PaletteList = []ColorSource{
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
        KittyPalette,
        LanternPalette,
        LemmingPalette,
        MiamiVicePalette,
        NeonPalette,
        NightspellPalette,
        PastellPalette,
        SeashorePalette,
        WayyouPalette,
        WepartedPalette,
    }
)

var (
    // In diesem Block werden die Paletten konkret erstellt.
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
    KittyPalette = NewGradientPaletteByList("Kitty", true, kittyColorList...)
    LanternPalette = NewGradientPaletteByList("Lantern", true, lanternColorList...)
    LemmingPalette = NewGradientPaletteByList("Lemming", true, lemmingColorList...)
    MiamiVicePalette = NewGradientPaletteByList("MiamiVice", true, miamiViceColorList...)
    NeonPalette = NewGradientPaletteByList("Neon", true, neonColorList...)
    NightspellPalette = NewGradientPaletteByList("Nightspell", true, nightspellColorList...)
    PastellPalette = NewGradientPalette("Pastell", pastellGradient...)
    SeashorePalette = NewGradientPalette("Seashore", seashoreGradient...)
    WayyouPalette = NewGradientPaletteByList("Wayyou", true, wayyouColorList...)
    WepartedPalette = NewGradientPaletteByList("Weparted", true, wepartedColorList...)
)

var (
    ColorList = []ColorSource{
        AliceBlueColor,
        AntiqueWhiteColor,
        AquaColor,
        AquamarineColor,
        AzureColor,
        BeigeColor,
        BisqueColor,
        BlackColor,
        BlanchedAlmondColor,
        BlueColor,
        BlueVioletColor,
        BrownColor,
        BurlyWoodColor,
        CadetBlueColor,
        ChartreuseColor,
        ChocolateColor,
        CoralColor,
        CornflowerBlueColor,
        CornsilkColor,
        CrimsonColor,
        CyanColor,
        DarkBlueColor,
        DarkCyanColor,
        DarkGoldenrodColor,
        DarkGrayColor,
        DarkGreenColor,
        DarkGreyColor,
        DarkKhakiColor,
        DarkMagentaColor,
        DarkOliveGreenColor,
        DarkOrangeColor,
        DarkOrchidColor,
        DarkRedColor,
        DarkSalmonColor,
        DarkSeaGreenColor,
        DarkSlateBlueColor,
        DarkSlateGrayColor,
        DarkSlateGreyColor,
        DarkTurquoiseColor,
        DarkVioletColor,
        DeepPinkColor,
        DeepSkyBlueColor,
        DimGrayColor,
        DimGreyColor,
        DodgerBlueColor,
        FireBrickColor,
        FloralWhiteColor,
        ForestGreenColor,
        FuchsiaColor,
        GainsboroColor,
        GhostWhiteColor,
        GoldColor,
        GoldenrodColor,
        GrayColor,
        GreenColor,
        GreenYellowColor,
        GreyColor,
        HoneydewColor,
        HotPinkColor,
        IndianRedColor,
        IndigoColor,
        IvoryColor,
        KhakiColor,
        LavenderColor,
        LavenderBlushColor,
        LawnGreenColor,
        LemonChiffonColor,
        LightBlueColor,
        LightCoralColor,
        LightCyanColor,
        LightGoldenrodYellowColor,
        LightGrayColor,
        LightGreenColor,
        LightGreyColor,
        LightPinkColor,
        LightSalmonColor,
        LightSeaGreenColor,
        LightSkyBlueColor,
        LightSlateGrayColor,
        LightSlateGreyColor,
        LightSteelBlueColor,
        LightYellowColor,
        LimeColor,
        LimeGreenColor,
        LinenColor,
        MagentaColor,
        MaroonColor,
        MediumAquamarineColor,
        MediumBlueColor,
        MediumOrchidColor,
        MediumPurpleColor,
        MediumSeaGreenColor,
        MediumSlateBlueColor,
        MediumSpringGreenColor,
        MediumTurquoiseColor,
        MediumVioletRedColor,
        MidnightBlueColor,
        MintCreamColor,
        MistyRoseColor,
        MoccasinColor,
        NavajoWhiteColor,
        NavyColor,
        OldLaceColor,
        OliveColor,
        OliveDrabColor,
        OrangeColor,
        OrangeRedColor,
        OrchidColor,
        PaleGoldenrodColor,
        PaleGreenColor,
        PaleTurquoiseColor,
        PaleVioletRedColor,
        PapayaWhipColor,
        PeachPuffColor,
        PeruColor,
        PinkColor,
        PlumColor,
        PowderBlueColor,
        PurpleColor,
        RedColor,
        RosyBrownColor,
        RoyalBlueColor,
        SaddleBrownColor,
        SalmonColor,
        SandyBrownColor,
        SeaGreenColor,
        SeashellColor,
        SiennaColor,
        SilverColor,
        SkyBlueColor,
        SlateBlueColor,
        SlateGrayColor,
        SlateGreyColor,
        SnowColor,
        SpringGreenColor,
        SteelBlueColor,
        TanColor,
        TealColor,
        ThistleColor,
        TomatoColor,
        TurquoiseColor,
        VioletColor,
        WheatColor,
        WhiteColor,
        WhiteSmokeColor,
        YellowColor,
        YellowGreenColor,
    }
)

var (
    // In diesem Block werden die uniformen Paletten erstellt.
    AliceBlueColor = NewUniformPalette("AliceBlue", colornames.AliceBlue)
    AntiqueWhiteColor = NewUniformPalette("AntiqueWhite", colornames.AntiqueWhite)
    AquaColor = NewUniformPalette("Aqua", colornames.Aqua)
    AquamarineColor = NewUniformPalette("Aquamarine", colornames.Aquamarine)
    AzureColor = NewUniformPalette("Azure", colornames.Azure)
    BeigeColor = NewUniformPalette("Beige", colornames.Beige)
    BisqueColor = NewUniformPalette("Bisque", colornames.Bisque)
    BlackColor = NewUniformPalette("Black", colornames.Black)
    BlanchedAlmondColor = NewUniformPalette("BlanchedAlmond", colornames.BlanchedAlmond)
    BlueColor = NewUniformPalette("Blue", colornames.Blue)
    BlueVioletColor = NewUniformPalette("BlueViolet", colornames.BlueViolet)
    BrownColor = NewUniformPalette("Brown", colornames.Brown)
    BurlyWoodColor = NewUniformPalette("BurlyWood", colornames.BurlyWood)
    CadetBlueColor = NewUniformPalette("CadetBlue", colornames.CadetBlue)
    ChartreuseColor = NewUniformPalette("Chartreuse", colornames.Chartreuse)
    ChocolateColor = NewUniformPalette("Chocolate", colornames.Chocolate)
    CoralColor = NewUniformPalette("Coral", colornames.Coral)
    CornflowerBlueColor = NewUniformPalette("CornflowerBlue", colornames.CornflowerBlue)
    CornsilkColor = NewUniformPalette("Cornsilk", colornames.Cornsilk)
    CrimsonColor = NewUniformPalette("Crimson", colornames.Crimson)
    CyanColor = NewUniformPalette("Cyan", colornames.Cyan)
    DarkBlueColor = NewUniformPalette("DarkBlue", colornames.DarkBlue)
    DarkCyanColor = NewUniformPalette("DarkCyan", colornames.DarkCyan)
    DarkGoldenrodColor = NewUniformPalette("DarkGoldenrod", colornames.DarkGoldenrod)
    DarkGrayColor = NewUniformPalette("DarkGray", colornames.DarkGray)
    DarkGreenColor = NewUniformPalette("DarkGreen", colornames.DarkGreen)
    DarkGreyColor = NewUniformPalette("DarkGrey", colornames.DarkGrey)
    DarkKhakiColor = NewUniformPalette("DarkKhaki", colornames.DarkKhaki)
    DarkMagentaColor = NewUniformPalette("DarkMagenta", colornames.DarkMagenta)
    DarkOliveGreenColor = NewUniformPalette("DarkOliveGreen", colornames.DarkOliveGreen)
    DarkOrangeColor = NewUniformPalette("DarkOrange", colornames.DarkOrange)
    DarkOrchidColor = NewUniformPalette("DarkOrchid", colornames.DarkOrchid)
    DarkRedColor = NewUniformPalette("DarkRed", colornames.DarkRed)
    DarkSalmonColor = NewUniformPalette("DarkSalmon", colornames.DarkSalmon)
    DarkSeaGreenColor = NewUniformPalette("DarkSeaGreen", colornames.DarkSeaGreen)
    DarkSlateBlueColor = NewUniformPalette("DarkSlateBlue", colornames.DarkSlateBlue)
    DarkSlateGrayColor = NewUniformPalette("DarkSlateGray", colornames.DarkSlateGray)
    DarkSlateGreyColor = NewUniformPalette("DarkSlateGrey", colornames.DarkSlateGrey)
    DarkTurquoiseColor = NewUniformPalette("DarkTurquoise", colornames.DarkTurquoise)
    DarkVioletColor = NewUniformPalette("DarkViolet", colornames.DarkViolet)
    DeepPinkColor = NewUniformPalette("DeepPink", colornames.DeepPink)
    DeepSkyBlueColor = NewUniformPalette("DeepSkyBlue", colornames.DeepSkyBlue)
    DimGrayColor = NewUniformPalette("DimGray", colornames.DimGray)
    DimGreyColor = NewUniformPalette("DimGrey", colornames.DimGrey)
    DodgerBlueColor = NewUniformPalette("DodgerBlue", colornames.DodgerBlue)
    FireBrickColor = NewUniformPalette("FireBrick", colornames.FireBrick)
    FloralWhiteColor = NewUniformPalette("FloralWhite", colornames.FloralWhite)
    ForestGreenColor = NewUniformPalette("ForestGreen", colornames.ForestGreen)
    FuchsiaColor = NewUniformPalette("Fuchsia", colornames.Fuchsia)
    GainsboroColor = NewUniformPalette("Gainsboro", colornames.Gainsboro)
    GhostWhiteColor = NewUniformPalette("GhostWhite", colornames.GhostWhite)
    GoldColor = NewUniformPalette("Gold", colornames.Gold)
    GoldenrodColor = NewUniformPalette("Goldenrod", colornames.Goldenrod)
    GrayColor = NewUniformPalette("Gray", colornames.Gray)
    GreenColor = NewUniformPalette("Green", colornames.Green)
    GreenYellowColor = NewUniformPalette("GreenYellow", colornames.GreenYellow)
    GreyColor = NewUniformPalette("Grey", colornames.Grey)
    HoneydewColor = NewUniformPalette("Honeydew", colornames.Honeydew)
    HotPinkColor = NewUniformPalette("HotPink", colornames.HotPink)
    IndianRedColor = NewUniformPalette("IndianRed", colornames.IndianRed)
    IndigoColor = NewUniformPalette("Indigo", colornames.Indigo)
    IvoryColor = NewUniformPalette("Ivory", colornames.Ivory)
    KhakiColor = NewUniformPalette("Khaki", colornames.Khaki)
    LavenderColor = NewUniformPalette("Lavender", colornames.Lavender)
    LavenderBlushColor = NewUniformPalette("LavenderBlush", colornames.LavenderBlush)
    LawnGreenColor = NewUniformPalette("LawnGreen", colornames.LawnGreen)
    LemonChiffonColor = NewUniformPalette("LemonChiffon", colornames.LemonChiffon)
    LightBlueColor = NewUniformPalette("LightBlue", colornames.LightBlue)
    LightCoralColor = NewUniformPalette("LightCoral", colornames.LightCoral)
    LightCyanColor = NewUniformPalette("LightCyan", colornames.LightCyan)
    LightGoldenrodYellowColor = NewUniformPalette("LightGoldenrodYellow", colornames.LightGoldenrodYellow)
    LightGrayColor = NewUniformPalette("LightGray", colornames.LightGray)
    LightGreenColor = NewUniformPalette("LightGreen", colornames.LightGreen)
    LightGreyColor = NewUniformPalette("LightGrey", colornames.LightGrey)
    LightPinkColor = NewUniformPalette("LightPink", colornames.LightPink)
    LightSalmonColor = NewUniformPalette("LightSalmon", colornames.LightSalmon)
    LightSeaGreenColor = NewUniformPalette("LightSeaGreen", colornames.LightSeaGreen)
    LightSkyBlueColor = NewUniformPalette("LightSkyBlue", colornames.LightSkyBlue)
    LightSlateGrayColor = NewUniformPalette("LightSlateGray", colornames.LightSlateGray)
    LightSlateGreyColor = NewUniformPalette("LightSlateGrey", colornames.LightSlateGrey)
    LightSteelBlueColor = NewUniformPalette("LightSteelBlue", colornames.LightSteelBlue)
    LightYellowColor = NewUniformPalette("LightYellow", colornames.LightYellow)
    LimeColor = NewUniformPalette("Lime", colornames.Lime)
    LimeGreenColor = NewUniformPalette("LimeGreen", colornames.LimeGreen)
    LinenColor = NewUniformPalette("Linen", colornames.Linen)
    MagentaColor = NewUniformPalette("Magenta", colornames.Magenta)
    MaroonColor = NewUniformPalette("Maroon", colornames.Maroon)
    MediumAquamarineColor = NewUniformPalette("MediumAquamarine", colornames.MediumAquamarine)
    MediumBlueColor = NewUniformPalette("MediumBlue", colornames.MediumBlue)
    MediumOrchidColor = NewUniformPalette("MediumOrchid", colornames.MediumOrchid)
    MediumPurpleColor = NewUniformPalette("MediumPurple", colornames.MediumPurple)
    MediumSeaGreenColor = NewUniformPalette("MediumSeaGreen", colornames.MediumSeaGreen)
    MediumSlateBlueColor = NewUniformPalette("MediumSlateBlue", colornames.MediumSlateBlue)
    MediumSpringGreenColor = NewUniformPalette("MediumSpringGreen", colornames.MediumSpringGreen)
    MediumTurquoiseColor = NewUniformPalette("MediumTurquoise", colornames.MediumTurquoise)
    MediumVioletRedColor = NewUniformPalette("MediumVioletRed", colornames.MediumVioletRed)
    MidnightBlueColor = NewUniformPalette("MidnightBlue", colornames.MidnightBlue)
    MintCreamColor = NewUniformPalette("MintCream", colornames.MintCream)
    MistyRoseColor = NewUniformPalette("MistyRose", colornames.MistyRose)
    MoccasinColor = NewUniformPalette("Moccasin", colornames.Moccasin)
    NavajoWhiteColor = NewUniformPalette("NavajoWhite", colornames.NavajoWhite)
    NavyColor = NewUniformPalette("Navy", colornames.Navy)
    OldLaceColor = NewUniformPalette("OldLace", colornames.OldLace)
    OliveColor = NewUniformPalette("Olive", colornames.Olive)
    OliveDrabColor = NewUniformPalette("OliveDrab", colornames.OliveDrab)
    OrangeColor = NewUniformPalette("Orange", colornames.Orange)
    OrangeRedColor = NewUniformPalette("OrangeRed", colornames.OrangeRed)
    OrchidColor = NewUniformPalette("Orchid", colornames.Orchid)
    PaleGoldenrodColor = NewUniformPalette("PaleGoldenrod", colornames.PaleGoldenrod)
    PaleGreenColor = NewUniformPalette("PaleGreen", colornames.PaleGreen)
    PaleTurquoiseColor = NewUniformPalette("PaleTurquoise", colornames.PaleTurquoise)
    PaleVioletRedColor = NewUniformPalette("PaleVioletRed", colornames.PaleVioletRed)
    PapayaWhipColor = NewUniformPalette("PapayaWhip", colornames.PapayaWhip)
    PeachPuffColor = NewUniformPalette("PeachPuff", colornames.PeachPuff)
    PeruColor = NewUniformPalette("Peru", colornames.Peru)
    PinkColor = NewUniformPalette("Pink", colornames.Pink)
    PlumColor = NewUniformPalette("Plum", colornames.Plum)
    PowderBlueColor = NewUniformPalette("PowderBlue", colornames.PowderBlue)
    PurpleColor = NewUniformPalette("Purple", colornames.Purple)
    RedColor = NewUniformPalette("Red", colornames.Red)
    RosyBrownColor = NewUniformPalette("RosyBrown", colornames.RosyBrown)
    RoyalBlueColor = NewUniformPalette("RoyalBlue", colornames.RoyalBlue)
    SaddleBrownColor = NewUniformPalette("SaddleBrown", colornames.SaddleBrown)
    SalmonColor = NewUniformPalette("Salmon", colornames.Salmon)
    SandyBrownColor = NewUniformPalette("SandyBrown", colornames.SandyBrown)
    SeaGreenColor = NewUniformPalette("SeaGreen", colornames.SeaGreen)
    SeashellColor = NewUniformPalette("Seashell", colornames.Seashell)
    SiennaColor = NewUniformPalette("Sienna", colornames.Sienna)
    SilverColor = NewUniformPalette("Silver", colornames.Silver)
    SkyBlueColor = NewUniformPalette("SkyBlue", colornames.SkyBlue)
    SlateBlueColor = NewUniformPalette("SlateBlue", colornames.SlateBlue)
    SlateGrayColor = NewUniformPalette("SlateGray", colornames.SlateGray)
    SlateGreyColor = NewUniformPalette("SlateGrey", colornames.SlateGrey)
    SnowColor = NewUniformPalette("Snow", colornames.Snow)
    SpringGreenColor = NewUniformPalette("SpringGreen", colornames.SpringGreen)
    SteelBlueColor = NewUniformPalette("SteelBlue", colornames.SteelBlue)
    TanColor = NewUniformPalette("Tan", colornames.Tan)
    TealColor = NewUniformPalette("Teal", colornames.Teal)
    ThistleColor = NewUniformPalette("Thistle", colornames.Thistle)
    TomatoColor = NewUniformPalette("Tomato", colornames.Tomato)
    TurquoiseColor = NewUniformPalette("Turquoise", colornames.Turquoise)
    VioletColor = NewUniformPalette("Violet", colornames.Violet)
    WheatColor = NewUniformPalette("Wheat", colornames.Wheat)
    WhiteColor = NewUniformPalette("White", colornames.White)
    WhiteSmokeColor = NewUniformPalette("WhiteSmoke", colornames.WhiteSmoke)
    YellowColor = NewUniformPalette("Yellow", colornames.Yellow)
    YellowGreenColor = NewUniformPalette("YellowGreen", colornames.YellowGreen)
)
