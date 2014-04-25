package colour

import (
//	"log"
)

type HarmonyStrategy int

const (
	Monochromatic      HarmonyStrategy = iota
	Complementary                      = iota
	Triad                              = iota
	SplitComplementary                 = iota
	Tetradic                           = iota
	Analogic                           = iota
	AccentedAnalogic                   = iota
	Square                             = iota
)

func GetHarmonyStrategy(strategy string) HarmonyStrategy {
	switch strategy {
	case "Monochromatic":
		return Monochromatic
	case "Complementary":
		return Complementary
	case "Triad":
		return Triad
	case "SplitComplementary":
		return SplitComplementary
	case "Tetradic":
		return Tetradic
	case "Analogic":
		return Analogic
	case "AccentedAnalogic":
		return AccentedAnalogic
	case "Square":
		return Square
	}
	return Monochromatic
}

type ColourSchemeGenerator struct {
	strategy HarmonyStrategy
	angle    float64
	tint     float64
	shade    float64
}

func NewColourSchemeGenerator() *ColourSchemeGenerator {
	return &ColourSchemeGenerator{SplitComplementary, 30.0, 0.3, 0.3}
}

type ColourType int

const (
	PrimaryColour       ColourType = iota
	SecondaryColour                = iota
	ComplementaryColour            = iota
)

type SchemeColour struct {
	Colour     Colour
	ColourType ColourType
	Tints      []Colour
	Shades     []Colour
}

func (c *ColourSchemeGenerator) SetStrategy(strategy HarmonyStrategy) {
	c.strategy = strategy

	// these strategies have specific angles
	switch strategy {
	case Triad:
	case Square:
		c.angle = 60.0
		break
	}
}

func (c *ColourSchemeGenerator) getClampedVal(val, min, max float64) float64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func (c *ColourSchemeGenerator) SetAngle(angle float64) {
	c.angle = c.getClampedVal(angle, -90.0, 90.00)
}

func (c *ColourSchemeGenerator) SetTint(tint float64) {
	c.tint = c.getClampedVal(tint, 0.0, 1.0)
}

func (c *ColourSchemeGenerator) SetShade(shade float64) {
	c.shade = c.getClampedVal(shade, 0.0, 1.0)
}

func (c *ColourSchemeGenerator) GetScheme(color Colour) []SchemeColour {
	switch c.strategy {
	case Monochromatic:
		return c.GenerateMonochromatic(color)
	case Complementary:
		return c.GenerateComplementary(color)
	case Triad:
	case SplitComplementary:
		return c.GenerateTriad(color)
	case Square:
	case Tetradic:
		return c.GenerateTetradic(color)
	case Analogic:
		return c.GenerateAnalogic(color)
	case AccentedAnalogic:
		return c.GenerateAccentedAnalogic(color)
	}
	return nil
}

func (c *ColourSchemeGenerator) GenerateTints(color Colour) []Colour {
	distanceToWhite := 1.0 - color.Value
	valStepUp := distanceToWhite * c.tint * .5
	satStepDown := color.Saturation * c.tint * .5

	colors := make([]Colour, 0, 2)
	colors = append(colors, NewColour(color.Hue, color.Saturation-satStepDown, color.Value+valStepUp))
	colors = append(colors, NewColour(color.Hue, color.Saturation-(2*satStepDown), color.Value+valStepUp))
	return colors
}

func (c *ColourSchemeGenerator) GenerateShades(color Colour) []Colour {
	valStepDown := color.Value * c.shade * .5
	distanceToPale := 1.0 - color.Saturation
	satStepUp := distanceToPale * c.shade * .5

	colors := make([]Colour, 0, 2)
	colors = append(colors, NewColour(color.Hue, color.Saturation+satStepUp, color.Value-valStepDown))
	colors = append(colors, NewColour(color.Hue, color.Saturation+satStepUp, color.Value-(1.5*valStepDown)))
	return colors
}

func (c *ColourSchemeGenerator) GenerateMonochromatic(colour Colour) []SchemeColour {
	colours := make([]SchemeColour, 0, 1)
	primary := SchemeColour{colour, PrimaryColour, nil, nil}
	primary.Tints = c.GenerateTints(primary.Colour)
	primary.Shades = c.GenerateShades(primary.Colour)

	colours = append(colours, primary)
	return colours
}

func (c *ColourSchemeGenerator) GenerateComplementary(colour Colour) []SchemeColour {
	colours := make([]SchemeColour, 0, 2)

	primary := SchemeColour{colour, PrimaryColour, nil, nil}
	primary.Tints = c.GenerateTints(primary.Colour)
	primary.Shades = c.GenerateShades(primary.Colour)

	complementaryColour := NewColour(colour.Hue, colour.Saturation, colour.Value)
	complementaryColour = complementaryColour.AddHue(180)
	complementary := SchemeColour{complementaryColour, ComplementaryColour, nil, nil}
	complementary.Tints = c.GenerateTints(complementary.Colour)
	complementary.Shades = c.GenerateShades(complementary.Colour)

	colours = append(colours, primary)
	colours = append(colours, complementary)
	return colours
}

