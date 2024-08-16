package ledgrid

// Verschiedene Objekte sollten Namen haben koennen, die man beispielsweise in
// GUIs oder TUIs anzeigen kann. Dieses Interface implementieren also alle
// benennbaren Objekte.
type Nameable interface {
	Name() string
	SetName(name string)
}

// Das entsprechende Embedable.
type NameableEmbed struct {
	name string
}

func (n *NameableEmbed) Init(name string) {
	n.SetName(name)
}

func (n *NameableEmbed) Name() string {
	return n.name
}

func (n *NameableEmbed) SetName(name string) {
	n.name = name
}
