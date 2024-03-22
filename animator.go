package ledgrid

import "time"

//----------------------------------------------------------------------------

type Animator struct {
	ticker       *time.Ticker
	t0, stopTime time.Time
	grid         *LedGrid
	client       *PixelClient
	animList     []Animatable
	drawList     []Drawable
}

func NewAnimator(grid *LedGrid, client *PixelClient) *Animator {
	a := &Animator{}
	a.grid = grid
	a.client = client
	a.animList = make([]Animatable, 0)
	a.drawList = make([]Drawable, 0)

	a.ticker = time.NewTicker(frameRefresh)
	// a.ticker = time.NewTicker(time.Duration(frameRefreshMs) * time.Millisecond)
	go func() {
		a.t0 = time.Now()
		for t := range a.ticker.C {
			t1 := t.Sub(a.t0).Seconds()
			for _, obj := range a.animList {
				obj.Update(t1)
			}
			a.grid.Clear(Black)
			for _, obj := range a.drawList {
				obj.Draw(a.grid)
			}
			a.client.Draw(a.grid)
		}
	}()

	return a
}

func (a *Animator) AddObject(obj any) {
	if v, ok := obj.(Animatable); ok {
		a.animList = append(a.animList, v)
	}
	if v, ok := obj.(Drawable); ok {
		a.drawList = append(a.drawList, v)
	}
}

func (a *Animator) Stop() {
	if !a.stopTime.IsZero() {
		return
	}
	a.ticker.Stop()
	a.stopTime = time.Now()
}

func (a *Animator) Reset() {
	if a.stopTime.IsZero() {
		return
	}
	a.t0 = a.t0.Add(time.Since(a.stopTime))
	a.ticker.Reset(frameRefresh)
	a.stopTime = time.Time{}
}
