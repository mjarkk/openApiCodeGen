package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type Route struct {
	HasBodyData   bool
	URL           string
	ReadableRoute string
	Method        string
	URLParms      []string
}

type SwaggerResponse struct {
	Ref string `json:"$ref"`
}

type SwaggerParam struct {
	Type     string `json:"type"`
	Format   string `json:"format"`
	Name     string `json:"name"`
	In       string `json:"in"`
	Required bool   `json:"required"`
}

type SwaggerPath struct {
	Tags       []string                   `json:"tags"`
	Parameters []SwaggerParam             `json:"parameters"`
	Responses  map[string]SwaggerResponse `json:"responses"`
}

type swaggerFile struct {
	BasePath string                            `json:"basePath"`
	Paths    map[string]map[string]SwaggerPath `json:"paths"` // The first map is the path (/some/route, /blog/{id}), the second map is the method (GET, POST, PATCH, ...)
}

var in = flag.String("in", "swagger.json", "The input swagger file (only supports json)")

func init() {
	flag.Parse()
}

func main() {
	// Open the open api file and parse it
	swaggerData, err := ioutil.ReadFile(*in)
	if err != nil {
		log.Fatal(err)
	}
	var data swaggerFile
	err = json.Unmarshal(swaggerData, &data)
	if err != nil {
		log.Fatal(err)
	}

	// Parse the output of the json file
	routes := []Route{}
	for path, methods := range data.Paths {
		for method, route := range methods {
			params, readableRoute, jsPath := getParamsFromPath(data.BasePath, path)

			newRoute := Route{
				Method:        strings.ToUpper(method),
				URL:           jsPath,
				URLParms:      params,
				HasBodyData:   hasBodyData(route.Parameters),
				ReadableRoute: readableRoute,
			}

			routes = append(routes, newRoute)
		}
	}

	// Print some header information
	fmt.Println(`import { apiFetcher } from "./apiUtil.js";`)
	fmt.Println()

	// Loop over the routes and print them out as javascript
	for _, route := range routes {
		fmt.Println()
		fmt.Print("export const " + route.ReadableRoute + route.Method + " = (")

		jsParams := ""
		for _, param := range route.URLParms {
			if len(jsParams) > 0 {
				jsParams += ", "
			}
			jsParams += "param" + firstLetterUpper(param)
		}
		fmt.Print(jsParams)

		if route.HasBodyData {
			if len(jsParams) > 0 {
				fmt.Print(", ")
			}
			fmt.Print("body")
		}

		fmt.Println(") => apiFetcher({")
		fmt.Println("  params: [" + jsParams + "],")
		fmt.Println("  method: \"" + route.Method + "\",")
		fmt.Println("  url: \"" + route.URL + "\",")
		if route.HasBodyData {
			fmt.Println("  body,")
		}
		fmt.Println("})")
	}
}

func hasBodyData(in []SwaggerParam) bool {
	for _, item := range in {
		if strings.ToLower(item.In) == "body" {
			return true
		}
	}
	return false
}

func getParamsFromPath(base, path string) (routeParams []string, readableRouteName, jsRoute string) {
	pathParts := strings.Split(path, "/")

	params := []string{}
	readableRoute := []string{}
	outPath := []string{}

	for _, part := range pathParts {
		open := strings.Index(part, "{")
		close := strings.Index(part, "}")
		if open < 0 || close < 0 || close < open {
			outPath = append(outPath, part)
			readableRoute = append(readableRoute, firstLetterUpper(part))
			continue
		}

		params = append(params, part[open+1:close])
		outPath = append(outPath, part[:open]+"$"+part[open:])
		readableRoute = append(readableRoute, firstLetterUpper(part[open+1:close]))
	}

	return params, strings.Join(readableRoute, ""), joinRoutes(base, strings.Join(outPath, "/"))
}

func joinRoutes(r1, r2 string) string {
	if len(r1) == 0 || len(r2) == 0 || (r1[len(r1)-1] == '/' && r2[0] != '/') || (r1[len(r1)-1] != '/' && r2[0] == '/') {
		return r1 + r2
	}

	if r1[len(r1)-1] == '/' && r2[0] == '/' {
		return r1 + r2[1:]
	}

	return r1 + "/" + r2
}

func firstLetterUpper(in string) string {
	if len(in) == 0 {
		return in
	}

	return strings.ToUpper(string(in[0])) + in[1:]
}
