package main

type SampleType string

const (
	single_stranded_dna SampleType = "Single Stranded DNA"
	double_stranded_dna SampleType = "Double Stranded DNA"
)

// C – Concentration of the nucleic acid in the sample.
// ​
// A260  – The maximum absorbance as indicated by the spectrophotometric reading. This usually occurs at the wavelength of 260 nm, but it may change depending on the nucleotide. So, if you wondered, why is 260 nm used for DNA?, this is the answer.
//
// l – Pathlength, and more precisely, the length of the cuvette used. The standard value is 1 cm, but your instrument may use a different size.
//
// DF – Dilution factor. It applies only when the sample is diluted. For instance, if you diluted 1 liter of sample in 50 liters of H2O, the dilution factor would be 50. The dilution factor calculator can help you determine the right value.
//
// CF – Conversion factor, which depends on the sample type:
//
// 33 µg/mL for single-stranded DNA (ssDNA).
//
// 50 µg/mL for double-stranded DNA (dsDNA).
//
// 40 µg/mL for RNA.
func dna_concentration(sample_type SampleType, conversion_factor float64, absorbance_at_max float64, pathlength float64, dilution_factor float64, concentration float64) float64 {

	if sample_type == nil && conversion_factor == nil {
		panic("Please provide at least of of these two arguments: Sample type | Conversion factor")
	}

	return (absorbance_at_max / pathlength) * dilution_factor * conversion_factor
}
