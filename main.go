// Use cobra cli or bubbletea
package main

import (
	"fmt"
)

func high_pass_filter_calculator() {
	fmt.Println("\n[High Pass Filter Calculator]")
	fmt.Println("Resistance (R)")
}

func fraction_to_percent_calculator() {
	fmt.Println("Fraction to percent")
}

func decimal_to_percent_calculator() {
	fmt.Println("Decimal to percent")
}

func main() {
	// Create a dictionary of domain/calculators
	calculators := make(map[string][]string)
	calculators["Physics"] = []string{
		"high_pass_filter_calculator",
	}
	calculators["Biology"] = []string{
		"mass",
		"gravity",
	}
	i := 1
	fmt.Println("Select a domain:")
	for key := range calculators {
		fmt.Printf("%d - %s\n", i, key)
		i++
	}

	var choice string

	fmt.Scanln(&choice)

	if choice == "1" {
		choice = "Physics"
	}
	if choice == "2" {
		choice = "Biology"
	}

	fmt.Printf("Choose a %s formula to calculate:\n", choice)
	for key, value := range calculators[choice] {
		fmt.Printf("%d - %s\n", key+1, value)
	}
	var formula_idx int
	for {
		_, err := fmt.Scanln(&formula_idx)
		if err == nil && formula_idx > 0 && formula_idx < 2 {
			break
		}
		fmt.Println("Please give a proper answer:")
	}

	formula := calculators[choice][formula_idx-1]

	switch formula {
	case calculators[choice][0]:
		high_pass_filter_calculator()
	}
}
