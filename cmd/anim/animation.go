package main

import (
	"encoding/gob"
	"math"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/colornames"
)

// Fuer Animationen, die endlos wiederholt weren sollen, kann diese Konstante
// fuer die Anzahl Wiederholungen verwendet werden.
const (
	AnimationRepeatForever = -1
)

// Mit dem Funktionstyp [AnimationCurve] kann der Verlauf einer Animation
// beeinflusst werden. Der Parameter [t] ist ein Wert im Intervall [0,1]
// und zeigt an, wo sich die Animation gerade befindet (t=0: Animation
// wurde eben gestartet; t=1: Animation ist zu Ende). Der Rueckgabewert
// ist ebenfalls ein Wert im Intervall [0,1] und hat die gleiche Bedeutung
// wie [t].
type AnimationCurve func(t float64) float64

// Bezeichnet eine lineare Animation.
func AnimationLinear(t float64) float64 {
	return t
}

// Beginnt langsam und nimmt immer mehr an Fahrt auf.
func AnimationEaseIn(t float64) float64 {
	return t * t
}

// Beginnt schnell und bremst zum Ende immer mehr ab.
func AnimationEaseOut(t float64) float64 {
	return t * (2 - t)
}

// Anfang und Ende der Animation werden abgebremst.
func AnimationEaseInOut(t float64) float64 {
	if t <= 0.5 {
		return t * t * 2
	}
	return -1 + (4-t*2)*t
}

// Alternative Funktion mit einer kubischen Funktion:
// f(x) = 3x^2 - 2x^3
func AnimationEaseInOutNew(t float64) float64 {
	return (3 - 2*t) * t * t
}

// Dies ist ein etwas unbeholfener Versuch, die Zielwerte bestimmter
// Animationen dynamisch berechnen zu lassen.
type ColorFuncType func() ledgrid.LedColor
type PointFuncType func() geom.Point
type FloatFuncType func() float64
type AlphaFuncType func() uint8

func RandColor() ColorFuncType {
	return func() ledgrid.LedColor {
		return colornames.RandColor()
	}
}
func RandPoint(r geom.Rectangle) PointFuncType {
	return func() geom.Point {
		fx := rand.Float64()
		fy := rand.Float64()
		return r.RelPos(fx, fy)
	}
}
func RandPointTrunc(r geom.Rectangle, t float64) PointFuncType {
	return func() geom.Point {
		fx := rand.Float64()
		fy := rand.Float64()
		p := r.RelPos(fx, fy)
		p.X = t * math.Round(p.X/t)
		p.Y = t * math.Round(p.Y/t)
		return p
	}
}
func RandSize(s1, s2 geom.Point) PointFuncType {
	return func() geom.Point {
		t := rand.Float64()
		return s1.Interpolate(s2, t)
	}
}
func RandFloat(a, b float64) FloatFuncType {
	return func() float64 {
		return a + (b-a)*rand.Float64()
	}
}
func RandAlpha(a, b uint8) AlphaFuncType {
	return func() uint8 {
		return a + uint8(rand.UintN(uint(b-a)))
	}
}

// Registriert alle Animationsarten, um sie exportieren oder importieren
// zu koennen.
func init() {
	gob.Register(&Group{})
	gob.Register(&Sequence{})
	gob.Register(&Timeline{})
	gob.Register(&AnimationEmbed{})
	gob.Register(&ColorAnimation{})
	gob.Register(&PaletteAnimation{})
	// gob.Register(&PositionAnimation{})
	gob.Register(&FloatAnimation{})
	gob.Register(&PathAnimation{})
}

// Dieses Interface ist von allen Typen zu implementieren, welche
// Animationen ausfuehren sollen/wollen.
type Animator interface {
	AddAnim(anims ...Animation)
	DelAnim(anim Animation)
	DelAllAnim()
	Stop()
	Continue()
}

