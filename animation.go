//go:build !tinygo

package ledgrid

import (
	"image"
	"math"
	"math/rand/v2"
	"sync"
	"time"

	"golang.org/x/image/math/fixed"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid/colors"
)

var (
	// Weil es nur einen (1) Animations-Kontroller geben kann, ist in diesem
	// Package eine globale Variable dafuer vorgesehen.
	AnimCtrl *AnimationController
)

const (
	// Fuer Animationen, die endlos wiederholt weren sollen, kann diese Konstante
	// fuer die Anzahl Wiederholungen verwendet werden.
	AnimationRepeatForever = -1
	// Mit refreshRate wird die Zeit in Millisekunden angegeben, die zwischen
	// den einzelnen Aktualisierungen gewartet wird.
	refreshRate = 30 * time.Millisecond
)

// Das Herzstueck der ganzen Animationen ist der AnimationController. Er
// sorgt dafuer, dass alle 30 ms (siehe Variable refreshRate) die
// Update-Methoden aller Animationen aufgerufen werden und veranlasst im
// Anschluss, dass alle darstellbaren Objekte neu gezeichnet werden und sendet
// das Bild schliesslich dem PixelController (oder dem PixelEmulator).
type AnimationController struct {
	AnimList   []Animation
	animMutex  *sync.RWMutex
	ticker     *time.Ticker
	quit       bool
	animPit    time.Time
	stopwatch  *Stopwatch
	numThreads int
	stop       time.Time
	delay      time.Duration
	isRunning  bool
	syncChan   chan bool
}

