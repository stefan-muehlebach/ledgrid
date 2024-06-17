package ledgrid

import (
	"fyne.io/fyne/v2/data/binding"
)

type NumericParameterType interface {
	~int | ~float64
}
type ParameterType interface {
	~bool | NumericParameterType | ~string | interface{}
}

//----------------------------------------------------------------------------

type DataListener interface {
	DataChanged()
}

type Parameter interface {
    Nameable
    SetCallback(fnc func(p Parameter))
    AddListener(binding.DataListener)
    RemoveListener(binding.DataListener)
}

//----------------------------------------------------------------------------

type baseParam[T ParameterType] struct {
    NameableEmbed
    val, init T
    callback func(p Parameter)
    	listeners []binding.DataListener
}
func newBaseParam[T ParameterType](name string, init T) *baseParam[T] {
    p := &baseParam[T]{}
    p.name = name
    p.val = init
    p.init = init
	p.listeners = make([]binding.DataListener, 0)
    return p
}

func (p *baseParam[T]) SetCallback(fnc func(p Parameter)) {
    p.callback = fnc
    p.callback(p)
}

func (p *baseParam[T]) AddListener(l binding.DataListener) {
	p.listeners = append(p.listeners, l)
}
func (p *baseParam[T]) RemoveListener(l binding.DataListener) {
	for idx, lsnr := range p.listeners {
		if l == lsnr {
            	p.listeners = append(p.listeners[:idx], p.listeners[idx+1:]...)
			break
		}
	}
}
func (p *baseParam[T]) Get() (T, error) {
    return p.val, nil
}
func (p *baseParam[T]) Set(v T) error {
    p.set(v)
    return nil
}

func (p *baseParam[T]) Val() T {
    return p.val
}
func (p *baseParam[T]) SetVal(v T) {
    p.set(v)
}
func (p *baseParam[T]) Reset() {
    p.set(p.init)
}
func (p *baseParam[T]) set(v T) {
	p.val = v
	// if b.valPtr != nil {
	// 	*b.valPtr = b.val
	// }
	if p.callback != nil {
		p.callback(p)
	}
	for _, lsnr := range p.listeners {
		if lsnr != nil {
			lsnr.DataChanged()
		}
	}
}



//----------------------------------------------------------------------------

type numericParam[T NumericParameterType] struct {
    baseParam[T]
    min, max, step T
}
func newNumericParam[T NumericParameterType](name string, init, min, max, step T) *numericParam[T] {
    p := &numericParam[T]{}
    p.name = name
    p.val = init
    p.init = init
    p.min = min
    p.max = max
    p.step = step
    return p
}
func (p *numericParam[T]) Range() (min, max, step T) {
    return p.min, p.max, p.step
}
func (p *numericParam[T]) Min() (T) {
    return p.min
}
func (p *numericParam[T]) Max() (T) {
    return p.max
}
func (p *numericParam[T]) Step() (T) {
    return p.step
}

//----------------------------------------------------------------------------

type BoolParameter interface {
    Parameter
    Get() (bool, error)
    Set(v bool) error
    Val() bool
    SetVal(v bool)
    Reset()
}
func NewBoolParameter(name string, init bool) BoolParameter {
    return newBaseParam(name, init)
}

type IntParameter interface {
    Parameter
    Get() (int, error)
    Set(v int) error
    Val() int
    SetVal(v int)
    Reset()
    Range() (min, max, step int)
    Min() int
    Max() int
    Step() int
}
func NewIntParameter(name string, init, min, max, step int) IntParameter {
    return newNumericParam(name, init, min, max, step)
}

type FloatParameter interface {
    Parameter
    Get() (float64, error)
    Set(v float64) error
    Val() float64
    SetVal(v float64)
    Reset()
    Range() (min, max, step float64)
    Min() float64
    Max() float64
    Step() float64
}
func NewFloatParameter(name string, init, min, max, step float64) FloatParameter {
    return newNumericParam(name, init, min, max, step)
}

type StringParameter interface {
    Parameter
    Get() (string, error)
    Set(v string) error
    Val() string
    SetVal(v string)
    Reset()
}
func NewStringParameter(name string, init string) StringParameter {
    return newBaseParam(name, init)
}

type PaletteParameter interface {
    Parameter
    Get() (ColorSource, error)
    Set(v ColorSource) error
    Val() ColorSource
    SetVal(v ColorSource)
    Reset()
}
func NewPaletteParameter(name string, init ColorSource) PaletteParameter {
    return newBaseParam(name, init)
}
