module github.com/stefan-muehlebach/ledgrid/cmd/tinygo

go 1.25.5

replace github.com/stefan-muehlebach/ledgrid => ../..

replace github.com/stefan-muehlebach/gg => ../../../gg

require (
	github.com/stefan-muehlebach/ledgrid v0.0.0-00010101000000-000000000000
	golang.org/x/image v0.34.0
)

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/stefan-muehlebach/gg v1.3.4 // indirect
	golang.org/x/text v0.32.0 // indirect
	periph.io/x/conn/v3 v3.7.2 // indirect
	periph.io/x/host/v3 v3.8.5 // indirect
)
