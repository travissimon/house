package colour

import (
	"log"
	"testing"
)

func Test_ColourSchemeGeneratorHookup(t *testing.T) {
	log.Println("Hookup Succeeded, testing color scheme generator")
}

func Test_Analogic(t *testing.T) {
	generator := NewColourSchemeGenerator()
	generator.SetStrategy(Analogic)

	primary := NewColour(0.0, 1.0, 1.0)

	schemeColours := generator.GetScheme(primary)
	if len(schemeColours) != 3 {
		t.Errorf("Expected 3 colours, actual: %v", len(schemeColours))
	}
}