// Das Interface fuer jede Art von Animation (bis jetzt zumindest).
type Animation interface {
	Duration() time.Duration
	SetDuration(dur time.Duration)
	Start()
	Stop()
	Continue()
	IsStopped() bool
	Update(t time.Time) bool
}

// Jede konkrete Animation (Farben, Positionen, Groessen, etc.) muss das
// Interface AnimationImpl implementieren.
type AnimationImpl interface {
	// Init wird vom Animationsframework aufgerufen, wenn diese Animation
	// gestartet wird. Wichtig: Wiederholungen und Umkehrungen (AutoReverse)
	// zaehlen nicht als Start!
	Init()
	// In Tick schliesslich ist die eigentliche Animationslogik enthalten.
	// Der Parameter t gibt an, wie weit die Animation bereits gelaufen ist.
	// t=0: Animation wurde eben gestartet
	// t=1: Die Animation ist fertig gelaufen.
	Tick(t float64)
}

type DurationEmbed struct {
    duration time.Duration
}

func (d *DurationEmbed) Duration() time.Duration {
	return d.duration
}

func (d *DurationEmbed) SetDuration(dur time.Duration) {
	d.duration = dur
}

// Eine Gruppe dient dazu, eine Anzahl von Animationen gleichzeitig zu starten.
// Die Laufzeit der Gruppe ist gleich der laengsten Laufzeit ihrer Animationen
// oder einer festen Dauer (je nach dem was groesser ist).
// Die Animationen, welche ueber eine Gruppe gestartet werden, sollten keine
// Endlos-Animationen sein, da sonst die Laufzeit der Gruppe ebenfalls
// unendlich wird.
type Group struct {
    DurationEmbed
	// Gibt an, wie oft diese Gruppe wiederholt werden soll.
	RepeatCount int

	animList         []Animation
	start, stop, end time.Time
	repeatsLeft      int
	running          bool
}

// Erstellt eine neue Gruppe, welche die Animationen in [anims] zusammen
// startet. Per Default ist die Laufzeit der Gruppe gleich der laengsten
// Laufzeit der hinzugefuegten Animationen.
func NewGroup(anims ...Animation) *Group {
	a := &Group{}
	a.duration = 0
	a.RepeatCount = 0
	a.Add(anims...)
	AnimCtrl.AddAnim(a)
	return a
}

// Fuegt der Gruppe weitere Animationen hinzu.
func (a *Group) Add(anims ...Animation) {
	for _, anim := range anims {
		if anim.Duration() > a.duration {
			a.duration = anim.Duration()
		}
		a.animList = append(a.animList, anim)
	}
}

// Startet die Gruppe.
func (a *Group) Start() {
	if a.running {
		return
	}
	a.start = time.Now()
	a.end = a.start.Add(a.duration)
	a.repeatsLeft = a.RepeatCount
	a.running = true
	for _, anim := range a.animList {
		anim.Start()
	}
}

// Unterbricht die Ausfuehrung der Gruppe.
func (a *Group) Stop() {
	if !a.running {
		return
	}
	a.stop = time.Now()
	a.running = false
}

// Setzt die Ausfuehrung der Gruppe fort.
func (a *Group) Continue() {
	if a.running {
		return
	}
	dt := time.Now().Sub(a.stop)
	a.start = a.start.Add(dt)
	a.end = a.end.Add(dt)
	a.running = true
}

// Liefert den Status der Gruppe zurueck.
func (a *Group) IsStopped() bool {
	return !a.running
}

func (a *Group) Update(t time.Time) bool {
	for _, anim := range a.animList {
		if !anim.IsStopped() {
			return true
		}
	}
	if t.After(a.end) {
		if a.repeatsLeft == 0 {
			a.running = false
			return false
		} else if a.repeatsLeft > 0 {
			a.repeatsLeft--
		}
		a.start = a.end
		a.end = a.start.Add(a.duration)
		for _, anim := range a.animList {
			anim.Start()
		}
	}
	return true
}

