package battlenetdev

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// NewRequest function
func NewRequest(path string) (map[string]interface{}, error) {
	resp, err := http.Get("https://develop.battlenet.com.cn/api/data/" + path + ".json")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("response error: %s", resp.Status)
	}
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return result, nil
}

// GetOpenAPISpecs function
func GetOpenAPISpecs(cardPage map[string]interface{}) Specification {
	spec := Specification{
		OpenAPI: oasVersion,
		Info: Info{
			Version:     apiVersion,
			Title:       cardPage["title"].(string),
			Description: cardPage["cardDescription"].(string),
			Contact: Contact{
				Name:  "James Zhu",
				Email: "fatindeed@hotmail.com",
			},
		},
		ExternalDocs: ExternalDocs{
			URL:         "https://develop.battlenet.com.cn/" + cardPage["path"].(string),
			Description: cardPage["cardTitle"].(string),
		},
		Servers: []Server{
			Server{
				URL: "http://{region}.api.blizzard.com",
				Variables: map[string]Variable{
					"region": Variable{
						Enum:    []string{"us", "eu", "kr", "tw"},
						Default: "us",
					},
				},
			},
		},
		Paths: map[string]map[string]Operation{},
		Components: Component{
			Parameters: map[string]Parameter{
				"namespace": Parameter{
					In:          "query",
					Name:        "namespace",
					Schema:      Schema{Type: "string"},
					Required:    true,
					Description: "The namespace to use to locate this document.",
					Example:     "static-classic-us",
				},
				"locale": Parameter{
					In:          "query",
					Name:        "locale",
					Schema:      Schema{Type: "string"},
					Required:    false,
					Description: "The locale to reflect in localized data.",
					Example:     "en_US",
				},
			},
		},
		// Tags: []Tag{},
	}
	// Load api documentation
	api, err := NewRequest("content/" + cardPage["path"].(string))
	if err != nil {
		log.Fatal(err)
	}
	resources := api["resources"].([]interface{})
	for _, resource := range resources {
		resource := resource.(map[string]interface{})
		spec.Tags = append(spec.Tags, Tag{Name: resource["name"].(string)})
		methods := resource["methods"].([]interface{})
		for _, method := range methods {
			method := method.(map[string]interface{})
			operation := Operation{
				Tags:        []string{resource["name"].(string)},
				Summary:     method["name"].(string),
				Description: method["description"].(string),
				OperationID: strings.Replace(method["name"].(string), " ", "", -1),
				// Parameters:  []Parameter{},
				Responses: map[string]Response{
					"200": Response{Description: "OK"},
				},
			}
			parameters := method["parameters"].([]interface{})
			for _, param := range parameters {
				param := param.(map[string]interface{})
				var parameter Parameter
				paramName := param["name"].(string)
				if _, ok := spec.Components.Parameters[paramName]; ok {
					parameter.Ref = "#/components/parameters/" + paramName
				} else {
					if strings.Index(method["path"].(string), paramName) > 0 {
						parameter.In = "path"
					} else {
						parameter.In = "query"
					}
					parameter.Name = strings.Trim(paramName, "{}")
					parameter.Schema.Type = param["type"].(string)
					parameter.Required = param["required"].(bool)
					parameter.Description = param["description"].(string)
					parameter.Example = param["defaultValue"]
				}
				operation.Parameters = append(operation.Parameters, parameter)
			}
			if spec.Paths[method["path"].(string)] == nil {
				spec.Paths[method["path"].(string)] = map[string]Operation{}
			}
			httpMethod := strings.ToLower(method["httpMethod"].(string))
			spec.Paths[method["path"].(string)][httpMethod] = operation
		}
	}
	return spec
}
