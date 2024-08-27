//go:build ignore

package ledgrid

import (
	"image"
	"image/color"
	"runtime"
	"time"

	"golang.org/x/image/draw"
)

var (
	theAnimator *Animator
)

func init() {
	theAnimator = newAnimator()
}

// Die Animation und die Darstellung von Objekten erfolgt zentral ueber einen
// Animator. Es ist sichergestellt, dass nur ein (1) Objekt von diesem Typ
// existiert.
type Animator struct {
	lg               *LedGrid
	client           PixelClient
	ticker           *time.Ticker
	animList         []Animation
	bgList, fgList   [2]Visual
	bgAnim, fgAnim   Animation
	bgAnimT, fgAnimT float64
}

// Erstellt einen neuen Animator, welcher fuer die Animation und die
// Darstellung aller Objekte auf dem LedGrid zustaendig ist.
func newAnimator() *Animator {
	if theAnimator != nil {
		return theAnimator
	}
	a := &Animator{}
	a.animList = make([]Animation, 0)

	return a
}

// Erstellt einen neuen Animator, welcher fuer die Animation und die
// Darstellung aller Objekte auf dem LedGrid zustaendig ist.
func NewAnimator(lg *LedGrid, client PixelClient) *Animator {
	theAnimator.lg = lg
	theAnimator.client = client
	theAnimator.coordinator()

	return theAnimator
}

func (a *Animator) SetBackground(obj Visual, fadeTime time.Duration) {
	if a.bgAnim != nil && !a.bgAnim.IsStopped() {
		return
	}
	obj.SetVisible(true)
	if a.bgList[0] == nil {
		a.bgList[0] = obj
		return
	}
	a.bgList[1] = obj
	a.bgAnim = NewNormAnimation(fadeTime, func(t float64) {
		if t == 1.0 {
			a.bgList[0].SetVisible(false)
			a.bgList[0], a.bgList[1] = a.bgList[1], nil
			a.bgAnimT = 0.0
		} else {
			a.bgAnimT = t
		}
	})
	a.bgAnim.Start()
}

func (a *Animator) SetForeground(obj Visual, fadeTime time.Duration) {
	if a.fgAnim != nil && !a.fgAnim.IsStopped() {
		return
	}
	obj.SetVisible(true)
	if a.fgList[0] == nil {
		a.fgList[0] = obj
		return
	}
	a.fgList[1] = obj
	a.fgAnim = NewNormAnimation(fadeTime, func(t float64) {
		if t == 1.0 {
			a.fgList[0].SetVisible(false)
			a.fgList[0], a.fgList[1] = a.fgList[1], nil
			a.fgAnimT = 0.0
		} else {
			a.fgAnimT = t
		}
	})
	a.fgAnim.Start()
}

func (a *Animator) Stop() {
	a.ticker.Stop()
}

func (a *Animator) addAnim(anim Animation) {
	a.animList = append(a.animList, anim)
}

func (a *Animator) Animations() []Animation {
	l := make([]Animation, len(a.animList))
	copy(l, a.animList)
	return l
}

type animationJob struct {
	anim Animation
	pt   time.Time
}

