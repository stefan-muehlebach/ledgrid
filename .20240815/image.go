//go:build ignore

package ledgrid

import (
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"golang.org/x/image/draw"
)

//----------------------------------------------------------------------------

// type Uniform struct {
//     pal ColorSource
// 	C LedColor
// }

// func NewUniform(pal ColorSource) *Uniform {
// 	u := &Uniform{}
//     u.pal = pal
// 	u.C = pal.Color(0)
// 	return u
// }

// func (u *Uniform) ColorModel() color.Model {
// 	return LedColorModel
// }

// func (u *Uniform) Bounds() image.Rectangle {
// 	return image.Rect(math.MinInt, math.MinInt, math.MaxInt, math.MaxInt)
// }

// func (u *Uniform) At(x, y int) color.Color {
// 	return u.pal.Color(0)
// }

// func (u *Uniform) Set(x, y int, c color.Color) {
// }

//----------------------------------------------------------------------------

type Image struct {
	VisualEmbed
	lg    *LedGrid
	img   draw.Image
	rect  image.Rectangle
}

func NewImageFromFile(lg *LedGrid, fileName string) *Image {
	i := &Image{}
	i.VisualEmbed.Init(fmt.Sprintf("%s (Image)", fileName))
	i.lg = lg
	fh, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Couldn't open file: %v", err)
	}
	tmp, err := png.Decode(fh)
	if err != nil {
		log.Fatalf("Couldn't decode file: %v", err)
	}
	i.img = tmp.(draw.Image)
	i.rect = i.img.Bounds()
	return i
}

func (i *Image) ColorModel() color.Model {
	return LedColorModel
}

func (i *Image) Bounds() image.Rectangle {
	return i.rect
}

func (i *Image) At(x, y int) color.Color {
	return i.img.At(x, y)
}

//----------------------------------------------------------------------------

type Uniform struct {
    VisualEmbed
    lg *LedGrid
    pal PaletteParameter
}

func NewUniform(lg *LedGrid, pal ColorSource) *Uniform {
	u := &Uniform{}
    u.VisualEmbed.Init("Uniform (Image)")
    u.lg = lg
    u.pal = NewPaletteParameter("Color", NewPaletteFader(pal))
	return u
}

func (u *Uniform) PaletteParam() PaletteParameter {
	return u.pal
}

func (u *Uniform) Palette() ColorSource {
	return u.pal.Val()
}

func (u *Uniform) SetPalette(pal ColorSource, fadeTime time.Duration) {
	u.pal.Val().(*PaletteFader).StartFade(pal, fadeTime)
}

func (u *Uniform) ColorModel() color.Model {
	return LedColorModel
}

func (u *Uniform) Bounds() image.Rectangle {
    return u.lg.Bounds()
	// return image.Rect(math.MinInt, math.MinInt, math.MaxInt, math.MaxInt)
}

func (u *Uniform) At(x, y int) color.Color {
	return u.pal.Val().Color(0)
}

//----------------------------------------------------------------------------


// func NewImageFromColor(lg *LedGrid, color *UniformPalette) *Image {
// 	i := &Image{}
// 	i.VisualEmbed.Init("Uniform Color (" + color.Name() + ")")
// 	i.lg = lg
//     i.pal = NewPaletteParameter("Color", NewPaletteFader(color))
// 	i.img = NewUniform(i.pal.Val())
// 	i.rect = i.lg.Bounds()
// 	return i
// }

// func (i *Image) Scale(dst draw.Image, dr image.Rectangle, src image.Image, sr image.Rectangle, op draw.Op, opts *draw.Options) {
// 	// draw.Draw(dst, dr, src, image.Point{}, op)
// }

func NewImageFromBlinken(lg *LedGrid, blk *BlinkenFile, fn int) *Image {
	var c color.Color

	i := &Image{}
	i.VisualEmbed.Init(fmt.Sprintf("Blinken [%d]", fn))
	i.lg = lg
	i.img = image.NewRGBA(image.Rect(0, 0, blk.Width, blk.Height))
	colorScale := uint8(255 / ((1 << blk.Bits) - 1))
	for row := range blk.Height {
		for col := range blk.Width {
			idxFrom := col * blk.Channels
			idxTo := idxFrom + blk.Channels
			src := blk.Frames[fn].Values[row][idxFrom:idxTo:idxTo]
			switch blk.Channels {
			case 1:
				v := colorScale * src[0]
				if v == 0 {
					c = color.RGBA{0, 0, 0, 0}
				} else {
					c = color.RGBA{v, v, v, 0xff}
				}
			case 3:
				r, g, b := colorScale*src[0], colorScale*src[1], colorScale*src[2]
				if r == 0 && g == 0 && b == 0 {
					c = color.RGBA{0, 0, 0, 0}
				} else {
					c = color.RGBA{r, g, b, 0xff}
				}
			}
			i.img.Set(col, row, c)
		}
	}
	i.rect = i.img.Bounds()
	return i
}

