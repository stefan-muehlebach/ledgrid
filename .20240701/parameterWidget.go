package ledgrid

import (
    	"fyne.io/fyne/v2"
    	"fyne.io/fyne/v2/widget"
   	"fyne.io/fyne/v2/data/binding"
)

type SliderParameter struct {
    label *widget.Label
    slider *widget.Slider
}

func NewSliderParameter(param *Bounded[float64]) *SliderParameter {
    p := &SliderParameter{}
    p.label = widget.NewLabelWithData(binding.FloatToStringWithFormat(param, param.Name()+" (%.3f)"))
    p.label.Alignment = fyne.TextAlignTrailing
    p.label.TextStyle.Bold = true
	p.slider = widget.NewSliderWithData(param.Min(), param.Max(), param)
	p.slider.Step = param.Step()
	p.slider.SetValue(param.Val())
    return p
}

func (p *SliderParameter) Label() fyne.Widget {
    return p.label
}

func (p *SliderParameter) Control() fyne.Widget {
    return p.slider
}

