package ledgrid

import (
	"encoding/gob"
	"image"
	"log"
	"math"
	"math/rand/v2"
	"os"
	"runtime"
	"sync"
	"time"

	"golang.org/x/image/math/fixed"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid/color"
)

var (
	AnimCtrl *AnimationController
)

// Fuer Animationen, die endlos wiederholt weren sollen, kann diese Konstante
// fuer die Anzahl Wiederholungen verwendet werden.
const (
	AnimationRepeatForever = -1
	refreshRate            = 30 * time.Millisecond
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

// Beginnt langsam und nimmt immer mehr an Fahrt auf.
func AnimationLazeIn(t float64) float64 {
	return t * t * t
}

// Beginnt schnell und bremst zum Ende immer mehr ab.
func AnimationEaseOut(t float64) float64 {
	return t * (2 - t)
}

// Beginnt langsam und nimmt immer mehr an Fahrt auf.
func AnimationLazeOut(t float64) float64 {
	return t * (t*(t-3) + 3)
}

// Anfang und Ende der Animation werden abgebremst.
func AnimationEaseInOut(t float64) float64 {
	if t <= 0.5 {
		return 2 * t * t
	}
	return (4-2*t)*t - 1
}

// Anfang und Ende der Animation werden abgebremst.
func AnimationLazeInOut(t float64) float64 {
	if t <= 0.5 {
		return 4 * t * t * t
	}
	return 4*(t-1)*(t-1)*(t-1) + 1
}

// Alternative Funktion mit einer kubischen Funktion:
// f(x) = 3x^2 - 2x^3
func AnimationEaseInOutNew(t float64) float64 {
	return (3 - 2*t) * t * t
}

// Dies ist ein etwas unbeholfener Versuch, die Zielwerte bestimmter
// Animationen dynamisch berechnen zu lassen. Alle XXXFuncType sind
// Funktionstypen fuer einen bestimmten Datentyp, der in den Animationen
// verwendet wird und dessen dynamische Berechnung Sinn macht.
type PaletteFuncType func() ColorSource

// type AlphaFuncType func() uint8

var (
	palId int = 0
)

// SeqPalette liefert eine Funktion, die bei jedem Aufruf die naechste Palette
// als Resultat zurueckgibt.
func SeqPalette() PaletteFuncType {
	return func() ColorSource {
		name := PaletteNames[palId]
		palId = (palId + 1) % len(PaletteNames)
		return PaletteMap[name]
	}
}

// RandPalette liefert eine Funktion, die bei jedem Aufruf eine zufaellig
// gewaehlte Palette retourniert.
func RandPalette() PaletteFuncType {
	return func() ColorSource {
		name := PaletteNames[rand.IntN(len(PaletteNames))]
		return PaletteMap[name]
	}
}

// Liefert bei jedem Aufruf einen zufaellig gewaehlten Punkt innerhalb des
// Rechtecks r.
func RandPoint(r geom.Rectangle) AnimValueFunc[geom.Point] {
	return func() geom.Point {
		fx := rand.Float64()
		fy := rand.Float64()
		return r.RelPos(fx, fy)
	}
}

// Wie RandPoint, sorgt jedoch dafuer dass die Koordinaten auf ein Vielfaches
// von t abgeschnitten werden.
func RandPointTrunc(r geom.Rectangle, t float64) AnimValueFunc[geom.Point] {
	return func() geom.Point {
		fx := rand.Float64()
		fy := rand.Float64()
		p := r.RelPos(fx, fy)
		p.X = t * math.Round(p.X/t)
		p.Y = t * math.Round(p.Y/t)
		return p
	}
}

// Macht eine Interpolation zwischen den Groessen s1 und s2. Der Interpolations-
// punkt wird zufaellig gewaehlt.
func RandSize(s1, s2 geom.Point) AnimValueFunc[geom.Point] {
	return func() geom.Point {
		t := rand.Float64()
		return s1.Interpolate(s2, t)
	}
}

// Liefert eine zufaellig gewaehlte Fliesskommazahl im Interval [a,b).
func RandFloat(a, b float64) AnimValueFunc[float64] {
	return func() float64 {
		return a + (b-a)*rand.Float64()
	}
}

// Liefert eine zufaellig gewaehlte natuerliche Zahl im Interval [a,b).
func RandAlpha(a, b uint8) AnimValueFunc[uint8] {
	return func() uint8 {
		return a + uint8(rand.UintN(uint(b-a)))
	}
}

// Das Herzstueck der ganzen Animationen ist der AnimationController. Er
// sorgt dafuer, dass alle 30 ms (siehe Variable refreshRate) die
// Update-Methoden aller Animationen aufgerufen werden und veranlasst im
// Anschluss, dass alle darstellbaren Objekte neu gezeichnet werden und sendet
// das Bild schliesslich dem PixelController (oder dem PixelEmulator).
type AnimationController struct {
	AnimList  []Animation
	animMutex *sync.RWMutex
	Canvas    *Canvas
	// Filter     Filter
	ledGrid    *LedGrid
	ticker     *time.Ticker
	quit       bool
	animPit    time.Time
	animWatch  *Stopwatch
	numThreads int
	stop       time.Time
	delay      time.Duration
	running    bool
	syncChan   chan bool
}

func NewAnimationController(canvas *Canvas, ledGrid *LedGrid) *AnimationController {
	if AnimCtrl != nil {
		return AnimCtrl
	}
	a := &AnimationController{}
	a.AnimList = make([]Animation, 0)
	a.animMutex = &sync.RWMutex{}
	a.Canvas = canvas
	// a.Filter = NewFilterIdent(canvas)
	a.ledGrid = ledGrid
	a.ticker = time.NewTicker(refreshRate)
	a.animWatch = NewStopwatch()
	a.numThreads = runtime.NumCPU()
	a.delay = time.Duration(0)
	a.syncChan = make(chan bool)

	AnimCtrl = a
	go a.backgroundThread()
	a.running = true

	canvas.StartRefresh(a.syncChan, ledGrid.syncChan)

	return a
}

// Fuegt weitere Animationen hinzu. Der Zugriff auf den entsprechenden Slice
// wird synchronisiert, da die Bearbeitung der Animationen durch den
// Background-Thread ebenfalls relativ haeufig auf den Slice zugreift.
func (a *AnimationController) Add(anims ...Animation) {
	a.animMutex.Lock()
	a.AnimList = append(a.AnimList, anims...)
	a.animMutex.Unlock()
}

// Loescht eine einzelne Animation.
func (a *AnimationController) Del(anim Animation) {
	a.animMutex.Lock()
	defer a.animMutex.Unlock()
	for idx, obj := range a.AnimList {
		if obj == anim {
			obj.Suspend()
			a.AnimList[idx] = nil
			return
		}
	}
}

func (a *AnimationController) DelAt(idx int) {
	a.animMutex.Lock()
	a.AnimList[idx].Suspend()
	a.AnimList[idx] = nil
	a.animMutex.Unlock()
}

// Loescht alle Animationen.
func (a *AnimationController) Purge() {
	a.animMutex.Lock()
	for _, anim := range a.AnimList {
		if anim == nil {
			continue
		}
		anim.Suspend()
	}
	a.AnimList = a.AnimList[:0]
	a.animMutex.Unlock()
}

// Mit Stop koennen die Animationen und die Darstellung auf der Hardware
// unterbunden werden.
func (a *AnimationController) Suspend() {
	if !a.running {
		return
	}
	a.running = false
	a.ticker.Stop()
	a.stop = time.Now()
}

// Setzt die Animationen wieder fort.
// TO DO: Die Fortsetzung sollte fuer eine:n Beobachter:in nahtlos erfolgen.
// Im Moment tut es das nicht - man muesste sich bei den Methoden und Ideen
// von AnimationEmbed bedienen.
func (a *AnimationController) Continue() {
	if a.running {
		return
	}
	a.delay += time.Since(a.stop)
	a.ticker.Reset(refreshRate)
	a.running = true
}

func (a *AnimationController) IsRunning() bool {
	return a.running
}

func (a *AnimationController) backgroundThread() {
	var wg sync.WaitGroup

	startChan := make(chan int)

	for range a.numThreads {
		go a.animationUpdater(startChan, &wg)
	}

	for pit := range a.ticker.C {
		if a.quit {
			break
		}
		a.animPit = pit.Add(-a.delay)

		a.animWatch.Start()
		wg.Add(a.numThreads)
		for id := range a.numThreads {
			startChan <- id
		}
		wg.Wait()
		a.animWatch.Stop()

		a.syncChan <- true
		<-a.syncChan
	}
	close(startChan)
}

func (a *AnimationController) animationUpdater(startChan <-chan int, wg *sync.WaitGroup) {
	for id := range startChan {
		a.animMutex.RLock()
		for i := id; i < len(a.AnimList); i += a.numThreads {
			anim := a.AnimList[i]
			if anim == nil || !anim.IsRunning() {
				continue
			}
			if !anim.Update(a.animPit) {
			}
		}
		a.animMutex.RUnlock()
		wg.Done()
	}
}

func (a *AnimationController) Save(fileName string) {
	fh, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Couldn't create file: %v", err)
	}
	defer fh.Close()
	gobEncoder := gob.NewEncoder(fh)
	err = gobEncoder.Encode(a)
	if err != nil {
		log.Fatalf("Couldn't encode data: %v", err)
	}
}

