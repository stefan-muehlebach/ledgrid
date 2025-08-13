module github.com/stefan-muehlebach/ledgrid/cmd/colorEdit

go 1.24.5

replace github.com/stefan-muehlebach/ledgrid => ../..

replace github.com/stefan-muehlebach/gg => ../../../gg

require (
	github.com/gbin/goncurses v0.0.0-20240517145248-be6a464272ae
	github.com/stefan-muehlebach/ledgrid v1.4.2
)

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/stefan-muehlebach/gg v1.4.1 // indirect
	golang.org/x/image v0.29.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	periph.io/x/conn/v3 v3.7.2 // indirect
	periph.io/x/host/v3 v3.8.5 // indirect
)
