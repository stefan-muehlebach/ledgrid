package ledgrid

import (
	"encoding/xml"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"golang.org/x/image/draw"
)

//----------------------------------------------------------------------------

type Image struct {
	VisualEmbed
	lg     *LedGrid
	img    image.Image
	scaler draw.Scaler
}

func NewImageFromFile(lg *LedGrid, fileName string) *Image {
	i := &Image{}
	i.VisualEmbed.Init("Image (" + fileName + ")")
	i.lg = lg
	fh, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Couldn't open file: %v", err)
	}
	i.img, err = png.Decode(fh)
	if err != nil {
		log.Fatalf("Couldn't decode file: %v", err)
	}
	i.scaler = draw.BiLinear.NewScaler(lg.Bounds().Dx(), lg.Bounds().Dy(),
		i.img.Bounds().Dx(), i.img.Bounds().Dy())
	return i
}

func (i *Image) Draw() {
	i.scaler.Scale(i.lg, i.lg.Bounds(), i.img, i.img.Bounds(), draw.Src, nil)
}

//----------------------------------------------------------------------------

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

// func (b *BlinkenFile) Write(fileName string) {
// 	var strBuild strings.Builder

// 	xmlFile, err := os.Create(fileName)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer xmlFile.Close()

// 	for i, frame := range b.Frames {
// 		for j, row := range frame.Values {
// 			strBuild.Reset()
// 			for _, v := range row {
// 				strBuild.WriteString(strconv.FormatUint(uint64(v), 16))
// 			}
// 			b.Frames[i].Rows[j] = []byte(strBuild.String())
// 		}
// 	}

// 	byteValue, err := xml.MarshalIndent(b, "", "    ")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	_, err = xmlFile.Write(byteValue)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