func (a *AnimationController) Load(fileName string) {
	fh, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Couldn't create file: %v", err)
	}
	defer fh.Close()
	gobDecoder := gob.NewDecoder(fh)
	err = gobDecoder.Decode(a)
	if err != nil {
		log.Fatalf("Couldn't decode data: %v", err)
	}
}

func (a *AnimationController) Watch() *Stopwatch {
	return a.animWatch
}

func (a *AnimationController) Now() time.Time {
	delay := a.delay
	if !a.running {
		delay += time.Since(a.stop)
	}
	return time.Now().Add(delay)
}

// Die folgenden Interfaces geben einen guten Ueberblick ueber die Arten
// von Hintergrundaktivitaeten.
//
// Ein Task hat nur die Moeglichkeit, gestartet zu werden. Anschl. lauft er
// asynchron ohne weitere Moeglichkeiten der Einflussnahme. Es ist sinnvoll,
// wenn der Code hinter einem Task so kurz wie moeglich gehalten wird.
type Task interface {
	Start()
}

// Ein Job besitzt alle Eigenschaften und Moeglichkeiten eines Tasks, bietet
// jedoch zusaetzlich die Moeglichkeit, von aussen bewusst gestoppt zu werden.
type Job interface {
	Task
	IsRunning() bool
	// Stop()
}

