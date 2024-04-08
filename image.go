package ledgrid

import (
	"golang.org/x/image/draw"
	"image/png"
	"encoding/xml"
	"image"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

//----------------------------------------------------------------------------

type Picture struct {
	DrawableEmbed
	lg  *LedGrid
	img image.Image
    scaler draw.Scaler
}

func NewPicture(lg *LedGrid, fileName string) *Picture {
	p := &Picture{}
	p.lg = lg
	fh, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Couldn't open file: %v", err)
	}
	p.img, err = png.Decode(fh)
	if err != nil {
		log.Fatalf("Couldn't decode file: %v", err)
	}
    p.scaler = draw.BiLinear.NewScaler(10, 10, p.img.Bounds().Dx(),
        p.img.Bounds().Dy())
	return p
}

func (p *Picture) Draw() {
    p.scaler.Scale(p.lg, p.lg.Bounds(), p.img, p.img.Bounds(), draw.Src, nil)
}

//----------------------------------------------------------------------------

type Image struct {
	DrawableEmbed
	lg  *LedGrid
	pal Colorable
	img []uint8
}

func NewImage(lg *LedGrid, pal Colorable) *Image {
	i := &Image{}
	i.DrawableEmbed.Init()
	i.lg = lg
	i.pal = pal
	i.img = make([]uint8, lg.Rect.Dx()*lg.Rect.Dy())
	return i
}

func (i *Image) Draw() {
	for idx, v := range i.img {
		row := idx / i.lg.Rect.Dx()
		col := idx % i.lg.Rect.Dx()
		fg := i.pal.Color(float64(v))
		bg := i.lg.LedColorAt(col, row)
		i.lg.SetLedColor(col, row, fg.Mix(bg, Blend))
	}
}

func (i *Image) SetPixels(pix [][]uint8) {
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
	imageList []*Image
	timeList  []time.Duration
	Idx       int
	Cycle     bool
}

func NewImageAnimation(lg *LedGrid) *ImageAnimation {
	i := &ImageAnimation{}
	i.VisualizableEmbed.Init("ImageAnimation")
	i.lg = lg
	i.imageList = make([]*Image, 0)
	i.timeList = make([]time.Duration, 0)
	i.Idx = 0
	i.Cycle = true
	return i
}

func (i *ImageAnimation) AddImage(img *Image, dur time.Duration) {
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

//----------------------------------------------------------------------------

type BlinkenFile struct {
	XMLName xml.Name       `xml:"blm"`
	Width   int            `xml:"width,attr"`
	Height  int            `xml:"height,attr"`
	Bits    int            `xml:"bits,attr"`
	Header  BlinkenHeader  `xml:"header"`
	Frames  []BlinkenFrame `xml:"frame"`
}

type BlinkenHeader struct {
	XMLName  xml.Name `xml:"header"`
	Title    string   `xml:"title"`
	Author   string   `xml:"author"`
	Email    string   `xml:"email"`
	Duration int      `xml:"duration,omitempty"`
}

type BlinkenFrame struct {
	XMLName  xml.Name  `xml:"frame"`
	Duration int       `xml:"duration,attr"`
	Rows     [][]byte  `xml:"row"`
	Values   [][]uint8 `xml:"-"`
}

func OpenBlinkenFile(fileName string) *BlinkenFile {
	b := &BlinkenFile{}

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

	for i, frame := range b.Frames {
		b.Frames[i].Values = make([][]uint8, b.Height)
		for j, row := range frame.Rows {
			b.Frames[i].Values[j] = make([]uint8, b.Width)
			for k, val := range row {
				v, err := strconv.ParseUint(string(val), 32, 8)
				if err != nil {
					log.Fatalf("'%s' not parseable: %v", string(val), err)
				}
				b.Frames[i].Values[j][k] = uint8(v)
			}
		}
	}
	return b
}

func (b *BlinkenFile) Write(fileName string) {
	var strBuild strings.Builder

	xmlFile, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()

	for i, frame := range b.Frames {
		for j, row := range frame.Values {
			strBuild.Reset()
			for _, v := range row {
				strBuild.WriteString(strconv.FormatUint(uint64(v), 32))
			}
			b.Frames[i].Rows[j] = []byte(strBuild.String())
		}
	}

	byteValue, err := xml.MarshalIndent(b, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	_, err = xmlFile.Write(byteValue)
	if err != nil {
		log.Fatal(err)
	}
}

func (b *BlinkenFile) MakeImageAnimation(lg *LedGrid, pal Colorable) *ImageAnimation {
	i := NewImageAnimation(lg)

	for _, frame := range b.Frames {
		img := NewImage(lg, pal)
		img.SetPixels(frame.Values)
		i.AddImage(img, time.Duration(frame.Duration)*time.Millisecond)
	}
	return i
}

//----------------------------------------------------------------------------

// func Main() {
// 	var blinkenFile *BlinkenFile
// 	var pal *SlicePalette

// 	blinkenFile = ReadBlinkenFile("alien.bml")

// 	frame := blinkenFile.Frames[20]
// 	fmt.Printf("%v\n", frame)

// 	pal = NewSlicePalette()
// 	i := NewImage(image.Point{10, 10}, pal)
// 	i.SetPixels(frame.Rows)
// }
