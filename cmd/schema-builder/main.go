package main

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/mxcd/broke/pkg/config"
)

func main() {
	jsonSchema := jsonschema.Reflect(&config.BrokeConfig{})
	data, err := jsonSchema.MarshalJSON()
	if err != nil {
		panic(err)
	}

	var obj map[string]interface{}
	json.Unmarshal(data, &obj)

	indentedData, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		panic(err)
	}

	file, err := os.Create("pkg/config/config-schema.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	w.Write(indentedData)
	w.Flush()
}
