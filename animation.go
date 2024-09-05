package ledgrid

import (
	"encoding/gob"
	"image"
	"image/draw"
	"log"
	"math"
	"math/rand/v2"
	"os"
	"runtime"
	"slices"
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
// Animationen dynamisch berechnen zu lassen. Alle XXXFuncType sind
// Funktionstypen fuer einen bestimmten Datentyp, der in den Animationen
// verwendet wird und dessen dynamische Berechnung Sinn macht.
type PaletteFuncType func() ColorSource
type AlphaFuncType func() uint8

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
func RandPoint(r geom.Rectangle) PointFuncType {
	return func() geom.Point {
		fx := rand.Float64()
		fy := rand.Float64()
		return r.RelPos(fx, fy)
	}
}

// Wie RandPoint, sorgt jedoch dafuer dass die Koordinaten auf ein Vielfaches
// von t abgeschnitten werden.
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

// Macht eine Interpolation zwischen den Groessen s1 und s2. Der Interpolations-
// punkt wird zufaellig gewaehlt.
func RandSize(s1, s2 geom.Point) PointFuncType {
	return func() geom.Point {
		t := rand.Float64()
		return s1.Interpolate(s2, t)
	}
}

// Liefert eine zufaellig gewaehlte Fliesskommazahl im Interval [a,b).
func RandFloat(a, b float64) FloatFuncType {
	return func() float64 {
		return a + (b-a)*rand.Float64()
	}
}

// Liefert eine zufaellig gewaehlte natuerliche Zahl im Interval [a,b).
func RandAlpha(a, b uint8) AlphaFuncType {
	return func() uint8 {
		return a + uint8(rand.UintN(uint(b-a)))
	}
}

// ----------------------------------------------------------------------------

// Dieses Interface ist von allen Typen zu implementieren, welche
// Animationen ausfuehren sollen/wollen. Aktuell gibt es nur einen Typ, der
// dieses Interface implementiert (AnimationController) und es ist auch
// fraglich, ob es je weitere Typen geben wird.
type Animator interface {
	Add(anims ...Animation)
	Del(anim Animation)
	Purge()
	Suspend()
	Continue()
	IsRunning() bool
}

// Das Herzstueck der ganzen Animationen ist der AnimationController. Er
// sorgt dafuer, dass alle 30 ms (siehe Variable refreshRate) die
// Update-Methoden aller Animationen aufgerufen werden und veranlasst im
// Anschluss, dass alle darstellbaren Objekte neu gezeichnet werden und sendet
// das Bild schliesslich dem PixelController (oder dem PixelEmulator).
type AnimationController struct {
	AnimList   []Animation
	animMutex  *sync.RWMutex
	Canvas     *Canvas
	ledGrid    *LedGrid
	ticker     *time.Ticker
	quit       bool
	animPit    time.Time
	animWatch  *Stopwatch
	numThreads int
	stop       time.Time
	delay      time.Duration
	running    bool
}

func NewAnimationController(canvas *Canvas, ledGrid *LedGrid /*, pixClient GridClient*/) *AnimationController {
	if AnimCtrl != nil {
		return AnimCtrl
	}
	a := &AnimationController{}
	a.AnimList = make([]Animation, 0)
	a.animMutex = &sync.RWMutex{}
	a.Canvas = canvas
	a.ledGrid = ledGrid
	a.ticker = time.NewTicker(refreshRate)
	a.animWatch = NewStopwatch()
	a.numThreads = runtime.NumCPU()
	a.delay = time.Duration(0)

	AnimCtrl = a
	go a.backgroundThread()
	a.running = true
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
			// a.AnimList = slices.Delete(a.AnimList, idx, idx+1)
			return
		}
	}
}