// Mit einer Sequence lassen sich eine Reihe von Animationen hintereinander
// ausfuehren. Dabei wird eine nachfolgende Animation erst dann gestartet,
// wenn die vorangehende beendet wurde.
type Sequence struct {
    DurationEmbed
	// Gibt an, wie oft diese Sequenz wiederholt werden soll.
	RepeatCount int

	animList         []Animation
	currAnim         int
	start, stop, end time.Time
	repeatsLeft      int
	running          bool
}

// Erstellt eine neue Sequenz welche die Animationen in [anims] hintereinander
// ausfuehrt.
func NewSequence(anims ...Animation) *Sequence {
	a := &Sequence{}
	a.duration = 0
	a.RepeatCount = 0
	a.Add(anims...)
	AnimCtrl.AddAnim(a)
	return a
}

// Fuegt der Sequenz weitere Animationen hinzu.
func (a *Sequence) Add(anims ...Animation) {
	for _, anim := range anims {
		a.duration = a.duration + anim.Duration()
		a.animList = append(a.animList, anim)
	}
}

// Startet die Sequenz.
func (a *Sequence) Start() {
	if a.running {
		return
	}
	a.start = time.Now()
	a.end = a.start.Add(a.duration)
	a.currAnim = 0
	a.repeatsLeft = a.RepeatCount
	a.running = true
	a.animList[a.currAnim].Start()
}

// Unterbricht die Ausfuehrung der Sequenz.
func (a *Sequence) Stop() {
	if !a.running {
		return
	}
	a.stop = time.Now()
	a.running = false
}

// Setzt die Ausfuehrung der Sequenz fort.
func (a *Sequence) Continue() {
	if a.running {
		return
	}
	dt := time.Now().Sub(a.stop)
	a.start = a.start.Add(dt)
	a.end = a.end.Add(dt)
	a.running = true
}

// Liefert den Status der Sequenz zurueck.
func (a *Sequence) IsStopped() bool {
	return !a.running
}

// Wird durch den Controller periodisch aufgerufen, prueft ob Animationen
// dieser Sequenz noch am Laufen sind und startet ggf. die naechste.
func (a *Sequence) Update(t time.Time) bool {
	if a.currAnim < len(a.animList) {
		if !a.animList[a.currAnim].IsStopped() {
			return true
		}
		a.currAnim++
	}
	if a.currAnim >= len(a.animList) {
		if t.After(a.end) {
			if a.repeatsLeft == 0 {
				a.running = false
				return false
			} else if a.repeatsLeft > 0 {
				a.repeatsLeft--
			}
			a.start = a.end
			a.end = a.start.Add(a.duration)
			a.currAnim = 0
			a.animList[a.currAnim].Start()
		}
		return true
	}
	a.animList[a.currAnim].Start()
	return true
}

// Mit einer Timeline koennen einzelne oder mehrere Animationen zu
// bestimmten Zeiten gestartet werden. Die Zeit ist relativ zur Startzeit
// der Timeline selber zu verstehen. Nach dem Start werden die Animationen
// nicht mehr weiter kontrolliert.
type Timeline struct {
    DurationEmbed
	// Gibt an, wie oft diese Timeline wiederholt werden soll.
	RepeatCount int

	posList          []*timelinePos
	nextPos          int
	start, stop, end time.Time
	repeatsLeft      int
	running          bool
}

// Interner Typ, mit dem Ausfuehrungszeitpunkt und Animationen festgehalten
// werden koennen.
type timelinePos struct {
	dt    time.Duration
	anims []Animation
}

// Erstellt eine neue Timeline mit Ausfuehrungsdauer d. Als d kann auch Null
// angegeben werden, dann ist die Laufzeit der Timeline gleich dem groessten
// Ausfuehrungszeitpunkt der hinterlegten Animationen.
func NewTimeline(d time.Duration) *Timeline {
	a := &Timeline{}
	a.duration = d
	a.RepeatCount = 0
	a.posList = make([]*timelinePos, 0)
	AnimCtrl.AddAnim(a)
	return a
}

