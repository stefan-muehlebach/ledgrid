module github.com/stefan-muehlebach/ledgrid

go 1.22.3

replace github.com/stefan-muehlebach/gg => ../gg

require (
	fyne.io/fyne/v2 v2.4.5
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/stefan-muehlebach/gg v0.0.0-00010101000000-000000000000
	gocv.io/x/gocv v0.36.1
	golang.org/x/image v0.15.0
	periph.io/x/conn/v3 v3.7.0
	periph.io/x/host/v3 v3.8.2
)

require (
	github.com/fredbi/uri v1.0.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