// Animationen schliesslich sind Jobs, koennen also gestoppt und gestartet
// werden, haben jedoch zusaetzlich klare Laufzeiten und werden periodisch
// durch den AnimationController aufgerufen, damit sie ihren Animations-Kram
// erledigen koennen. Man kann sie ausserdem Suspendieren und Fortsetzen.
type Animation interface {
	Job
	Suspend()
	Continue()
	Update(t time.Time) bool
}

type TimedAnimation interface {
	Animation
	Duration() time.Duration
	SetDuration(dur time.Duration)
}

// Jede konkrete Animation (Farben, Positionen, Groessen, etc.) muss das
// Interface AnimationImpl implementieren.
type NormAnimation interface {
	TimedAnimation
	// Init wird vom Animationsframework aufgerufen, wenn diese Animation
	// gestartet wird. Wichtig: Wiederholungen und Umkehrungen (AutoReverse)
	// zaehlen nicht als Start!
	Init()
	// In Tick schliesslich ist die eigentliche Animationslogik enthalten.
	// Der Parameter t gibt an, wie weit die Animation bereits gelaufen ist.
	// t=0: Animation wurde eben gestartet
	// t=1: Die Animation ist fertig
	Tick(t float64)
}

// Haben Animationen eine Dauer, so koennen sie dieses Embeddable einbinden
// und erhalten somit die klassischen Methoden fuer das Setzen und Abfragen
// der Dauer.
type DurationEmbed struct {
	duration time.Duration
	// Dummy    int
}

func (d *DurationEmbed) Duration() time.Duration {
	return d.duration
}

func (d *DurationEmbed) SetDuration(dur time.Duration) {
	d.duration = dur
}

// Animationen werden nur selten direkt dem AnimationController hinzugefuegt.
// Meistens stehen Animationen in einer Beziehung zu weiteren Animationen
// bspw. wenn eine Gruppe von Animationen gleichzeitig oder direkt
// hintereinander gestartet werden sollen. Die Typen 'Group', 'Sequence' und
// 'Timeline' sind solche Steuer-Animationen, implementieren jedoch alle
// das Animation-Interface und koennen somit ineinander verschachtelt werden.

