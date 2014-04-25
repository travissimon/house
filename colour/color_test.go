package colour

import (
	"log"
	"math"
	"testing"
)

func Test_ColourHookup(t *testing.T) {
	log.Println("Hookup Succeeded, testing colors")
}

var testValues = []struct {
	red        uint8
	green      uint8
	blue       uint8
	hue        float64
	saturation float64
	value      float64
	hex        string
}{
	{255, 255, 255, 0.0, 0.0, 1.0, "#ffffff"},
	{128, 255, 255, 180.0, 0.5, 1.0, "#80ffff"},
	{255, 128, 255, 300.0, 0.5, 1.0, "#ff80ff"},
	{255, 255, 128, 60.0, 0.5, 1.0, "#ffff80"},
	{128, 128, 128, 0.0, 0.0, 0.5, "#808080"},
	{255, 0, 0, 0.0, 1.0, 1.0, "#ff0000"},
	{0, 255, 0, 120.0, 1.0, 1.0, "#00ff00"},
	{0, 0, 255, 240.0, 1.0, 1.0, "#0000ff"},
	{0, 255, 255, 180.0, 1.0, 1.0, "#00ffff"},
}

func Test_HsvToRgb(t *testing.T) {
	for _, row := range testValues {
		hue := row.hue
		sat := row.saturation
		val := row.value

		expR := row.red
		expG := row.green
		expB := row.blue

		r, g, b := hsvToRgb(hue, sat, val)
		if r != expR {
			t.Errorf("Red, expected: %v, actual: %v", expR, r)
		}
		if g != expG {
			t.Errorf("Green, expected: %v, actual: %v", expG, g)
		}
		if b != expB {
			t.Errorf("Blue, expected: %v, actual: %v", expB, b)
		}
	}
}

func Test_RgbToHsv(t *testing.T) {
	for _, row := range testValues {
		expH := row.hue
		expS := row.saturation
		expV := row.value

		h, s, v := rgbToHsv(row.red, row.green, row.blue)
		if !ApproximatelyEqual(h, expH) {
			t.Errorf("Hue, expected: %v, actual: %v", expH, h)
		}
		if !ApproximatelyEqual(s, expS) {
			t.Errorf("Saturation, expected: %v, actual: %v", expS, s)
		}
		if !ApproximatelyEqual(v, expV) {
			t.Errorf("Value, expected: %v, actual: %v", expV, v)
		}
	}
}

func Test_Hex(t *testing.T) {
	for _, row := range testValues {
		expHex := row.hex

		c := NewColour(row.hue, row.saturation, row.value)
		hex := c.ToHex()

		if expHex != hex {
			t.Errorf("Expected hex '%v', actual: '%v'", expHex, hex)
		}

		fromHex, _ := FromHex(expHex)
		if !ApproximatelyEqual(fromHex.Hue, row.hue) {
			t.Errorf("fromHex hue, expected %v, actual: %v", row.hue, fromHex.Hue)
		}

		if !ApproximatelyEqual(fromHex.Saturation, row.saturation) {
			t.Errorf("fromHex saturation, expected %v, actual: %v", row.saturation, fromHex.Saturation)
		}

		if !ApproximatelyEqual(fromHex.Value, row.value) {
			t.Errorf("fromHex value, expected %v, actual %v", row.value, fromHex.Value)
		}

	}
}

func Test_HueAddition(t *testing.T) {
	c := NewColour(0.0, 1.0, 1.0)

	for i := 0.0; i < 900.0; i++ {
		mod := math.Mod(i, 360.0)
		c2 := c.AddHue(i)
		if c2.Hue != mod {
			t.Errorf("Could not add %v to Hue", i)
		}
	}

	for i := 0.0; i > -900.0; i-- {
		mod := math.Abs(math.Mod(i, 360))
		c2 := c.AddHue(i)
		if c2.Hue != mod {
			t.Errorf("Error adding Hue (%v + %v = %v), expected: %v", c.Hue, i, c2.Hue, mod)
		}
	}
}

func Test_SaturationAndValueAddition(t *testing.T) {
	c := NewColour(0.0, 0.0, 0.0)

	for i := -2.0; i < 2.0; i += 0.1 {
		sat := c.AddSaturation(i)
		val := c.AddValue(i)

		if i < 0.0 {
			if sat.Saturation != 0.0 {
				t.Errorf("Saturation: exp 0.0, actual: %v", sat.Saturation)
			}
			if val.Value != 0.0 {
				t.Errorf("Value: exp 0.0, actual: %v", val.Value)
			}
		} else if i > 1.0 {
			if sat.Saturation != 1.0 {
				t.Errorf("Saturation: exp 1.0, actual: %v", sat.Saturation)
			}
			if val.Value != 1.0 {
				t.Errorf("Value: exp 1.0, actual: %v", val.Value)
			}
		} else {
			if sat.Saturation != i {
				t.Errorf("Saturation: exp %v, actual: %v", i, sat.Saturation)
			}
			if val.Value != i {
				t.Errorf("Value: exp %v, actual: %v", i, val.Value)
			}
		}
	}
}

func ApproximatelyEqual(v1, v2 float64) bool {
	return math.Abs(v1-v2) < 0.01
}