// Fuegt der Timeline die Animation anim hinzu mit Ausfuehrungszeitpunkt
// dt nach Start der Timeline. Im Moment muessen die Animationen noch in
// der Reihenfolge ihres Ausfuehrungszeitpunktes hinzugefuegt werden.
func (a *Timeline) Add(pit time.Duration, anims ...Animation) {
	var i int

	if pit > a.duration {
		a.duration = pit
	}

	for i = 0; i < len(a.posList); i++ {
		pos := a.posList[i]
		if pos.dt == pit {
			pos.anims = append(pos.anims, anims...)
			return
		}
		if pos.dt > pit {
			break
		}
	}
	a.posList = slices.Insert(a.posList, i, &timelinePos{pit, anims})
}

// Startet die Timeline.
func (a *Timeline) Start() {
	if a.running {
		return
	}
	a.start = time.Now()
	a.end = a.start.Add(a.duration)
	a.repeatsLeft = a.RepeatCount
	a.nextPos = 0
	a.running = true
}

// Unterbricht die Ausfuehrung der Timeline.
func (a *Timeline) Stop() {
	if !a.running {
		return
	}
	a.stop = time.Now()
	a.running = false
}

// Setzt die Ausfuehrung der Timeline fort.
func (a *Timeline) Continue() {
	if a.running {
		return
	}
	dt := time.Now().Sub(a.stop)
	a.start = a.start.Add(dt)
	a.end = a.end.Add(dt)
	a.running = true
}

// Retourniert den Status der Timeline.
func (a *Timeline) IsStopped() bool {
	return !a.running
}

// Wird periodisch durch den Controller aufgerufen und aktualisiert die
// Timeline.
func (a *Timeline) Update(t time.Time) bool {
	if a.nextPos >= len(a.posList) {
		if t.After(a.end) {
			if a.repeatsLeft == 0 {
				a.running = false
				return false
			} else if a.repeatsLeft > 0 {
				a.repeatsLeft--
			}
			a.start = a.end
			a.end = a.start.Add(a.duration)
			a.nextPos = 0
		}
		return true
	}
	pos := a.posList[a.nextPos]
	if t.Sub(a.start) >= pos.dt {
		for _, anim := range pos.anims {
			anim.Start()
		}
		a.nextPos++
	}
	return true
}

// Embeddable mit in allen Animationen benoetigen Variablen und Methoden.
// Erleichert das Erstellen von neuen Animationen gewaltig.
type AnimationEmbed struct {
    DurationEmbed
	// Falls true, wird die Animation einmal vorwaerts und einmal rueckwerts
	// abgespielt.
	AutoReverse bool
	// Curve bezeichnet eine Interpolationsfunktion, welche einen beliebigen
	// Verlauf der Animation erlaubt (Beschleunigung am Anfang, Abbremsen
	// gegen Schluss, etc).
	Curve AnimationCurve
	// Bezeichnet die Anzahl Wiederholungen dieser Animation.
	RepeatCount int

	// duration         time.Duration
	wrapper          Animation
	reverse          bool
	start, stop, end time.Time
	total            float64
	repeatsLeft      int
	running          bool
	tick             func(t float64)
	init             func()
}

// Muss beim Erstellen einer Animation aufgerufen werden, welche dieses
// Embeddable einbindet.
func (a *AnimationEmbed) ExtendAnimation(wrapper Animation) {
	a.AutoReverse = false
	a.Curve = AnimationEaseInOut
	a.RepeatCount = 0
	a.wrapper = wrapper
	a.running = false
	a.tick = wrapper.(AnimationImpl).Tick
	a.init = wrapper.(AnimationImpl).Init
	AnimCtrl.AddAnim(wrapper)
}

func (a *AnimationEmbed) Duration() time.Duration {
	factor := 1
	if a.RepeatCount > 0 {
		factor += a.RepeatCount
	}
	if a.AutoReverse {
		factor *= 2
	}
	return time.Duration(factor) * a.duration
}

