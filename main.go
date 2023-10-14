package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/ananrafs/descartes/core"
	"github.com/ananrafs/descartes/law"
)

func main() {
	folderLocation := flag.String("folder", "./dump/law_2", "folder location")
	factFile := flag.String("fact", "fact", "fact file name")
	lawFile := flag.String("law", "law", "law file name")
	outFile := flag.String("out", "output", "output file name")

	// Parse the command-line arguments
	flag.Parse()

	core.InitRule(core.WithDefaultRules)
	core.InitEvaluator(core.WithDefaultEvaluators)
	core.InitActions(core.WithDefaultActions)

	l, err := law.CreateLaw(getStringFromFile(fmt.Sprintf("%s/%s.json", *folderLocation, *lawFile)))
	if err != nil {
		panic(err)
	}
	err = core.Register(l)
	if err != nil {
		panic(err)
	}

	f, err := law.CreateMultipleFact(getStringFromFile(fmt.Sprintf("%s/%s.json", *folderLocation, *factFile)))
	if err != nil {
		panic(err)
	}
	var responses []interface{}
	for _, fact := range f {
		res, err := core.Eval(fact)
		if err != nil {
			responses = append(responses, err)
			continue
		}
		responses = append(responses, res)
	}

	writeToFile(fmt.Sprintf("%s/%s.json", *folderLocation, *outFile), responses)

}

func getFromFile(fileLocation string) ([]byte, error) {
	jsonFile, err := os.Open(fileLocation)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	return io.ReadAll(jsonFile)

}

func getStringFromFile(fileLocation string) string {
	strByte, err := getFromFile(fileLocation)
	if err != nil {
		panic(err)
	}

	return string(strByte)
}

func writeToFile(fileLocation string, obj interface{}) {
	// Marshal the interface{} to JSON.
	jsonData, err := json.MarshalIndent(obj, "", " ")
	if err != nil {
		fmt.Println("Error marshaling data to JSON:", err)
		return
	}

	// Write the JSON data to a file.
	err = os.WriteFile(fileLocation, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return
	}
}
