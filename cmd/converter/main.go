package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/mourice12/hl7-to-fhir/internal/converter"
	"github.com/mourice12/hl7-to-fhir/internal/hl7"
)

func main() {
	//sample HL7 message

	inputFile := flag.String("input", "", "input HL7FilePath")
	outputFile := flag.String("output", "", "Output FHIR JSON File Path")
	flag.Parse()

	//validate input
	if *inputFile == "" {
		fmt.Println("Usage: converter input <file.hl7> [-output <file.json]")
		os.Exit(1)
	}

	//read input file
	data, err := os.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	//Parse HL7
	msg, err := hl7.Parse(string(data))
	if err != nil {
		fmt.Printf("Error Parsing Hl7: %v\n", err)
		os.Exit(1)
	}

	//Convert to bundle instead of just patient
	bundle, err := converter.ConvertToBundle(msg)
	if err != nil {
		fmt.Printf("Error converting: %v\n", err)
		os.Exit(1)
	}
	//convert to FHIR
	//patient, err := converter.ConvertToPatient(msg)
	//if err != nil {
	//	fmt.Printf("Error converting JSON: %v\n", err)
	///	os.Exit(1)
	//}

	//Marshal JSON
	jsonBytes, err := json.MarshalIndent(bundle, "", " ")
	if err != nil {
		fmt.Printf("Error Marshaling JSON: %v\n", err)
	}

	//Output
	if *outputFile != "" {
		err = os.WriteFile(*outputFile, jsonBytes, 0644)
		if err != nil {
			fmt.Printf("Error writing file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("FHIR output written to %s\n", *outputFile)
	} else {
		fmt.Println(string(jsonBytes))
	}
}
