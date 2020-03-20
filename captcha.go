package gocaptcha

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"math"
	"math/rand"
	"path"
	"time"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

const (
	defaultChar = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type Options struct {
	CharPreset      string      // 验证码预设的字符
	Length          int         // 验证码位数
	Width           int         // 验证码宽
	Height          int         // 验证码高
	Curve           int         // 曲线数 0 为不设曲线,默认0
	BackgroundColor color.Color // 背景色,默认黑色全透明
	FontDPI         float64     // 字体DPI,默认72.0
	FontScale       float64     // 字体比例系数,默认1.0
	Noise           float64     // 噪点密集度系数,默认1.0
}

type Data struct {
	Text string
	img  *image.NRGBA
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func defaultOptions(options *Options) *Options {
	opts := *options
	if opts.CharPreset == "" {
		opts.CharPreset = defaultChar
	}
	if opts.Length == 0 {
		opts.Length = 6
	}
	if opts.Width == 0 {
		opts.Width = 240
	}
	if opts.Height == 0 {
		opts.Height = 80
	}
	if opts.BackgroundColor == nil {
		opts.BackgroundColor = color.Transparent
	}
	if opts.FontDPI == 0 {
		opts.FontDPI = 72.0
	}
	if opts.FontScale == 0 {
		opts.FontScale = 1.0
	}
	if opts.Noise == 0 {
		opts.Noise = 1.0
	}

	return &Options{
		CharPreset:      opts.CharPreset,
		Length:          opts.Length,
		Width:           opts.Width,
		Height:          opts.Height,
		Curve:           opts.Curve,
		BackgroundColor: opts.BackgroundColor,
		FontDPI:         opts.FontDPI,
		FontScale:       opts.FontScale,
		Noise:           opts.Noise,
	}
}

// New 生成验证码
func New(options *Options) (*Data, error) {
	opts := defaultOptions(options)

	text := randText(opts)

	img := image.NewNRGBA(image.Rect(0, 0, opts.Width, opts.Height))

	draw.Draw(img, img.Bounds(), &image.Uniform{opts.BackgroundColor}, image.Point{}, draw.Src)
	drawNoise(img, opts)  // 噪点
	drawCurves(img, opts) // 曲线

	err := drawText(text, img, opts)
	if err != nil {
		return nil, err
	}

	return &Data{Text: text, img: img}, nil

}

// WriteImage 输出图片
func (data *Data) WriteImage(w io.Writer) error {
	return png.Encode(w, data.img)
}

// randText 随机取opts.Length位字作为验证码
func randText(opts *Options) (text string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	textNum := len(opts.CharPreset)
	for i := 0; i < opts.Length; i++ {
		text += string(opts.CharPreset[r.Intn(textNum)])
	}

	return text
}

// drawNoise 添加噪点
func drawNoise(img *image.NRGBA, opts *Options) {
	noiseCount := opts.Width * opts.Height / int(28.0/opts.Noise)
	for i := 0; i < noiseCount; i++ {
		x := rand.Intn(opts.Width)
		y := rand.Intn(opts.Height)
		img.Set(x, y, randColor())
	}
}

// randColor 噪点颜色
func randColor() color.RGBA {
	red := r.Intn(256)
	green := r.Intn(256)
	blue := r.Intn(256)

	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(255)}
}

// drawCurves 添加opts.Curve条曲线
func drawCurves(img *image.NRGBA, opts *Options) {
	for i := 0; i < opts.Curve; i++ {
		drawSineCurve(img, opts)
	}
}
func drawSineCurve(img *image.NRGBA, opts *Options) {
	var xStart, xEnd int
	if opts.Width <= 40 {
		xStart, xEnd = 1, opts.Width-1
	} else {
		xStart = r.Intn(opts.Width/10) + 1
		xEnd = opts.Width - r.Intn(opts.Width/10) - 1
	}
	curveHeight := float64(r.Intn(opts.Height/6) + opts.Height/6)
	yStart := r.Intn(opts.Height*2/3) + opts.Height/6
	angle := 1.0 + r.Float64()
	yFlip := 1.0
	if r.Intn(2) == 0 {
		yFlip = -1.0
	}
	curveColor := randMainColor(opts)

	for x1 := xStart; x1 <= xEnd; x1++ {
		y := math.Sin(math.Pi*angle*float64(x1)/float64(opts.Width)) * curveHeight * yFlip
		img.Set(x1, int(y)+yStart, curveColor)
		img.Set(x1, int(y)+yStart+1, curveColor) // 加粗
	}
}

// randMainColor 随机主要文字颜色,包括曲线
func randMainColor(opts *Options) color.Color {
	baseLightness := getLightness(opts.BackgroundColor)
	var value float64
	if baseLightness >= 0.5 {
		value = baseLightness - 0.3 - r.Float64()*0.2
	} else {
		value = baseLightness + 0.3 + r.Float64()*0.2
	}
	hue := float64(r.Intn(361)) / 360
	saturation := 0.6 + r.Float64()*0.2

	return hsva{h: hue, s: saturation, v: value, a: 255}
}

func getLightness(colour color.Color) float64 {
	r, g, b, a := colour.RGBA()
	// transparent
	if a == 0 {
		return 1.0
	}
	max := maxColor(r, g, b)
	min := minColor(r, g, b)

	l := (float64(max) + float64(min)) / (2 * 255)

	return l
}
func maxColor(numList ...uint32) (max uint32) {
	for _, num := range numList {
		colorVal := num & 255
		if colorVal > max {
			max = colorVal
		}
	}

	return max
}

func minColor(numList ...uint32) (min uint32) {
	min = 255
	for _, num := range numList {
		colorVal := num & 255
		if colorVal < min {
			min = colorVal
		}
	}

	return min
}

// 随机取一个字体
func randFontFamily() (*truetype.Font, error) {
	fontFiles, err := AssetDir("fonts")
	if err != nil {
		return &truetype.Font{}, err
	}

	fontFile := path.Join("fonts", fontFiles[r.Intn(len(fontFiles))])

	// fontBytes, err := ioutil.ReadFile(fontFile)
	fontBytes, err := Asset(fontFile)
	if err != nil {
		return &truetype.Font{}, err
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return &truetype.Font{}, err
	}
	return f, nil
}

// 在图片上写text
func drawText(text string, img *image.NRGBA, opts *Options) error {

	ttfFont, err := randFontFamily()
	if err != nil {
		return err
	}

	ctx := freetype.NewContext()
	ctx.SetDPI(opts.FontDPI)
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)
	ctx.SetHinting(font.HintingFull)
	ctx.SetFont(ttfFont)

	fontSpacing := opts.Width / len(text)
	fontOffset := r.Intn(fontSpacing / 3)

	for idx, char := range text {
		fontScale := 0.9 + r.Float64()*0.4
		fontSize := float64(opts.Height) / fontScale * opts.FontScale
		ctx.SetFontSize(fontSize)
		ctx.SetSrc(image.NewUniform(randMainColor(opts)))
		x := fontSpacing*idx + fontOffset
		y := opts.Height/5 + r.Intn(opts.Height/3) + int(fontSize/2)
		// y := int(-fontSize/6)+ctx.PointToFixed(fontSize).Ceil()
		pt := freetype.Pt(x, y)
		if _, err := ctx.DrawString(string(char), pt); err != nil {
			return err
		}
	}

	return nil
}