func init() {
	gob.Register(&Group{})
	gob.Register(&Sequence{})
	gob.Register(&Timeline{})
	gob.Register(&SimpleTask{})
	gob.Register(&HideShowAnimation{})
	// gob.Register(&StartStopAnimation{})
	gob.Register(&FadeAnimation{})
	gob.Register(&ColorAnimation{})
	gob.Register(&PaletteAnimation{})
	gob.Register(&PaletteFadeAnimation{})
	gob.Register(&FloatAnimation{})
	// gob.Register(&PathAnimation{})
	gob.Register(&FixedPosAnimation{})
	gob.Register(&IntegerPosAnimation{})
	gob.Register(&ShaderAnimation{})
	gob.Register(&NormAnimationEmbed{})
}

// Mit einem Task koennen beliebige Funktionsaufrufe in die
// Animationsketten aufgenommen werden. Sie koennen beliebig oft gestartet
// werden, haben als Dauer stets 0 ms und werden immer als 'gestoppt'
// ausgewiesen. Es empfiehlt sich, nur kurze Aktionen damit zu realisieren.
type SimpleTask struct {
	fn func()
}

func NewTask(fn func()) *SimpleTask {
	a := &SimpleTask{fn}
	return a
}
func (a *SimpleTask) Start() {
	a.fn()
}

// Mit einer ShowHideAnimation kann die Eigenschaft IsHidden von
// CanvasObjekten beeinflusst werden. Mit jedem Aufruf wird dieser Switch
// umgestellt (also Hidden -> Shown -> Hidden -> etc.).
type HideShowAnimation struct {
	obj CanvasObject
}

func NewHideShowAnimation(obj CanvasObject) *HideShowAnimation {
	a := &HideShowAnimation{obj: obj}
	return a
}

func (a *HideShowAnimation) Start() {
	if a.obj.IsVisible() {
		a.obj.Hide()
	} else {
		a.obj.Show()
	}
}

func (a *HideShowAnimation) IsRunning() bool {
	return false
}

// Analog zu HideShowAnimation dient SuspContAnimation dazu, die Eigenschaft
// IsRunning von Animation-Objekten zu beeinflussen. Jeder Aufruf dieser
// Animation wechselt die Eigenschaft (Stopped -> Started -> Stopped -> etc.)
type SuspContAnimation struct {
	anim Animation
}

func NewSuspContAnimation(anim Animation) *SuspContAnimation {
	a := &SuspContAnimation{anim: anim}
	return a
}
func (a *SuspContAnimation) Start() {
	if a.anim.IsRunning() {
		a.anim.Suspend()
	} else {
		a.anim.Continue()
	}
}
func (a *SuspContAnimation) IsRunning() bool {
	return false
}

// Embeddable mit in allen Animationen benoetigen Variablen und Methoden.
// Erleichert das Erstellen von neuen Animationen gewaltig.
type NormAnimationEmbed struct {
	wrapper NormAnimation
	// Falls true, wird die Animation einmal vorwaerts und einmal rueckwerts
	// abgespielt.
	AutoReverse bool
	// Curve bezeichnet eine Interpolationsfunktion, welche einen beliebigen
	// Verlauf der Animation erlaubt (Beschleunigung am Anfang, Abbremsen
	// gegen Schluss, etc).
	Curve AnimationCurve
	// Bezeichnet die Anzahl Wiederholungen dieser Animation.
	RepeatCount int
	DurationEmbed
	reverse          bool
	start, stop, end time.Time
	total            float64
	repeatsLeft      int
	running          bool
}

// Muss beim Erstellen einer Animation aufgerufen werden, welche dieses
// Embeddable einbindet.
func (a *NormAnimationEmbed) Extend(wrapper NormAnimation) {
	a.wrapper = wrapper
	a.AutoReverse = false
	a.Curve = AnimationEaseInOut
	a.RepeatCount = 0
	a.running = false
	AnimCtrl.Add(a.wrapper)
}

func (a *NormAnimationEmbed) Duration() time.Duration {
	factor := 1
	if a.RepeatCount > 0 {
		factor += a.RepeatCount
	}
	if a.AutoReverse {
		factor *= 2
	}
	return time.Duration(factor) * a.duration
}

