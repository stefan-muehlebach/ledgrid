package ledgrid

import (
	"github.com/stefan-muehlebach/gg/color"
)

var (
    // ColorList ist ein Slice mit allen vorhandenen Paletten.
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

    // ColorMap ist ein Map um Paletten mit ihrem Namen anzusprechen.
    ColorMap = map[string]ColorSource{
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
    // In diesem Block werden die uniformen Paletten erstellt.
    AliceBlueColor = NewUniformPalette("AliceBlueColor", color.AliceBlue)
    AntiqueWhiteColor = NewUniformPalette("AntiqueWhiteColor", color.AntiqueWhite)
    AquaColor = NewUniformPalette("AquaColor", color.Aqua)
    AquamarineColor = NewUniformPalette("AquamarineColor", color.Aquamarine)
    AzureColor = NewUniformPalette("AzureColor", color.Azure)
    BeigeColor = NewUniformPalette("BeigeColor", color.Beige)
    BisqueColor = NewUniformPalette("BisqueColor", color.Bisque)
    BlackColor = NewUniformPalette("BlackColor", color.Black)
    BlanchedAlmondColor = NewUniformPalette("BlanchedAlmondColor", color.BlanchedAlmond)
    BlueColor = NewUniformPalette("BlueColor", color.Blue)
    BlueVioletColor = NewUniformPalette("BlueVioletColor", color.BlueViolet)
    BrownColor = NewUniformPalette("BrownColor", color.Brown)
    BurlyWoodColor = NewUniformPalette("BurlyWoodColor", color.BurlyWood)
    CadetBlueColor = NewUniformPalette("CadetBlueColor", color.CadetBlue)
    ChartreuseColor = NewUniformPalette("ChartreuseColor", color.Chartreuse)
    ChocolateColor = NewUniformPalette("ChocolateColor", color.Chocolate)
    CoralColor = NewUniformPalette("CoralColor", color.Coral)
    CornflowerBlueColor = NewUniformPalette("CornflowerBlueColor", color.CornflowerBlue)
    CornsilkColor = NewUniformPalette("CornsilkColor", color.Cornsilk)
    CrimsonColor = NewUniformPalette("CrimsonColor", color.Crimson)
    CyanColor = NewUniformPalette("CyanColor", color.Cyan)
    DarkBlueColor = NewUniformPalette("DarkBlueColor", color.DarkBlue)
    DarkCyanColor = NewUniformPalette("DarkCyanColor", color.DarkCyan)
    DarkGoldenrodColor = NewUniformPalette("DarkGoldenrodColor", color.DarkGoldenrod)
    DarkGrayColor = NewUniformPalette("DarkGrayColor", color.DarkGray)
    DarkGreenColor = NewUniformPalette("DarkGreenColor", color.DarkGreen)
    DarkGreyColor = NewUniformPalette("DarkGreyColor", color.DarkGrey)
    DarkKhakiColor = NewUniformPalette("DarkKhakiColor", color.DarkKhaki)
    DarkMagentaColor = NewUniformPalette("DarkMagentaColor", color.DarkMagenta)
    DarkOliveGreenColor = NewUniformPalette("DarkOliveGreenColor", color.DarkOliveGreen)
    DarkOrangeColor = NewUniformPalette("DarkOrangeColor", color.DarkOrange)
    DarkOrchidColor = NewUniformPalette("DarkOrchidColor", color.DarkOrchid)
    DarkRedColor = NewUniformPalette("DarkRedColor", color.DarkRed)
    DarkSalmonColor = NewUniformPalette("DarkSalmonColor", color.DarkSalmon)
    DarkSeaGreenColor = NewUniformPalette("DarkSeaGreenColor", color.DarkSeaGreen)
    DarkSlateBlueColor = NewUniformPalette("DarkSlateBlueColor", color.DarkSlateBlue)
    DarkSlateGrayColor = NewUniformPalette("DarkSlateGrayColor", color.DarkSlateGray)
    DarkSlateGreyColor = NewUniformPalette("DarkSlateGreyColor", color.DarkSlateGrey)
    DarkTurquoiseColor = NewUniformPalette("DarkTurquoiseColor", color.DarkTurquoise)
    DarkVioletColor = NewUniformPalette("DarkVioletColor", color.DarkViolet)
    DeepPinkColor = NewUniformPalette("DeepPinkColor", color.DeepPink)
    DeepSkyBlueColor = NewUniformPalette("DeepSkyBlueColor", color.DeepSkyBlue)
    DimGrayColor = NewUniformPalette("DimGrayColor", color.DimGray)
    DimGreyColor = NewUniformPalette("DimGreyColor", color.DimGrey)
    DodgerBlueColor = NewUniformPalette("DodgerBlueColor", color.DodgerBlue)
    FireBrickColor = NewUniformPalette("FireBrickColor", color.FireBrick)
    FloralWhiteColor = NewUniformPalette("FloralWhiteColor", color.FloralWhite)
    ForestGreenColor = NewUniformPalette("ForestGreenColor", color.ForestGreen)
    FuchsiaColor = NewUniformPalette("FuchsiaColor", color.Fuchsia)
    GainsboroColor = NewUniformPalette("GainsboroColor", color.Gainsboro)
    GhostWhiteColor = NewUniformPalette("GhostWhiteColor", color.GhostWhite)
    GoldColor = NewUniformPalette("GoldColor", color.Gold)
    GoldenrodColor = NewUniformPalette("GoldenrodColor", color.Goldenrod)
    GrayColor = NewUniformPalette("GrayColor", color.Gray)
    GreenColor = NewUniformPalette("GreenColor", color.Green)
    GreenYellowColor = NewUniformPalette("GreenYellowColor", color.GreenYellow)
    GreyColor = NewUniformPalette("GreyColor", color.Grey)
    HoneydewColor = NewUniformPalette("HoneydewColor", color.Honeydew)
    HotPinkColor = NewUniformPalette("HotPinkColor", color.HotPink)
    IndianRedColor = NewUniformPalette("IndianRedColor", color.IndianRed)
    IndigoColor = NewUniformPalette("IndigoColor", color.Indigo)
    IvoryColor = NewUniformPalette("IvoryColor", color.Ivory)
    KhakiColor = NewUniformPalette("KhakiColor", color.Khaki)
    LavenderBlushColor = NewUniformPalette("LavenderBlushColor", color.LavenderBlush)
    LavenderColor = NewUniformPalette("LavenderColor", color.Lavender)
    LawnGreenColor = NewUniformPalette("LawnGreenColor", color.LawnGreen)
    LemonChiffonColor = NewUniformPalette("LemonChiffonColor", color.LemonChiffon)
    LightBlueColor = NewUniformPalette("LightBlueColor", color.LightBlue)
    LightCoralColor = NewUniformPalette("LightCoralColor", color.LightCoral)
    LightCyanColor = NewUniformPalette("LightCyanColor", color.LightCyan)
    LightGoldenrodYellowColor = NewUniformPalette("LightGoldenrodYellowColor", color.LightGoldenrodYellow)
    LightGrayColor = NewUniformPalette("LightGrayColor", color.LightGray)
    LightGreenColor = NewUniformPalette("LightGreenColor", color.LightGreen)
    LightGreyColor = NewUniformPalette("LightGreyColor", color.LightGrey)
    LightPinkColor = NewUniformPalette("LightPinkColor", color.LightPink)
    LightSalmonColor = NewUniformPalette("LightSalmonColor", color.LightSalmon)
    LightSeaGreenColor = NewUniformPalette("LightSeaGreenColor", color.LightSeaGreen)
    LightSkyBlueColor = NewUniformPalette("LightSkyBlueColor", color.LightSkyBlue)
    LightSlateGrayColor = NewUniformPalette("LightSlateGrayColor", color.LightSlateGray)
    LightSlateGreyColor = NewUniformPalette("LightSlateGreyColor", color.LightSlateGrey)
    LightSteelBlueColor = NewUniformPalette("LightSteelBlueColor", color.LightSteelBlue)
    LightYellowColor = NewUniformPalette("LightYellowColor", color.LightYellow)
    LimeColor = NewUniformPalette("LimeColor", color.Lime)
    LimeGreenColor = NewUniformPalette("LimeGreenColor", color.LimeGreen)
    LinenColor = NewUniformPalette("LinenColor", color.Linen)
    MagentaColor = NewUniformPalette("MagentaColor", color.Magenta)
    MaroonColor = NewUniformPalette("MaroonColor", color.Maroon)
    MediumAquamarineColor = NewUniformPalette("MediumAquamarineColor", color.MediumAquamarine)
    MediumBlueColor = NewUniformPalette("MediumBlueColor", color.MediumBlue)
    MediumOrchidColor = NewUniformPalette("MediumOrchidColor", color.MediumOrchid)
    MediumPurpleColor = NewUniformPalette("MediumPurpleColor", color.MediumPurple)
    MediumSeaGreenColor = NewUniformPalette("MediumSeaGreenColor", color.MediumSeaGreen)
    MediumSlateBlueColor = NewUniformPalette("MediumSlateBlueColor", color.MediumSlateBlue)
    MediumSpringGreenColor = NewUniformPalette("MediumSpringGreenColor", color.MediumSpringGreen)
    MediumTurquoiseColor = NewUniformPalette("MediumTurquoiseColor", color.MediumTurquoise)
    MediumVioletRedColor = NewUniformPalette("MediumVioletRedColor", color.MediumVioletRed)
    MidnightBlueColor = NewUniformPalette("MidnightBlueColor", color.MidnightBlue)
    MintCreamColor = NewUniformPalette("MintCreamColor", color.MintCream)
    MistyRoseColor = NewUniformPalette("MistyRoseColor", color.MistyRose)
    MoccasinColor = NewUniformPalette("MoccasinColor", color.Moccasin)
    NavajoWhiteColor = NewUniformPalette("NavajoWhiteColor", color.NavajoWhite)
    NavyColor = NewUniformPalette("NavyColor", color.Navy)
    OldLaceColor = NewUniformPalette("OldLaceColor", color.OldLace)
    OliveColor = NewUniformPalette("OliveColor", color.Olive)
    OliveDrabColor = NewUniformPalette("OliveDrabColor", color.OliveDrab)
    OrangeColor = NewUniformPalette("OrangeColor", color.Orange)
    OrangeRedColor = NewUniformPalette("OrangeRedColor", color.OrangeRed)
    OrchidColor = NewUniformPalette("OrchidColor", color.Orchid)
    PaleGoldenrodColor = NewUniformPalette("PaleGoldenrodColor", color.PaleGoldenrod)
    PaleGreenColor = NewUniformPalette("PaleGreenColor", color.PaleGreen)
    PaleTurquoiseColor = NewUniformPalette("PaleTurquoiseColor", color.PaleTurquoise)
    PaleVioletRedColor = NewUniformPalette("PaleVioletRedColor", color.PaleVioletRed)
    PapayaWhipColor = NewUniformPalette("PapayaWhipColor", color.PapayaWhip)
    PeachPuffColor = NewUniformPalette("PeachPuffColor", color.PeachPuff)
    PeruColor = NewUniformPalette("PeruColor", color.Peru)
    PinkColor = NewUniformPalette("PinkColor", color.Pink)
    PlumColor = NewUniformPalette("PlumColor", color.Plum)
    PowderBlueColor = NewUniformPalette("PowderBlueColor", color.PowderBlue)
    PurpleColor = NewUniformPalette("PurpleColor", color.Purple)
    RedColor = NewUniformPalette("RedColor", color.Red)
    RosyBrownColor = NewUniformPalette("RosyBrownColor", color.RosyBrown)
    RoyalBlueColor = NewUniformPalette("RoyalBlueColor", color.RoyalBlue)
    SaddleBrownColor = NewUniformPalette("SaddleBrownColor", color.SaddleBrown)
    SalmonColor = NewUniformPalette("SalmonColor", color.Salmon)
    SandyBrownColor = NewUniformPalette("SandyBrownColor", color.SandyBrown)
    SeaGreenColor = NewUniformPalette("SeaGreenColor", color.SeaGreen)
    SeashellColor = NewUniformPalette("SeashellColor", color.Seashell)
    SiennaColor = NewUniformPalette("SiennaColor", color.Sienna)
    SilverColor = NewUniformPalette("SilverColor", color.Silver)
    SkyBlueColor = NewUniformPalette("SkyBlueColor", color.SkyBlue)
    SlateBlueColor = NewUniformPalette("SlateBlueColor", color.SlateBlue)
    SlateGrayColor = NewUniformPalette("SlateGrayColor", color.SlateGray)
    SlateGreyColor = NewUniformPalette("SlateGreyColor", color.SlateGrey)
    SnowColor = NewUniformPalette("SnowColor", color.Snow)
    SpringGreenColor = NewUniformPalette("SpringGreenColor", color.SpringGreen)
    SteelBlueColor = NewUniformPalette("SteelBlueColor", color.SteelBlue)
    TanColor = NewUniformPalette("TanColor", color.Tan)
    TealColor = NewUniformPalette("TealColor", color.Teal)
    ThistleColor = NewUniformPalette("ThistleColor", color.Thistle)
    TomatoColor = NewUniformPalette("TomatoColor", color.Tomato)
    TurquoiseColor = NewUniformPalette("TurquoiseColor", color.Turquoise)
    VioletColor = NewUniformPalette("VioletColor", color.Violet)
    WheatColor = NewUniformPalette("WheatColor", color.Wheat)
    WhiteColor = NewUniformPalette("WhiteColor", color.White)
    WhiteSmokeColor = NewUniformPalette("WhiteSmokeColor", color.WhiteSmoke)
    YellowColor = NewUniformPalette("YellowColor", color.Yellow)
    YellowGreenColor = NewUniformPalette("YellowGreenColor", color.YellowGreen)
)
