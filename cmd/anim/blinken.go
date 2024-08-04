package main

import (
	"github.com/stefan-muehlebach/gg/geom"
	"encoding/xml"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

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
		log.Fatalf("Couldn't open file '%s': %v", fileName, err)
	}
	defer xmlFile.Close()

	byteValue, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		log.Fatalf("Couldn't read content of file: %v", err)
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
						log.Fatalf("Cannot parse '%s': %v", string(val), err)
					}
					idx = k*b.Channels + l
					b.Frames[i].Values[j][idx] = uint8(v)
				}
			}
		}
	}
	return b
}

func (b *BlinkenFile) Image(idx int) *Image {
	var c color.Color

	i := &Image{}
	i.img = image.NewRGBA(image.Rect(0, 0, b.Width, b.Height))
	colorScale := uint8(255 / ((1 << b.Bits) - 1))
	for row := range b.Height {
		for col := range b.Width {
			idxFrom := col * b.Channels
			idxTo := idxFrom + b.Channels
			src := b.Frames[idx].Values[row][idxFrom:idxTo:idxTo]
			switch b.Channels {
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
	i.Size = ConvertSize(geom.NewPointIMG(i.img.Bounds().Size()))
	return i
}