// func (a *AnimationEmbed) SetDuration(dur time.Duration) {
// 	a.duration = dur
// }

// Startet die Animation mit jenen Parametern, die zum Startzeitpunkt
// aktuell sind. Ist die Animaton bereits am Laufen ist diese Methode
// ein no-op.
func (a *AnimationEmbed) Start() {
	if a.running {
		return
	}
	a.start = time.Now()
	a.end = a.start.Add(a.duration)
	a.total = a.end.Sub(a.start).Seconds()
	a.repeatsLeft = a.RepeatCount
	a.reverse = false
	a.running = true
	a.init()
}

// Haelt die Animation an, laesst sie jedoch in der Animation-Queue der
// Applikation. Mit [Continue] kann eine gestoppte Animation wieder
// fortgesetzt werden.
func (a *AnimationEmbed) Stop() {
	if !a.running {
		return
	}
	a.stop = time.Now()
	a.running = false
}

// Setzt eine mit [Stop] angehaltene Animation wieder fort.
func (a *AnimationEmbed) Continue() {
	if a.running {
		return
	}
	dt := time.Now().Sub(a.stop)
	a.start = a.start.Add(dt)
	a.end = a.end.Add(dt)
	a.running = true
}

// Liefert true, falls die Animation mittels [Stop] angehalten wurde oder
// falls die Animation zu Ende ist.
func (a *AnimationEmbed) IsStopped() bool {
	return !a.running
}

// Diese Methode ist fuer die korrekte Abwicklung (Beachtung von Reverse und
// RepeatCount, etc) einer Animation zustaendig. Wenn die Animation zu Ende
// ist, retourniert Update false. Der Parameter [t] ist ein "Point in Time".
func (a *AnimationEmbed) Update(t time.Time) bool {
	if t.After(a.end) {
		if a.reverse {
			a.tick(a.Curve(0.0))
			if a.repeatsLeft == 0 {
				a.running = false
				return false
			}
			a.reverse = false
		} else {
			a.tick(a.Curve(1.0))
			if a.AutoReverse {
				a.reverse = true
			}
		}
		if !a.reverse {
			if a.repeatsLeft == 0 {
				a.running = false
				return false
			}
			if a.repeatsLeft > 0 {
				a.repeatsLeft--
			}
		}
		a.start = a.end
		a.end = a.start.Add(a.duration)
		a.init()
		return true
	}

	delta := t.Sub(a.start).Seconds()
	val := delta / a.total
	if a.reverse {
		a.tick(a.Curve(1.0 - val))
	} else {
		a.tick(a.Curve(val))
	}
	return true
}

// Animation fuer einen Verlauf zwischen zwei Farben.
type ColorAnimation struct {
	AnimationEmbed
	Cont       bool
	ValPtr     *ledgrid.LedColor
	Val1, Val2 ledgrid.LedColor
	ValFunc    ColorFuncType
}

func NewColorAnimation(valPtr *ledgrid.LedColor, val2 ledgrid.LedColor, dur time.Duration) *ColorAnimation {
	a := &ColorAnimation{}
	a.AnimationEmbed.ExtendAnimation(a)
	a.SetDuration(dur)
	a.ValPtr = valPtr
	a.Val1 = *valPtr
	a.Val2 = val2
	return a
}

func (a *ColorAnimation) Init() {
	if a.Cont {
		a.Val1 = *a.ValPtr
	}
	if a.ValFunc != nil {
		a.Val2 = a.ValFunc()
	}
}

func (a *ColorAnimation) Tick(t float64) {
    alpha := (*a.ValPtr).A
	*a.ValPtr = a.Val1.Interpolate(a.Val2, t)
    (*a.ValPtr).A = alpha
}

// Will man allerdings nur die Durchsichtigkeit (den Alpha-Wert) einer Farbe
// veraendern und kennt beispielsweise die Farbe selber gar nicht, dann ist
// die AlphaAnimation genau das Richtige.
type AlphaAnimation struct {
	AnimationEmbed
	Cont       bool
	ValPtr     *uint8
	Val1, Val2 uint8
	ValFunc    AlphaFuncType
}

