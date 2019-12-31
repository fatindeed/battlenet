<?php

class ApiDoc
{
    /**
     * API version.
     *
     * @var string
     */
    private $version = '1.0.0';

    /**
     * The OpenApi Specification (OAS) version.
     *
     * @var string
     */
    private $oas = '3.0.1';

    /**
     * @var string
     */
    private $content;

    /**
     * Generate api doc.
     *
     * @param stdClass $page
     *
     * @return void
     */
    public function generate(stdClass $page): void
    {
        // Init schema
        $data = [
            'openapi' => $this->oas,
            'info' => [
                'version' => $this->version,
                'title' => $page->title,
                'description' => $page->cardDescription,
                'contact' => [
                    'name' => 'James Zhu',
                    'email' => 'fatindeed@hotmail.com'
                ]
            ],
            'externalDocs' => [
                'url' => 'https://develop.battlenet.com.cn/' . $page->path,
                'description' => $page->cardTitle
            ],
            'servers' => [
                [
                    'url' => 'http://{region}.api.blizzard.com',
                    'variables' => [
                        'region' => [
                            'enum' => ['us', 'eu', 'kr', 'tw'],
                            'default' => 'us'
                        ]
                    ]
                ],
                // TODO: cn domain
            ],
            'paths' => [],
            'components' => [
                'parameters' => [
                    // 'region' => [
                    //     'in' => 'query',
                    //     'name' => ':region',
                    //     'schema' => ['type' => 'string'],
                    //     'required' => true,
                    //     'description' => 'The region of the data to retrieve.',
                    //     'example' => 'us',
                    // ],
                    'namespace' => [
                        'in' => 'query',
                        'name' => 'namespace',
                        'schema' => ['type' => 'string'],
                        'required' => true,
                        'description' => 'The namespace to use to locate this document.',
                        'example' => 'static-classic-us',
                    ],
                    'locale' => [
                        'in' => 'query',
                        'name' => 'locale',
                        'schema' => ['type' => 'string'],
                        'description' => 'The locale to reflect in localized data.',
                        'example' => 'en_US',
                    ]
                ]
            ],
            'tags' => []
        ];
        // Load api documentation
        $json = file_get_contents('https://develop.battlenet.com.cn/api/data/content/' . $page->path . '.json');
        if (! $json) {
            trigger_error('Failed to read ' . $page->path, E_USER_ERROR);
        }
        $schema = json_decode($json);
        foreach ($schema->resources as $i => $resource) {
            $data['tags'][$i] = ['name' => $resource->name];
            foreach ($resource->methods as $method) {
                $schema = [
                    'tags' => [$resource->name],
                    'summary' => $method->name,
                    'description' => $method->description,
                    'operationId' => lcfirst(str_replace(' ', '', $method->name)),
                    'parameters' => [],
                    'responses' => [
                        200 => ['description' => 'OK']
                    ]
                ];
                foreach ($method->parameters as $j => $parameter) {
                    if (isset($data['components']['parameters'][$parameter->name])) {
                        $schema['parameters'][$j] = ['$ref' => '#/components/parameters/' . $parameter->name];
                    } else {
                        $schema['parameters'][$j] = [
                            'in' => (strpos($method->path, $parameter->name) !== false) ? 'path' : 'query',
                            'name' => trim($parameter->name, '{}'),
                            'schema' => ['type' => $parameter->type],
                            'required' => $parameter->required,
                            'description' => $parameter->description
                        ];
                        if (property_exists($parameter, 'defaultValue')) {
                            $schema['parameters'][$j]['example'] = $parameter->defaultValue;
                        }
                    }
                }
                $data['paths'][$method->path][strtolower($method->httpMethod)] = $schema;
            }
        }
        $this->content = trim(trim(yaml_emit($data)), '-.');
    }

    /**
     * Publish api to swagger hub.
     * 
     * @param string $owner
     * @param string $api
     *
     * @return void
     *
     * @see https://app.swaggerhub.com/apis-docs/swagger-hub/registry-api/1.0.47
     */
    public function publish(string $owner, string $api): void
    {
        $query = [
            'isPrivate' => 'false',
            'version' => $this->version,
            'oas' => $this->oas,
            'force' => 'true'
        ];
        $ch = curl_init();
        curl_setopt($ch, CURLOPT_URL, 'https://api.swaggerhub.com/apis/' . $owner. '/' . $api. '?' . http_build_query($query));
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        curl_setopt($ch, CURLOPT_POST, true);
        curl_setopt($ch, CURLOPT_POSTFIELDS, $this->content);
        curl_setopt($ch, CURLOPT_HTTPHEADER, [
            'Authorization: ' . getenv('API_KEY'),
            'Content-Type: application/yaml',
            'Accept: application/json'
        ]);
        $content = curl_exec($ch);
        if (curl_errno($ch)) {
            trigger_error('Curl Error#' . curl_errno($ch) . ': ' . curl_error($ch), E_USER_ERROR);
        }
        $code = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        if ($code >= 400) {
            $result = json_decode($content);
            trigger_error('Swaggerhub Response Error#' . $result->code . ': ' . $result->message, E_USER_ERROR);
        }
        curl_close($ch);
    }
}
