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
	lg           *LedGrid
	client       *PixelClient
	ticker       *time.Ticker
	t0, stopTime time.Time
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
    a.objList = make([]any, 0)

	a.ticker = time.NewTicker(frameRefresh)
	a.t0 = time.Now()
    a.stopTime = time.Time{}
	go func() {
		for _ = range a.ticker.C {
			dt := frameRefresh
			for _, val := range a.objList {
                if obj, ok := val.(Animatable); ok {
                    if obj.Alive() {
				        obj.Update(dt)
                    }
                }
			}
			a.lg.Clear(BlackColor)
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

// Gibt den Status der Animation zurueck. true bedeutet, dass die Animation
// laeuft; false bedeutet, dass die Animation gestoppt ist.
func (a *Animator) IsRunning() (bool) {
    return a.stopTime.IsZero()
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
