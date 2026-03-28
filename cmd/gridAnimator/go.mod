module github.com/stefan-muehlebach/ledgrid/cmd/gridAnimator

go 1.26.1

replace github.com/stefan-muehlebach/ledgrid => ../..

replace github.com/stefan-muehlebach/gg => ../../../gg

require (
	github.com/korandiz/v4l v1.1.0
	github.com/stefan-muehlebach/gg v1.4.1
	github.com/stefan-muehlebach/ledgrid v1.4.3
	github.com/vladimirvivien/go4vl v0.3.0
	gocv.io/x/gocv v0.43.0
	golang.org/x/image v0.37.0
)

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.35.0 // indirect
	periph.io/x/conn/v3 v3.7.2 // indirect
	periph.io/x/host/v3 v3.8.5 // indirect
)
