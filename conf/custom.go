package conf

// ---------------------------------------------------------------------------
// In order to have some non-default module configurations to test the conf
// package on one hand or the grid emulator on the other, this section
// contains some special, yet weird but working configurations. Some of the
// panels even have holes in them - but they work!
var (
	TetrisTile = ModuleConfig{}
	tetrSpec   = []ModuleSpec{
		{0, 0, ModLR000},
		{1, 0, ModRL090},
		{1, 1, ModLR000},
		{2, 1, ModLR000},
	}

	LowerCurve = ModuleConfig{}
	lowCurSpec = []ModuleSpec{
		{0, 0, ModLR000},
		{1, 0, ModRL090},
		{1, 1, ModLR000},
		{2, 1, ModLR000},
		{3, 1, ModLR000},
		{3, 0, ModRL270},
		{4, 0, ModLR000},
	}

	SquareWithHole = ModuleConfig{}
	squWiHolSpec   = []ModuleSpec{
		{0, 0, ModLR000},
		{1, 0, ModLR000},
		{2, 0, ModRL090},
		{2, 1, ModRL090},
		{2, 2, ModRL090},
		{1, 2, ModLR180},
		{0, 2, ModRL270},
		{0, 1, ModRL270},
	}

	SmallChessBoard = ModuleConfig{}
	smChBrdSpec     = []ModuleSpec{
		{0, 0, ModRL180},
		{1, 1, ModLR000},
	}

	ChessBoard = ModuleConfig{}
	chBrdSpec  = []ModuleSpec{
		{1, 0, ModRL180},
		{2, 1, ModLR000},
		{3, 0, ModRL180},
		{4, 1, ModRL090},
		{3, 2, ModLR270},
		{4, 3, ModRL090},
		{3, 4, ModRL000},
		{2, 3, ModLR180},
		{1, 4, ModRL000},
		{0, 3, ModRL270},
		{1, 2, ModLR090},
		{0, 1, ModRL270},
	}
)

func init() {
	TetrisTile.AddMods(tetrSpec)
	LowerCurve.AddMods(lowCurSpec)
	SquareWithHole.AddMods(squWiHolSpec)
	SmallChessBoard.AddMods(smChBrdSpec)
	ChessBoard.AddMods(chBrdSpec)
}
