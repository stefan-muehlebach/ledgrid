module github.com/stefan-muehlebach/ledgrid/cmd/colorEdit

go 1.26.2

replace github.com/stefan-muehlebach/ledgrid => ../..

replace github.com/stefan-muehlebach/gg => ../../../gg

require (
	github.com/gbin/goncurses v0.0.0-20251113135420-86371713952c
	github.com/stefan-muehlebach/gg v1.4.1
	github.com/stefan-muehlebach/ledgrid v1.4.3
)

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	golang.org/x/image v0.38.0 // indirect
	golang.org/x/text v0.35.0 // indirect
	periph.io/x/conn/v3 v3.7.3 // indirect
	periph.io/x/host/v3 v3.8.5 // indirect
)
