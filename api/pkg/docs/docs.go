package docs

import (
    "encoding/json"
    "fmt"
    "reflect"
    "strings"
)

type RouteInfo struct {
    Method      string   `json:"method"`
    Path        string   `json:"path"`
    Handler     string   `json:"handler"`
    Middlewares []string `json:"middlewares"`
    Description string   `json:"description"`
}

var routes []RouteInfo

func RegisterRoute(method, path, handler string, middlewares []string, desc string) {
    routes = append(routes, RouteInfo{
        Method:      method,
        Path:        path,
        Handler:     handler,
        Middlewares: middlewares,
        Description: desc,
    })
}

func GenerateMarkdown() string {
    var builder strings.Builder
    builder.WriteString("# API Documentation\n\n")
    
    for _, route := range routes {
        builder.WriteString(fmt.Sprintf("## %s %s\n", route.Method, route.Path))
        builder.WriteString(fmt.Sprintf("**Handler:** %s\n\n", route.Handler))
        builder.WriteString(fmt.Sprintf("**Description:** %s\n\n", route.Description))
        if len(route.Middlewares) > 0 {
            builder.WriteString(fmt.Sprintf("**Middlewares:** %s\n\n", strings.Join(route.Middlewares, ", ")))
        }
        builder.WriteString("---\n\n")
    }
    
    return builder.String()
}