//go:build cameraOpenCV || cameraV4L2

package main

import (
	"context"

	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
)

func init() {
	// programList.AddTitle("Camera Animations")
	programList.Add("Ordinary camera", "Camera", OrdinaryCamera)
	programList.Add("Differential camera", "Camera", DiffCamera)

}

func OrdinaryCamera(ctx context.Context, canv *ledgrid.Canvas) {
	pos := geom.Point{}
	size := geom.Point{float64(width), float64(height)}

	cam := NewCamera(pos, size, ctx)
	canv.Add(cam)
	cam.Start()
}

func DiffCamera(ctx context.Context, canv *ledgrid.Canvas) {
	pos := geom.Point{}
	size := geom.Point{float64(width), float64(height)}

	cam := NewHistCamera(pos, size, ctx, colors.GoTeal)
	canv.Add(cam)
	cam.Start()
}
