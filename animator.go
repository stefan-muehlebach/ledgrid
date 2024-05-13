package ledgrid

import (
	"runtime"
	"time"
)

var (
	theAnimator *Animator
)

// Die Animation und die Darstellung von Objekten erfolgt zentral ueber einen
// Animator. Es ist sichergestellt, dass nur ein (1) Objekt von diesem Typ
// existiert.
type Animator struct {
	lg       *LedGrid
	client   PixelClient
	ticker   *time.Ticker
	objList  []Visual
	animList []Animation
	backObjs [2]Visual
	backAnim Animation
	backT    float64
	foreObjs [2]Visual
	foreAnim Animation
	foreT    float64
}

// Erstellt einen neuen Animator, welcher fuer die Animation und die
// Darstellung aller Objekte auf dem LedGrid zustaendig ist.
func NewAnimator(lg *LedGrid, client PixelClient) *Animator {
	if theAnimator != nil {
		return theAnimator
	}
	a := &Animator{}
	a.lg = lg
	a.client = client
	a.objList = make([]Visual, 0)
	a.animList = make([]Animation, 0)

	a.coordinator()

	theAnimator = a
	return a
}

func (a *Animator) SetBackground(obj Visual, fadeTime time.Duration) {
	if a.backAnim != nil && !a.backAnim.IsStopped() {
		return
	}
    obj.SetVisible(true)
	if a.backObjs[0] == nil {
		a.backObjs[0] = obj
		return
	}
	a.backObjs[1] = obj
	a.backAnim = NewNormAnimation(fadeTime, func(t float64) {
        if t == 1.0 {
            a.backObjs[0].SetVisible(false)
            a.backObjs[0], a.backObjs[1] = a.backObjs[1], nil
            a.backT = 0.0
        } else {
            a.backT = t
        }
	})
    a.backAnim.Start()
}

// Fuegt ein neues Objekt dem Animator hinzu.
func (a *Animator) AddObjects(objs ...Visual) {
	a.objList = append(a.objList, objs...)
}

// Retourniert alle Objekte.
func (a *Animator) Objects() []Visual {
	l := make([]Visual, len(a.objList))
	copy(l, a.objList)
	return l
}

func (r *Animator) addAnim(anim Animation) {
	r.animList = append(r.animList, anim)
}

func (r *Animator) Animations() []Animation {
	l := make([]Animation, len(r.animList))
	copy(l, r.animList)
	return l
}

func (a *Animator) coordinator() {
	numCores := runtime.NumCPU()
	objChan := make(chan Animation, 2*numCores)
	doneChan := make(chan bool, 2*numCores)
	for range numCores {
		go a.animUpdater(objChan, doneChan)
	}
	a.ticker = time.NewTicker(frameRefresh)
	go func() {
		for range a.ticker.C {
			numObjs := 0
			for _, obj := range a.animList {
				if obj.IsStopped() {
					continue
				}
				objChan <- obj
				numObjs++
			}
			for range numObjs {
				<-doneChan
			}
			a.lg.Clear(Black)
			// for _, obj := range a.objList {
			// 	if obj.Visible() {
			// 		obj.Draw()
			// 	}
			// }
            if a.backObjs[0] == nil {
                continue
            }
            a.backObjs[0].Draw()
			a.client.Draw(a.lg)
		}
	}()
}

func (a *Animator) animUpdater(objChan <-chan Animation, doneChan chan<- bool) {
	for obj := range objChan {
		obj.update(time.Now())
		doneChan <- true
	}
}

// Es gibt mehrere Implementationen von Animations-Objekten, alle haben jedoch
// das folgende Interface zu implementieren, damit sie vom zentralen
// Animator bedient werden koennen.
type Animation interface {
	Start()
	Stop()
	IsStopped() bool
	update(t time.Time) bool
}

//----------------------------------------------------------------------------

type AnimationCurve func(t float64) float64

func LinearAnimationCurve(t float64) float64 {
	return t
}

