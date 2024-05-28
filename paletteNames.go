
// ACHTUNG: dieses File wird automatisch durch das Tool 'gen' in diesem
// Verzeichnis erzeugt! Manuelle Anpassungen koennen verloren gehen.

package ledgrid

import (
    "github.com/stefan-muehlebach/gg/colornames"
)

var (
    // PaletteList ist ein Slice mit allen vorhandenen Paletten.
    PaletteList = []ColorSource{
        BlackAndWhitePalette,
        DarkJunglePalette,
        DarkPalette,
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
        LavenderBlushColor,
        LavenderColor,
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

    // PaletteMap ist ein Map um Paletten mit ihrem Namen anzusprechen.
    PaletteMap = map[string]ColorSource{
        "BlackAndWhitePalette": BlackAndWhitePalette,
        "DarkJunglePalette": DarkJunglePalette,
        "DarkPalette": DarkPalette,
        "DarkerPalette": DarkerPalette,
        "DefaultPalette": DefaultPalette,
        "EarthAndSkyPalette": EarthAndSkyPalette,
        "FadeBluePalette": FadeBluePalette,
        "FadeCyanPalette": FadeCyanPalette,
        "FadeGreenPalette": FadeGreenPalette,
        "FadeMagentaPalette": FadeMagentaPalette,
        "FadeRedPalette": FadeRedPalette,
        "FadeYellowPalette": FadeYellowPalette,
        "FirePalette": FirePalette,
        "HipsterPalette": HipsterPalette,
        "HotAndColdPalette": HotAndColdPalette,
        "HotFirePalette": HotFirePalette,
        "KittyPalette": KittyPalette,
        "LanternPalette": LanternPalette,
        "LemmingPalette": LemmingPalette,
        "MiamiVicePalette": MiamiVicePalette,
        "NeonPalette": NeonPalette,
        "NightspellPalette": NightspellPalette,
        "PastellPalette": PastellPalette,
        "SeashorePalette": SeashorePalette,
        "WayyouPalette": WayyouPalette,
        "WepartedPalette": WepartedPalette,
        "AliceBlueColor": AliceBlueColor,
        "AntiqueWhiteColor": AntiqueWhiteColor,
        "AquaColor": AquaColor,
        "AquamarineColor": AquamarineColor,
        "AzureColor": AzureColor,
        "BeigeColor": BeigeColor,
        "BisqueColor": BisqueColor,
        "BlackColor": BlackColor,
        "BlanchedAlmondColor": BlanchedAlmondColor,
        "BlueColor": BlueColor,
        "BlueVioletColor": BlueVioletColor,
        "BrownColor": BrownColor,
        "BurlyWoodColor": BurlyWoodColor,
        "CadetBlueColor": CadetBlueColor,
        "ChartreuseColor": ChartreuseColor,
        "ChocolateColor": ChocolateColor,
        "CoralColor": CoralColor,
        "CornflowerBlueColor": CornflowerBlueColor,
        "CornsilkColor": CornsilkColor,
        "CrimsonColor": CrimsonColor,
        "CyanColor": CyanColor,
        "DarkBlueColor": DarkBlueColor,
        "DarkCyanColor": DarkCyanColor,
        "DarkGoldenrodColor": DarkGoldenrodColor,
        "DarkGrayColor": DarkGrayColor,
        "DarkGreenColor": DarkGreenColor,
        "DarkGreyColor": DarkGreyColor,
        "DarkKhakiColor": DarkKhakiColor,
        "DarkMagentaColor": DarkMagentaColor,
        "DarkOliveGreenColor": DarkOliveGreenColor,
        "DarkOrangeColor": DarkOrangeColor,
        "DarkOrchidColor": DarkOrchidColor,
        "DarkRedColor": DarkRedColor,
        "DarkSalmonColor": DarkSalmonColor,
        "DarkSeaGreenColor": DarkSeaGreenColor,
        "DarkSlateBlueColor": DarkSlateBlueColor,
        "DarkSlateGrayColor": DarkSlateGrayColor,
        "DarkSlateGreyColor": DarkSlateGreyColor,
        "DarkTurquoiseColor": DarkTurquoiseColor,
        "DarkVioletColor": DarkVioletColor,
        "DeepPinkColor": DeepPinkColor,
        "DeepSkyBlueColor": DeepSkyBlueColor,
        "DimGrayColor": DimGrayColor,
        "DimGreyColor": DimGreyColor,
        "DodgerBlueColor": DodgerBlueColor,
        "FireBrickColor": FireBrickColor,
        "FloralWhiteColor": FloralWhiteColor,
        "ForestGreenColor": ForestGreenColor,
        "FuchsiaColor": FuchsiaColor,
        "GainsboroColor": GainsboroColor,
        "GhostWhiteColor": GhostWhiteColor,
        "GoldColor": GoldColor,
        "GoldenrodColor": GoldenrodColor,
        "GrayColor": GrayColor,
        "GreenColor": GreenColor,
        "GreenYellowColor": GreenYellowColor,
        "GreyColor": GreyColor,
        "HoneydewColor": HoneydewColor,
        "HotPinkColor": HotPinkColor,
        "IndianRedColor": IndianRedColor,
        "IndigoColor": IndigoColor,
        "IvoryColor": IvoryColor,
        "KhakiColor": KhakiColor,
        "LavenderBlushColor": LavenderBlushColor,
        "LavenderColor": LavenderColor,
        "LawnGreenColor": LawnGreenColor,
        "LemonChiffonColor": LemonChiffonColor,
        "LightBlueColor": LightBlueColor,
        "LightCoralColor": LightCoralColor,
        "LightCyanColor": LightCyanColor,
        "LightGoldenrodYellowColor": LightGoldenrodYellowColor,
        "LightGrayColor": LightGrayColor,
        "LightGreenColor": LightGreenColor,
        "LightGreyColor": LightGreyColor,
        "LightPinkColor": LightPinkColor,
        "LightSalmonColor": LightSalmonColor,
        "LightSeaGreenColor": LightSeaGreenColor,
        "LightSkyBlueColor": LightSkyBlueColor,
        "LightSlateGrayColor": LightSlateGrayColor,
        "LightSlateGreyColor": LightSlateGreyColor,
        "LightSteelBlueColor": LightSteelBlueColor,
        "LightYellowColor": LightYellowColor,
        "LimeColor": LimeColor,
        "LimeGreenColor": LimeGreenColor,
        "LinenColor": LinenColor,
        "MagentaColor": MagentaColor,
        "MaroonColor": MaroonColor,
        "MediumAquamarineColor": MediumAquamarineColor,
        "MediumBlueColor": MediumBlueColor,
        "MediumOrchidColor": MediumOrchidColor,
        "MediumPurpleColor": MediumPurpleColor,
        "MediumSeaGreenColor": MediumSeaGreenColor,
        "MediumSlateBlueColor": MediumSlateBlueColor,
        "MediumSpringGreenColor": MediumSpringGreenColor,
        "MediumTurquoiseColor": MediumTurquoiseColor,
        "MediumVioletRedColor": MediumVioletRedColor,
        "MidnightBlueColor": MidnightBlueColor,
        "MintCreamColor": MintCreamColor,
        "MistyRoseColor": MistyRoseColor,
        "MoccasinColor": MoccasinColor,
        "NavajoWhiteColor": NavajoWhiteColor,
        "NavyColor": NavyColor,
        "OldLaceColor": OldLaceColor,
        "OliveColor": OliveColor,
        "OliveDrabColor": OliveDrabColor,
        "OrangeColor": OrangeColor,
        "OrangeRedColor": OrangeRedColor,
        "OrchidColor": OrchidColor,
        "PaleGoldenrodColor": PaleGoldenrodColor,
        "PaleGreenColor": PaleGreenColor,
        "PaleTurquoiseColor": PaleTurquoiseColor,
        "PaleVioletRedColor": PaleVioletRedColor,
        "PapayaWhipColor": PapayaWhipColor,
        "PeachPuffColor": PeachPuffColor,
        "PeruColor": PeruColor,
        "PinkColor": PinkColor,
        "PlumColor": PlumColor,
        "PowderBlueColor": PowderBlueColor,
        "PurpleColor": PurpleColor,
        "RedColor": RedColor,
        "RosyBrownColor": RosyBrownColor,
        "RoyalBlueColor": RoyalBlueColor,
        "SaddleBrownColor": SaddleBrownColor,
        "SalmonColor": SalmonColor,
        "SandyBrownColor": SandyBrownColor,
        "SeaGreenColor": SeaGreenColor,
        "SeashellColor": SeashellColor,
        "SiennaColor": SiennaColor,
        "SilverColor": SilverColor,
        "SkyBlueColor": SkyBlueColor,
        "SlateBlueColor": SlateBlueColor,
        "SlateGrayColor": SlateGrayColor,
        "SlateGreyColor": SlateGreyColor,
        "SnowColor": SnowColor,
        "SpringGreenColor": SpringGreenColor,
        "SteelBlueColor": SteelBlueColor,
        "TanColor": TanColor,
        "TealColor": TealColor,
        "ThistleColor": ThistleColor,
        "TomatoColor": TomatoColor,
        "TurquoiseColor": TurquoiseColor,
        "VioletColor": VioletColor,
        "WheatColor": WheatColor,
        "WhiteColor": WhiteColor,
        "WhiteSmokeColor": WhiteSmokeColor,
        "YellowColor": YellowColor,
        "YellowGreenColor": YellowGreenColor,
    }

)

var (
    // In diesem Block werden die Paletten konkret erstellt.
    BlackAndWhitePalette = NewGradientPalette("BlackAndWhitePalette", blackAndWhiteGradient...)
    DarkJunglePalette = NewGradientPaletteByList("DarkJunglePalette", true, darkJungleColorList...)
    DarkPalette = NewGradientPalette("DarkPalette", darkGradient...)
    DarkerPalette = NewGradientPalette("DarkerPalette", darkerGradient...)
    DefaultPalette = NewGradientPalette("DefaultPalette", defaultGradient...)
    EarthAndSkyPalette = NewGradientPalette("EarthAndSkyPalette", earthAndSkyGradient...)
    FadeBluePalette = NewGradientPaletteByList("FadeBluePalette", false, fadeBlueColorListNonCyc...)
    FadeCyanPalette = NewGradientPaletteByList("FadeCyanPalette", false, fadeCyanColorListNonCyc...)
    FadeGreenPalette = NewGradientPaletteByList("FadeGreenPalette", false, fadeGreenColorListNonCyc...)
    FadeMagentaPalette = NewGradientPaletteByList("FadeMagentaPalette", false, fadeMagentaColorListNonCyc...)
    FadeRedPalette = NewGradientPaletteByList("FadeRedPalette", false, fadeRedColorListNonCyc...)
    FadeYellowPalette = NewGradientPaletteByList("FadeYellowPalette", false, fadeYellowColorListNonCyc...)
    FirePalette = NewGradientPalette("FirePalette", fireGradient...)
    HipsterPalette = NewGradientPaletteByList("HipsterPalette", true, hipsterColorList...)
    HotAndColdPalette = NewGradientPalette("HotAndColdPalette", hotAndColdGradient...)
    HotFirePalette = NewGradientPaletteByList("HotFirePalette", true, hotFireColorList...)
    KittyPalette = NewGradientPaletteByList("KittyPalette", true, kittyColorList...)
    LanternPalette = NewGradientPaletteByList("LanternPalette", true, lanternColorList...)
    LemmingPalette = NewGradientPaletteByList("LemmingPalette", true, lemmingColorList...)
    MiamiVicePalette = NewGradientPaletteByList("MiamiVicePalette", true, miamiViceColorList...)
    NeonPalette = NewGradientPaletteByList("NeonPalette", true, neonColorList...)
    NightspellPalette = NewGradientPaletteByList("NightspellPalette", true, nightspellColorList...)
    PastellPalette = NewGradientPalette("PastellPalette", pastellGradient...)
    SeashorePalette = NewGradientPalette("SeashorePalette", seashoreGradient...)
    WayyouPalette = NewGradientPaletteByList("WayyouPalette", true, wayyouColorList...)
    WepartedPalette = NewGradientPaletteByList("WepartedPalette", true, wepartedColorList...)
    AliceBlueColor = NewUniformPalette("AliceBlueColor", colornames.AliceBlue)
    AntiqueWhiteColor = NewUniformPalette("AntiqueWhiteColor", colornames.AntiqueWhite)
    AquaColor = NewUniformPalette("AquaColor", colornames.Aqua)
    AquamarineColor = NewUniformPalette("AquamarineColor", colornames.Aquamarine)
    AzureColor = NewUniformPalette("AzureColor", colornames.Azure)
    BeigeColor = NewUniformPalette("BeigeColor", colornames.Beige)
    BisqueColor = NewUniformPalette("BisqueColor", colornames.Bisque)
    BlackColor = NewUniformPalette("BlackColor", colornames.Black)
    BlanchedAlmondColor = NewUniformPalette("BlanchedAlmondColor", colornames.BlanchedAlmond)
    BlueColor = NewUniformPalette("BlueColor", colornames.Blue)
    BlueVioletColor = NewUniformPalette("BlueVioletColor", colornames.BlueViolet)
    BrownColor = NewUniformPalette("BrownColor", colornames.Brown)
    BurlyWoodColor = NewUniformPalette("BurlyWoodColor", colornames.BurlyWood)
    CadetBlueColor = NewUniformPalette("CadetBlueColor", colornames.CadetBlue)
    ChartreuseColor = NewUniformPalette("ChartreuseColor", colornames.Chartreuse)
    ChocolateColor = NewUniformPalette("ChocolateColor", colornames.Chocolate)
    CoralColor = NewUniformPalette("CoralColor", colornames.Coral)
    CornflowerBlueColor = NewUniformPalette("CornflowerBlueColor", colornames.CornflowerBlue)
    CornsilkColor = NewUniformPalette("CornsilkColor", colornames.Cornsilk)
    CrimsonColor = NewUniformPalette("CrimsonColor", colornames.Crimson)
    CyanColor = NewUniformPalette("CyanColor", colornames.Cyan)
    DarkBlueColor = NewUniformPalette("DarkBlueColor", colornames.DarkBlue)
    DarkCyanColor = NewUniformPalette("DarkCyanColor", colornames.DarkCyan)
    DarkGoldenrodColor = NewUniformPalette("DarkGoldenrodColor", colornames.DarkGoldenrod)
    DarkGrayColor = NewUniformPalette("DarkGrayColor", colornames.DarkGray)
    DarkGreenColor = NewUniformPalette("DarkGreenColor", colornames.DarkGreen)
    DarkGreyColor = NewUniformPalette("DarkGreyColor", colornames.DarkGrey)
    DarkKhakiColor = NewUniformPalette("DarkKhakiColor", colornames.DarkKhaki)
    DarkMagentaColor = NewUniformPalette("DarkMagentaColor", colornames.DarkMagenta)
    DarkOliveGreenColor = NewUniformPalette("DarkOliveGreenColor", colornames.DarkOliveGreen)
    DarkOrangeColor = NewUniformPalette("DarkOrangeColor", colornames.DarkOrange)
    DarkOrchidColor = NewUniformPalette("DarkOrchidColor", colornames.DarkOrchid)
    DarkRedColor = NewUniformPalette("DarkRedColor", colornames.DarkRed)
    DarkSalmonColor = NewUniformPalette("DarkSalmonColor", colornames.DarkSalmon)
    DarkSeaGreenColor = NewUniformPalette("DarkSeaGreenColor", colornames.DarkSeaGreen)
    DarkSlateBlueColor = NewUniformPalette("DarkSlateBlueColor", colornames.DarkSlateBlue)
    DarkSlateGrayColor = NewUniformPalette("DarkSlateGrayColor", colornames.DarkSlateGray)
    DarkSlateGreyColor = NewUniformPalette("DarkSlateGreyColor", colornames.DarkSlateGrey)
    DarkTurquoiseColor = NewUniformPalette("DarkTurquoiseColor", colornames.DarkTurquoise)
    DarkVioletColor = NewUniformPalette("DarkVioletColor", colornames.DarkViolet)
    DeepPinkColor = NewUniformPalette("DeepPinkColor", colornames.DeepPink)
    DeepSkyBlueColor = NewUniformPalette("DeepSkyBlueColor", colornames.DeepSkyBlue)
    DimGrayColor = NewUniformPalette("DimGrayColor", colornames.DimGray)
    DimGreyColor = NewUniformPalette("DimGreyColor", colornames.DimGrey)
    DodgerBlueColor = NewUniformPalette("DodgerBlueColor", colornames.DodgerBlue)
    FireBrickColor = NewUniformPalette("FireBrickColor", colornames.FireBrick)
    FloralWhiteColor = NewUniformPalette("FloralWhiteColor", colornames.FloralWhite)
    ForestGreenColor = NewUniformPalette("ForestGreenColor", colornames.ForestGreen)
    FuchsiaColor = NewUniformPalette("FuchsiaColor", colornames.Fuchsia)
    GainsboroColor = NewUniformPalette("GainsboroColor", colornames.Gainsboro)
    GhostWhiteColor = NewUniformPalette("GhostWhiteColor", colornames.GhostWhite)
    GoldColor = NewUniformPalette("GoldColor", colornames.Gold)
    GoldenrodColor = NewUniformPalette("GoldenrodColor", colornames.Goldenrod)
    GrayColor = NewUniformPalette("GrayColor", colornames.Gray)
    GreenColor = NewUniformPalette("GreenColor", colornames.Green)
    GreenYellowColor = NewUniformPalette("GreenYellowColor", colornames.GreenYellow)
    GreyColor = NewUniformPalette("GreyColor", colornames.Grey)
    HoneydewColor = NewUniformPalette("HoneydewColor", colornames.Honeydew)
    HotPinkColor = NewUniformPalette("HotPinkColor", colornames.HotPink)
    IndianRedColor = NewUniformPalette("IndianRedColor", colornames.IndianRed)
    IndigoColor = NewUniformPalette("IndigoColor", colornames.Indigo)
    IvoryColor = NewUniformPalette("IvoryColor", colornames.Ivory)
    KhakiColor = NewUniformPalette("KhakiColor", colornames.Khaki)
    LavenderBlushColor = NewUniformPalette("LavenderBlushColor", colornames.LavenderBlush)
    LavenderColor = NewUniformPalette("LavenderColor", colornames.Lavender)
    LawnGreenColor = NewUniformPalette("LawnGreenColor", colornames.LawnGreen)
    LemonChiffonColor = NewUniformPalette("LemonChiffonColor", colornames.LemonChiffon)
    LightBlueColor = NewUniformPalette("LightBlueColor", colornames.LightBlue)
    LightCoralColor = NewUniformPalette("LightCoralColor", colornames.LightCoral)
    LightCyanColor = NewUniformPalette("LightCyanColor", colornames.LightCyan)
    LightGoldenrodYellowColor = NewUniformPalette("LightGoldenrodYellowColor", colornames.LightGoldenrodYellow)
    LightGrayColor = NewUniformPalette("LightGrayColor", colornames.LightGray)
    LightGreenColor = NewUniformPalette("LightGreenColor", colornames.LightGreen)
    LightGreyColor = NewUniformPalette("LightGreyColor", colornames.LightGrey)
    LightPinkColor = NewUniformPalette("LightPinkColor", colornames.LightPink)
    LightSalmonColor = NewUniformPalette("LightSalmonColor", colornames.LightSalmon)
    LightSeaGreenColor = NewUniformPalette("LightSeaGreenColor", colornames.LightSeaGreen)
    LightSkyBlueColor = NewUniformPalette("LightSkyBlueColor", colornames.LightSkyBlue)
    LightSlateGrayColor = NewUniformPalette("LightSlateGrayColor", colornames.LightSlateGray)
    LightSlateGreyColor = NewUniformPalette("LightSlateGreyColor", colornames.LightSlateGrey)
    LightSteelBlueColor = NewUniformPalette("LightSteelBlueColor", colornames.LightSteelBlue)
    LightYellowColor = NewUniformPalette("LightYellowColor", colornames.LightYellow)
    LimeColor = NewUniformPalette("LimeColor", colornames.Lime)
    LimeGreenColor = NewUniformPalette("LimeGreenColor", colornames.LimeGreen)
    LinenColor = NewUniformPalette("LinenColor", colornames.Linen)
    MagentaColor = NewUniformPalette("MagentaColor", colornames.Magenta)
    MaroonColor = NewUniformPalette("MaroonColor", colornames.Maroon)
    MediumAquamarineColor = NewUniformPalette("MediumAquamarineColor", colornames.MediumAquamarine)
    MediumBlueColor = NewUniformPalette("MediumBlueColor", colornames.MediumBlue)
    MediumOrchidColor = NewUniformPalette("MediumOrchidColor", colornames.MediumOrchid)
    MediumPurpleColor = NewUniformPalette("MediumPurpleColor", colornames.MediumPurple)
    MediumSeaGreenColor = NewUniformPalette("MediumSeaGreenColor", colornames.MediumSeaGreen)
    MediumSlateBlueColor = NewUniformPalette("MediumSlateBlueColor", colornames.MediumSlateBlue)
    MediumSpringGreenColor = NewUniformPalette("MediumSpringGreenColor", colornames.MediumSpringGreen)
    MediumTurquoiseColor = NewUniformPalette("MediumTurquoiseColor", colornames.MediumTurquoise)
    MediumVioletRedColor = NewUniformPalette("MediumVioletRedColor", colornames.MediumVioletRed)
    MidnightBlueColor = NewUniformPalette("MidnightBlueColor", colornames.MidnightBlue)
    MintCreamColor = NewUniformPalette("MintCreamColor", colornames.MintCream)
    MistyRoseColor = NewUniformPalette("MistyRoseColor", colornames.MistyRose)
    MoccasinColor = NewUniformPalette("MoccasinColor", colornames.Moccasin)
    NavajoWhiteColor = NewUniformPalette("NavajoWhiteColor", colornames.NavajoWhite)
    NavyColor = NewUniformPalette("NavyColor", colornames.Navy)
    OldLaceColor = NewUniformPalette("OldLaceColor", colornames.OldLace)
    OliveColor = NewUniformPalette("OliveColor", colornames.Olive)
    OliveDrabColor = NewUniformPalette("OliveDrabColor", colornames.OliveDrab)
    OrangeColor = NewUniformPalette("OrangeColor", colornames.Orange)
    OrangeRedColor = NewUniformPalette("OrangeRedColor", colornames.OrangeRed)
    OrchidColor = NewUniformPalette("OrchidColor", colornames.Orchid)
    PaleGoldenrodColor = NewUniformPalette("PaleGoldenrodColor", colornames.PaleGoldenrod)
    PaleGreenColor = NewUniformPalette("PaleGreenColor", colornames.PaleGreen)
    PaleTurquoiseColor = NewUniformPalette("PaleTurquoiseColor", colornames.PaleTurquoise)
    PaleVioletRedColor = NewUniformPalette("PaleVioletRedColor", colornames.PaleVioletRed)
    PapayaWhipColor = NewUniformPalette("PapayaWhipColor", colornames.PapayaWhip)
    PeachPuffColor = NewUniformPalette("PeachPuffColor", colornames.PeachPuff)
    PeruColor = NewUniformPalette("PeruColor", colornames.Peru)
    PinkColor = NewUniformPalette("PinkColor", colornames.Pink)
    PlumColor = NewUniformPalette("PlumColor", colornames.Plum)
    PowderBlueColor = NewUniformPalette("PowderBlueColor", colornames.PowderBlue)
    PurpleColor = NewUniformPalette("PurpleColor", colornames.Purple)
    RedColor = NewUniformPalette("RedColor", colornames.Red)
    RosyBrownColor = NewUniformPalette("RosyBrownColor", colornames.RosyBrown)
    RoyalBlueColor = NewUniformPalette("RoyalBlueColor", colornames.RoyalBlue)
    SaddleBrownColor = NewUniformPalette("SaddleBrownColor", colornames.SaddleBrown)
    SalmonColor = NewUniformPalette("SalmonColor", colornames.Salmon)
    SandyBrownColor = NewUniformPalette("SandyBrownColor", colornames.SandyBrown)
    SeaGreenColor = NewUniformPalette("SeaGreenColor", colornames.SeaGreen)
    SeashellColor = NewUniformPalette("SeashellColor", colornames.Seashell)
    SiennaColor = NewUniformPalette("SiennaColor", colornames.Sienna)
    SilverColor = NewUniformPalette("SilverColor", colornames.Silver)
    SkyBlueColor = NewUniformPalette("SkyBlueColor", colornames.SkyBlue)
    SlateBlueColor = NewUniformPalette("SlateBlueColor", colornames.SlateBlue)
    SlateGrayColor = NewUniformPalette("SlateGrayColor", colornames.SlateGray)
    SlateGreyColor = NewUniformPalette("SlateGreyColor", colornames.SlateGrey)
    SnowColor = NewUniformPalette("SnowColor", colornames.Snow)
    SpringGreenColor = NewUniformPalette("SpringGreenColor", colornames.SpringGreen)
    SteelBlueColor = NewUniformPalette("SteelBlueColor", colornames.SteelBlue)
    TanColor = NewUniformPalette("TanColor", colornames.Tan)
    TealColor = NewUniformPalette("TealColor", colornames.Teal)
    ThistleColor = NewUniformPalette("ThistleColor", colornames.Thistle)
    TomatoColor = NewUniformPalette("TomatoColor", colornames.Tomato)
    TurquoiseColor = NewUniformPalette("TurquoiseColor", colornames.Turquoise)
    VioletColor = NewUniformPalette("VioletColor", colornames.Violet)
    WheatColor = NewUniformPalette("WheatColor", colornames.Wheat)
    WhiteColor = NewUniformPalette("WhiteColor", colornames.White)
    WhiteSmokeColor = NewUniformPalette("WhiteSmokeColor", colornames.WhiteSmoke)
    YellowColor = NewUniformPalette("YellowColor", colornames.Yellow)
    YellowGreenColor = NewUniformPalette("YellowGreenColor", colornames.YellowGreen)
)
