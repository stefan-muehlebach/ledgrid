package ledgrid

// import (
// 	"testing"
// )

// var (
// 	param []Parameter = make([]Parameter, 5)
// )

// func TestParameter(t *testing.T) {
// 	param[0] = NewBoolParam("Checkbox", false)
// 	param[1] = NewFloatParam("Float", 0.0, -10.0, 10.0, 2.0)
// 	param[2] = NewIntParam("Integer", 2, 0, 10, 1)
// 	param[3] = NewStringParam("Textbox", "Something")
// 	param[4] = NewPaletteParam("Color", WhiteSmokeColor)
// 	for i, p := range param {
// 		switch o := p.(type) {
// 		case BoolParameter:
// 			t.Logf("[%d] parameter \"%s\", value: %t", i, p.Name(), o.Get())
// 		case IntParameter:
// 			t.Logf("[%d] parameter \"%s\", value: %d", i, p.Name(), o.Get())
// 		case FloatParameter:
// 			t.Logf("[%d] parameter \"%s\", value: %.3f", i, p.Name(), o.Get())
// 		case StringParameter:
// 			t.Logf("[%d] parameter \"%s\", value: '%s'", i, p.Name(), o.Get())
// 		case PaletteParameter:
// 			t.Logf("[%d] parameter \"%s\", value: %v", i, p.Name(), o.Get())
// 		}
// 	}
// }
