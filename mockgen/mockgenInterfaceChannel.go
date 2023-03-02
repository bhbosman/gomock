package main

import (
	"fmt"
	"github.com/golang/mock/mockgen/model"
	"strings"
)

func (g *generator) GenerateChannelInterface(intf *model.Interface, outputPackagePath string) error {
	//longTp, shortTp := g.formattedTypeParams(intf, outputPackagePath)

	//mockType := g.mockName(intf.Name)

	g.p("")
	g.p("// Interface %v ", intf.Comment)
	g.p("// Interface %v ", outputPackagePath)
	g.p("// Interface %v ", intf.Name)
	_ = g.GenerateChannelMethods(intf, outputPackagePath)
	_ = g.GenerateChannelReceivedEvent(intf, outputPackagePath)
	return nil
}
func (g *generator) GenerateChannelReceivedEvent(intf *model.Interface, outputPackagePath string) error {

	g.p("func ChannelEventsFor%v(next %v, event interface{}) (bool, error){",
		intf.Name,
		intf.NameType.String(g.packageMap, outputPackagePath))
	g.p("switch v := event.(type) {")

	for _, m := range intf.Methods {
		argNames := g.getArgNames("v.inData.", "...", m)

		argOutNames := g.getOutArgNames(m)

		name := fmt.Sprintf("%v%v", intf.Name, m.Name)
		g.p("case *%v:", name)
		g.p("data := %vOut{}", name)
		if len(argOutNames) == 0 {
			g.p("next.%v(%v)", m.Name, strings.Join(argNames, ","))

		} else {
			g.p("%v = next.%v(%v)", strings.Join(argOutNames, ","), m.Name, strings.Join(argNames, ","))
		}

		g.p("if v.outDataChannel != nil{")
		g.p("v.outDataChannel <- data")
		g.p("}")
		g.p("return true, nil")
	}
	g.p("default:")
	g.p("return false, nil")
	g.p("}")
	g.p("}")
	g.p("")
	return nil
}

func (g *generator) GenerateChannelMethods(intf *model.Interface, outputPackagePath string) error {
	for _, m := range intf.Methods {
		g.p("// Interface %v, Method: %v ", intf.Name, m.Name)
		name := fmt.Sprintf("%v%v", intf.Name, m.Name)

		argNames := g.getArgNames("", "", m)
		argTypes := g.getArgTypes(m, outputPackagePath)
		argString := makeArgString(argNames, argTypes)

		rets := make([]string, len(m.Out))
		for i, p := range m.Out {
			rets[i] = p.Type.String(g.packageMap, outputPackagePath)
		}
		retString := strings.Join(rets, ", ")
		if len(rets) > 1 {
			retString = "(" + retString + ")"
		}
		if retString != "" {
			retString = " " + retString
		}

		_ = g.GenerateStructIn(intf, m, outputPackagePath, name, argString, argNames, argTypes)
		_ = g.GenerateStructOut(intf, m, outputPackagePath, name, argString, argNames, argTypes)
		_ = g.GenerateError(intf, m, outputPackagePath, name, argString, argNames, argTypes)
		_ = g.GenerateStruct(intf, m, outputPackagePath, name, argString, argNames, argTypes)
		_ = g.GenerateCallFunction(intf, m, outputPackagePath, name, argString, argNames, argTypes)
	}
	return nil
}

func (g *generator) GenerateCallFunction(
	intf *model.Interface,
	m *model.Method,
	outputPackagePath string,
	name string,
	argString string,
	argNames []string,
	argTypes []string) error {

	g.p("func Call%v(context context.Context, channel chan<- interface{}, waitToComplete bool, %v)(%vOut, error) {", name, argString, name)
	g.p("if context != nil && context.Err() != nil {")
	g.p("return %vOut{}, context.Err()", name)
	g.p("}")
	if len(argNames) == 0 {
		g.p("data := New%v(waitToComplete)", name)
	} else {
		sss := ""
		for i, argName := range argNames {
			s := argTypes[i]
			if len(s) > 3 && s[0:3] == "..." {
				sss = fmt.Sprintf("%v, %v...", sss, argName)
			} else {
				sss = fmt.Sprintf("%v, %v", sss, argName)

			}
		}

		g.p("data := New%v(waitToComplete%v)", name, sss)

	}

	g.p("if waitToComplete {")
	g.p("defer func(data *%v) {", name)
	g.p("err := data.Close()")
	g.p("if err != nil {")
	g.p("}")
	g.p("}(data)")
	g.p("}")

	g.p("if context != nil && context.Err() != nil {")
	g.p("return %vOut{}, context.Err()", name)
	g.p("}")
	g.p("channel <- data")
	g.p("var err error")
	g.p("var v %vOut", name)
	g.p("if waitToComplete{")
	g.p("v, err = data.Wait(func(interfaceName string, methodName string, err error) error {")
	g.p("return err")
	g.p("})")
	g.p("}else{")
	g.p("err =errors.NoWaitOperationError")
	g.p("}")
	g.p("if err != nil {")
	g.p("return %vOut{}, err", name)
	g.p("}")
	g.p("return v, nil")
	g.p("}")
	g.p("")
	return nil
}

