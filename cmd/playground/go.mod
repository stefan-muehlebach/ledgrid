module github.com/stefan-muehlebach/ledgrid/cmd/playground

go 1.24.1

replace github.com/stefan-muehlebach/ledgrid => ../..

replace github.com/stefan-muehlebach/gg => ../../../gg

require (
	github.com/stefan-muehlebach/gg v1.3.4
	github.com/stefan-muehlebach/ledgrid v0.0.0-00010101000000-000000000000
	golang.org/x/image v0.23.0
	golang.org/x/term v0.23.0
)

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	golang.org/x/sys v0.23.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	periph.io/x/conn/v3 v3.7.1 // indirect
	periph.io/x/host/v3 v3.8.3 // indirect
)
