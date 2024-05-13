//go:build ignore

package ledgrid

import (
	"image"
	"image/png"
	"log"
	"os"
	"time"
)


type PictureAnimation struct {
	VisualEmbed
	lg       *LedGrid
	pictList []*Picture
	timeList []time.Duration
	Idx      int
	Cycle    bool
}

func NewPictureAnimation(lg *LedGrid) *PictureAnimation {
	i := &PictureAnimation{}
	i.VisualEmbed.Init("PictAnim")
	i.lg = lg
	i.pictList = make([]*Picture, 0)
	i.timeList = make([]time.Duration, 0)
	i.Idx = 0
	i.Cycle = true
	return i
}

func (i *PictureAnimation) AddPicture(pict *Picture, dur time.Duration) {
	i.pictList = append(i.pictList, pict)
	if len(i.timeList) > 0 {
		dur += i.timeList[len(i.timeList)-1]
	}
	i.timeList = append(i.timeList, dur)
}

func (i *PictureAnimation) Update(dt time.Duration) bool {
	// i.AnimatableEmbed.Update(dt)
	t := i.t0 % i.timeList[len(i.timeList)-1]
	for idx, v := range i.timeList {
		if t < v {
			i.Idx = idx
			return true
		}
	}
	return true
}

func (i *PictureAnimation) Draw() {
	i.pictList[i.Idx].Draw()
}

//----------------------------------------------------------------------------

type PixelImage struct {
	VisualEmbed
	lg  *LedGrid
	pal ColorSource
	img []uint8
}

func NewPixelImage(lg *LedGrid, pal ColorSource) *PixelImage {
	i := &PixelImage{}
	i.DrawableEmbed.Init()
	i.lg = lg
	i.pal = pal
	i.img = make([]uint8, lg.Rect.Dx()*lg.Rect.Dy())
	return i
}

func (i *PixelImage) Draw() {
	for idx, v := range i.img {
		row := idx / i.lg.Rect.Dx()
		col := idx % i.lg.Rect.Dx()
		fg := i.pal.Color(float64(v))
		bg := i.lg.LedColorAt(col, row)
		i.lg.SetLedColor(col, row, fg.Mix(bg, Blend))
	}
}

func (i *PixelImage) SetPixels(pix [][]uint8) {
	for row, data := range pix {
		for col, v := range data {
			i.img[row*i.lg.Rect.Dx()+col] = v
		}
	}
}

//----------------------------------------------------------------------------

type ImageAnimation struct {
	VisualizableEmbed
	lg        *LedGrid
	imageList []*PixelImage
	timeList  []time.Duration
	Idx       int
	Cycle     bool
}

func NewImageAnimation(lg *LedGrid) *ImageAnimation {
	i := &ImageAnimation{}
	i.VisualizableEmbed.Init("ImageAnim")
	i.lg = lg
	i.imageList = make([]*PixelImage, 0)
	i.timeList = make([]time.Duration, 0)
	i.Idx = 0
	i.Cycle = true
	return i
}

func (i *ImageAnimation) AddImage(img *PixelImage, dur time.Duration) {
	i.imageList = append(i.imageList, img)
	if len(i.timeList) > 0 {
		dur += i.timeList[len(i.timeList)-1]
	}
	i.timeList = append(i.timeList, dur)
}

func (i *ImageAnimation) Update(dt time.Duration) bool {
	i.AnimatableEmbed.Update(dt)
	t := i.t0 % i.timeList[len(i.timeList)-1]
	for idx, v := range i.timeList {
		if t < v {
			i.Idx = idx
			return true
		}
	}
	return true
}

func (i *ImageAnimation) Draw() {
	i.imageList[i.Idx].Draw()
}
