package conf

// ---------------------------------------------------------------------------
// In order to have some non-default module configurations to test the conf
// package on one hand or the grid emulator on the other, this section
// contains some special, yet weird but working configurations. Some of the
// panels even have holes in them - but they work!
var (
    CustomConf = ModuleConfig{
        {Col: 0, Row: 0, Mod: ModLR000},
        {Col: 1, Row: 0, Mod: ModRL090},
        {Col: 1, Row: 1, Mod: ModRL090},
        {Col: 0, Row: 1, Mod: ModLR180},
    }

	TetrisTile = ModuleConfig{
		{Col: 0, Row: 0, Mod: ModLR000},
		{Col: 1, Row: 0, Mod: ModRL090},
		{Col: 1, Row: 1, Mod: ModLR000},
		{Col: 2, Row: 1, Mod: ModLR000},
	}

	LowerCurve = ModuleConfig{
		{Col: 0, Row: 0, Mod: ModLR000},
		{Col: 1, Row: 0, Mod: ModRL090},
		{Col: 1, Row: 1, Mod: ModLR000},
		{Col: 2, Row: 1, Mod: ModLR000},
		{Col: 3, Row: 1, Mod: ModLR000},
		{Col: 3, Row: 0, Mod: ModRL270},
		{Col: 4, Row: 0, Mod: ModLR000},
	}

	SquareWithHole = ModuleConfig{
		{Col: 0, Row: 0, Mod: ModLR000},
		{Col: 1, Row: 0, Mod: ModLR000},
		{Col: 2, Row: 0, Mod: ModRL090},
		{Col: 2, Row: 1, Mod: ModRL090},
		{Col: 2, Row: 2, Mod: ModRL090},
		{Col: 1, Row: 2, Mod: ModLR180},
		{Col: 0, Row: 2, Mod: ModRL270},
		{Col: 0, Row: 1, Mod: ModRL270},
	}

	SmallChessBoard = ModuleConfig{
		{Col: 0, Row: 0, Mod: ModRL180},
		{Col: 1, Row: 1, Mod: ModLR000},
	}

	ChessBoard = ModuleConfig{
		{Col: 1, Row: 0, Mod: ModRL180},
		{Col: 2, Row: 1, Mod: ModLR000},
		{Col: 3, Row: 0, Mod: ModRL180},
		{Col: 4, Row: 1, Mod: ModRL090},
		{Col: 3, Row: 2, Mod: ModLR270},
		{Col: 4, Row: 3, Mod: ModRL090},
		{Col: 3, Row: 4, Mod: ModRL000},
		{Col: 2, Row: 3, Mod: ModLR180},
		{Col: 1, Row: 4, Mod: ModRL000},
		{Col: 0, Row: 3, Mod: ModRL270},
		{Col: 1, Row: 2, Mod: ModLR090},
		{Col: 0, Row: 1, Mod: ModRL270},
	}
)