func (a *Animator) coordinator() {
	numCores := runtime.NumCPU()
	jobChan := make(chan animationJob, 2*numCores)
	doneChan := make(chan bool, 2*numCores)
	drawMask := image.NewUniform(color.Alpha{0xff})
	drawOpts := &draw.Options{
		SrcMask: drawMask,
	}
	scaler := draw.BiLinear

	for range numCores {
		go a.animUpdater(jobChan, doneChan)
	}
	a.ticker = time.NewTicker(frameRefresh)
	go func() {
		var srcRect, dstRect image.Rectangle
		var srcRatio, dstRatio float64

		dstRect = a.lg.Bounds()
		dstRatio = float64(dstRect.Dy()) / float64(dstRect.Dx())

		for pt := range a.ticker.C {
			numObjs := 0
			for _, anim := range a.animList {
				if anim.IsStopped() {
					continue
				}
				jobChan <- animationJob{anim, pt}
				numObjs++
			}
			for range numObjs {
				<-doneChan
			}
			a.lg.Clear(Black)

			if bg := a.bgList[0]; bg != nil {
				alpha := uint8((1 - a.bgAnimT) * 255.0)
				drawMask.C = color.Alpha{alpha}
				scaler.Scale(a.lg, a.lg.Bounds(), bg, bg.Bounds(), draw.Over, drawOpts)
				if bg := a.bgList[1]; bg != nil {
					alpha := uint8(a.bgAnimT * 255.0)
					drawMask.C = color.Alpha{alpha}
					scaler.Scale(a.lg, a.lg.Bounds(), bg, bg.Bounds(), draw.Over, drawOpts)
				}
			}
			if fg := a.fgList[0]; fg != nil {
				alpha := uint8((1 - a.fgAnimT) * 255.0)
				drawMask.C = color.Alpha{alpha}
				srcRect = fg.Bounds()
				dstRect = a.lg.Bounds()
				srcRatio = float64(srcRect.Dy()) / float64(srcRect.Dx())
				// log.Printf("srcRatio: %f, dstRatio: %f", srcRatio, dstRatio)
				if dstRatio > srcRatio {
					// Destination hoeher als Source
					h := int(srcRatio * float64(dstRect.Dx()))
					m := (dstRect.Dy() - h) / 2
					dstRect.Min.Y = m
					dstRect.Max.Y = m + h
				} else if dstRatio < srcRatio {
					// Destination flacher als Source
					w := int(float64(dstRect.Dy()) / srcRatio)
					m := (dstRect.Dx() - w) / 2
					dstRect.Min.X = m
					dstRect.Max.X = m + w
				}
				scaler.Scale(a.lg, dstRect, fg, srcRect, draw.Over, drawOpts)
				if fg := a.fgList[1]; fg != nil {
					alpha := uint8(a.fgAnimT * 255.0)
					drawMask.C = color.Alpha{alpha}
					srcRect = fg.Bounds()
					dstRect = a.lg.Bounds()
					srcRatio = float64(srcRect.Dy()) / float64(srcRect.Dx())
					// log.Printf("srcRatio: %f, dstRatio: %f", srcRatio, dstRatio)
					if dstRatio > srcRatio {
						// Destination hoeher als Source
						h := int(srcRatio * float64(dstRect.Dx()))
						m := (dstRect.Dy() - h) / 2
						dstRect.Min.Y = m
						dstRect.Max.Y = m + h
					} else if dstRatio < srcRatio {
						// Destination flacher als Source
						w := int(float64(dstRect.Dy()) / srcRatio)
						m := (dstRect.Dx() - w) / 2
						dstRect.Min.X = m
						dstRect.Max.X = m + w
					}
					scaler.Scale(a.lg, dstRect, fg, srcRect, draw.Over, drawOpts)
				}
			}
			a.client.Send(a.lg)
		}
	}()
}

func (a *Animator) animUpdater(jobChan <-chan animationJob, doneChan chan<- bool) {
	for job := range jobChan {
		job.anim.update(job.pt)
		doneChan <- true
	}
}

// Es gibt mehrere Implementationen von Animations-Objekten, alle haben jedoch
// das folgende Interface zu implementieren, damit sie vom zentralen
// Animator bedient werden koennen.
type Animation interface {
	Start()
	Stop()
	Cont()
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

	reverse          bool
	start, stop, end time.Time
	total            int64
	repeatsLeft      int
	running          bool
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

// Startet die Animation, resp. fuehrt einen Restart durch, falls die Animation
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
	if !a.running {
		return
	}
	a.stop = time.Now()
	a.running = false
}

func (a *NormAnimation) Cont() {
	if a.running {
		return
	}
	dt := time.Now().Sub(a.stop)
	a.start = a.start.Add(dt)
	a.end = a.end.Add(dt)
	a.running = true
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
	Tick        func(t float64)
	start, stop time.Time
	running     bool
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
	if !a.running {
		return
	}
	a.stop = time.Now()
	a.running = false
}

func (a *InfAnimation) Cont() {
	if a.running {
		return
	}
	a.start = a.start.Add(time.Now().Sub(a.stop))
	a.running = true
}

func (a *InfAnimation) IsStopped() bool {
	return !a.running
}

func (a *InfAnimation) update(t time.Time) bool {
	delta := t.Sub(a.start).Seconds()
	a.Tick(delta)
	return true
}
