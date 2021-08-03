package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
)

func loadConfig(configFile string) (config ConfigModel, err error) {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal([]byte(file), &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func inStringSlice(item string, slice []string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func getFilePath(baseDir string, filePrefix string, separator string, serviceName string) string {
	prefix := ""
	if filePrefix != "" {
		prefix = fmt.Sprintf("%s%s", filePrefix, separator)
	}

	fileName := fmt.Sprintf("%s%s.json", prefix, serviceName)
	fPath := path.Join(path.Dir(baseDir), fileName)
	return fPath
}

func doRequest(url string, response interface{}) (err error) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	_ = json.Unmarshal(body, &response)
	return err
}

func transformSwaggerResponse(serviceConfig Service, response SwaggerResponse) (paths []KrakendPath) {
	var headerToPass []string
	for _, v := range response.SecurityDefinitions {
		if v.In == HttpAdditionalInformation {
			headerToPass = append(headerToPass, v.Name)
		}
	}
	for url, swaggerPaths := range response.Paths {
		if inStringSlice(url, serviceConfig.ExcludedEndpoints) {
			continue
		}
		for method, swaggerMethod := range swaggerPaths {
			var qParams []string
			for _, param := range swaggerMethod.Parameters {
				if param.In == SwaggerQueryStringKey {
					qParams = append(qParams, param.Name)
				}
			}

			kPath := KrakendPath{
				Endpoint:          url,
				OutputEncoding:    serviceConfig.KrakendOutputEncoding,
				QuerystringParams: qParams,
				Method:            strings.ToUpper(method),
				HeadersToPass:     headerToPass,
				Backend: []KrakendBackend{
					{
						Encoding:   serviceConfig.KrakendBackendEncoding,
						URLPattern: url,
						Host:       []string{serviceConfig.InternalAddress},
					},
				},
			}
			paths = append(paths, kPath)
		}
	}
	return paths
}

func crossMapping(serviceConfig Service) (paths []KrakendPath) {
	for _, router := range serviceConfig.CrossMapping {
		from := strings.ReplaceAll(router.From, ServiceReplacementKey, serviceConfig.Name)
		to := strings.ReplaceAll(router.To, ServiceReplacementKey, serviceConfig.Name)
		kPath := KrakendPath{
			Endpoint:       from,
			OutputEncoding: serviceConfig.KrakendOutputEncoding,
			Method:         router.Method,
			HeadersToPass:  router.HeadersToPass,
			Backend: []KrakendBackend{
				{
					Encoding:   serviceConfig.KrakendBackendEncoding,
					URLPattern: to,
					Host:       []string{router.InternalAddress},
				},
			},
		}
		paths = append(paths, kPath)
	}
	return paths
}

func main() {

	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, service := range config.Services {
		var paths []KrakendPath

		response := SwaggerResponse{}
		err := doRequest(service.Url, &response)
		if err != nil {
			log.Println(service.Name, err.Error())
			continue
		}
		swaggerPaths := transformSwaggerResponse(service, response)
		paths = append(paths, swaggerPaths...)

		crossPaths := crossMapping(service)
		paths = append(paths, crossPaths...)

		jsonResult, err := json.Marshal(paths)
		if err != nil {
			log.Fatal("json file creation error ", err)
		}

		endpointScope := string(jsonResult)[1:]
		endpointScope = endpointScope[:len(endpointScope)-1]

		filePath := getFilePath(config.BaseDir, config.FilePrefix, config.PrefixSeparator, service.Name)

		err = ioutil.WriteFile(filePath, []byte(endpointScope), 0644)
		if err != nil {
			log.Panic(err)
		}

	}

}