func NewAlphaAnimation(valPtr *uint8, val2 uint8, dur time.Duration) *AlphaAnimation {
	a := &AlphaAnimation{}
	a.AnimationEmbed.ExtendAnimation(a)
	a.SetDuration(dur)
	a.ValPtr = valPtr
	a.Val1 = *valPtr
	a.Val2 = val2
	return a
}

func (a *AlphaAnimation) Init() {
	if a.Cont {
		a.Val1 = *a.ValPtr
	}
	if a.ValFunc != nil {
		a.Val2 = a.ValFunc()
	}
}

func (a *AlphaAnimation) Tick(t float64) {
	*a.ValPtr = uint8((1.0-t)*float64(a.Val1) + t*float64(a.Val2))
}

// Animation fuer einen Farbverlauf ueber die Farben einer Palette.
type PaletteAnimation struct {
	AnimationEmbed
	ValPtr *ledgrid.LedColor
	pal    ledgrid.ColorSource
}

func NewPaletteAnimation(valPtr *ledgrid.LedColor, pal ledgrid.ColorSource, dur time.Duration) *PaletteAnimation {
	a := &PaletteAnimation{}
	a.AnimationEmbed.ExtendAnimation(a)
	a.SetDuration(dur)
	a.Curve = AnimationLinear
	a.ValPtr = valPtr
	a.pal = pal
	return a
}

func (a *PaletteAnimation) Init() {}

func (a *PaletteAnimation) Tick(t float64) {
	*a.ValPtr = a.pal.Color(t)
}

// Animation fuer einen Wechsel der Position, resp. Veraenderung der Groesse.
// type PositionAnimation struct {
// 	AnimationEmbed
// 	Cont       bool
// 	ValPtr     *geom.Point
// 	Val1, Val2 geom.Point
// 	ValFunc    PointFuncType
// }

// func NewPositionAnimationOld(valPtr *geom.Point, val2 geom.Point, dur time.Duration) *PositionAnimation {
// 	a := &PositionAnimation{}
// 	a.AnimationEmbed.ExtendAnimation(a)
//     a.SetDuration(dur)
// 	a.ValPtr = valPtr
// 	a.Val1 = *valPtr
// 	a.Val2 = val2
// 	return a
// }

// func (a *PositionAnimation) Init() {
// 	if a.Cont {
// 		a.Val1 = *a.ValPtr
// 	}
// 	if a.ValFunc != nil {
// 		a.Val2 = a.ValFunc()
// 	}
// }

// func (a *PositionAnimation) Tick(t float64) {
// 	*a.ValPtr = a.Val1.Interpolate(a.Val2, t)
// }

// Da Positionen und Groessen mit dem gleichen Objekt aus geom realisiert
// werden (geom.Point), ist die Animation einer Groesse und einer Position
// im Wesentlichen das gleiche. Die Funktion NewSizeAnimation ist als
// syntaktische Vereinfachung zu verstehen.

// Animation fuer einen Verlauf zwischen zwei Fliesskommazahlen.
type FloatAnimation struct {
	AnimationEmbed
	Cont       bool
	ValPtr     *float64
	Val1, Val2 float64
	ValFunc    FloatFuncType
}

func NewFloatAnimation(valPtr *float64, val2 float64, dur time.Duration) *FloatAnimation {
	a := &FloatAnimation{}
	a.AnimationEmbed.ExtendAnimation(a)
	a.SetDuration(dur)
	a.ValPtr = valPtr
	a.Val1 = *valPtr
	a.Val2 = val2
	return a
}

func (a *FloatAnimation) Init() {
	if a.Cont {
		a.Val1 = *a.ValPtr
	}
	if a.ValFunc != nil {
		a.Val2 = a.ValFunc()
	}
}

func (a *FloatAnimation) Tick(t float64) {
	*a.ValPtr = (1-t)*a.Val1 + t*a.Val2
}

