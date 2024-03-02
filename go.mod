module github.com/stefan-muehlebach/ledgrid

go 1.22.0

replace github.com/stefan-muehlebach/gg => ../gg

require (
	github.com/gbin/goncurses v0.0.0-20240205203827-f1026e55db44
	github.com/stefan-muehlebach/gg v0.0.0-00010101000000-000000000000
	periph.io/x/conn/v3 v3.7.0
	periph.io/x/host/v3 v3.8.2
)

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	golang.org/x/image v0.15.0 // indirect
)
