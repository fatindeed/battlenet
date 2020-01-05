/*----------------------------------------------------------------------------------------
 * Copyright (c) Microsoft Corporation. All rights reserved.
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 *---------------------------------------------------------------------------------------*/

package main

import (
	"battlenetdev"
	"log"
	"os"
	"strings"
	"swaggerhub"

	"gopkg.in/yaml.v3"
)

func main() {
	nav, err := battlenetdev.NewRequest("navigation/documentation")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	children := nav["children"].([]interface{})[0].(map[string]interface{})["children"].([]interface{})
	for _, child := range children {
		sections := child.(map[string]interface{})["page"].(map[string]interface{})["content"].(map[string]interface{})["sections"].([]interface{})
		for _, section := range sections {
			cardPages := section.(map[string]interface{})["cardPages"].([]interface{})
			for _, cardPage := range cardPages {
				cardPage := cardPage.(map[string]interface{})
				if cardPage["contentType"] != nil && cardPage["contentType"].(string) == "api-reference" {
					spec := battlenetdev.GetOpenAPISpecs(cardPage)
					yaml, err := yaml.Marshal(spec)
					if err != nil {
						log.Fatal(err)
						os.Exit(1)
					}
					apiName := strings.Replace(cardPage["path"].(string), "/", "-", -1)
					err = swaggerhub.SaveDefinition("battlenet", apiName, yaml)
					if err != nil {
						log.Fatal(err)
						os.Exit(1)
					}
					log.Println(apiName + " done")
				}
			}
		}
	}
}