func NewAnimationController(syncChan chan bool) *AnimationController {
	if AnimCtrl != nil {
		return AnimCtrl
	}
	a := &AnimationController{}
	a.AnimList = make([]Animation, 0)
	a.animMutex = &sync.RWMutex{}
	a.ticker = time.NewTicker(refreshRate)
	a.stopwatch = NewStopwatch()
	a.numThreads = 1
	a.delay = time.Duration(0)
	a.syncChan = syncChan

	AnimCtrl = a
	go a.backgroundThread()
	a.isRunning = true

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

// Mit Suspend koennen die Animationen und die Darstellung auf der Hardware
// unterbunden werden.
func (a *AnimationController) Suspend() {
	if !a.isRunning {
		return
	}
	a.isRunning = false
	a.ticker.Stop()
	a.stop = time.Now()
}

// Setzt die Animationen wieder fort.
// TO DO: Die Fortsetzung sollte fuer eine:n Beobachter:in nahtlos erfolgen.
// Im Moment tut es das nicht - man muesste sich bei den Methoden und Ideen
// von AnimationEmbed bedienen.
func (a *AnimationController) Continue() {
	if a.isRunning {
		return
	}
	a.delay += time.Since(a.stop)
	a.ticker.Reset(refreshRate)
	a.isRunning = true
}

func (a *AnimationController) IsRunning() bool {
	return a.isRunning
}

// Diese Funktion wird im Hintergrund gestartet und ist fuer die Koordination
// von Animation, Zeichnen und Sender der Daten zur Hardware verantwortlich.
// TO DO: in Zukunft sollte die Koordination durch das Objekt LedGrid
// vorgenommen werden. Davon ist in jedem Programm nur eine Instanz vorhanden
// AnimationController's koennte es grundsaetzlich mehrere geben.
func (a *AnimationController) backgroundThread() {
	for pit := range a.ticker.C {
		if a.quit {
			break
		}
		a.animPit = pit.Add(-a.delay)

		a.stopwatch.Start()
		a.animationUpdater(0)
		a.stopwatch.Stop()

		a.syncChan <- true
		<-a.syncChan
	}
}

// Von dieser Funktion werden pro Core eine Go-Routine gestartet. Sie werden
// durch eine Message ueber den Kanal statChan aktiviert und uebernehmen die
// Aktualisierung ihrer Animationen.
func (a *AnimationController) animationUpdater(id int) {
	a.animMutex.RLock()
	for id, anim := range a.AnimList {
		if anim == nil || !anim.IsRunning() {
			continue
		}
		a.animMutex.RUnlock()
		if !anim.Update(a.animPit) {
			a.AnimList[id] = nil
		}
		a.animMutex.RLock()
	}
	a.animMutex.RUnlock()
}

func (a *AnimationController) Stopwatch() *Stopwatch {
	return a.stopwatch
}

func (a *AnimationController) Now() time.Time {
	return a.animPit
}

// Mit dem Funktionstyp [AnimationCurve] kann der Verlauf einer Animation
// beeinflusst werden. Der Parameter [t] ist ein Wert im Intervall [0,1]
// und zeigt an, wo sich die Animation gerade befindet (t=0: Animation
// wurde eben gestartet; t=1: Animation ist zu Ende). Der Rueckgabewert
// ist ebenfalls ein Wert im Intervall [0,1] und hat die gleiche Bedeutung
// wie [t].
type AnimationCurve func(t float64) float64

// Linearer Verlauf zwischen 0.0 und 1.0.
func AnimationLinear(t float64) float64 {
	return t
}

// Die folgenden Funktionen dienen als Basis fuer eine ganze Schar von
// Animationskurven. Es sind im Grunde parametrisierbare Parabeln; mit dem
// Parameter (einem Fliesskommawert in [0,8]) wird der Exponent bestimmt.
// Man kann also Parabeln 1.25-ter oder 0.65-ter Ordnung erstellen.
// Die Namen (...In, ...Out, ...InOut) bezeichnen, an welcher Stelle der
// Animation (Anfang, Ende, Anfang-und-Ende) eine Verzoegerung erzeugt wird.
type genericCurve func(t, a float64) float64

func genericIn(t, a float64) float64 {
	return math.Pow(t, a)
}

func genericOut(t, a float64) float64 {
	return -math.Pow(1-t, a) + 1
}

func genericInOut(t, a float64) float64 {
	if t < 0.5 {
		return genericIn(2*t, a) / 2
	} else {
		return (genericOut(2*t-1, a) + 1) / 2
	}
}

// Mit dieser Methode schliesslich, werden konkrete Animationskurven erstellt.
// Der erste Parameter bezeichnet eine der generischen Funktionen und mit
// dem Parameter a wird der gewuenschte Exponent fix hinterlegt. Das Resultat
// ist eine Funktion des Typs [AnimationCurve], der im Feld [Curve] einer
// Animation hinterlegt werden kann.
func NewAnimationCurve(fnc genericCurve, a float64) AnimationCurve {
	return func(t float64) float64 {
		return fnc(t, a)
	}
}

var (
	// Die Animationskurven, welche auf quadratischen und kubischen Parabeln
	// basieren, sind bereits vorgefertigt.
	AnimationEaseIn     = NewAnimationCurve(genericIn, 2.0)
	AnimationEaseOut    = NewAnimationCurve(genericOut, 2.0)
	AnimationEaseInOut  = NewAnimationCurve(genericInOut, 2.0)
	AnimationCubicIn    = NewAnimationCurve(genericIn, 3.0)
	AnimationCubicOut   = NewAnimationCurve(genericOut, 3.0)
	AnimationCubicInOut = NewAnimationCurve(genericInOut, 3.0)
)

// Dies ist ein etwas unbeholfener Versuch, die Zielwerte bestimmter
// Animationen dynamisch berechnen zu lassen. Alle XXXFuncType sind
// Funktionstypen fuer einen bestimmten Datentyp, der in den Animationen
// verwendet wird und dessen dynamische Berechnung Sinn macht.
type PaletteFuncType func() ColorSource

// SeqPalette liefert eine Funktion, die bei jedem Aufruf die naechste Palette
// als Resultat zurueckgibt.
func SeqPalette() PaletteFuncType {
	var palId int = 0
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

var (
	randColor, randGroupColor colors.LedColor
)

func RandColor(new bool) AnimValueFunc[colors.LedColor] {
	return func() colors.LedColor {
		if new {
			randColor = colors.RandColor()
		}
		return randColor
	}
}

func RandGroupColor(group colors.ColorGroup, new bool) AnimValueFunc[colors.LedColor] {
	return func() colors.LedColor {
		if new {
			randGroupColor = colors.RandGroupColor(group)
		}
		return randGroupColor
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

// Die folgenden Interfaces geben einen guten Ueberblick ueber die Arten
// von Hintergrundaktivitaeten.
//
// Ein Task hat nur die Moeglichkeit, gestartet zu werden. Anschl. lauft er
// asynchron ohne weitere Moeglichkeiten der Einflussnahme. Es ist sinnvoll,
// wenn der Code hinter einem Task so kurz wie moeglich gehalten wird.
type Task interface {
	Start()
	StartAt(t time.Time)
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
// Interface NormAnimation implementieren.
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
}

func (d *DurationEmbed) Duration() time.Duration {
	return d.duration
}

func (d *DurationEmbed) SetDuration(dur time.Duration) {
	d.duration = dur
}

// Mit einem Task koennen beliebige Funktionsaufrufe in die
// Animationsketten aufgenommen werden. Sie koennen beliebig oft gestartet
// werden. Es empfiehlt sich, nur kurze Aktionen damit zu realisieren
// (bspw. Setzen von Variablen)
type SimpleTask struct {
	fn func()
}

func NewTask(fn func()) *SimpleTask {
	a := &SimpleTask{fn}
	return a
}
func (a *SimpleTask) StartAt(t time.Time) {
	a.fn()
}
func (a *SimpleTask) Start() {
	a.StartAt(AnimCtrl.Now())
}

// Dieses Embeddable wird von allen Animationen verwendet, welche eine
// Animation implementieren, die folgende Kriterien aufweist:
//   - sie hat eine begrenzte Laufzeit (ohne Beruecksichtiung von Umkehrungen
//     und Wiederholungen!)
type NormAnimationEmbed struct {
	// Falls true, wird die Animation einmal vorwaerts und einmal rueckwerts
	// abgespielt.
	AutoReverse bool
	// Curve bezeichnet eine Interpolationsfunktion, welche einen beliebigen
	// Verlauf der Animation erlaubt (Beschleunigung am Anfang, Abbremsen
	// gegen Schluss, etc).
	Curve AnimationCurve
	// Bezeichnet die Anzahl Wiederholungen dieser Animation.
	RepeatCount int
	// Ueber dieses Embeddable werden Variablen und Methdoden zum Setzen und
	// Abfragen der Laufzeit importiert.
	DurationEmbed
	// Falls die Animation nicht am Anfang starten soll, sondern zu einem
	// beliebigen Punkt, setzt man dieses Feld auf einen Wert zwischen 0 und
	// 1.
	Pos float64

	wrapper          NormAnimation
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
	a.Curve = AnimationEaseInOut
}

// Mit Duration wird die gesamte Laufzeit der Animation (als inkl. Umkehrungen
// und Wiederholungen) ermittelt. Falls die Anzahl Wiederholgungen auf -1
// (d.h. Endloswiederholgung) gesetzt ist, liefert Duration die Laufzeit eines
// Animationszyklus.
func (a *NormAnimationEmbed) Duration() time.Duration {
	factor := 1
	startDiff := time.Duration(0)
	if a.RepeatCount > 0 {
		factor += a.RepeatCount
	}
	if a.AutoReverse {
		factor *= 2
	}
	if a.Pos > 0.0 {
		if a.AutoReverse {
			startDiff = time.Duration(a.Pos * 2.0 * float64(a.duration))
		} else {
			startDiff = time.Duration(a.Pos * float64(a.duration))
		}
	}

	return time.Duration(factor)*a.duration - startDiff
}

// Retourniert einige Werte, welche das Timing der Animation betreffen. Wird
// wohl eher fuer Debugging verwendet. Ausserdem ist der Zugriff nicht
// synchronisiert! Passt man nicht auf, kriegt man inkonsistente Angaben.
func (a *NormAnimationEmbed) TimeInfo() (start, end time.Time, total float64) {
	return a.start, a.end, a.total
}

// Startet die Animation mit jenen Parametern, die zum Startzeitpunkt
// aktuell sind. Ist die Animaton bereits am Laufen ist diese Methode
// ein no-op.
func (a *NormAnimationEmbed) StartAt(t time.Time) {
	if a.running {
		return
	}
	a.start = t
	a.reverse = false
	if a.Pos > 0.0 {
		if a.AutoReverse {
			a.Pos *= 2.0
			if a.Pos >= 1.0 {
				a.reverse = true
				a.Pos -= 1.0
			}
		}
		a.start = a.start.Add(-time.Duration(a.Pos * float64(a.duration)))
		a.Pos = 0.0
	}
	a.end = a.start.Add(a.duration)
	a.total = a.end.Sub(a.start).Seconds()
	a.repeatsLeft = a.RepeatCount
	a.wrapper.Init()
	a.running = true
	AnimCtrl.Add(a.wrapper)
}

func (a *NormAnimationEmbed) Start() {
	a.StartAt(AnimCtrl.Now())
}

// Haelt die Animation an, laesst sie jedoch in der Animation-Queue der
// Applikation. Mit [Continue] kann eine gestoppte Animation wieder
// fortgesetzt werden.
func (a *NormAnimationEmbed) Suspend() {
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
// ist, retourniert Update false. Der Parameter t ist ein fortlaufender
// "Point in Time", der fuer das gesamte Animationsframework konsistent
// ermittelt wird. Nur wenn der AnimationController angehalten wird, stoppt
// auch diese Zeitbasis.
func (a *NormAnimationEmbed) Update(t time.Time) bool {
	if t.After(a.end) {
		if a.reverse {
			a.wrapper.Tick(a.Curve(0.0))
			if a.repeatsLeft == 0 {
				a.running = false
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
		// a.wrapper.Init()
	} else {
		delta := t.Sub(a.start).Seconds()
		val := delta / a.total
		if a.reverse {
			a.wrapper.Tick(a.Curve(1.0 - val))
		} else {
			a.wrapper.Tick(a.Curve(val))
		}
	}
	return true
}

// Ein Delay kann als ganz normale Animation ueberall dort eingefuegt werden,
// wo ein unmittelbares Fortfahren der Animationen nicht gewuenscht ist und
// wo der Einsatz der Timeline zu aufwaendig ist.
type Delay struct {
	NormAnimationEmbed
}

func NewDelay(d time.Duration) *Delay {
	a := &Delay{}
	a.SetDuration(d)
	a.NormAnimationEmbed.Extend(a)
	return a
}

func (a *Delay) Init()          {}
func (a *Delay) Tick(t float64) {}

// ---------------------------------------------------------------------------

// Folgende Datentypen lassen sich 'animieren', d.h. ueber einen bestimmten
// Zeitraum von einem Wert A zu einem zweiten Wert B interpolieren. Mit
// dieser Sammlung von Typen und Funktionen habe ich wohl mein Gesellenstueck
// in Go geschrieben. Als besonderes Schmankerl ist das ganze auch noch
// generisch umgesetzt.

// Dies sind die animierbaren skalaren Datentypen.
type AnimNumbers interface {
	~float64 | ~uint8 | ~int
}
// AnimPoints sind die animierbaren 2D-Vektoren von denen es 3 Typen gibt.
type AnimPoints interface {
	geom.Point | image.Point | fixed.Point26_6
}
// AnimColors enthaelt bloss einen Datentyp fuer animierbare Farben.
type AnimColors interface {
	colors.LedColor
}
// AnimValue schliesslich ist der Zusammenschluss aller animierbaren Typen.
type AnimValue interface {
	AnimNumbers | AnimPoints | AnimColors
}
// Die eigentliche Arbeit erbringen Funktionen vom Typ AnimValueFunc, welche
// ohne Parameter aufgerufen werden und als Rueckgabewert einen Wert des
// Typs AnimValue haben.
type AnimValueFunc[T AnimValue] func() T

// Nicht alles muss immer animiert werden, soll sich aber dem Interface
// moeglichst anpassen. Dazu kann mit der Funktion Const eine AnimValueFunc
// erzeugt werden, die bei jedem Aufruf den selben Wert zurueckliefert.
func Const[T AnimValue](v T) AnimValueFunc[T] {
	return func() T { return v }
}

// Damit der Code zu den Animationen so knapp und uebersichtlich wie moeglich
// wird, ist ein Grossteil generisch implementiert. GenericAnimation wird als
// Grund-Typ von (nahezu) allen Animationen verwendet. Die Instanziierung
// erfolgt ueber einen der Datenypen, welche in AnimValue zusammengefasst sind.
type GenericAnimation[T AnimValue] struct {
	// NormAnimationEmbed wird eingebunden, um Methoden wie Start(), Suspend(),
	// Continue() nur einmal implementieren zu muessen.
	NormAnimationEmbed
	// ValPtr zeigt auf jenes Feld eines Objektes, welches durch diese
    // Animation veraendert werden soll. Die Animation kennt also den
    // Grund-Typ (Rectangle, Circle, Pixel, etc) nicht, sondern nur die zu
    // animierende Eigenschaft.
	ValPtr *T
	// Val1 und Val2 enthalten Funktionen, mit welchen der Start-, resp. End-
	// wert einer Animation zum Startzeitpunkt ermittelt werden. Damit lassen
	// sich die Animationen sehr dynamisch gestalten.
	Val1, Val2 AnimValueFunc[T]
	// Ist Cont (Continue) auf true, dann wird val1 zum Startzeitpunkt mit
	// dem aktuellen Wert hinter ValPtr belegt und nicht mit dem Funktionswert
	// aus Val1.
	Cont       bool
	val1, val2 T
}

func (a *GenericAnimation[T]) InitAnim(valPtr *T, val2 T, dur time.Duration) {
	a.ValPtr = valPtr
	a.Val1 = Const(*valPtr)
	a.Val2 = Const(val2)
	a.Cont = true
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

type Fadable interface {
	AlphaPtr() *uint8
}

type FadeEmbed struct {
	colPtr *colors.LedColor
}

func (e *FadeEmbed) Init(c *colors.LedColor) {
	e.colPtr = c
}
func (e *FadeEmbed) AlphaPtr() *uint8 {
	return &e.colPtr.A
}

type FadeAnimation struct {
	GenericAnimation[uint8]
}

func NewFadeAnim(obj Fadable, fade FadeType, dur time.Duration) *FadeAnimation {
	a := &FadeAnimation{}
	a.InitAnim(obj.AlphaPtr(), uint8(fade), dur)
	a.NormAnimationEmbed.Extend(a)
	return a
}

func (a *FadeAnimation) Tick(t float64) {
	*a.ValPtr = uint8((1.0-t)*float64(a.val1) + t*float64(a.val2))
}

// type AlphaEmbed struct {
// 	Alpha uint8
// }

// func (e *AlphaEmbed) AlphaPtr() *uint8 {
// 	return &e.Alpha
// }

// type AlphaAnimation struct {
// 	GenericAnimation[uint8]
// }

// func NewAlphaAnim(obj Fadable, fade FadeType, dur time.Duration) *AlphaAnimation {
// 	a := &AlphaAnimation{}
// 	a.InitAnim(obj.AlphaPtr(), uint8(fade), dur)
// 	a.NormAnimationEmbed.Extend(a)
// 	return a
// }

// func (a *AlphaAnimation) Tick(t float64) {
// 	*a.ValPtr = uint8((1.0-t)*float64(a.val1) + t*float64(a.val2))
// }

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

type Rotateable interface {
	AnglePtr() *float64
}

type StrokeWidtheable interface {
	StrokeWidthPtr() *float64
}

// Animation fuer einen Verlauf zwischen zwei Fliesskommazahlen. Kann fuer
// verschiedene konkrete Animationen verwendet werden.
type FloatAnimation struct {
	GenericAnimation[float64]
}

func NewAngleAnim(obj Rotateable, val2 float64, dur time.Duration) *FloatAnimation {
	a := &FloatAnimation{}
	a.InitAnim(obj.AnglePtr(), val2, dur)
	a.NormAnimationEmbed.Extend(a)
	return a
}

func NewStrokeWidthAnim(obj StrokeWidtheable, val2 float64, dur time.Duration) *FloatAnimation {
	a := &FloatAnimation{}
	a.InitAnim(obj.StrokeWidthPtr(), val2, dur)
	a.NormAnimationEmbed.Extend(a)
	return a
}

func (a *FloatAnimation) Tick(t float64) {
	*a.ValPtr = (1-t)*a.val1 + t*a.val2
}

type Colorable interface {
	ColorPtr() *colors.LedColor
}

type ColorFillable interface {
	Colorable
	FillColorPtr() *colors.LedColor
}

// type ColorStrokable interface {
//     StrokeColorPtr() *colors.LedColor
// }

// Animation fuer einen Verlauf zwischen zwei Farben.
type ColorAnimation struct {
	GenericAnimation[colors.LedColor]
}

func NewColorAnim(obj Colorable, val2 colors.LedColor, dur time.Duration) *ColorAnimation {
	a := &ColorAnimation{}
	a.InitAnim(obj.ColorPtr(), val2, dur)
	a.NormAnimationEmbed.Extend(a)
	return a
}

func NewFillColorAnim(obj ColorFillable, val2 colors.LedColor, dur time.Duration) *ColorAnimation {
	a := &ColorAnimation{}
	a.InitAnim(obj.FillColorPtr(), val2, dur)
	a.NormAnimationEmbed.Extend(a)
	return a
}

func (a *ColorAnimation) Tick(t float64) {
	alpha := (*a.ValPtr).A
	*a.ValPtr = a.val1.Interpolate(a.val2, t)
	(*a.ValPtr).A = alpha
}

// Animation fuer einen Farbverlauf ueber die Farben einer Palette.
type PaletteAnimation struct {
	GenericAnimation[colors.LedColor]
	pal ColorSource
}

func NewPaletteAnim(obj Colorable, pal ColorSource, dur time.Duration) *PaletteAnimation {
	a := &PaletteAnimation{}
	a.InitAnim(obj.ColorPtr(), colors.Black, dur)
	a.NormAnimationEmbed.Extend(a)
	a.Curve = AnimationLinear
	a.pal = pal
	return a
}

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

func NewPaletteFadeAnim(fader *PaletteFader, pal2 ColorSource, dur time.Duration) *PaletteFadeAnimation {
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

// Animation fuer das Fahren entlang eines Pfades. Mit fnc kann eine konkrete,
// Pfad-generierende Funktion angegeben werden. Siehe auch [PathFunction]
type Positionable interface {
	PosPtr() *geom.Point
}

type PathAnimation struct {
	GenericAnimation[geom.Point]
	Path Path
}

func NewPathAnim(obj Positionable, path *GeomPath, size geom.Point, dur time.Duration) *PathAnimation {
	a := &PathAnimation{}
	a.InitAnim(obj.PosPtr(), geom.Point{}, dur)
	a.NormAnimationEmbed.Extend(a)
	a.Path = path
	a.Val2 = func() geom.Point {
		return a.Val1().Add(size)
	}
	return a
}

func NewPolyPathAnim(obj Positionable, path *PolygonPath, dur time.Duration) *PathAnimation {
	a := &PathAnimation{}
	a.InitAnim(obj.PosPtr(), geom.Point{}, dur)
	a.NormAnimationEmbed.Extend(a)
	a.Path = path
	a.Val2 = func() geom.Point {
		return a.Val1().AddXY(1, 1)
	}
	return a
}

func NewPositionAnim(obj Positionable, val2 geom.Point, dur time.Duration) *PathAnimation {
	a := &PathAnimation{}
	a.InitAnim(obj.PosPtr(), val2, dur)
	a.NormAnimationEmbed.Extend(a)
	a.Path = LinearPath
	return a
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

type Sizeable interface {
	SizePtr() *geom.Point
}

type SizeAnimation struct {
	GenericAnimation[geom.Point]
	Path Path
}

func NewSizeAnim(obj Sizeable, val2 geom.Point, dur time.Duration) *SizeAnimation {
	a := &SizeAnimation{}
	a.InitAnim(obj.SizePtr(), val2, dur)
	a.NormAnimationEmbed.Extend(a)
	a.Path = LinearPath
	return a
}

func (a *SizeAnimation) Tick(t float64) {
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
type FixedPositionable interface {
	PosPtr() *fixed.Point26_6
}

type FixedPosAnimation struct {
	GenericAnimation[fixed.Point26_6]
}

// func NewFixedPosAnimation(valPtr *fixed.Point26_6, val2 fixed.Point26_6, dur time.Duration) *FixedPosAnimation {
// 	a := &FixedPosAnimation{}
// 	a.InitAnim(valPtr, val2, dur)
// 	a.NormAnimationEmbed.Extend(a)
// 	return a
// }

func NewFixedPosAnim(obj FixedPositionable, val2 fixed.Point26_6, dur time.Duration) *FixedPosAnimation {
	a := &FixedPosAnimation{}
	a.InitAnim(obj.PosPtr(), val2, dur)
	a.NormAnimationEmbed.Extend(a)
	return a
}

func (a *FixedPosAnimation) Tick(t float64) {
	*a.ValPtr = a.val1.Mul(float2fix(1.0 - t)).Add(a.val2.Mul(float2fix(t)))
}

func float2fix(x float64) fixed.Int26_6 {
	return fixed.Int26_6(math.Round(x * 64))
}

type IntegerPositionable interface {
	PosPtr() *image.Point
}
type IntegerPosAnimation struct {
	GenericAnimation[image.Point]
}

func NewIntegerPosAnim(obj IntegerPositionable, val2 image.Point, dur time.Duration) *IntegerPosAnimation {
	a := &IntegerPosAnimation{}
	a.InitAnim(obj.PosPtr(), val2, dur)
	a.NormAnimationEmbed.Extend(a)
	return a
}

func (a *IntegerPosAnimation) Tick(t float64) {
	v1 := geom.NewPointIMG(a.val1)
	v2 := geom.NewPointIMG(a.val2)
	np := v1.Mul(1.0 - t).Add(v2.Mul(t))
	*a.ValPtr = np.Int()
}

// Fuer den klassischen Shader wird pro Pixel folgende Animation gestartet.
type Shader TimedAnimation

type ColorShaderFunc func(t, x, y, z float64, idx, nPix int) colors.LedColor

type ColorShaderAnim struct {
	ValPtr      *colors.LedColor
	X, Y, Z     float64
	Idx, NPix   int
	Fnc         ColorShaderFunc
	start, stop time.Time
	running     bool
}

func NewColorShaderAnim(obj Colorable, x, y, z float64, idx, nPix int, fnc ColorShaderFunc) *ColorShaderAnim {
	a := &ColorShaderAnim{}
	a.ValPtr = obj.ColorPtr()
	a.X, a.Y, a.Z = x, y, z
	a.Idx = idx
	a.NPix = nPix
	a.Fnc = fnc
	// AnimCtrl.Add(0, a)
	return a
}

func (a *ColorShaderAnim) Duration() time.Duration {
	return time.Duration(0)
}

func (a *ColorShaderAnim) SetDuration(d time.Duration) {}

// Startet die Animation.
func (a *ColorShaderAnim) StartAt(t time.Time) {
	if a.running {
		return
	}
	a.start = t
	a.running = true
	AnimCtrl.Add(a)
}

func (a *ColorShaderAnim) Start() {
	a.StartAt(AnimCtrl.Now())
}

// Unterbricht die Ausfuehrung der Animation.
func (a *ColorShaderAnim) Suspend() {
	if !a.running {
		return
	}
	a.stop = AnimCtrl.Now()
	a.running = false
}

// Setzt die Ausfuehrung der Animation fort.
func (a *ColorShaderAnim) Continue() {
	if a.running {
		return
	}
	dt := AnimCtrl.Now().Sub(a.stop)
	a.start = a.start.Add(dt)
	a.running = true
}

// Liefert den Status der Animation zurueck.
func (a *ColorShaderAnim) IsRunning() bool {
	return a.running
}

func (a *ColorShaderAnim) Update(t time.Time) bool {
	*a.ValPtr = a.Fnc(0.6*t.Sub(a.start).Seconds(), a.X, a.Y, a.Z, a.Idx, a.NPix)
	return true
}

type NormShaderFunc func(t, x, y float64) float64

type ShaderAnimation struct {
	ValPtr      *colors.LedColor
	Pal         ColorSource
	X, Y        float64
	Fnc         NormShaderFunc
	start, stop time.Time
	running     bool
}

func NewShaderAnim(obj Colorable, pal ColorSource, x, y float64,
	fnc NormShaderFunc) *ShaderAnimation {
	a := &ShaderAnimation{}
	a.ValPtr = obj.ColorPtr()
	a.Pal = pal
	a.X, a.Y = x, y
	a.Fnc = fnc
	// AnimCtrl.Add(0, a)
	return a
}

func (a *ShaderAnimation) Duration() time.Duration {
	return time.Duration(0)
}

func (a *ShaderAnimation) SetDuration(d time.Duration) {}

// Startet die Animation.
func (a *ShaderAnimation) StartAt(t time.Time) {
	if a.running {
		return
	}
	a.start = t
	a.running = true
	AnimCtrl.Add(a)
}

func (a *ShaderAnimation) Start() {
	a.StartAt(AnimCtrl.Now())
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
	*a.ValPtr = a.Pal.Color(a.Fnc(t.Sub(a.start).Seconds(), a.X, a.Y))
	return true
}
