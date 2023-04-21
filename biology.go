package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
)

type SampleType string

const (
	ssDNA SampleType = "Single Stranded DNA"
	dsDNA SampleType = "Double Stranded DNA"
	rna   SampleType = "RNA"
)

type DnaConcentration struct {
	sample_type       SampleType
	conversion_factor float64
	absorbance_at_max float64
	pathlength        float64
	dilution_factor   float64
	// concentration     float64
}

/*
Calculate DNA Concentration

C – Concentration of the nucleic acid in the sample.
​
A260  – The maximum absorbance as indicated by the spectrophotometric reading. This usually occurs at the wavelength of 260 nm, but it may change depending on the nucleotide. So, if you wondered, why is 260 nm used for DNA?, this is the answer.

l – Pathlength, and more precisely, the length of the cuvette used. The standard value is 1 cm, but your instrument may use a different size.

DF – Dilution factor. It applies only when the sample is diluted. For instance, if you diluted 1 liter of sample in 50 liters of H2O, the dilution factor would be 50. The dilution factor calculator can help you determine the right value.

CF – Conversion factor, which depends on the sample type:

	33 µg/mL for single-stranded DNA (ssDNA).

	50 µg/mL for double-stranded DNA (dsDNA).

	40 µg/mL for RNA.
*/
func DNAConcentration[T any](p DnaConcentration) float64 {

	if p.sample_type == "" && p.conversion_factor == 0 {
		panic("Please provide at least of of these two arguments: Sample type | Conversion factor")
	} else if p.sample_type == "Single Stranded DNA" {
		p.conversion_factor = 33
	} else if p.sample_type == "Double Stranded DNA" {
		p.conversion_factor = 50
	} else if p.sample_type == "RNA" {
		p.conversion_factor = 40
	}

	return (p.absorbance_at_max / p.pathlength) * p.dilution_factor * p.conversion_factor
}

// Test function
type sTest struct {
	a int
	b int
}

func test(m model) {
	var t textinput.Model

	t = textinput.New()
	t.CursorStyle = cursorStyle

	t.Placeholder = "Example argument 1"
	t.Focus()
	t.PromptStyle = focusedStyle
	t.TextStyle = focusedStyle
	t.CharLimit = 64

	m.Inputs = append(m.Inputs, t)

	t = textinput.New()
	t.CursorStyle = cursorStyle

	t.Placeholder = "Example argument 1"
	t.Focus()
	t.PromptStyle = focusedStyle
	t.TextStyle = focusedStyle
	t.CharLimit = 64

	m.Inputs = append(m.Inputs, t)

	fmt.Println(m.Inputs)
}

// func main() {
// 	param := DNAConcentration{sample_type: "Single Stranded DNA", absorbance_at_max: 5, pathlength: 1, dilution_factor: 1}
// 	var result = dna_concentration(param)
// 	fmt.Println(result)
// }
