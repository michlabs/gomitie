package main 

import (
	"log"
	"flag"
	"fmt"
	"os"

	"github.com/michlabs/gomitie/ner"
)

var modelFP string
var trainingFP string
var fileFormat string
var outputFP string

func init() {
	flag.StringVar(&modelFP, "model", "", "path to total word feature model file")
	flag.StringVar(&trainingFP, "input", "", "path to the training file")
	flag.StringVar(&fileFormat, "format", "stanford")
	flag.StringVar(&outputFP, "output", "./ner_model.dat", "where to save output model")
	flag.Parse()
}

func main() {
	if modelFP == "" {
		log.Println("Error: Model file is required but empty")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if trainingFP == "" {
		log.Println("Error: Training file is required but empty")
		flag.PrintDefaults()
		os.Exit(1)
	}

	
}