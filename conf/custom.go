package conf

// ---------------------------------------------------------------------------
// In order to have some non-default module configurations to test the conf
// package on one hand or the grid emulator on the other, this section
// contains some special, yet weird but working configurations. Some of the
// panels even have holes in them - but they work!
var (
	TetrisTile = ModuleConfig{
		ModulePosition{Col: 0, Row: 0, Idx: 0,   Mod: ModLR000},
		ModulePosition{Col: 1, Row: 0, Idx: 100, Mod: ModRL090},
		ModulePosition{Col: 1, Row: 1, Idx: 200, Mod: ModLR000},
		ModulePosition{Col: 2, Row: 1, Idx: 300, Mod: ModLR000},
	}

	LowerCurve = ModuleConfig{
		ModulePosition{Col: 0, Row: 0, Idx: 0,   Mod: ModLR000},
		ModulePosition{Col: 1, Row: 0, Idx: 100, Mod: ModRL090},
		ModulePosition{Col: 1, Row: 1, Idx: 200, Mod: ModLR000},
		ModulePosition{Col: 2, Row: 1, Idx: 300, Mod: ModLR000},
		ModulePosition{Col: 3, Row: 1, Idx: 400, Mod: ModLR000},
		ModulePosition{Col: 3, Row: 0, Idx: 500, Mod: ModRL270},
		ModulePosition{Col: 4, Row: 0, Idx: 600, Mod: ModLR000},
	}

	SquareWithHole = ModuleConfig{
		ModulePosition{Col: 0, Row: 0, Idx: 0,   Mod: ModLR000},
		ModulePosition{Col: 1, Row: 0, Idx: 100, Mod: ModLR000},
		ModulePosition{Col: 2, Row: 0, Idx: 200, Mod: ModRL090},
		ModulePosition{Col: 2, Row: 1, Idx: 300, Mod: ModRL090},
		ModulePosition{Col: 2, Row: 2, Idx: 400, Mod: ModRL090},
		ModulePosition{Col: 1, Row: 2, Idx: 500, Mod: ModLR180},
		ModulePosition{Col: 0, Row: 2, Idx: 600, Mod: ModRL270},
		ModulePosition{Col: 0, Row: 1, Idx: 700, Mod: ModRL270},
	}

	ChessBoard = ModuleConfig{
		ModulePosition{Col: 1, Row: 0, Idx: 0, Mod: ModRL180},
		ModulePosition{Col: 2, Row: 1, Idx: 100, Mod: ModLR000},
		ModulePosition{Col: 3, Row: 0, Idx: 200, Mod: ModRL180},
		ModulePosition{Col: 4, Row: 1, Idx: 300, Mod: ModRL090},
		ModulePosition{Col: 3, Row: 2, Idx: 400, Mod: ModLR270},
		ModulePosition{Col: 4, Row: 3, Idx: 500, Mod: ModRL090},
		ModulePosition{Col: 3, Row: 4, Idx: 600, Mod: ModRL000},
		ModulePosition{Col: 2, Row: 3, Idx: 700, Mod: ModLR180},
		ModulePosition{Col: 1, Row: 4, Idx: 800, Mod: ModRL000},
		ModulePosition{Col: 0, Row: 3, Idx: 900, Mod: ModRL270},
		ModulePosition{Col: 1, Row: 2, Idx: 1000, Mod: ModLR090},
		ModulePosition{Col: 0, Row: 1, Idx: 1100, Mod: ModRL270},
	}
)
