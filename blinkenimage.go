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

type Image struct {
	VisualEmbed
	lg     *LedGrid
	img    draw.Image
	scaler draw.Scaler
}

func NewImageFromFile(lg *LedGrid, fileName string) *Image {
	i := &Image{}
	i.VisualEmbed.Init(fmt.Sprintf("Image '%s'", fileName))
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
	i.scaler = draw.BiLinear.NewScaler(lg.Bounds().Dx(), lg.Bounds().Dy(),
		i.img.Bounds().Dx(), i.img.Bounds().Dy())
	return i
}

func NewImageFromBlinken(lg *LedGrid, blk *BlinkenFile, fn int) *Image {
	var c color.Color

	i := &Image{}
	i.VisualEmbed.Init(fmt.Sprintf("Blinken [%d]", fn))
	i.lg = lg
	switch blk.Channels {
	case 1:
		i.img = image.NewGray(image.Rect(0, 0, blk.Width, blk.Height))
	case 3:
		i.img = image.NewRGBA(image.Rect(0, 0, blk.Width, blk.Height))
	default:
		log.Fatal("Only grayscale or RGB images are supported!")
	}
	colorScale := uint8(255 / ((1 << blk.Bits) - 1))
	for row := range blk.Height {
		for col := range blk.Width {
			idxFrom := col * blk.Channels
			idxTo := idxFrom + blk.Channels
			src := blk.Frames[fn].Values[row][idxFrom:idxTo:idxTo]
			switch blk.Channels {
			case 1:
				c = color.Gray{colorScale * src[0]}
			case 3:
				c = color.RGBA{R: colorScale * src[0], G: colorScale * src[1], B: colorScale * src[2], A: 255}
			}
			i.img.Set(col, row, c)
		}
	}
	i.scaler = draw.BiLinear.NewScaler(lg.Bounds().Dx(), lg.Bounds().Dy(),
		i.img.Bounds().Dx(), i.img.Bounds().Dy())
	return i
}

func (i *Image) Draw() {
	i.scaler.Scale(i.lg, i.lg.Bounds(), i.img, i.img.Bounds(), draw.Src, nil)
}

//----------------------------------------------------------------------------

type ImageAnimation struct {
	VisualEmbed
	lg       *LedGrid
	imgIdx   int
	total    float64
	imgList  []*Image
	timeList []float64
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

func (a *ImageAnimation) Draw() {
	a.imgList[a.imgIdx].Draw()
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
	theAnimator.AddAnimations(a.anim)
	return a
}
