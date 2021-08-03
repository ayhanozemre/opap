package main

const HttpAdditionalInformation = "header"
const ServiceReplacementKey = "<service_name>"
const SwaggerQueryStringKey = "path"

type Service struct {
	Name                   string   `json:"name"`
	Url                    string   `json:"swagger_url"`
	InternalAddress        string   `json:"internal_address"`
	KrakendOutputEncoding  string   `json:"output_encoding"`
	KrakendBackendEncoding string   `json:"backend_encoding"`
	ExcludedEndpoints      []string `json:"excluded_endpoints"`
	CrossMapping           []Router `json:"cross_mapping"`
}

type Router struct {
	From            string   `json:"from"`
	To              string   `json:"to"`
	Method          string   `json:"method"`
	InternalAddress string   `json:"internal_address"`
	HeadersToPass   []string `json:"headers_to_pass"`
}

type ConfigModel struct {
	Services        []Service `json:"services"`
	FilePrefix      string    `json:"file_prefix"`
	PrefixSeparator string    `json:"prefix_separator"`
	BaseDir         string    `json:"base_dir"`
}

type SwaggerMethodParameter struct {
	Name string `json:"name"`
	In   string `json:"in"`
}

type SwaggerMethod struct {
	Security   []map[string]interface{} `json:"security"`
	Parameters []SwaggerMethodParameter `json:"parameters"`
}

type SwaggerSecurityDefinition struct {
	Name string `json:"name"`
	In   string `json:"in"`
}

type SwaggerResponse struct {
	Paths               map[string]map[string]SwaggerMethod  `json:"paths"`
	SecurityDefinitions map[string]SwaggerSecurityDefinition `json:"securityDefinitions"`
}

type KrakendBackend struct {
	Encoding   string   `json:"encoding,omitempty"`
	URLPattern string   `json:"url_pattern"`
	Host       []string `json:"host"`
}

type KrakendPath struct {
	Endpoint          string           `json:"endpoint,omitempty"`
	HeadersToPass     []string         `json:"headers_to_pass,omitempty"`
	QuerystringParams []string         `json:"querystring_params,omitempty"`
	OutputEncoding    string           `json:"output_encoding,omitempty"`
	Method            string           `json:"method,omitempty"`
	Backend           []KrakendBackend `json:"backend"`
}
