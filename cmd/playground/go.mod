module github.com/stefan-muehlebach/ledgrid/cmd/playground

go 1.26.1

replace github.com/stefan-muehlebach/ledgrid => ../..

replace github.com/stefan-muehlebach/gg => ../../../gg

require (
	github.com/stefan-muehlebach/gg v1.4.1
	github.com/stefan-muehlebach/ledgrid v1.4.3
	golang.org/x/image v0.37.0
	golang.org/x/term v0.41.0
)

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.35.0 // indirect
	periph.io/x/conn/v3 v3.7.2 // indirect
	periph.io/x/host/v3 v3.8.5 // indirect
)
