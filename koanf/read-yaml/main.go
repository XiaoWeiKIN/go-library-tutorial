package main

import (
	"fmt"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/maps"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

var k = koanf.New(".")

func main() {
	k.Load(file.Provider("mock/mock1.yaml"), yaml.Parser())

	for k, v := range maps.Unflatten(k.All(), ".") {
		fmt.Println(k)
		fmt.Println(v)
		fmt.Println("-------")
	}
}
