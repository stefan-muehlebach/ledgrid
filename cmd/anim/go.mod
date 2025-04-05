module github.com/stefan-muehlebach/ledgrid/cmd/anim

go 1.24.2

replace github.com/stefan-muehlebach/ledgrid => ../..

replace github.com/stefan-muehlebach/gg => ../../../gg

require (
	github.com/korandiz/v4l v1.1.0
	github.com/stefan-muehlebach/gg v1.4.0
	github.com/stefan-muehlebach/ledgrid v1.4.1
	github.com/vladimirvivien/go4vl v0.0.5
	gocv.io/x/gocv v0.41.0
	golang.org/x/image v0.25.0
)

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	periph.io/x/conn/v3 v3.7.2 // indirect
	periph.io/x/host/v3 v3.8.4 // indirect
)
