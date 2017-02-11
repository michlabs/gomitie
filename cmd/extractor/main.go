package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/michlabs/gomitie"
	"github.com/michlabs/gomitie/ner"
)

var modelFP string
var inputFP string

func init() {
	// /usr/local/share/MITIE-models/english/ner_model.dat
	flag.StringVar(&modelFP, "model", "", "path to MITITE model file")
	flag.StringVar(&inputFP, "input", "", "path to the input file")
	flag.Parse()
}

func main() {
	if modelFP == "" {
		log.Println("Error: Model file is required but empty")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if inputFP == "" {
		log.Println("Error: Input file is required but empty")
		flag.PrintDefaults()
		os.Exit(1)
	}

	log.Println("Loading model...")
	ext, err := ner.NewExtractor(modelFP)
	if err != nil {
		log.Fatal(err)
	}
	defer ext.Free()

	log.Println("Tags: %+v", ext.Tags())

	log.Println("Reading input file...")
	txt, err := ioutil.ReadFile(inputFP)
	if err != nil {
		log.Fatal(err)
	}

	tokens := gomitie.Tokenize(string(txt))

	log.Println("Extracting entities...")
	es, err := ext.Extract(tokens)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range es {
		log.Printf("%+v", v)
	}
}