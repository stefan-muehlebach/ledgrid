//go:build cameraV4L2

package main

import (
	"context"
	"image"

	"golang.org/x/image/draw"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"

	"github.com/stefan-muehlebach/gg/colors"
	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
)

var (
	HistCameraCtrls = []CameraCtrl{
		{v4l2.CtrlRotate, 180},
		{v4l2.CtrlHFlip, 1},
		{v4l2.CtrlColorFX, v4l2.CtrlValue(v4l2.ColorFXBlackWhite)},
		{v4l2.CtrlCameraSceneMode, 8},
		{v4l2.CtrlCameraIsoSensitivityAuto, 1},
	}
)

type HistCamera struct {
	camera
	color *image.Uniform
	img   *image.RGBA
	//imgRaw  *image.Gray
	imgIn   []*image.RGBA
	imgMask *image.Alpha
}

func NewHistCamera(pos, size geom.Point, ctx context.Context,
	col colors.RGBA) *HistCamera {
	c := &HistCamera{}
	c.CanvasObjectEmbed.Extend(c)
	c.Init(pos, size, ctx)
	c.ctrlList = HistCameraCtrls

	camRect := image.Rect(0, 0, camWidth, camHeight)

	c.color = image.NewUniform(col)
	c.img = image.NewRGBA(camRect)
	//c.imgRaw = image.NewGray(camRect)
	c.imgIn = make([]*image.RGBA, 2)
	for i := range 2 {
		c.imgIn[i] = image.NewRGBA(camRect)
	}
	c.imgMask = image.NewAlpha(camRect)
	c.capture = func(frame *device.Frame) {
		c.imgIn[0], c.imgIn[1] = c.imgIn[1], c.imgIn[0]
		copyRGB(c.imgIn[0], frame.Data)
		diffRGB(c.imgMask, c.imgIn[0], c.imgIn[1])
	}

	return c
}

func (c *HistCamera) Draw(canv *ledgrid.Canvas) {
	c.scaler.Scale(canv.Img, c.dstRect, c.color, c.srcRect, draw.Over,
		&draw.Options{
			SrcMask: c.imgMask,
		})
}