// Animation fuer das Fahren entlang eines Pfades. Mit fnc kann eine konkrete,
// Pfad-generierende Funktion angegeben werden. Siehe auch [PathFunction]
type PathAnimation struct {
	AnimationEmbed
	Cont       bool
	ValPtr     *geom.Point
	Val1, Val2 geom.Point
	Size       geom.Point
	ValFunc    PointFuncType
	PathFunc   PathFunctionType
}

func NewPathAnimation(valPtr *geom.Point, pathFunc PathFunctionType, size geom.Point, dur time.Duration) *PathAnimation {
	a := &PathAnimation{}
	a.AnimationEmbed.ExtendAnimation(a)
	a.SetDuration(dur)
	a.ValPtr = valPtr
	a.Val1 = *valPtr
	a.Size = size
	a.PathFunc = pathFunc
	return a
}

func NewPositionAnimation(valPtr *geom.Point, val2 geom.Point, dur time.Duration) *PathAnimation {
	a := &PathAnimation{}
	a.AnimationEmbed.ExtendAnimation(a)
	a.SetDuration(dur)
	a.ValPtr = valPtr
	a.Val1 = *valPtr
    a.Val2 = val2
	a.PathFunc = LinearPath
	return a
}

var (
	NewSizeAnimation = NewPositionAnimation
)

func (a *PathAnimation) Init() {
	if a.Cont {
		a.Val1 = *a.ValPtr
	}
	if a.ValFunc != nil {
        if !a.Size.Eq(geom.Point{}) {
            a.Size = a.ValFunc()
        } else {
		    a.Val2 = a.ValFunc()
        }
	}
}

func (a *PathAnimation) Tick(t float64) {
    var dp geom.Point

    if !a.Size.Eq(geom.Point{}) {
        dp = a.PathFunc(t).Scale(a.Size)
    } else {
        	dp = a.PathFunc(t).Scale(a.Val2.Sub(a.Val1))
    }
	*a.ValPtr = a.Val1.Add(dp)
}

//----------------------------------------------------------------------------

func NewPolyPathAnimation(valPtr *geom.Point, polyPath *PolygonPath, dur time.Duration) *PathAnimation {
	a := &PathAnimation{}
	a.AnimationEmbed.ExtendAnimation(a)
	a.SetDuration(dur)
	a.ValPtr = valPtr
	a.Val1 = *valPtr
	a.Size = geom.Point{1, 1}
	a.PathFunc = polyPath.RelPoint
	return a
}

// Neben den vorhandenen Pfaden (Kreise, Halbkreise, Viertelkreise) koennen
// Positions-Animationen auch entlang komplett frei definierten Pfaden
// erfolgen. Der Schluessel dazu ist der Typ PolygonPath.
type PolygonPath struct {
	rect     geom.Rectangle
	stopList []polygonStop
}

type polygonStop struct {
	len float64
	pos geom.Point
}

// Erstellt ein neues PolygonPath-Objekt und verwendet die Punkte in points
// als Eckpunkte eines offenen Polygons. Punkte koennen nur beim Erstellen
// angegeben werden.
func NewPolygonPath(points ...geom.Point) *PolygonPath {
	p := &PolygonPath{}
	p.stopList = make([]polygonStop, len(points))

	origin := geom.Point{}
	for i, point := range points {
		if i == 0 {
			origin = point
			p.stopList[i] = polygonStop{0.0, geom.Point{0, 0}}
			continue
		}
		pt := point.Sub(origin)
		len := p.stopList[i-1].len + pt.Distance(p.stopList[i-1].pos)
		p.stopList[i] = polygonStop{len, pt}

		p.rect.Min = p.rect.Min.Min(pt)
		p.rect.Max = p.rect.Max.Max(pt)
	}
	return p
}

