package ledgrid

import (
	"encoding/gob"
	"image"
	"math"
	"math/rand/v2"
	"sync"
	"time"

	"golang.org/x/image/math/fixed"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid/color"
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
	AnimList   [NumLayers][]Animation
	animMutex  [NumLayers]*sync.RWMutex
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

func NewAnimationController(syncChan chan bool) *AnimationController {
	if AnimCtrl != nil {
		return AnimCtrl
	}
	a := &AnimationController{}
	for i := range NumLayers {
		a.AnimList[i] = make([]Animation, 0)
		a.animMutex[i] = &sync.RWMutex{}
	}
	a.ticker = time.NewTicker(refreshRate)
	a.animWatch = NewStopwatch()
	a.numThreads = 1
	// a.numThreads = runtime.NumCPU()
	a.delay = time.Duration(0)
	a.syncChan = syncChan

	AnimCtrl = a
	go a.backgroundThread()
	a.running = true

	return a
}

// Fuegt weitere Animationen hinzu. Der Zugriff auf den entsprechenden Slice
// wird synchronisiert, da die Bearbeitung der Animationen durch den
// Background-Thread ebenfalls relativ haeufig auf den Slice zugreift.
func (a *AnimationController) Add(layer int, anims ...Animation) {
	a.animMutex[layer].Lock()
	a.AnimList[layer] = append(a.AnimList[layer], anims...)
	a.animMutex[layer].Unlock()
}

// Loescht eine einzelne Animation.
func (a *AnimationController) Del(layer int, anim Animation) {
	a.animMutex[layer].Lock()
	defer a.animMutex[layer].Unlock()
	for idx, obj := range a.AnimList[layer] {
		if obj == anim {
			obj.Suspend()
			a.AnimList[layer][idx] = nil
			return
		}
	}
}

func (a *AnimationController) DelAt(layer int, idx int) {
	a.animMutex[layer].Lock()
	a.AnimList[layer][idx].Suspend()
	a.AnimList[layer][idx] = nil
	a.animMutex[layer].Unlock()
}

// Loescht alle Animationen.
func (a *AnimationController) Purge(layer int) {
	a.animMutex[layer].Lock()
	for _, anim := range a.AnimList[layer] {
		if anim == nil {
			continue
		}
		anim.Suspend()
	}
	a.AnimList[layer] = a.AnimList[layer][:0]
	a.animMutex[layer].Unlock()
}

func (a *AnimationController) PurgeAll() {
	for layer := range NumLayers {
		a.Purge(layer)
	}
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

// Diese Funktion wird im Hintergrund gestartet und ist fuer die Koordination
// von Animation, Zeichnen und Sender der Daten zur Hardware verantwortlich.
// TO DO: in Zukunft sollte die Koordination durch das Objekt LedGrid
// vorgenommen werden. Davon ist in jedem Programm nur eine Instanz vorhanden
// AnimationController's koennte es grundsaetzlich mehrere geben.
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

// Von dieser Funktion werden pro Core eine Go-Routine gestartet. Sie werden
// durch eine Message ueber den Kanal statChan aktiviert und uebernehmen die
// Aktualisierung ihrer Animationen.
func (a *AnimationController) animationUpdater(startChan <-chan int, wg *sync.WaitGroup) {
	for id := range startChan {
		for listIdx := range a.AnimList {
			a.animMutex[listIdx].RLock()
			for i := id; i < len(a.AnimList[listIdx]); i += a.numThreads {
				anim := a.AnimList[listIdx][i]
				if anim == nil || !anim.IsRunning() {
					continue
				}
				if !anim.Update(a.animPit) {
					// a.AnimList[0][i] = nil
				}
			}
			a.animMutex[listIdx].RUnlock()
		}
		wg.Done()
	}
}

// func (a *AnimationController) Save(fileName string) {
// 	fh, err := os.Create(fileName)
// 	if err != nil {
// 		log.Fatalf("Couldn't create file: %v", err)
// 	}
// 	defer fh.Close()
// 	gobEncoder := gob.NewEncoder(fh)
// 	err = gobEncoder.Encode(a)
// 	if err != nil {
// 		log.Fatalf("Couldn't encode data: %v", err)
// 	}
// }

// func (a *AnimationController) Load(fileName string) {
// 	fh, err := os.Open(fileName)
// 	if err != nil {
// 		log.Fatalf("Couldn't create file: %v", err)
// 	}
// 	defer fh.Close()
// 	gobDecoder := gob.NewDecoder(fh)
// 	err = gobDecoder.Decode(a)
// 	if err != nil {
// 		log.Fatalf("Couldn't decode data: %v", err)
// 	}
// }

func (a *AnimationController) Watch() *Stopwatch {
	return a.animWatch
}

func (a *AnimationController) Now() time.Time {
	return a.animPit
	// delay := a.delay
	// if !a.running {
	// 	delay += time.Since(a.stop)
	// }
	// return time.Now().Add(delay)
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

// Beginnt langsam und nimmt immer mehr an Fahrt auf
// (quadratische Grundlage).
func AnimationEaseIn(t float64) float64 {
	return t * t
}

// Beginnt schnell und bremst zum Ende immer mehr ab
// (quadratische Grundlage).
func AnimationEaseOut(t float64) float64 {
	return t * (2 - t)
}

func AnimationMiddleStop(t float64) float64 {
    if t <= 0.5 {
        return -2*(t-0.5)*(t-0.5) + 0.5
    } else {
        return 2*(t-0.5)*(t-0.5) + 0.5
    }
}

// Anfang und Ende der Animation werden abgebremst
// (quadratische Grundlage, stueckweise Funktion).
func AnimationEaseInOut(t float64) float64 {
	if t <= 0.5 {
		return 2 * t * t
	}
	return (4-2*t)*t - 1
}

// Beginnt langsam und nimmt immer mehr an Fahrt auf.
// (kubische Grundlage).
func AnimationLazeIn(t float64) float64 {
	return t * t * t
}

// Beginnt langsam und nimmt immer mehr an Fahrt auf.
// (kubische Grundlage).
func AnimationLazeOut(t float64) float64 {
	return t * (t*(t-3) + 3)
}

// Anfang und Ende der Animation werden abgebremst.
// (kubische Grundlage, stueckweise Funktion).
func AnimationLazeInOut(t float64) float64 {
	if t <= 0.5 {
		return 4 * t * t * t
	}
	return 4*(t-1)*(t-1)*(t-1) + 1
}

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
	randColor, randGroupColor color.LedColor
)

func RandColor(new bool) AnimValueFunc[color.LedColor] {
	return func() color.LedColor {
		if new {
			randColor = color.RandColor()
		}
		return randColor
	}
}

func RandGroupColor(group color.ColorGroup, new bool) AnimValueFunc[color.LedColor] {
	return func() color.LedColor {
		if new {
			randGroupColor = color.RandGroupColor(group)
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

func (a *HideShowAnimation) StartAt(t time.Time) {
	if a.obj.IsVisible() {
		a.obj.Hide()
	} else {
		a.obj.Show()
	}
}
func (a *HideShowAnimation) Start() {
	a.StartAt(AnimCtrl.Now())
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
func (a *SuspContAnimation) StartAt(t time.Time) {
	if a.anim.IsRunning() {
		a.anim.Suspend()
	} else {
		a.anim.Continue()
	}
}
func (a *SuspContAnimation) Start() {
	a.StartAt(AnimCtrl.Now())
}

func (a *SuspContAnimation) IsRunning() bool {
	return false
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
	AnimCtrl.Add(0, a.wrapper)
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

// func (a *NormAnimationEmbed) SetAnimPos(r float64) {
// 	a.startPos = r
// }

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

// ---------------------------------------------------------------------------

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

// Folgende Datentypen lassen sich 'animieren', d.h. ueber einen bestimmten
// Zeitraum von einem Wert A zu einem zweiten Wert B interpolieren.
type AnimNumbers interface {
	~float64 | ~uint8 | ~int
}
type AnimPoints interface {
	geom.Point | image.Point | fixed.Point26_6
}
type AnimColors interface {
	color.LedColor
}
type AnimValue interface {
	AnimNumbers | AnimPoints | AnimColors
}

type AnimValueFunc[T AnimValue] func() T

func Const[T AnimValue](v T) AnimValueFunc[T] {
	return func() T { return v }
}

// Damit der Code zu den Animationen so knapp und uebersichtlich wie moeglich
// wird, ist ein Grossteil generisch implementiert. GenericAnimation wird als
// Grund-Typ von (nahezu) allen Animationen verwendet. Die Instanziierung
// erfolgt ueber einen der Datenypen, welche in AnimValue zusammengefasst sind.
type GenericAnimation[T AnimValue] struct {
	// NormAnimationEmbed wird eingebunden, um Methoden wie Start(), Suspend(),
	// Continue() nur einmal implementiern zu muessen.
	NormAnimationEmbed
	// ValPtr zeigt auf jenes Feld eines Record, welches durch diese Animation
	// veraendert werden soll. Die Animation kenn also den Grund-Typ (Rectangle,
	// Circle, Pixel, etc) nicht, sondern nur den zu animierenden Wert.
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
	colPtr *color.LedColor
}

func (e *FadeEmbed) Init(c *color.LedColor) {
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
	ColorPtr() *color.LedColor
}

type ColorFillable interface {
	Colorable
	FillColorPtr() *color.LedColor
}

// type ColorStrokable interface {
//     StrokeColorPtr() *color.LedColor
// }

// Animation fuer einen Verlauf zwischen zwei Farben.
type ColorAnimation struct {
	GenericAnimation[color.LedColor]
}

func NewColorAnim(obj Colorable, val2 color.LedColor, dur time.Duration) *ColorAnimation {
	a := &ColorAnimation{}
	a.InitAnim(obj.ColorPtr(), val2, dur)
	a.NormAnimationEmbed.Extend(a)
	return a
}

func NewFillColorAnim(obj ColorFillable, val2 color.LedColor, dur time.Duration) *ColorAnimation {
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

// func NewIntegerPosAnimation(valPtr *image.Point, val2 image.Point, dur time.Duration) *IntegerPosAnimation {
// 	a := &IntegerPosAnimation{}
// 	a.InitAnim(valPtr, val2, dur)
// 	a.NormAnimationEmbed.Extend(a)
// 	return a
// }

func (a *IntegerPosAnimation) Tick(t float64) {
	v1 := geom.NewPointIMG(a.val1)
	v2 := geom.NewPointIMG(a.val2)
	np := v1.Mul(1.0 - t).Add(v2.Mul(t))
	*a.ValPtr = np.Int()
}

// Animation fuer einen Farbverlauf ueber die Farben einer Palette.
type PaletteAnimation struct {
	GenericAnimation[color.LedColor]
	pal ColorSource
}

func NewPaletteAnim(obj Colorable, pal ColorSource, dur time.Duration) *PaletteAnimation {
	a := &PaletteAnimation{}
	a.InitAnim(obj.ColorPtr(), color.Black, dur)
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
type Shader TimedAnimation

type ColorShaderFunc func(t, x, y, z float64, idx, nPix int) color.LedColor

type ColorShaderAnim struct {
	ValPtr      *color.LedColor
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
	AnimCtrl.Add(0, a)
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
	ValPtr      *color.LedColor
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
	AnimCtrl.Add(0, a)
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