func (g *generator) GenerateStructIn(
	intf *model.Interface,
	m *model.Method,
	outputPackagePath string,
	name string,
	argString string,
	argNames []string,
	argTypes []string) error {
	g.p("type %vIn struct{", name)
	//argStrings := makeArgStrings(true, argNames, argTypes)
	//for _, s := range argStrings {
	//	g.p("%v", s)
	//}
	for i, argName := range argNames {
		s := argTypes[i]
		if len(s) > 3 && s[0:3] == "..." {
			s = fmt.Sprintf("[]%v", s[3:])
		}
		g.p("%v %v", argName, s)

	}
	g.p("}")
	g.p("")

	return nil
}

func (g *generator) GenerateStructOut(
	intf *model.Interface,
	m *model.Method,
	outputPackagePath string,
	name string,
	argString string,
	argNames []string,
	argTypes []string) error {

	g.p("type %vOut  struct{", name)
	for i, o := range m.Out {
		g.p("Args%v %v", i, o.Type.String(g.packageMap, outputPackagePath))
	}
	g.p("}")

	return nil
}

func (g *generator) GenerateStruct(
	intf *model.Interface,
	m *model.Method,
	outputPackagePath string,
	name string,
	argString string,
	argNames []string,
	argTypes []string) error {

	g.p("type %v struct{", name)
	g.in()
	g.p("inData %vIn", name)
	g.p("outDataChannel chan %vOut", name)
	g.out()
	g.p("}")
	if len(argString) == 0 {
		g.p("func New%v(waitToComplete bool) *%v{", name, name)
	} else {
		g.p("func New%v(waitToComplete bool, %v) *%v{", name, argString, name)
	}
	g.in()
	g.p("var outDataChannel chan %vOut", name)
	g.p("if waitToComplete {")
	g.p("outDataChannel = make(chan %vOut)", name)
	g.p("}else{")
	g.p("outDataChannel = nil")
	g.p("}")
	g.p("return &%v{", name)
	g.in()
	g.p("inData: %vIn{", name)
	g.in()
	for _, s := range argNames {
		g.p("%v: %v,", s, s)
	}
	g.out()
	g.p("},")
	g.p("outDataChannel: outDataChannel,")
	g.out()
	g.p("}")
	g.out()
	g.p("}")
	g.p("")
	g.p("func (self * %v) Wait(onError func(interfaceName string, methodName string, err error )error) (%vOut, error){ ", name, name)
	g.in()
	g.p("data, ok := <-self.outDataChannel")
	g.p("if !ok{")
	g.in()
	g.p("generatedError :=  &%vError{", name)
	g.p("InterfaceName: \"%v\",", intf.Name)
	g.p("MethodName:    \"%v\",", m.Name)
	g.p("Reason:        \"Channel for %v::%v returned false\",", intf.Name, m.Name)
	g.p("}")
	g.p("if onError != nil{")
	g.in()
	g.p("err := onError(\"%v\", \"%v\", generatedError)", intf.Name, m.Name)
	g.p("return %vOut{}, err", name)
	g.out()
	g.p("}else{")
	g.p("return %vOut{}, generatedError", name)
	g.p("}")
	g.out()
	g.p("}")
	g.p("return data, nil")
	g.out()
	g.p("}")
	g.p("")
	g.p("func (self * %v) Close() error{", name)
	g.in()
	g.p("close(self.outDataChannel)")
	g.p("return nil")
	g.out()
	g.p("}")

	return nil
}

func (g *generator) GenerateError(
	intf *model.Interface,
	m *model.Method,
	outputPackagePath string,
	name string,
	argString string,
	argNames []string,
	argTypes []string) error {

	g.p("type %vError struct{", name)
	g.p("InterfaceName string")
	g.p("MethodName string")
	g.p("Reason string")
	g.p("}")
	g.p("")
	g.p("func (self *%vError) Error() string{", name)
	g.p("return fmt.Sprintf(\"error in data coming back from %%v::%%v. Reason: %%v\", self.InterfaceName, self.MethodName, self.Reason)")
	g.p("}")
	g.p("")

	return nil
}