// Diese Methode ist bei der Erstellung einer Pfad-Animation als Parameter
// fnc anzugeben.
func (p *PolygonPath) RelPoint(t float64) geom.Point {
	dist := t * p.stopList[len(p.stopList)-1].len
	for i, stop := range p.stopList[1:] {
		if dist < stop.len {
			p1 := p.stopList[i].pos
			p2 := stop.pos
			relDist := dist - p.stopList[i].len
			f := relDist / (stop.len - p.stopList[i].len)
			return p1.Interpolate(p2, f)
		}
	}
	return p.stopList[len(p.stopList)-1].pos
}

// Die PathFunctionType muss folgende Bedingungen erfuellen:
//  1. t ist in [0,1]
//  2. f(0) = (0,0)
//  3. max(f(t).X) - min(f(t).X) = 1.0 und
//     max(f(t).Y) - min(f(t).Y) = 1.0
type PathFunctionType func(t float64) geom.Point

// Beschreibt eine Gerade
func LinearPath(t float64) geom.Point {
	return geom.Point{t, t}
}

// Beschreibt ein Rechteck im Uhrzeigersinn.
// Startpunkt ist auf 12 Uhr.
func RectanglePathA(t float64) geom.Point {
	switch {
	case t < 1.0/8.0:
		return geom.Point{0.5 * 8.0 * t, 0.0}
	case t < 3.0/8.0:
		return geom.Point{0.5, 4.0 * (t - 1.0/8.0)}
	case t < 5.0/8.0:
		return geom.Point{0.5 - 4.0*(t-3.0/8.0), 1.0}
	case t < 7.0/8.0:
		return geom.Point{-0.5, 1.0 - 4.0*(t-5.0/8.0)}
	default:
		return geom.Point{-0.5 + 0.5*8.0*(t-7.0/8.0), 0.0}
	}
}

// Beschreibt ein Rechteck im Uhrzeigersinn.
// Startpunkt ist auf 9 Uhr.
func RectanglePathB(t float64) geom.Point {
	switch {
	case t < 1.0/8.0:
		return geom.Point{0.0, -0.5 * 8.0 * t}
	case t < 3.0/8.0:
		return geom.Point{4.0 * (t - 1.0/8.0), -0.5}
	case t < 5.0/8.0:
		return geom.Point{1.0, 4.0*(t-3.0/8.0) - 0.5}
	case t < 7.0/8.0:
		return geom.Point{1.0 - 4.0*(t-5.0/8.0), 0.5}
	default:
		return geom.Point{0, 0.5 * (1.0 - 8.0*(t-7.0/8.0))}
	}
}

// Beschreibt einen Kreis oder Ellipse im Uhrzeigersinn.
// Startpunkt ist auf 12 Uhr.
func FullCirclePathA(t float64) geom.Point {
	phi := t * 2 * math.Pi
	return geom.Point{0.5 * math.Sin(phi), 0.5 - 0.5*math.Cos(phi)}
}

// Startpunkt ist auf 9 Uhr.
func FullCirclePathB(t float64) geom.Point {
	phi := t * 2 * math.Pi
	return geom.Point{0.5 - 0.5*math.Cos(phi), -(0.5 * math.Sin(phi))}
}

// Beschreibt einen Halbkreis.
func HalfCirclePathA(t float64) geom.Point {
	phi := t * math.Pi
	return geom.Point{math.Sin(phi), (1.0 - math.Cos(phi)) / 2.0}
}

func HalfCirclePathB(t float64) geom.Point {
	phi := t * math.Pi
	return geom.Point{(1.0 - math.Cos(phi)) / 2.0, math.Sin(phi)}
}

// Beschreibt einen Viertelkreis.
// Horizontaler Start.
func QuarterCirclePathA(t float64) geom.Point {
	phi := t * math.Pi / 2.0
	return geom.Point{math.Sin(phi), 1.0 - math.Cos(phi)}
}

// Vertikaler Start.
func QuarterCirclePathB(t float64) geom.Point {
	phi := t * math.Pi / 2.0
	return geom.Point{1.0 - math.Cos(phi), math.Sin(phi)}
}