// Startet die Animation mit jenen Parametern, die zum Startzeitpunkt
// aktuell sind. Ist die Animaton bereits am Laufen ist diese Methode
// ein no-op.
func (a *NormAnimationEmbed) Start() {
	if a.running {
		return
	}
	a.start = AnimCtrl.Now()
	a.end = a.start.Add(a.duration)
	a.total = a.end.Sub(a.start).Seconds()
	a.repeatsLeft = a.RepeatCount
	a.reverse = false
	a.wrapper.Init()
	a.running = true
}

// Haelt die Animation an, laesst sie jedoch in der Animation-Queue der
// Applikation. Mit [Continue] kann eine gestoppte Animation wieder
// fortgesetzt werden.
func (a *NormAnimationEmbed) Suspend() {
	// fmt.Printf("Von aussen gestoppt...\n")
	if !a.running {
		return
	}
	a.stop = AnimCtrl.Now()
	a.running = false
}

// Setzt eine mit [Stop] angehaltene Animation wieder fort.
func (a *NormAnimationEmbed) Continue() {
	if a.running {
		return
	}
	dt := AnimCtrl.Now().Sub(a.stop)
	a.start = a.start.Add(dt)
	a.end = a.end.Add(dt)
	a.running = true
}

// Liefert true, falls die Animation mittels [Stop] angehalten wurde oder
// falls die Animation zu Ende ist.
func (a *NormAnimationEmbed) IsRunning() bool {
	return a.running
}

// Diese Methode ist fuer die korrekte Abwicklung (Beachtung von Reverse und
// RepeatCount, etc) einer Animation zustaendig. Wenn die Animation zu Ende
// ist, retourniert Update false. Der Parameter [t] ist ein "Point in Time".
func (a *NormAnimationEmbed) Update(t time.Time) bool {
	if t.After(a.end) {
		if a.reverse {
			a.wrapper.Tick(a.Curve(0.0))
			if a.repeatsLeft == 0 {
				a.running = false
				// fmt.Printf("erster ausstieg...\n  %+v\n", a)

				return false
			}
			a.reverse = false
		} else {
			a.wrapper.Tick(a.Curve(1.0))
			if a.AutoReverse {
				a.reverse = true
			}
		}
		if !a.reverse {
			if a.repeatsLeft == 0 {
				a.running = false
				// fmt.Printf("zweiter ausstieg...\n  %+v\n", a)
				return false
			}
			if a.repeatsLeft > 0 {
				a.repeatsLeft--
			}
		}
		a.start = a.end
		a.end = a.start.Add(a.duration)
		a.wrapper.Init()
		return true
	}

	delta := t.Sub(a.start).Seconds()
	val := delta / a.total
	if a.reverse {
		a.wrapper.Tick(a.Curve(1.0 - val))
	} else {
		a.wrapper.Tick(a.Curve(val))
	}
	return true
}

// ---------------------------------------------------------------------------

type AnimValue interface {
	~float64 | ~uint8 | ~int | geom.Point | image.Point | fixed.Point26_6 | color.LedColor
}

type AnimValueFunc[T AnimValue] func() T

func Const[T AnimValue](v T) AnimValueFunc[T] {
	return func() T { return v }
}

type GenericAnimation[T AnimValue] struct {
	NormAnimationEmbed
	ValPtr     *T
	val1, val2 T
	Val1, Val2 AnimValueFunc[T]
	Cont       bool
}

func (a *GenericAnimation[T]) InitAnim(valPtr *T, val2 T, dur time.Duration) {
	a.ValPtr = valPtr
	a.Val1 = Const(*valPtr)
	a.Val2 = Const(val2)
	a.SetDuration(dur)
}

func (a *GenericAnimation[T]) Init() {
	if a.Cont {
		a.Val1 = Const(*a.ValPtr)
	}
	a.val1 = a.Val1()
	a.val2 = a.Val2()
}

// ---------------------------------------------------------------------------

type FadeType int

const (
	FadeOut FadeType = 0x00
	FadeIn           = 0xff
)

type FadeAnimation struct {
	GenericAnimation[uint8]
}

func NewFadeAnimation(valPtr *uint8, fade FadeType, dur time.Duration) *FadeAnimation {
	a := &FadeAnimation{}
	a.InitAnim(valPtr, uint8(fade), dur)
	a.NormAnimationEmbed.Extend(a)
	return a
}

func (a *FadeAnimation) Tick(t float64) {
	*a.ValPtr = uint8((1.0-t)*float64(a.val1) + t*float64(a.val2))
}

