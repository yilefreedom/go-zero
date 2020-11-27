package javagen

import (
	"fmt"
	"io"
	"path"
	"strings"
	"text/template"

	"github.com/yilefreedom/go-zero/tools/goctl/api/spec"
	apiutil "github.com/yilefreedom/go-zero/tools/goctl/api/util"
	"github.com/yilefreedom/go-zero/tools/goctl/util"
)

const (
	componentTemplate = `// Code generated by goctl. DO NOT EDIT.
package com.xhb.logic.http.packet.{{.packet}}.model;

import com.xhb.logic.http.DeProguardable;

{{.componentType}}
`
)

func genComponents(dir, packetName string, api *spec.ApiSpec) error {
	types := apiutil.GetSharedTypes(api)
	if len(types) == 0 {
		return nil
	}
	for _, ty := range types {
		if err := createComponent(dir, packetName, ty); err != nil {
			return err
		}
	}

	return nil
}

func createComponent(dir, packetName string, ty spec.Type) error {
	modelFile := util.Title(ty.Name) + ".java"
	filename := path.Join(dir, modelDir, modelFile)
	if err := util.RemoveOrQuit(filename); err != nil {
		return err
	}

	fp, created, err := apiutil.MaybeCreateFile(dir, modelDir, modelFile)
	if err != nil {
		return err
	}
	if !created {
		return nil
	}
	defer fp.Close()

	tys, err := buildType(ty)
	if err != nil {
		return err
	}

	t := template.Must(template.New("componentType").Parse(componentTemplate))
	return t.Execute(fp, map[string]string{
		"componentType": tys,
		"packet":        packetName,
	})
}

func buildType(ty spec.Type) (string, error) {
	var builder strings.Builder
	if err := writeType(&builder, ty); err != nil {
		return "", apiutil.WrapErr(err, "Type "+ty.Name+" generate error")
	}
	return builder.String(), nil
}

func writeType(writer io.Writer, tp spec.Type) error {
	fmt.Fprintf(writer, "public class %s implements DeProguardable {\n", util.Title(tp.Name))
	for _, member := range tp.Members {
		if err := writeProperty(writer, member, 1); err != nil {
			return err
		}
	}
	genGetSet(writer, tp, 1)
	fmt.Fprintf(writer, "}\n")
	return nil
}
