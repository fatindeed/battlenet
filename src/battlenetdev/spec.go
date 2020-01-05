package battlenetdev

const (
	// API version.
	apiVersion = "1.0.0"
	// The OpenApi Specification (OAS) version.
	oasVersion = "3.0.2"
)

// Specification type
type Specification struct {
	OpenAPI      string                          `yaml:"openapi"`
	Info         Info                            `yaml:"info"`
	ExternalDocs ExternalDocs                    `yaml:"externalDocs"`
	Servers      []Server                        `yaml:"servers"`
	Paths        map[string]map[string]Operation `yaml:"paths"`
	Components   Component                       `yaml:"components"`
	Tags         []Tag                           `yaml:"tags"`
}

// Info type
type Info struct {
	Version        string  `yaml:"version,omitempty"`
	Title          string  `yaml:"title,omitempty"`
	Description    string  `yaml:"description,omitempty"`
	Contact        Contact `yaml:"contact,omitempty"`
	TermsOfService string  `yaml:"termsOfService,omitempty"`
	// License        License `yaml:"license,omitempty"`
}

// Contact type
type Contact struct {
	Name  string `yaml:"name,omitempty"`
	Email string `yaml:"email,omitempty"`
	URL   string `yaml:"url,omitempty"`
}

// ExternalDocs type
type ExternalDocs struct {
	URL         string `yaml:"url,omitempty"`
	Description string `yaml:"description,omitempty"`
}

// Server type
type Server struct {
	URL       string              `yaml:"url,omitempty"`
	Variables map[string]Variable `yaml:"variables,omitempty"`
}

// Variable type
type Variable struct {
	Enum    []string `yaml:"enum,omitempty"`
	Default string   `yaml:"default,omitempty"`
}

// Operation type
type Operation struct {
	Tags        []string            `yaml:"tags,omitempty"`
	Summary     string              `yaml:"summary,omitempty"`
	Description string              `yaml:"description,omitempty"`
	OperationID string              `yaml:"operationId,omitempty"`
	Parameters  []Parameter         `yaml:"parameters,omitempty"`
	Responses   map[string]Response `yaml:"responses,omitempty"`
}

// Parameter type
type Parameter struct {
	Ref         string      `yaml:"$ref,omitempty"`
	In          string      `yaml:"in,omitempty"`
	Name        string      `yaml:"name,omitempty"`
	Schema      Schema      `yaml:"schema,omitempty"`
	Required    bool        `yaml:"required,omitempty"`
	Description string      `yaml:"description,omitempty"`
	Example     interface{} `yaml:"example,omitempty"`
}

// Response type
type Response struct {
	Description string `yaml:"description,omitempty"`
}

// Component type
type Component struct {
	Parameters map[string]Parameter `yaml:"parameters"`
}

// Schema type
type Schema struct {
	Type string `yaml:"type"`
}

// Tag type
type Tag struct {
	Name         string       `yaml:"name,omitempty"`
	Description  string       `yaml:"description,omitempty"`
	ExternalDocs ExternalDocs `yaml:"externalDocs,omitempty"`
}