// func (i *Image) PaletteParam() PaletteParameter {
// 	return i.pal
// }

// func (i *Image) Palette() ColorSource {
// 	return i.pal.Val()
// }

// func (i *Image) SetPalette(pal ColorSource, fadeTime time.Duration) {
// 	i.pal.Val().(*PaletteFader).StartFade(pal, fadeTime)
// }


//----------------------------------------------------------------------------

type ImageAnimation struct {
	VisualEmbed
	lg       *LedGrid
	imgIdx   int
	total    float64
	imgList  []*Image
	timeList []float64
	rect     image.Rectangle
	anim     Animation
}

func NewImageAnimation(lg *LedGrid) *ImageAnimation {
	a := &ImageAnimation{}
	a.VisualEmbed.Init("Image Animation")
	a.lg = lg
	a.imgIdx = 0
	a.total = 0.0
	a.imgList = make([]*Image, 0)
	a.timeList = make([]float64, 0)
	return a
}

func (a *ImageAnimation) AddImage(img *Image, dur time.Duration) {
	a.imgList = append(a.imgList, img)
	a.total += dur.Seconds()
	a.timeList = append(a.timeList, a.total)
}

func (a *ImageAnimation) Frame() int {
	return a.imgIdx
}

func (a *ImageAnimation) SetFrame(i int) {
	a.imgIdx = i
}

func (a *ImageAnimation) Total() float64 {
	return a.total
}

func (a *ImageAnimation) SetVisible(vis bool) {
	if vis {
		a.anim.Start()
	} else {
		a.anim.Stop()
	}
	a.VisualEmbed.SetVisible(vis)
}

func (a *ImageAnimation) Update(t float64) {
	t = math.Mod(t, a.total)
	for i, ts := range a.timeList {
		if t <= ts {
			a.imgIdx = i
			return
		}
	}
}

func (a *ImageAnimation) ColorModel() color.Model {
	return LedColorModel
}

func (a *ImageAnimation) Bounds() image.Rectangle {
	return a.imgList[a.imgIdx].Bounds()
}

func (a *ImageAnimation) At(x, y int) color.Color {
	return a.imgList[a.imgIdx].At(x, y)
}

//----------------------------------------------------------------------------

type BlinkenFile struct {
	XMLName  xml.Name       `xml:"blm"`
	Width    int            `xml:"width,attr"`
	Height   int            `xml:"height,attr"`
	Bits     int            `xml:"bits,attr"`
	Channels int            `xml:"channels,attr"`
	Header   BlinkenHeader  `xml:"header"`
	Frames   []BlinkenFrame `xml:"frame"`
}

type BlinkenHeader struct {
	XMLName  xml.Name `xml:"header"`
	Title    string   `xml:"title"`
	Author   string   `xml:"author"`
	Email    string   `xml:"email"`
	Creator  string   `xml:"creator"`
	Duration int      `xml:"duration,omitempty"`
}

type BlinkenFrame struct {
	XMLName  xml.Name  `xml:"frame"`
	Duration int       `xml:"duration,attr"`
	Rows     [][]byte  `xml:"row"`
	Values   [][]uint8 `xml:"-"`
}

func ReadBlinkenFile(fileName string) *BlinkenFile {
	b := &BlinkenFile{Channels: 1}

	xmlFile, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()

	byteValue, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		log.Fatal(err)
	}

	err = xml.Unmarshal(byteValue, b)
	if err != nil {
		log.Fatal(err)
	}

	numberWidth := b.Bits / 4
	if b.Bits%4 != 0 {
		numberWidth++
	}
	for i, frame := range b.Frames {
		b.Frames[i].Values = make([][]uint8, b.Height)
		for j, row := range frame.Rows {
			b.Frames[i].Values[j] = make([]uint8, b.Width*b.Channels)
			for k := 0; k < b.Width; k++ {
				for l := range b.Channels {
					idx := k*numberWidth*b.Channels + l*numberWidth
					val := row[idx : idx+numberWidth]
					v, err := strconv.ParseUint(string(val), 16, b.Bits)
					if err != nil {
						log.Fatalf("'%s' not parseable: %v", string(val), err)
					}
					idx = k*b.Channels + l
					b.Frames[i].Values[j][idx] = uint8(v)
				}
			}
		}
	}
	return b
}

func (b *BlinkenFile) NewImageAnimation(lg *LedGrid) *ImageAnimation {
	a := NewImageAnimation(lg)
	a.SetName(b.Header.Title + " (BlinkenLight)")
	for i, frame := range b.Frames {
		img := NewImageFromBlinken(lg, b, i)
		a.AddImage(img, time.Duration(frame.Duration)*time.Millisecond)
	}
	a.anim = NewInfAnimation(a.Update)
	return a
}
