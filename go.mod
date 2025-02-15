module github.com/stefan-muehlebach/ledgrid

go 1.24.0

replace github.com/stefan-muehlebach/gg => ../gg

require (
	github.com/stefan-muehlebach/gg v1.3.4
	golang.org/x/image v0.23.0
	periph.io/x/conn/v3 v3.7.1
	periph.io/x/host/v3 v3.8.3
)

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	golang.org/x/text v0.21.0 // indirect
)
