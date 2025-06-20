//go:build !tinygo

package ledgrid

import (
	"slices"
	"time"
)

// Eine Gruppe dient dazu, eine Anzahl von Animationen gleichzeitig zu starten.
// Die Laufzeit der Gruppe ist gleich der laengsten Laufzeit ihrer Animationen
// oder einer festen Dauer (je nachdem, welche Dauer groesser ist).
// Die Animationen, welche ueber eine Gruppe gestartet werden, sollten keine
// Endlos-Animationen sein, da sonst die Laufzeit der Gruppe ebenfalls
// endlos wird.
type Group struct {
	DurationEmbed
	// Gibt an, wie oft diese Gruppe wiederholt werden soll.
	RepeatCount int
	// Liste, der durch diese Gruppe gestarteten Tasks.
	Tasks []Task

	start, stop, end time.Time
	repeatsLeft      int
	running          bool
}

// Erstellt eine neue Gruppe, welche die Animationen in [anims] zusammen
// startet. Per Default ist die Laufzeit der Gruppe gleich der laengsten
// Laufzeit der hinzugefuegten Animationen.
func NewGroup(tasks ...Task) *Group {
	a := &Group{}
	a.Add(tasks...)
	// AnimCtrl.Add(0, a)
	return a
}

// Fuegt der Gruppe weitere Animationen hinzu.
func (a *Group) Add(tasks ...Task) {
	for _, task := range tasks {
		a.Tasks = append(a.Tasks, task)
	}
}

func (a *Group) updateDuration() {
	a.duration = time.Duration(0)
	for _, task := range a.Tasks {
		if anim, ok := task.(TimedAnimation); ok {
			a.duration = max(a.duration, anim.Duration())
		}
	}
}

// Startet die Gruppe.
func (a *Group) StartAt(t time.Time) {
	if a.running {
		return
	}
	a.updateDuration()
	a.start = t
	a.end = a.start.Add(a.duration)
	a.repeatsLeft = a.RepeatCount
	a.running = true
	for _, task := range a.Tasks {
		task.StartAt(t)
	}
	AnimCtrl.Add(a)
}

func (a *Group) Start() {
	a.StartAt(AnimCtrl.Now())
}

// Unterbricht die Ausfuehrung der Gruppe.
func (a *Group) Suspend() {
	if !a.running {
		return
	}
	a.stop = AnimCtrl.Now()
	for _, task := range a.Tasks {
        if anim, ok := task.(Animation); ok {
            anim.Suspend()
        }
	}
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
	for _, task := range a.Tasks {
        if anim, ok := task.(Animation); ok {
            anim.Continue()
        }
	}
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
		a.updateDuration()
		a.start = a.end
		a.end = a.start.Add(a.duration)
		for _, task := range a.Tasks {
			task.StartAt(t)
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
	activeTask       int
	start, stop, end time.Time
	repeatsLeft      int
	running          bool
}

// Erstellt eine neue Sequenz welche die Animationen in [anims] hintereinander
// ausfuehrt.
func NewSequence(tasks ...Task) *Sequence {
	a := &Sequence{}
	a.Add(tasks...)
	// AnimCtrl.Add(0, a)
	return a
}

func (a *Sequence) TimeInfo() (start, end time.Time) {
	return a.start, a.end
}

// Fuegt der Sequenz weitere Animationen hinzu.
func (a *Sequence) Add(tasks ...Task) {
	for _, task := range tasks {
		a.Tasks = append(a.Tasks, task)
	}
}

func (a *Sequence) Put(tasks ...Task) {
	for _, task := range tasks {
		a.Tasks = append([]Task{task}, a.Tasks...)
	}
}

func (a *Sequence) updateDuration() {
	a.duration = time.Duration(0)
	for _, task := range a.Tasks {
		if anim, ok := task.(TimedAnimation); ok {
			a.duration = a.duration + anim.Duration()
		}
	}
}

// Startet die Sequenz.
func (a *Sequence) StartAt(t time.Time) {
	if a.running {
		return
	}
	a.updateDuration()
	a.start = t
	a.end = a.start.Add(a.duration)
	a.activeTask = 0
	a.repeatsLeft = a.RepeatCount
	a.running = true
	a.Tasks[a.activeTask].StartAt(t)
	AnimCtrl.Add(a)
}

func (a *Sequence) Start() {
	a.StartAt(AnimCtrl.Now())
}

// Unterbricht die Ausfuehrung der Sequenz.
func (a *Sequence) Suspend() {
	if !a.running {
		return
	}
	a.stop = AnimCtrl.Now()
	for _, task := range a.Tasks {
        if anim, ok := task.(Animation); ok {
            anim.Suspend()
        }
	}
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
	for _, task := range a.Tasks {
        if anim, ok := task.(Animation); ok {
            anim.Continue()
        }
	}
	a.running = true
}

// Liefert den Status der Sequenz zurueck.
func (a *Sequence) IsRunning() bool {
	return a.running
}

// Wird durch den Controller periodisch aufgerufen, prueft ob Animationen
// dieser Sequenz noch am Laufen sind und startet ggf. die naechste.
func (a *Sequence) Update(t time.Time) bool {
	if a.activeTask < len(a.Tasks) {
		if job, ok := a.Tasks[a.activeTask].(Job); ok {
			if job.IsRunning() {
				return true
			}
		}
		a.activeTask++
	}
	if a.activeTask >= len(a.Tasks) {
		if t.After(a.end) {
			if a.repeatsLeft == 0 {
				a.running = false
				return false
			} else if a.repeatsLeft > 0 {
				a.repeatsLeft--
			}
			a.updateDuration()
			a.start = a.end
			a.end = a.start.Add(a.duration)
			a.activeTask = 0
			a.Tasks[a.activeTask].StartAt(t)
		}
		return true
	}
	a.Tasks[a.activeTask].StartAt(t)
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
	a.Slots = make([]*TimelineSlot, 0)
	// AnimCtrl.Add(0, a)
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
func (a *Timeline) StartAt(t time.Time) {
	if a.running {
		return
	}
	a.start = t
	a.end = a.start.Add(a.duration)
	a.repeatsLeft = a.RepeatCount
	a.nextSlot = 0
	a.running = true
	AnimCtrl.Add(a)
}

func (a *Timeline) Start() {
	a.StartAt(AnimCtrl.Now())
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
			task.StartAt(t)
		}
		a.nextSlot++
	}
	return true
}
