package ledgrid

import (
	"encoding/xml"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

// Fuer die Einbindung von Animationen im BlinkenLight-Format sind hier erst
// mal die Typen für die Verarbeitung der XML-Dateien.
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

// Mit folgender Funktion wird eine Datei im BlinkenLight-Format eingelesen.
// Die Bilddaten werden dabei noch nicht decodiert, d.h. noch nicht in ein
// 'image'-Format umgewandelt (siehe dazu auch die Methode [Decode]).
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

// Retourniert die Anzahl Frames in der BlinkenLight-Animation.
func (b *BlinkenFile) NumFrames() int {
	return len(b.Frames)
}

// Retourniert die Anzeigedauer des Frames mit Index idx. Für eine
// Aufsummierung der Anzeigewerte zur Berechnung der Anzeigedauer der gesamten
// Bildfolge ist ein aufrufendes Programm verantwortlich.
func (b *BlinkenFile) Duration(idx int) time.Duration {
	return time.Duration(b.Frames[idx].Duration) * time.Millisecond
}

func (b *BlinkenFile) SetAllDuration(durMs int) {
    for i := range b.Frames {
        b.Frames[i].Duration = durMs
    }
}

// Hier schliesslich werden die Bilddaten des Frames mit Index idx decodiert
// und als image.RGBA-Struktur zurückgegeben.
func (b *BlinkenFile) Decode(idx int) draw.Image {
	var c color.Color

	img := image.NewRGBA(image.Rect(0, 0, b.Width, b.Height))
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
			img.Set(col, row, c)
		}
	}
	return img
}

// func (b *BlinkenFile) Image(idx int) *Image {
// 	var c color.Color

// 	i := &Image{}
// 	i.Img = image.NewRGBA(image.Rect(0, 0, b.Width, b.Height))
// 	colorScale := uint8(255 / ((1 << b.Bits) - 1))
// 	for row := range b.Height {
// 		for col := range b.Width {
// 			idxFrom := col * b.Channels
// 			idxTo := idxFrom + b.Channels
// 			src := b.Frames[idx].Values[row][idxFrom:idxTo:idxTo]
// 			switch b.Channels {
// 			case 1:
// 				v := colorScale * src[0]
// 				if v == 0 {
// 					c = color.RGBA{0, 0, 0, 0}
// 				} else {
// 					c = color.RGBA{v, v, v, 0xff}
// 				}
// 			case 3:
// 				r, g, b := colorScale*src[0], colorScale*src[1], colorScale*src[2]
// 				if r == 0 && g == 0 && b == 0 {
// 					c = color.RGBA{0, 0, 0, 0}
// 				} else {
// 					c = color.RGBA{r, g, b, 0xff}
// 				}
// 			}
// 			i.Img.Set(col, row, c)
// 		}
// 	}
// 	i.Size = ConvertSize(geom.NewPointIMG(i.Img.Bounds().Size()))
// 	return i
// }