// ---------------------------------------------------------------------------

type IntAnimation struct {
	GenericAnimation[int]
}

func NewIntAnimation(valPtr *int, val2 int, dur time.Duration) *IntAnimation {
	a := &IntAnimation{}
	a.InitAnim(valPtr, val2, dur)
	a.NormAnimationEmbed.Extend(a)
	return a
}

func (a *IntAnimation) Tick(t float64) {
	*a.ValPtr = int(math.Round((1.0-t)*float64(a.val1) + t*float64(a.val2)))
}

// Animation fuer einen Verlauf zwischen zwei Fliesskommazahlen.
type FloatAnimation struct {
	GenericAnimation[float64]
}

func NewFloatAnimation(valPtr *float64, val2 float64, dur time.Duration) *FloatAnimation {
	a := &FloatAnimation{}
	a.InitAnim(valPtr, val2, dur)
	a.NormAnimationEmbed.Extend(a)
	return a
}

func (a *FloatAnimation) Tick(t float64) {
	*a.ValPtr = (1-t)*a.val1 + t*a.val2
}

// Animation fuer einen Verlauf zwischen zwei Farben.
type ColorAnimation struct {
	GenericAnimation[color.LedColor]
}

func NewColorAnimation(valPtr *color.LedColor, val2 color.LedColor, dur time.Duration) *ColorAnimation {
	a := &ColorAnimation{}
	a.InitAnim(valPtr, val2, dur)
	a.NormAnimationEmbed.Extend(a)
	return a
}

func (a *ColorAnimation) Tick(t float64) {
	alpha := (*a.ValPtr).A
	*a.ValPtr = a.val1.Interpolate(a.val2, t)
	(*a.ValPtr).A = alpha
}


// Animation fuer das Fahren entlang eines Pfades. Mit fnc kann eine konkrete,
// Pfad-generierende Funktion angegeben werden. Siehe auch [PathFunction]
type PathAnimation struct {
	GenericAnimation[geom.Point]
	Path Path
}

func NewPathAnimation(valPtr *geom.Point, path *GeomPath, size geom.Point, dur time.Duration) *PathAnimation {
	a := &PathAnimation{}
	a.InitAnim(valPtr, geom.Point{}, dur)
	a.NormAnimationEmbed.Extend(a)
	a.Path = path
	a.Val2 = func() geom.Point {
		return a.Val1().Add(size)
	}
	return a
}

func NewPolyPathAnimation(valPtr *geom.Point, path *PolygonPath, dur time.Duration) *PathAnimation {
	a := &PathAnimation{}
	a.InitAnim(valPtr, (*valPtr).AddXY(1, 1), dur)
	a.NormAnimationEmbed.Extend(a)
	a.Path = path
	return a
}

func NewPositionAnimation(valPtr *geom.Point, val2 geom.Point, dur time.Duration) *PathAnimation {
	a := &PathAnimation{}
	a.InitAnim(valPtr, val2, dur)
	a.NormAnimationEmbed.Extend(a)
	a.Path = LinearPath
	return a
}

func NewSizeAnimation(valPtr *geom.Point, val2 geom.Point, dur time.Duration) *PathAnimation {
	return NewPositionAnimation(valPtr, val2, dur)
}

func (a *PathAnimation) Tick(t float64) {
	var dp geom.Point
	var s geom.Point

	dp = a.Path.Pos(t)
	s = a.val2.Sub(a.val1)
	dp.X *= s.X
	dp.Y *= s.Y
	*a.ValPtr = a.val1.Add(dp)
}

// Animation fuer eine Positionsveraenderung anhand des Fixed-Datentyps
// [fixed/Point26_6]. Dies wird insbesondere für die Positionierung von
// Schriften verwendet.
type FixedPosAnimation struct {
    GenericAnimation[fixed.Point26_6]
}

func NewFixedPosAnimation(valPtr *fixed.Point26_6, val2 fixed.Point26_6, dur time.Duration) *FixedPosAnimation {
	a := &FixedPosAnimation{}
    a.InitAnim(valPtr, val2, dur)
	a.NormAnimationEmbed.Extend(a)
	return a
}

func (a *FixedPosAnimation) Tick(t float64) {
	*a.ValPtr = a.val1.Mul(float2fix(1.0 - t)).Add(a.val2.Mul(float2fix(t)))
}

