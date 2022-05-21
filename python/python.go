package python

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/syzkrash/skol/codegen"
	"github.com/syzkrash/skol/parser"
)

// Python is a simple transpiler that turns a skol AST into Python code
type Python struct {
	parser *parser.Parser
}

func NewPython(fn string, src io.RuneScanner) codegen.Generator {
	gen := &Python{
		parser: parser.NewParser(fn, src),
	}
	gen.addEnv()
	return gen
}

func (p *Python) value(v parser.Node, output io.StringWriter) (err error) {
	switch v.Kind() {
	case parser.NdInteger:
		i := v.(*parser.IntegerNode)
		_, err = output.WriteString(fmt.Sprint(i.Int))
	case parser.NdFloat:
		f := v.(*parser.FloatNode)
		_, err = output.WriteString(fmt.Sprint(f.Float))
	case parser.NdString:
		s := v.(*parser.StringNode)
		_, err = output.WriteString(fmt.Sprintf("\"%s\"", s.Str))
	case parser.NdChar:
		c := v.(*parser.CharNode)
		_, err = output.WriteString(fmt.Sprintf("%d", c.Char))
	case parser.NdVarRef:
		r := v.(*parser.VarRefNode)
		_, err = output.WriteString(r.Var)
	case parser.NdFuncCall:
		c := v.(*parser.FuncCallNode)
		if err = p.funcCall(c, output); err != nil {
			return
		}
	default:
		err = fmt.Errorf("unexpected value node: %s", v.Kind())
	}
	return
}

var ops = map[string]string{
	"add": "+",
	"sub": "-",
	"mul": "*",
	"pow": "**",
	"div": "/",
	"mod": "%",
}

func (p *Python) funcCall(f *parser.FuncCallNode, output io.StringWriter) (err error) {
	if op, ok := ops[f.Func]; ok {
		if _, err = output.WriteString("("); err != nil {
			return
		}
		if err = p.value(f.Args[0], output); err != nil {
			return
		}
		if _, err = output.WriteString(op); err != nil {
			return
		}
		if err = p.value(f.Args[1], output); err != nil {
			return
		}
		if _, err = output.WriteString(")"); err != nil {
			return
		}
	} else if err = p.regularFuncCall(f, output); err != nil {
		return
	}
	return
}

func (p *Python) regularFuncCall(f *parser.FuncCallNode, output io.StringWriter) (err error) {
	if _, err = output.WriteString(f.Func); err != nil {
		return
	}
	if _, err = output.WriteString("("); err != nil {
		return
	}
	if len(f.Args) == 1 {
		if err = p.value(f.Args[0], output); err != nil {
			return
		}
	} else if len(f.Args) > 1 {
		for _, a := range f.Args[:len(f.Args)-1] {
			if err = p.value(a, output); err != nil {
				return
			}
			if _, err = output.WriteString(","); err != nil {
				return
			}
		}
		if err = p.value(f.Args[len(f.Args)-1], output); err != nil {
			return
		}
	}
	if _, err = output.WriteString(")"); err != nil {
		return
	}
	return
}

func (p *Python) internalGenerate(n parser.Node, output io.StringWriter) (err error) {
	switch n.Kind() {
	case parser.NdReturn:
		r := n.(*parser.ReturnNode)
		if _, err = output.WriteString("return "); err != nil {
			return
		}
		if err = p.value(r.Value, output); err != nil {
			return
		}
	case parser.NdVarDef:
		v := n.(*parser.VarDefNode)
		if _, err = output.WriteString(v.Var); err != nil {
			return
		}
		if _, err = output.WriteString("="); err != nil {
			return
		}
		if err = p.value(v.Value, output); err != nil {
			return
		}
		output.WriteString("\n")
	case parser.NdFuncDef:
		f := n.(*parser.FuncDefNode)
		if _, err = output.WriteString("def "); err != nil {
			return
		}
		if _, err = output.WriteString(f.Func); err != nil {
			return
		}
		if _, err = output.WriteString("("); err != nil {
			return
		}
		argNames := []string{}
		for n := range f.Arg {
			argNames = append(argNames, n)
		}
		if len(argNames) == 1 {
			if _, err = output.WriteString(argNames[0]); err != nil {
				return
			}
		} else if len(argNames) > 1 {
			for _, n := range argNames[:len(argNames)-1] {
				if _, err = output.WriteString(n); err != nil {
					return
				}
				if _, err = output.WriteString(","); err != nil {
					return
				}
			}
			if _, err = output.WriteString(argNames[len(argNames)-1]); err != nil {
				return
			}
		}
		if _, err = output.WriteString("):\n"); err != nil {
			return
		}
		for _, n := range f.Body {
			if _, err = output.WriteString("\t"); err != nil {
				return
			}
			if err = p.internalGenerate(n, output); err != nil {
				return
			}
		}
		output.WriteString("\n")
	case parser.NdFuncCall:
		f := n.(*parser.FuncCallNode)
		p.funcCall(f, output)
		output.WriteString("\n")

	default:
		return fmt.Errorf("unexpected top-level node: %s", n.Kind())
	}
	return nil
}

func (p *Python) Generate(output io.StringWriter) error {
	for {
		n, err := p.parser.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if err = p.internalGenerate(n, output); err != nil {
			return err
		}
	}
}

func (p *Python) addEnv() {
	p.parser.Scope.Funcs = map[string]*parser.Function{
		"print": {
			Name: "print",
			Args: map[string]parser.ValueType{
				"a": parser.VtAny,
			},
			Ret: parser.VtNothing,
		},
	}
	for fn := range ops {
		p.parser.Scope.Funcs[fn] = &parser.Function{
			Name: fn,
			Args: map[string]parser.ValueType{
				"a": parser.VtFloat,
				"b": parser.VtFloat,
			},
			Ret: parser.VtFloat,
		}
	}
}

func (*Python) CanRun() bool {
	return true
}

func (*Python) Run(fn string) error {
	cmd := exec.Command("py", fn)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