func (a *AnimationController) DelAt(idx int) {
	a.animMutex.Lock()
	a.AnimList[idx].Suspend()
	a.AnimList[idx] = nil
	// a.AnimList = slices.Delete(a.AnimList, idx, idx+1)
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

		a.Canvas.Refresh()
		draw.Draw(a.ledGrid, a.ledGrid.Bounds(), a.Canvas, image.Point{}, draw.Over)
        a.ledGrid.Show()
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
				// fmt.Printf("Anim %T ist zu ende..\n", anim)
				// fmt.Printf("  %+v\n", anim)
				// a.AnimList[i] = nil
			}
		}
		a.animMutex.RUnlock()
		wg.Done()
		// doneChan <- true
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
// Ein Task hat nur die Moeglichkeit, gestartet zu werden. Anschl. lauft er
// asynchron ohne weitere Moeglichkeiten der Einflussnahme. Einzig der aktuelle
// Status ('laeuft'/'laeuft nicht') kann noch ermittelt werden.
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
	Duration() time.Duration
	SetDuration(dur time.Duration)
	Suspend()
	Continue()
	Update(t time.Time) bool
}

// Jede konkrete Animation (Farben, Positionen, Groessen, etc.) muss das
// Interface AnimationImpl implementieren.
type NormAnimation interface {
	Animation
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

// Haben Animationen eine Dauer, so koennen sie dieses Embeddable einbinden
// und erhalten somit die klassischen Methoden fuer das Setzen und Abfragen
// der Dauer.
type DurationEmbed struct {
	duration time.Duration
	Dummy    int
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
	gob.Register(&BackgroundTask{})
	gob.Register(&HideShowAnimation{})
	// gob.Register(&StartStopAnimation{})
	gob.Register(&Uint8Animation{})
	gob.Register(&ColorAnimation{})
	gob.Register(&PaletteAnimation{})
	gob.Register(&PaletteFadeAnimation{})
	gob.Register(&FloatAnimation{})
	gob.Register(&PathAnimation{})
	gob.Register(&FixedPosAnimation{})
	gob.Register(&IntegerPosAnimation{})
	gob.Register(&ShaderAnimation{})
	gob.Register(&NormAnimationEmbed{})
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

	Tasks            []Task
	start, stop, end time.Time
	repeatsLeft      int
	running          bool
}

// Erstellt eine neue Gruppe, welche die Animationen in [anims] zusammen
// startet. Per Default ist die Laufzeit der Gruppe gleich der laengsten
// Laufzeit der hinzugefuegten Animationen.
func NewGroup(tasks ...Task) *Group {
	a := &Group{}
	a.duration = 0
	a.RepeatCount = 0
	a.Add(tasks...)
	AnimCtrl.Add(a)
	return a
}

// Fuegt der Gruppe weitere Animationen hinzu.
func (a *Group) Add(tasks ...Task) {
	for _, task := range tasks {
		if anim, ok := task.(Animation); ok {
			a.duration = max(a.duration, anim.Duration())
		}
		a.Tasks = append(a.Tasks, task)
	}
}

// Startet die Gruppe.
func (a *Group) Start() {
	if a.running {
		return
	}
	// fmt.Printf("Group: starting\n")
	a.start = AnimCtrl.Now()
	a.end = a.start.Add(a.duration)
	a.repeatsLeft = a.RepeatCount
	a.running = true
	for _, task := range a.Tasks {
		task.Start()
	}
	// AnimCtrl.Add(a)
}

// Unterbricht die Ausfuehrung der Gruppe.
func (a *Group) Suspend() {
	if !a.running {
		return
	}
	a.stop = AnimCtrl.Now()
	a.running = false
}

// Setzt die Ausfuehrung der Gruppe fort.
func (a *Group) Continue() {
	if a.running {
		return
	}
	dt := AnimCtrl.Now().Sub(a.stop)
	a.start = a.start.Add(dt)
	a.end = a.end.Add(dt)
	a.running = true
}

// Liefert den Status der Gruppe zurueck.
func (a *Group) IsRunning() bool {
	return a.running
}

func (a *Group) Update(t time.Time) bool {
	for _, task := range a.Tasks {
		if job, ok := task.(Job); ok {
			if job.IsRunning() {
				return true
			}
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
		for _, task := range a.Tasks {
			task.Start()
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

	Tasks            []Task
	currTask         int
	start, stop, end time.Time
	repeatsLeft      int
	running          bool
}

// Erstellt eine neue Sequenz welche die Animationen in [anims] hintereinander
// ausfuehrt.
func NewSequence(tasks ...Task) *Sequence {
	a := &Sequence{}
	a.duration = 0
	a.RepeatCount = 0
	a.Add(tasks...)
	AnimCtrl.Add(a)
	return a
}

// Fuegt der Sequenz weitere Animationen hinzu.
func (a *Sequence) Add(tasks ...Task) {
	for _, task := range tasks {
		if anim, ok := task.(Animation); ok {
			a.duration = a.duration + anim.Duration()
		}
		a.Tasks = append(a.Tasks, task)
	}
}

// Startet die Sequenz.
func (a *Sequence) Start() {
	if a.running {
		return
	}
	a.start = AnimCtrl.Now()
	a.end = a.start.Add(a.duration)
	a.currTask = 0
	a.repeatsLeft = a.RepeatCount
	a.running = true
	a.Tasks[a.currTask].Start()
}

// Unterbricht die Ausfuehrung der Sequenz.
func (a *Sequence) Suspend() {
	if !a.running {
		return
	}
	a.stop = AnimCtrl.Now()
	a.running = false
}

// Setzt die Ausfuehrung der Sequenz fort.
func (a *Sequence) Continue() {
	if a.running {
		return
	}
	dt := AnimCtrl.Now().Sub(a.stop)
	a.start = a.start.Add(dt)
	a.end = a.end.Add(dt)
	a.running = true
}

// Liefert den Status der Sequenz zurueck.
func (a *Sequence) IsRunning() bool {
	return a.running
}

// Wird durch den Controller periodisch aufgerufen, prueft ob Animationen
// dieser Sequenz noch am Laufen sind und startet ggf. die naechste.
func (a *Sequence) Update(t time.Time) bool {
	if a.currTask < len(a.Tasks) {
		if job, ok := a.Tasks[a.currTask].(Job); ok {
			if job.IsRunning() {
				return true
			}
		}
		a.currTask++
	}
	if a.currTask >= len(a.Tasks) {
		if t.After(a.end) {
			if a.repeatsLeft == 0 {
				a.running = false
				return false
			} else if a.repeatsLeft > 0 {
				a.repeatsLeft--
			}
			a.start = a.end
			a.end = a.start.Add(a.duration)
			a.currTask = 0
			a.Tasks[a.currTask].Start()
		}
		return true
	}
	a.Tasks[a.currTask].Start()
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

	Slots            []*TimelineSlot
	nextSlot         int
	start, stop, end time.Time
	repeatsLeft      int
	running          bool
}

// Interner Typ, mit dem Ausfuehrungszeitpunkt und Animationen festgehalten
// werden koennen.
type TimelineSlot struct {
	Duration time.Duration
	Tasks    []Task
}

// Erstellt eine neue Timeline mit Ausfuehrungsdauer d. Als d kann auch Null
// angegeben werden, dann ist die Laufzeit der Timeline gleich dem groessten
// Ausfuehrungszeitpunkt der hinterlegten Animationen.
func NewTimeline(d time.Duration) *Timeline {
	a := &Timeline{}
	a.duration = d
	a.RepeatCount = 0
	a.Slots = make([]*TimelineSlot, 0)
	AnimCtrl.Add(a)
	return a
}

// Fuegt der Timeline die Animation anim hinzu mit Ausfuehrungszeitpunkt
// dt nach Start der Timeline. Im Moment muessen die Animationen noch in
// der Reihenfolge ihres Ausfuehrungszeitpunktes hinzugefuegt werden.
func (a *Timeline) Add(pit time.Duration, tasks ...Task) {
	var i int

	if pit > a.duration {
		a.duration = pit
	}

	for i = 0; i < len(a.Slots); i++ {
		pos := a.Slots[i]
		if pos.Duration == pit {
			pos.Tasks = append(pos.Tasks, tasks...)
			return
		}
		if pos.Duration > pit {
			break
		}
	}
	a.Slots = slices.Insert(a.Slots, i, &TimelineSlot{pit, tasks})
}

// Startet die Timeline.
func (a *Timeline) Start() {
	if a.running {
		return
	}
	a.start = AnimCtrl.Now()
	a.end = a.start.Add(a.duration)
	a.repeatsLeft = a.RepeatCount
	a.nextSlot = 0
	a.running = true
}

// Unterbricht die Ausfuehrung der Timeline.
func (a *Timeline) Suspend() {
	if !a.running {
		return
	}
	a.stop = AnimCtrl.Now()
	a.running = false
}

// Setzt die Ausfuehrung der Timeline fort.
func (a *Timeline) Continue() {
	if a.running {
		return
	}
	dt := AnimCtrl.Now().Sub(a.stop)
	a.start = a.start.Add(dt)
	a.end = a.end.Add(dt)
	a.running = true
}

// Retourniert den Status der Timeline.
func (a *Timeline) IsRunning() bool {
	return a.running
}

// Wird periodisch durch den Controller aufgerufen und aktualisiert die
// Timeline.
func (a *Timeline) Update(t time.Time) bool {
	if a.nextSlot >= len(a.Slots) {
		if t.After(a.end) {
			if a.repeatsLeft == 0 {
				a.running = false
				return false
			} else if a.repeatsLeft > 0 {
				a.repeatsLeft--
			}
			a.start = a.end
			a.end = a.start.Add(a.duration)
			a.nextSlot = 0
		}
		return true
	}
	slot := a.Slots[a.nextSlot]
	if t.Sub(a.start) >= slot.Duration {
		for _, task := range slot.Tasks {
			task.Start()
		}
		a.nextSlot++
	}
	return true
}

// Mit einem BackgroundTask koennen beliebige Funktionsaufrufe in die
// Animationsketten aufgenommen werden. Sie koennen beliebig oft gestartet
// werden, haben als Dauer stets 0 ms und werden immer als 'gestoppt'
// ausgewiesen. Es empfiehlt sich, nur kurze Aktionen damit zu realisieren.
type BackgroundTask struct {
	fn func()
}

func NewBackgroundTask(fn func()) *BackgroundTask {
	a := &BackgroundTask{fn}
	return a
}
func (a *BackgroundTask) Start() {
	a.fn()
}
func (a *BackgroundTask) IsRunning() bool {
	return false
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
	if a.obj.IsHidden() {
		a.obj.Show()
	} else {
		a.obj.Hide()
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
	// fmt.Printf("Gestartet...\n")
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
	// AnimCtrl.Add(a.wrapper)
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
// Im folgenden Abschnitt befinden sich nun alle Typen, welche

// Will man allerdings nur die Durchsichtigkeit (den Alpha-Wert) einer Farbe
// veraendern und kennt beispielsweise die Farbe selber gar nicht, dann ist
// die AlphaAnimation genau das Richtige.
// type AlphaAnimation struct {
// 	NormAnimationEmbed
// 	ValPtr     *uint8
// 	Val1, Val2 uint8
// 	ValFunc    AlphaFuncType
// 	Cont       bool
// }

// func NewAlphaAnimation(valPtr *uint8, val2 uint8, dur time.Duration) *AlphaAnimation {
// 	a := &AlphaAnimation{ValPtr: valPtr, Val1: *valPtr, Val2: val2}
// 	a.NormAnimationEmbed.Extend(a)
// 	a.SetDuration(dur)
// 	return a
// }

// func (a *AlphaAnimation) Init() {
// 	if a.Cont {
// 		a.Val1 = *a.ValPtr
// 	}
// 	if a.ValFunc != nil {
// 		a.Val2 = a.ValFunc()
// 	}
// }

// func (a *AlphaAnimation) Tick(t float64) {
// 	*a.ValPtr = uint8((1.0-t)*float64(a.Val1) + t*float64(a.Val2))
// }

type Uint8FuncType func() uint8
type FloatFuncType func() float64
type ColorFuncType func() color.LedColor
type PointFuncType func() geom.Point

func ConstUint8(v uint8) Uint8FuncType {
	return func() uint8 { return v }
}
func ConstFloat(v float64) FloatFuncType {
	return func() float64 { return v }
}
func ConstColor(v color.LedColor) ColorFuncType {
	return func() color.LedColor { return v }
}
func ConstPoint(v geom.Point) PointFuncType {
	return func() geom.Point { return v }
}

type Uint8Animation struct {
	NormAnimationEmbed
	ValPtr     *uint8
	val1, val2 uint8
	Val1, Val2 Uint8FuncType
	Cont       bool
}

func NewUint8Animation(valPtr *uint8, val2 uint8, dur time.Duration) *Uint8Animation {
	a := &Uint8Animation{ValPtr: valPtr, Val1: ConstUint8(*valPtr), Val2: ConstUint8(val2)}
	a.NormAnimationEmbed.Extend(a)
	a.SetDuration(dur)
	return a
}

func (a *Uint8Animation) Init() {
	if a.Cont {
		a.Val1 = ConstUint8(*a.ValPtr)
	}
	a.val1 = a.Val1()
	a.val2 = a.Val2()
}

func (a *Uint8Animation) Tick(t float64) {
	*a.ValPtr = uint8((1.0-t)*float64(a.val1) + t*float64(a.val2))
}

// Da Positionen und Groessen mit dem gleichen Objekt aus geom realisiert
// werden (geom.Point), ist die Animation einer Groesse und einer Position
// im Wesentlichen das gleiche. Die Funktion NewSizeAnimation ist als
// syntaktische Vereinfachung zu verstehen.

// Animation fuer einen Verlauf zwischen zwei Fliesskommazahlen.
type FloatAnimation struct {
	NormAnimationEmbed
	ValPtr     *float64
	val1, val2 float64
	Val1, Val2 FloatFuncType
	Cont       bool
}

func NewFloatAnimation(valPtr *float64, val2 float64, dur time.Duration) *FloatAnimation {
	a := &FloatAnimation{ValPtr: valPtr, Val1: ConstFloat(*valPtr), Val2: ConstFloat(val2)}
	a.NormAnimationEmbed.Extend(a)
	a.SetDuration(dur)
	return a
}

func (a *FloatAnimation) Init() {
	if a.Cont {
		a.Val1 = ConstFloat(*a.ValPtr)
	}
	a.val1 = a.Val1()
	a.val2 = a.Val2()
}

func (a *FloatAnimation) Tick(t float64) {
	*a.ValPtr = (1-t)*a.val1 + t*a.val2
}

// Animation fuer einen Verlauf zwischen zwei Farben.
type ColorAnimation struct {
	NormAnimationEmbed
	ValPtr     *color.LedColor
	val1, val2 color.LedColor
	Val1, Val2 ColorFuncType
	Cont       bool
}

func NewColorAnimation(valPtr *color.LedColor, val2 color.LedColor, dur time.Duration) *ColorAnimation {
	a := &ColorAnimation{ValPtr: valPtr, Val1: ConstColor(*valPtr), Val2: ConstColor(val2)}
	a.NormAnimationEmbed.Extend(a)
	a.SetDuration(dur)
	return a
}

func (a *ColorAnimation) Init() {
	if a.Cont {
		a.Val1 = ConstColor(*a.ValPtr)
	}
	a.val1 = a.Val1()
	a.val2 = a.Val2()
}

func (a *ColorAnimation) Tick(t float64) {
	alpha := (*a.ValPtr).A
	*a.ValPtr = a.val1.Interpolate(a.val2, t)
	(*a.ValPtr).A = alpha
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

// Animation fuer das Fahren entlang eines Pfades. Mit fnc kann eine konkrete,
// Pfad-generierende Funktion angegeben werden. Siehe auch [PathFunction]
type PathAnimation struct {
	NormAnimationEmbed
	ValPtr     *geom.Point
	val1, val2 geom.Point
	Val1, Val2 PointFuncType
	PathFunc   PathFunctionType
	Cont       bool
}

func NewPathAnimation(valPtr *geom.Point, pathFunc PathFunctionType, size geom.Point, dur time.Duration) *PathAnimation {
	a := &PathAnimation{ValPtr: valPtr, Val1: ConstPoint(*valPtr), PathFunc: pathFunc}
	a.NormAnimationEmbed.Extend(a)
	a.SetDuration(dur)
	a.Val2 = func() geom.Point {
		return a.Val1().Add(size)
	}
	return a
}

func NewPositionAnimation(valPtr *geom.Point, val2 geom.Point, dur time.Duration) *PathAnimation {
	a := &PathAnimation{ValPtr: valPtr, Val1: ConstPoint(*valPtr), Val2: ConstPoint(val2), PathFunc: LinearPath}
	a.NormAnimationEmbed.Extend(a)
	a.SetDuration(dur)
	return a
}

func NewSizeAnimation(valPtr *geom.Point, val2 geom.Point, dur time.Duration) *PathAnimation {
    return NewPositionAnimation(valPtr, val2, dur)
}

func (a *PathAnimation) Init() {
	if a.Cont {
		a.Val1 = ConstPoint(*a.ValPtr)
	}
	a.val1 = a.Val1()
	a.val2 = a.Val2()
}

func (a *PathAnimation) Tick(t float64) {
	var dp geom.Point
	var s geom.Point

	dp = a.PathFunc(t)
	s = a.val2.Sub(a.val1)
	dp.X *= s.X
	dp.Y *= s.Y
	*a.ValPtr = a.val1.Add(dp)
}

//----------------------------------------------------------------------------

func NewPolyPathAnimation(valPtr *geom.Point, polyPath *PolygonPath, dur time.Duration) *PathAnimation {
	a := &PathAnimation{ValPtr: valPtr, Val1: ConstPoint(*valPtr),
		Val2: ConstPoint((*valPtr).AddXY(1, 1)), PathFunc: polyPath.RelPoint}
	a.NormAnimationEmbed.Extend(a)
	a.SetDuration(dur)
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

// Animation fuer eine Positionsveraenderung anhand des Fixed-Datentyps
// [fixed/Point26_6]. Dies wird insbesondere für die Positionierung von
// Schriften verwendet.
type FixedPosAnimation struct {
	NormAnimationEmbed
	ValPtr     *fixed.Point26_6
	Val1, Val2 fixed.Point26_6
	Cont       bool
}

func NewFixedPosAnimation(valPtr *fixed.Point26_6, val2 fixed.Point26_6, dur time.Duration) *FixedPosAnimation {
	a := &FixedPosAnimation{}
	a.NormAnimationEmbed.Extend(a)
	a.SetDuration(dur)
	a.ValPtr = valPtr
	a.Val1 = *valPtr
	a.Val2 = val2
	return a
}

func (a *FixedPosAnimation) Init() {
	if a.Cont {
		a.Val1 = *a.ValPtr
	}
}

func (a *FixedPosAnimation) Tick(t float64) {
	*a.ValPtr = a.Val1.Mul(float2fix(1.0 - t)).Add(a.Val2.Mul(float2fix(t)))
}

func float2fix(x float64) fixed.Int26_6 {
	return fixed.Int26_6(math.Round(x * 64))
}

type IntegerPosAnimation struct {
	NormAnimationEmbed
	ValPtr     *image.Point
	Val1, Val2 image.Point
	Cont       bool
}

func NewIntegerPosAnimation(valPtr *image.Point, val2 image.Point, dur time.Duration) *IntegerPosAnimation {
	a := &IntegerPosAnimation{}
	a.NormAnimationEmbed.Extend(a)
	a.SetDuration(dur)
	a.ValPtr = valPtr
	a.Val1 = *valPtr
	a.Val2 = val2
	return a
}

func (a *IntegerPosAnimation) Init() {
	if a.Cont {
		a.Val1 = *a.ValPtr
	}
}

func (a *IntegerPosAnimation) Tick(t float64) {
	v1 := geom.NewPointIMG(a.Val1)
	v2 := geom.NewPointIMG(a.Val2)
	np := v1.Mul(1.0 - t).Add(v2.Mul(t))
	*a.ValPtr = np.Int()
}

type ShaderFuncType func(x, y, t float64) float64

// Fuer den klassischen Shader wird pro Pixel folgende Animation gestartet.
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