func (c *ColourSchemeGenerator) GenerateTriad(colour Colour) []SchemeColour {
	colours := make([]SchemeColour, 0, 3)

	primary := SchemeColour{colour, PrimaryColour, nil, nil}
	primary.Tints = c.GenerateTints(primary.Colour)
	primary.Shades = c.GenerateShades(primary.Colour)

	secondaryColour1 := NewColour(colour.Hue, colour.Saturation, colour.Value)
	secondaryColour1 = secondaryColour1.AddHue(180.0 + c.angle)
	secondary1 := SchemeColour{secondaryColour1, SecondaryColour, nil, nil}
	secondary1.Tints = c.GenerateTints(secondary1.Colour)
	secondary1.Shades = c.GenerateShades(secondary1.Colour)

	secondaryColour2 := NewColour(colour.Hue, colour.Saturation, colour.Value)
	secondaryColour2 = secondaryColour2.AddHue(180.0 - c.angle)
	secondary2 := SchemeColour{secondaryColour2, SecondaryColour, nil, nil}
	secondary2.Tints = c.GenerateTints(secondary2.Colour)
	secondary2.Shades = c.GenerateShades(secondary2.Colour)

	colours = append(colours, primary)
	colours = append(colours, secondary1)
	colours = append(colours, secondary2)
	return colours
}

func (c *ColourSchemeGenerator) GenerateTetradic(colour Colour) []SchemeColour {
	colours := make([]SchemeColour, 0, 4)

	primary := SchemeColour{colour, PrimaryColour, nil, nil}
	primary.Tints = c.GenerateTints(primary.Colour)
	primary.Shades = c.GenerateShades(primary.Colour)

	secondaryColour1 := NewColour(colour.Hue, colour.Saturation, colour.Value)
	secondaryColour1 = secondaryColour1.AddHue(c.angle)
	secondary1 := SchemeColour{secondaryColour1, SecondaryColour, nil, nil}
	secondary1.Tints = c.GenerateTints(secondary1.Colour)
	secondary1.Shades = c.GenerateShades(secondary1.Colour)

	secondaryColour2 := NewColour(colour.Hue, colour.Saturation, colour.Value)
	secondaryColour2 = secondaryColour2.AddHue(180.0 - c.angle)
	secondary2 := SchemeColour{secondaryColour2, SecondaryColour, nil, nil}
	secondary2.Tints = c.GenerateTints(secondary2.Colour)
	secondary2.Shades = c.GenerateShades(secondary2.Colour)

	complementaryColour := NewColour(colour.Hue, colour.Saturation, colour.Value)
	complementaryColour = secondaryColour2.AddHue(180.0)
	complementary := SchemeColour{complementaryColour, ComplementaryColour, nil, nil}
	complementary.Tints = c.GenerateTints(complementary.Colour)
	complementary.Shades = c.GenerateShades(complementary.Colour)

	colours = append(colours, primary)
	colours = append(colours, secondary1)
	colours = append(colours, secondary2)
	colours = append(colours, complementary)
	return colours
}

func (c *ColourSchemeGenerator) GenerateAnalogic(colour Colour) []SchemeColour {
	colours := make([]SchemeColour, 0, 3)

	primary := SchemeColour{colour, PrimaryColour, nil, nil}
	primary.Tints = c.GenerateTints(primary.Colour)
	primary.Shades = c.GenerateShades(primary.Colour)

	secondaryColour1 := NewColour(colour.Hue, colour.Saturation, colour.Value)
	secondaryColour1 = secondaryColour1.AddHue(c.angle)
	secondary1 := SchemeColour{secondaryColour1, SecondaryColour, nil, nil}
	secondary1.Tints = c.GenerateTints(secondary1.Colour)
	secondary1.Shades = c.GenerateShades(secondary1.Colour)

	secondaryColour2 := NewColour(colour.Hue, colour.Saturation, colour.Value)
	secondaryColour2 = secondaryColour2.AddHue(-c.angle)
	secondary2 := SchemeColour{secondaryColour2, SecondaryColour, nil, nil}
	secondary2.Tints = c.GenerateTints(secondary2.Colour)
	secondary2.Shades = c.GenerateShades(secondary2.Colour)

	colours = append(colours, primary)
	colours = append(colours, secondary1)
	colours = append(colours, secondary2)
	return colours
}

func (c *ColourSchemeGenerator) GenerateAccentedAnalogic(colour Colour) []SchemeColour {
	colours := make([]SchemeColour, 0, 4)

	primary := SchemeColour{colour, PrimaryColour, nil, nil}
	primary.Tints = c.GenerateTints(primary.Colour)
	primary.Shades = c.GenerateShades(primary.Colour)

	secondaryColour1 := NewColour(colour.Hue, colour.Saturation, colour.Value)
	secondaryColour1 = secondaryColour1.AddHue(c.angle)
	secondary1 := SchemeColour{secondaryColour1, SecondaryColour, nil, nil}
	secondary1.Tints = c.GenerateTints(secondary1.Colour)
	secondary1.Shades = c.GenerateShades(secondary1.Colour)

	secondaryColour2 := NewColour(colour.Hue, colour.Saturation, colour.Value)
	secondaryColour2 = secondaryColour2.AddHue(-c.angle)
	secondary2 := SchemeColour{secondaryColour2, SecondaryColour, nil, nil}
	secondary2.Tints = c.GenerateTints(secondary2.Colour)
	secondary2.Shades = c.GenerateShades(secondary2.Colour)

	complementaryColour := NewColour(colour.Hue, colour.Saturation, colour.Value)
	complementaryColour = secondaryColour2.AddHue(180.0)
	complementary := SchemeColour{complementaryColour, ComplementaryColour, nil, nil}
	complementary.Tints = c.GenerateTints(complementary.Colour)
	complementary.Shades = c.GenerateShades(complementary.Colour)

	colours = append(colours, primary)
	colours = append(colours, secondary1)
	colours = append(colours, secondary2)
	colours = append(colours, complementary)
	return colours
}
