package ledgrid

import (
	"time"
)

var (
	theAnimator *Animator
)

// Die Animation und die Darstellung von Objekten erfolgt zentral ueber einen
// Animator. Es ist sichergestellt, dass nur ein (1) Objekt von diesem Typ
// existiert.
type Animator struct {
	ticker       *time.Ticker
	t0, stopTime time.Time
	// t1           time.Duration
	// Speedup      float64
	lg           *LedGrid
	client       *PixelClient
	// animList     []Animatable
	// drawList     []Drawable
    objList      []any
}

// Erstellt einen neuen Animator, welcher fuer die Animation und die
// Darstellung aller Objekte auf dem LedGrid zustaendig ist.
func NewAnimator(lg *LedGrid, client *PixelClient) *Animator {
	if theAnimator != nil {
		return theAnimator
	}
	a := &Animator{}
	a.lg = lg
	a.client = client
	// a.animList = make([]Animatable, 0)
	// a.drawList = make([]Drawable, 0)
    a.objList = make([]any, 0)

	// a.Speedup = 1.0
	a.ticker = time.NewTicker(frameRefresh)
	go func() {
		a.t0 = time.Now()
		// a.t1 = time.Duration(0)
		for _ = range a.ticker.C {
			dt := frameRefresh
			// dt := time.Duration(float64(frameRefresh) * a.Speedup)
			// a.t1 += dt
			for _, val := range a.objList {
                if obj, ok := val.(Animatable); ok {
                    if obj.Alive() {
				        obj.Update(dt)
                    }
                }
				// obj.Update(a.t1.Seconds(), dt.Seconds())
			}
			a.lg.Clear(Black)
			for _, val := range a.objList {
                if obj, ok := val.(Drawable); ok {
                    if obj.Visible() {
				        obj.Draw()
                    }
                }
			}
			a.client.Draw(a.lg)
		}
	}()

	theAnimator = a
	return a
}

// Fuegt ein neues Objekt dem Animator hinzu. Je nachdem, ob das Objekt nur
// animiert oder auch gezeichnet werden kann, wird es in eine der zentralen
// Listen am Ende angehaengt.
func (a *Animator) AddObjects(objs ...any) {
    a.objList = append(a.objList, objs...)
	// for _, obj := range objs {
	// 	if v, ok := obj.(Animatable); ok {
	// 		a.animList = append(a.animList, v)
	// 	}
	// 	if v, ok := obj.(Drawable); ok {
	// 		a.drawList = append(a.drawList, v)
	// 	}
	// }
}

// Unterbricht die Animation.
func (a *Animator) Stop() {
	if !a.stopTime.IsZero() {
		return
	}
	a.ticker.Stop()
	a.stopTime = time.Now()
}

// Setzt die Animation wieder fort.
func (a *Animator) Reset() {
	if a.stopTime.IsZero() {
		return
	}
	a.t0 = a.t0.Add(time.Since(a.stopTime))
	a.ticker.Reset(frameRefresh)
	a.stopTime = time.Time{}
}
