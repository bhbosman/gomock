package main

import (
	"fmt"
	"github.com/golang/mock/mockgen/model"
	"strings"
)

func (g *generator) GenerateMockMethodExtensions(mockType string, m *model.Method, pkgOverride, longTp, shortTp string) error {
	argNames := g.getArgNames("", "", m)
	g.p("// argNames: %v", argNames)

	defaultArgs := make([]string, len(argNames))
	for i := 0; i < len(defaultArgs); i++ {
		defaultArgs[i] = "gomock.Any()"
	}
	g.p("// defaultArgs: %v", defaultArgs)
	defaultArgsAsString := strings.Join(defaultArgs, ",")
	g.p("// defaultArgsAsString: %v", defaultArgsAsString)

	argTypes := g.getArgTypes(m, pkgOverride)
	g.p("// argTypes: %v", argTypes)
	argString := makeArgString(argNames, argTypes)
	g.p("// argString: %v", argString)

	rets := make([]string, len(m.Out))
	for i, p := range m.Out {
		rets[i] = p.Type.String(g.packageMap, pkgOverride)
	}
	g.p("// rets: %v", rets)
	retString := strings.Join(rets, ", ")
	g.p("// retString: %v", retString)

	if len(rets) > 1 {
		retString = "(" + retString + ")"
	}
	if retString != "" {
		retString = " " + retString
	}
	g.p("// retString: %v", retString)

	ia := newIdentifierAllocator(argNames)
	g.p("// ia: %v", ia)
	idRecv := ia.allocateIdentifier("mr")
	g.p("// idRecv: %v", idRecv)
	if len(argNames) > 0 {
		g.p("// 1")
		g.p("func (%s *%vMockRecorder%v) On%vDoAndReturn(\n\t%v interface{}, \n\tf func(%v)%v) *gomock.Call {",
			idRecv, mockType, shortTp, m.Name, strings.Join(argNames, ","), argString, retString)
		g.p("return %v.\n\t%v(%v).\n\tDoAndReturn(f)", idRecv, m.Name, strings.Join(argNames, ","))
		g.p("}")
		g.p("")

		g.p("// 1")
		g.p("func (%s *%vMockRecorder%v) On%vDo(\n\t%v interface{}, \n\tf func(%v)) *gomock.Call {",
			idRecv, mockType, shortTp, m.Name, strings.Join(argNames, ","), argString)
		g.p("return %v.\n\t%v(%v).\n\tDo(f)", idRecv, m.Name, strings.Join(argNames, ","))
		g.p("}")
		g.p("")

		g.p("// 1")
		g.p("func (%s *%vMockRecorder%v) On%vDoAndReturnDefault(\n\tf func(%v)%v) *gomock.Call {",
			idRecv, mockType, shortTp, m.Name, argString, retString)
		g.p("return %v.\n\t%v(%v).\n\tDoAndReturn(f)", idRecv, m.Name, defaultArgsAsString)
		g.p("}")
		g.p("")

		g.p("// 1")
		g.p("func (%s *%vMockRecorder%v) On%vDoDefault(\n\tf func(%v)) *gomock.Call {",
			idRecv, mockType, shortTp, m.Name, argString)
		g.p("return %v.\n\t%v(%v).\n\tDo(f)", idRecv, m.Name, defaultArgsAsString)
		g.p("}")
		g.p("")

		if len(rets) > 0 {
			retNames := make([]string, len(rets))
			retArgs := make([]string, len(rets))
			for i, v := range rets {
				retNames[i] = ia.allocateIdentifier(fmt.Sprintf("ret%d", i))
				retArgs[i] = fmt.Sprintf("%v %v", retNames[i], v)
			}

			g.p("// retNames: %v", retNames)
			g.p("// retArgs: %v", retArgs)
			g.p("// retArgs22: %v", strings.Join(retArgs, ","))

			g.p("// 1")
			g.p("func (%s *%vMockRecorder%v) On%vReturn(\n\t%v interface{}, \n\t%v) *gomock.Call {",
				idRecv, mockType, shortTp, m.Name, strings.Join(argNames, ","), strings.Join(retArgs, ","))
			g.p("return %v.\n\t%v(%v).\n\tReturn(%v)", idRecv, m.Name, strings.Join(argNames, ","), strings.Join(retNames, ","))
			g.p("}")

			g.p("// 1")
			g.p("func (%s *%vMockRecorder%v) On%vReturnDefault(\n\t%v) *gomock.Call {",
				idRecv, mockType, shortTp, m.Name, strings.Join(retArgs, ","))
			g.p("return %v.\n\t%v(%v).\n\tReturn(%v)", idRecv, m.Name, defaultArgsAsString, strings.Join(retNames, ","))
			g.p("}")

		}

	} else {
		g.p("// 0")
		g.p("func (%s *%vMockRecorder%v) On%vDoAndReturn(\n\tf func(%v)%v) *gomock.Call {",
			idRecv, mockType, shortTp, m.Name, argString, retString)
		g.p("return %v.\n\t%v().\n\tDoAndReturn(f)", idRecv, m.Name)
		g.p("}")

		g.p("// 0")
		g.p("func (%s *%vMockRecorder%v) On%vDo(\n\tf func(%v)) *gomock.Call {",
			idRecv, mockType, shortTp, m.Name, argString)
		g.p("return %v.\n\t%v().\n\tDoAndReturn(f)", idRecv, m.Name)
		g.p("}")

		if len(rets) > 0 {
			retNames := make([]string, len(rets))
			retArgs := make([]string, len(rets))
			for i, v := range rets {
				retNames[i] = ia.allocateIdentifier(fmt.Sprintf("ret%d", i))
				retArgs[i] = fmt.Sprintf("%v %v", retNames[i], v)
			}

			g.p("// retNames: %v", retNames)
			g.p("// retArgs: %v", retArgs)
			g.p("// retArgs22: %v", strings.Join(retArgs, ","))

			g.p("// 1")
			g.p("func (%s *%vMockRecorder%v) On%vReturn(%v) *gomock.Call {",
				idRecv, mockType, shortTp, m.Name, strings.Join(retArgs, ","))
			g.p("return %v.\n\t%v(%v).\n\tReturn(%v)", idRecv, m.Name, strings.Join(argNames, ","), strings.Join(retNames, ","))
			g.p("}")
		}

	}
	return nil
}

func (g *generator) GenerateMockExtensions(intf *model.Interface, outputPackagePath string) error {
	for _, method := range intf.Methods {
		mockType := g.mockName(intf.Name)
		longTp, shortTp := g.formattedTypeParams(intf, outputPackagePath)

		_ = g.GenerateMockMethodExtensions(mockType, method, outputPackagePath, longTp, shortTp)
	}

	return nil
}
