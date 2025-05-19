module github.com/stefan-muehlebach/ledgrid

go 1.24.3

replace github.com/stefan-muehlebach/gg => ../gg

require (
	github.com/stefan-muehlebach/gg v1.4.1
	golang.org/x/image v0.27.0
	periph.io/x/conn/v3 v3.7.2
	periph.io/x/host/v3 v3.8.5
)

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	golang.org/x/text v0.25.0 // indirect
)