func float2fix(x float64) fixed.Int26_6 {
	return fixed.Int26_6(math.Round(x * 64))
}

type IntegerPosAnimation struct {
    GenericAnimation[image.Point]
}

func NewIntegerPosAnimation(valPtr *image.Point, val2 image.Point, dur time.Duration) *IntegerPosAnimation {
	a := &IntegerPosAnimation{}
    a.InitAnim(valPtr, val2, dur)
	a.NormAnimationEmbed.Extend(a)
	return a
}

func (a *IntegerPosAnimation) Tick(t float64) {
	v1 := geom.NewPointIMG(a.val1)
	v2 := geom.NewPointIMG(a.val2)
	np := v1.Mul(1.0 - t).Add(v2.Mul(t))
	*a.ValPtr = np.Int()
}

// Animation fuer einen Farbverlauf ueber die Farben einer Palette.
type PaletteAnimation struct {
	NormAnimationEmbed
	ValPtr *color.LedColor
	pal    ColorSource
}

func NewPaletteAnimation(valPtr *color.LedColor, pal ColorSource, dur time.Duration) *PaletteAnimation {
	a := &PaletteAnimation{}
	a.NormAnimationEmbed.Extend(a)
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

// Dies schliesslich ist eine Animation, bei welcher stufenlos von einer
// Palette auf eine andere umgestellt wird.
type PaletteFadeAnimation struct {
	NormAnimationEmbed
	Fader   *PaletteFader
	Val2    ColorSource
	ValFunc PaletteFuncType
}

func NewPaletteFadeAnimation(fader *PaletteFader, pal2 ColorSource, dur time.Duration) *PaletteFadeAnimation {
	a := &PaletteFadeAnimation{}
	a.NormAnimationEmbed.Extend(a)
	a.SetDuration(dur)
	a.Fader = fader
	a.Val2 = pal2
	return a
}

func (a *PaletteFadeAnimation) Init() {
	a.Fader.T = 0.0
	if a.ValFunc != nil {
		a.Fader.Pals[1] = a.ValFunc()
	} else {
		a.Fader.Pals[1] = a.Val2
	}
}

func (a *PaletteFadeAnimation) Tick(t float64) {
	if t == 1.0 {
		a.Fader.Pals[0], a.Fader.Pals[1] = a.Fader.Pals[1], a.Fader.Pals[0]
		a.Fader.T = 0.0
	} else {
		a.Fader.T = t
	}
}

// Fuer den klassischen Shader wird pro Pixel folgende Animation gestartet.
type ShaderFuncType func(x, y, t float64) float64

type ShaderAnimation struct {
	ValPtr      *color.LedColor
	Pal         ColorSource
	X, Y        float64
	Fnc         ShaderFuncType
	start, stop time.Time
	running     bool
}

func NewShaderAnimation(valPtr *color.LedColor, pal ColorSource,
	x, y float64, fnc ShaderFuncType) *ShaderAnimation {
	a := &ShaderAnimation{}
	a.ValPtr = valPtr
	a.Pal = pal
	a.X, a.Y = x, y
	a.Fnc = fnc
	AnimCtrl.Add(a)
	return a
}

func (a *ShaderAnimation) Duration() time.Duration {
	return time.Duration(0)
}

func (a *ShaderAnimation) SetDuration(d time.Duration) {}

// Startet die Animation.
func (a *ShaderAnimation) Start() {
	if a.running {
		return
	}
	a.start = AnimCtrl.Now()
	a.running = true
}

// Unterbricht die Ausfuehrung der Animation.
func (a *ShaderAnimation) Suspend() {
	if !a.running {
		return
	}
	a.stop = AnimCtrl.Now()
	a.running = false
}

// Setzt die Ausfuehrung der Animation fort.
func (a *ShaderAnimation) Continue() {
	if a.running {
		return
	}
	dt := AnimCtrl.Now().Sub(a.stop)
	a.start = a.start.Add(dt)
	a.running = true
}

// Liefert den Status der Animation zurueck.
func (a *ShaderAnimation) IsRunning() bool {
	return a.running
}

func (a *ShaderAnimation) Update(t time.Time) bool {
	*a.ValPtr = a.Pal.Color(a.Fnc(a.X, a.Y, t.Sub(a.start).Seconds()))
	return true
}