func CubicAnimationCurve(t float64) float64 {
	return 3*t*t - 2*t*t*t
}

const (
	AnimationRepeatForever = -1
)

// Mit der 'normierten' Animation (daher der Name) kann eine recht flexible
// Animation erzeugt werden.
type NormAnimation struct {
	// Falls true, wird die Animation am Ende des Zeitraums umgekehrt.
	AutoReverse bool
	// Mit Curve kann eine modifizierende Funktion definiert werden (siehe
	// auch Funktionstyp AnimationCurve).
	Curve AnimationCurve
	// Legt die Dauer der Animation fest.
	Duration time.Duration
	// Legt die Anzahl Wiederholungen der Animation fest.
	RepeatCount int
	// Diese Funktion wird bei jedem Update aufgerufen.
	Tick func(t float64)

	reverse     bool
	start, end  time.Time
	total       int64
	repeatsLeft int
	running     bool
}

// Erstellt eine neue Animation, welche ueber den Zeitraum d die Funktion
// fn aufruft. Bei t=0 wird die Funktion mit dem Wert 0.0 aufgerufen und bei
// t=d mit 1.0
func NewNormAnimation(d time.Duration, fn func(float64)) *NormAnimation {
	a := &NormAnimation{}
	a.Curve = LinearAnimationCurve
	a.Duration = d
	a.Tick = fn
    theAnimator.addAnim(a)
	return a
}

// Startet die Animation, resp. fuehrt eine Restart durch, falls die Animation
// bereits am Laufen ist.
func (a *NormAnimation) Start() {
	a.start = time.Now()
	a.end = a.start.Add(a.Duration)
	a.total = a.end.Sub(a.start).Milliseconds()
	a.repeatsLeft = a.RepeatCount
	a.running = true
}

// Stoppt die Animation.
func (a *NormAnimation) Stop() {
	a.running = false
}

// Dient der Abfrage, ob die Animation noch am Laufen ist.
func (a *NormAnimation) IsStopped() bool {
	return !a.running
}

// Interne Funktion, welche durch das System pro Update aufgerufen wird.
func (a *NormAnimation) update(t time.Time) bool {
	if t.After(a.end) {
		if a.reverse {
			a.Tick(a.Curve(0.0))
			if a.repeatsLeft == 0 {
				a.running = false
				return false
			}
			a.reverse = false
		} else {
			a.Tick(a.Curve(1.0))
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
		a.start = t
		a.end = a.start.Add(a.Duration)
		return true
	}

	delta := t.Sub(a.start).Milliseconds()
	val := float64(delta) / float64(a.total)
	if a.reverse {
		a.Tick(a.Curve(1.0 - val))
	} else {
		a.Tick(a.Curve(val))
	}
	return true
}

// Fuer Animationen, die permanent laufen sollen, ist die Infinite-Animation
// (InfAnimation) gedacht. Sie kann auch gestartet, resp. gestoppt werden,
// es gibt aber keine max. Animationsdauer oder die Möglichkeit die Animation
// zyklisch auszufuehren.
type InfAnimation struct {
	Tick    func(t float64)
	start   time.Time
	running bool
}

// Erzeugt eine neue Infinite-Animation, welche bei jedem Refresh die Funktion
// fn aufruft und ihr die Anzahl Sekunden (Fliesskommazahl) seit dem Start
// uebergibt.
func NewInfAnimation(fn func(float64)) *InfAnimation {
	a := &InfAnimation{}
	a.Tick = fn
    theAnimator.addAnim(a)
	return a
}

func (a *InfAnimation) Start() {
	a.start = time.Now()
	a.running = true
}

func (a *InfAnimation) Stop() {
	a.running = false
}

func (a *InfAnimation) IsStopped() bool {
	return !a.running
}

func (a *InfAnimation) update(t time.Time) bool {
	delta := t.Sub(a.start).Seconds()
	a.Tick(delta)
	return true
}
