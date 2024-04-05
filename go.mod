module github.com/stefan-muehlebach/ledgrid

go 1.22.1

replace github.com/stefan-muehlebach/gg => ../gg

require (
	github.com/rthornton128/goncurses v0.0.0-20231014161942-82671379df88
	github.com/stefan-muehlebach/gg v0.0.0-00010101000000-000000000000
	golang.org/x/image v0.15.0
	periph.io/x/conn/v3 v3.7.0
	periph.io/x/host/v3 v3.8.2
)

require golang.org/x/text v0.14.0 // indirect
