package colour

import (
	"fmt"
	"math"
)

// Store as HSV
// 0 <= Hue <= 360
// 0 <= Saturation <= 1
// 0 <= Value <= 1
type Colour struct {
	Hue, Saturation, Value float64
}

func (c Colour) String() string {
	return fmt.Sprintf("(hsv: %1.2f, %1.2f, %1.2f)", c.Hue, c.Saturation, c.Value)
}

func NewColour(hue, sat, val float64) Colour {
	return Colour{clampVal(hue, 0, 360), clampVal(sat, 0, 1), clampVal(val, 0, 1)}
}

func FromHue(hue uint16, sat uint8, bri uint8) Colour {
	// get hsv as 360.0, 1.0, 1.0
	h := float64(hue) / 65534.0 * 360.0
	s := float64(sat) / 255.0
	v := float64(bri) / 255.0
	return NewColour(h, s, v)
}

func (c *Colour) ToHue() (hue uint16, sat uint8, bri uint8) {
	hTmp := c.Hue * 65534.0 / 360.0
	hue = uint16(hTmp)
	sat = uint8(c.Saturation * 255.0)
	bri = uint8(c.Value * 255.0)
	return
}

func clampVal(val, min, max float64) float64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func (c Colour) AddHue(amount float64) Colour {
	newHue := c.Hue + amount
	if newHue == 0.0 {
		c.Hue = 0.0
	} else {
		c.Hue = math.Abs(math.Mod(newHue, 360.0))
	}
	return c
}

func (c Colour) AddSaturation(amount float64) Colour {
	c.Saturation = clampVal(c.Saturation+amount, 0.0, 1.0)
	return c
}

func (c Colour) AddValue(amount float64) Colour {
	c.Value = clampVal(c.Value+amount, 0.0, 1.0)
	return c
}

func (c Colour) ToHex() string {
	// Add 0.5 for rounding
	r, g, b := hsvToRgb(c.Hue, c.Saturation, c.Value)
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func FromHex(hex string) (Colour, error) {
	format := "#%02x%02x%02x"
	if len(hex) == 4 {
		format = "#%x%x%x"
	}

	var red, green, blue uint8
	n, err := fmt.Sscanf(hex, format, &red, &green, &blue)
	if err != nil {
		return Colour{}, err
	}
	if n != 3 {
		return Colour{}, fmt.Errorf("color: '%v' is not a valid hex color", hex)
	}

	h, s, v := rgbToHsv(red, green, blue)
	return NewColour(h, s, v), nil
}

func (c Colour) ToRgb() (r, g, b uint8) {
	return hsvToRgb(c.Hue, c.Saturation, c.Value)
}

func FromRgb(r, g, b uint8) (c Colour) {
	h, s, v := rgbToHsv(r, g, b)
	return NewColour(h, s, v)
}

func hsvToRgb(hue, saturation, value float64) (r, g, b uint8) {
	huePrime := hue / 60.0
	chroma := value * saturation
	X := chroma * (1.0 - math.Abs(math.Mod(huePrime, 2.0)-1.0))

	m := value - chroma
	red, green, blue := 0.0, 0.0, 0.0

	switch {
	case 0.0 <= huePrime && huePrime < 1.0:
		red = chroma
		green = X
	case 1.0 <= huePrime && huePrime < 2.0:
		red = X
		green = chroma
	case 2.0 <= huePrime && huePrime < 3.0:
		green = chroma
		blue = X
	case 3.0 <= huePrime && huePrime < 4.0:
		green = X
		blue = chroma
	case 4.0 <= huePrime && huePrime < 5.0:
		red = X
		blue = chroma
	case 5.0 <= huePrime && huePrime < 6.0:
		red = chroma
		blue = X
	}

	return uint8((red+m)*255.0 + 0.5), uint8((green+m)*255 + 0.5), uint8((blue+m)*255 + 0.5)
}

func rgbToHsv(red, green, blue uint8) (hue, saturation, value float64) {
	r := float64(red) / 255.0
	g := float64(green) / 255.0
	b := float64(blue) / 255.0

	min := math.Min(math.Min(r, g), b)
	value = math.Max(math.Max(r, g), b)
	chroma := value - min

	saturation = 0.0
	if value != 0.0 {
		saturation = chroma / value
	}

	hue = 0.0
	if min != value {
		if value == r {
			hue = math.Mod((g-b)/chroma, 6.0)
		} else if value == g {
			hue = (b-r)/chroma + 2.0
		} else {
			hue = (r-g)/chroma + 4.0
		}

		hue *= 60.0
		if hue < 0.0 {
			hue += 360
		}
	}

	return
}
