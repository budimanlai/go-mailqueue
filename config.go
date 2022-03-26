package main

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/eqto/go-json"
)

var (
	config json.Object
)

func InitConfig() {
	var err1 error
	config, err1 = NewConfig(`config/main-local.json`)
	if err1 != nil {
		panic(err1)
	}
}

func NewConfig(file string) (json.Object, error) {
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		return nil, errors.New(`File '` + file + `' not found`)
	}

	jsonFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	obj, err := json.Parse(byteValue)
	if err != nil {
		return nil, err
	}

	return obj, nil
}
