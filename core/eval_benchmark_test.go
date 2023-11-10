package core_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/ananrafs/descartes/core"
	"github.com/ananrafs/descartes/law"
)

var tableTest = []struct {
	fileLocation string
}{
	{fileLocation: "dynamic"},
	{fileLocation: "dynamic cache"},
	{fileLocation: "static"},
	{fileLocation: "static cache"},
	{fileLocation: "test_slice"},
	{fileLocation: "test_slice_struct"},
	// {fileLocation: "cache"},
	// {fileLocation: "rule_random"},
	// {fileLocation: "test_actiongroup"},
}

func BenchmarkJudgeLaw_Test(b *testing.B) {
	core.InitFactory(core.WithDefaults())
	for _, tt := range tableTest {
		l, err := law.CreateLaw(getStringFromFile(fmt.Sprintf("../dump/%s/law.json", tt.fileLocation)))
		if err != nil {
			fmt.Println(err)
			continue
		}
		core.Register(l)

		f, err := law.CreateMultipleFact(getStringFromFile(fmt.Sprintf("../dump/%s/fact.json", tt.fileLocation)))
		if err != nil {
			fmt.Println(err)
			continue
		}

		b.Run(tt.fileLocation, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				core.Eval(f[0])
			}
		})
	}
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
