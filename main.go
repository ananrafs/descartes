package main

import (
	"fmt"
	"io"
	"os"

	"github.com/ananrafs/descartes/core"
	"github.com/ananrafs/descartes/law"
)

func main() {

	core.InitRule(core.WithDefaultRules)
	core.InitEvaluator(core.WithDefaultEvaluators)

	l, err := law.CreateLaw(getStringFromFile("./dump/rule_random/law.json"))
	if err != nil {
		panic(err)
	}
	err = core.Register(l)
	if err != nil {
		panic(err)
	}

	f, err := law.CreateFact(getStringFromFile("./dump/rule_random/fact_1.json"))
	if err != nil {
		panic(err)
	}
	res, err := core.Eval(f)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

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
