package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

type GasData struct {
	Fast        decimal.Decimal
	Fastest     decimal.Decimal
	SafeLow     decimal.Decimal
	Average     decimal.Decimal
	FastWait    decimal.Decimal
	FastestWait decimal.Decimal
	SafeLowWait decimal.Decimal
	AverageWait decimal.Decimal
	BlockTime   decimal.Decimal `json:"block_time"`
	BlockNum    uint64
	DateAdded   time.Time
}

var (
	lightGray color.Color = color.Gray{}
	teal      color.Color = color.RGBA{0, 200, 200, 255}
	red       color.Color = color.RGBA{200, 30, 30, 255}
	blue      color.Color = color.RGBA{0, 0, 200, 255}
	black     color.Color = color.Black
	fadedBlue color.Color = color.RGBA{0, 0, 200, 10}

	maxX    = 1024
	maxY    = 1024
	marginX = 40
	marginY = 40
)

// Interface for displaying Shapes

type displayShape interface {
	drawShape() *image.RGBA
}

// Struct for rectangle

type Rectangle struct {
	p             image.Point
	length, width int
}

// Struct for circle

type Circle struct {
	p image.Point
	r int
}

func (c *Circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *Circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *Circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}

// Rectangle Draw shape function

func (r Rectangle) drawShape() *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, r.length, r.width))
}

func drawlines(pos [4]float64, m *image.RGBA, col color.Color) {
	gc := draw2dimg.NewGraphicContext(m)
	gc.MoveTo(pos[0], pos[1])
	gc.LineTo(pos[2], pos[3])
	gc.SetStrokeColor(col)
	gc.SetLineWidth(3)
	gc.Stroke()
}

// func addLabel(img *image.RGBA, x, y int, label string) {
// 	col := color.Black
// 	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

// 	d := &font.Drawer{
// 		Dst:  img,
// 		Src:  image.NewUniform(col),
// 		Face: basicfont.Face7x13,
// 		Dot:  point,
// 	}
// 	d.DrawString(label)
// }

func max(gds *[]GasData) decimal.Decimal {
	var m decimal.Decimal
	for i, e := range *gds {
		if i == 0 || e.SafeLow.GreaterThan(m) {
			m = e.SafeLow
		}
	}
	return m
}

func buildMap(gds *[]GasData) *image.RGBA {
	maxSafeLow := max(gds)
	maxSafeLowGrid := maxSafeLow.DivRound(decimal.NewFromInt(10), 0)
	maxY := int(maxSafeLowGrid.IntPart())*180 + marginY + 90

	firstTime := (*gds)[0].DateAdded
	lastTime := (*gds)[len(*gds)-1].DateAdded
	duration := lastTime.Sub(firstTime)
	numHours := duration.Hours()

	blockstart := (*gds)[0].BlockNum
	blockend := (*gds)[len(*gds)-1].BlockNum
	rangeX := blockend - blockstart
	scale := float64(maxX-2*marginX) / float64(rangeX)

	surface := Rectangle{length: maxX, width: maxY}  // Surface to draw on
	rpainter := Rectangle{length: maxX, width: maxY} // Colored Mask Layer

	m := surface.drawShape()
	cr := rpainter.drawShape()

	draw.Draw(m, m.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)
	draw.Draw(cr, cr.Bounds(), &image.Uniform{red}, image.ZP, draw.Src)

	/** Draws Edges **/
	marginXf := float64(marginX)
	drawlines([4]float64{marginXf, 0, marginXf, float64(maxY - 1)}, m, black)
	drawlines([4]float64{0, float64(maxY - marginY), float64(maxX), float64(maxY - marginY)}, m, black)
	for j := float64(maxY - marginY); j > 0; j -= 180 {
		drawlines([4]float64{0, j, float64(maxX), j}, m, fadedBlue)
	}

	for j := float64(40); marginXf+j < 1023; j += 40 {
		drawlines([4]float64{marginXf + j, 10, marginXf + j, float64(maxY - 1)}, m, fadedBlue)
	}
	//gasPrices := []decimal.Decimal{decimal.NewFromFloat(1.0), decimal.NewFromFloat(26.3), decimal.NewFromFloat(26.3), decimal.NewFromFloat(26.3)}
	var prevX, prevY float64
	for n, t := range *gds {
		if n == 0 {
			prevY = float64(maxY - marginY - int(t.SafeLow.IntPart())*18)
			prevX = float64(marginX + int(float64(t.BlockNum-blockstart)*scale))
			continue
		}
		yPos := float64(maxY - marginY - int(t.SafeLow.IntPart())*18)
		xPos := float64(marginX + int(float64(t.BlockNum-blockstart)*scale))
		//draw.DrawMask(m, m.Bounds(), cr, image.ZP, &Circle{image.Point{xPos, yPos}, 3}, image.ZP, draw.Over)
		drawlines([4]float64{prevX, prevY, xPos, yPos}, m, red)
		prevX = xPos
		prevY = yPos
	}
	for j := int64(1); j <= maxSafeLowGrid.IntPart(); j++ {
		yNum := int(j * 10)
		drawNumbers(yNum, 10, maxY-marginY-yNum*18, m)
	}
	// drawNumbers(50, 10, maxY-marginY-50*18, m)
	// drawNumbers(40, 10, maxY-marginY-40*18, m)
	// drawNumbers(30, 10, maxY-marginY-30*18, m)
	// drawNumbers(20, 10, maxY-marginY-20*18, m)
	// drawNumbers(10, 10, maxY-marginY-10*18, m)
	//---->
	start := (*gds)[0].DateAdded.Hour()
	mult := 1.0
	increment := float64(maxX-marginX) / numHours
	if numHours > 40 {
		increment *= numHours / 40
		mult = numHours / 40
		numHours = 40
	}
	for j := float64(0); j < numHours; j++ {
		xPos := marginX + int(j*mult*increment)
		yPos := maxY - marginY + 12
		drawNumbers((start+int(j*mult))%24, xPos, yPos, m)
	}
	GWei(marginX+3, 10, m)
	Hours(maxX/2-50, maxY-marginY+28, m)
	//pixfont.DrawString(m, 0, maxY-marginY-50*18, "50", color.Black)
	//addLabel(m, 0, maxY-marginY-50*18, "50")
	return m
}

// func Show(name string) {
// 	command := "open"
// 	arg1 := "-a"
// 	arg2 := "/Applications/Preview.app"
// 	cmd := exec.Command(command, arg1, arg2, name)
// 	err := cmd.Run()

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
var start, stop string

func initViper() {
	viper.SetConfigFile("./config.json")
	if err := viper.ReadInConfig(); err != nil {
		viper.SetConfigFile("./config.json")
		log.Fatal("config.json", err)
	}
}
func main() {
	initViper()

	logName := viper.GetString("log")
	log.SetOutput(&lumberjack.Logger{
		Filename:   "./logs/" + logName,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
	})
	flag.StringVar(&start, "start", "", "beginning of plot")
	flag.StringVar(&stop, "stop", "", "end of plot")
	flag.Parse()
	fmt.Println(">", start, "<>", stop, "<")
	gds, err := loadGasDataBetween(start, stop)
	if err != nil {
		log.Fatal(err)
	}
	for _, gas := range gds {
		fmt.Println(gas.DateAdded, gas.BlockNum, gas.SafeLow.String())
	}
	w, _ := os.Create("blogmap.png")
	defer w.Close()
	m := buildMap(&gds)
	png.Encode(w, m) //Encode writes the Image m to w in PNG format.

	//Show(w.Name())
}
