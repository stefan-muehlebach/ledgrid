//go:build cameraV4L2

package main

import (
	"context"
	"image"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"golang.org/x/image/draw"
)

var (
	CameraCtrls = []CameraCtrl{
		{v4l2.CtrlRotate, 180},
		{v4l2.CtrlHFlip, 1},
		{v4l2.CtrlColorFX, v4l2.CtrlValue(v4l2.ColorFXSkinWhiten)},
		{v4l2.CtrlCameraSceneMode, 8},
		{v4l2.CtrlCameraIsoSensitivityAuto, 1},
	}
)

type Camera struct {
	camera
	imgOut *image.RGBA
}

func NewCamera(pos, size geom.Point, ctx context.Context) *Camera {
	c := &Camera{}
	c.CanvasObjectEmbed.Extend(c)
	c.Init(pos, size, ctx)
	c.ctrlList = CameraCtrls
	c.imgOut = image.NewRGBA(image.Rect(0, 0, camWidth, camHeight))
	c.capture = func(frame *device.Frame) {
		copyRGB(c.imgOut, frame.Data)
	}

	return c
}

func (c *Camera) Draw(canv *ledgrid.Canvas) {
	c.scaler.Scale(canv.Img, c.dstRect, c.imgOut, c.srcRect, draw.Src, nil)
}
