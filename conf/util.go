// Copyright (c) 2016, German Neuroinformatics Node (G-Node),
//                     Adrian Stoewer <adrian.stoewer@rz.ifi.lmu.de>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted under the terms of the BSD License. See
// LICENSE file in the root of the Project.

package conf

import (
	"fmt"
	"html/template"
	"net/url"
	"path"
)

// MakeUrl makes a URL for other web resources provided by gin-auth using
// the base url from the server config file.
func MakeUrl(pathFormat string, param ...interface{}) string {
	baseUrl := GetServerConfig().BaseURL
	for i, s := range param {
		switch s := s.(type) {
		case string:
			param[i] = url.QueryEscape(s)
		}
	}
	pathFormat = fmt.Sprintf(pathFormat, param...)
	return baseUrl + pathFormat
}

// MakeTemplate loads a template using the default layout and the given content template file.
func MakeTemplate(name string) *template.Template {
	layout := path.Join(resourcesPath, "templates", "layout.html")
	content := path.Join(resourcesPath, "templates", name)
	tmpl, err := template.ParseFiles(layout, content)
	if err != nil {
		panic(err)
	}
	return tmpl
}
