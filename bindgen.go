package bindgen

import (
	"bytes"
	"fmt"
	"go/token"
	"log"
	"text/template"

	"github.com/cznic/cc"
	"github.com/cznic/xc"
)

type TypeKey struct {
	IsPointer bool
	Kind      cc.Kind
}

// Declaration is a description of a C function declaration.
type Declaration struct {
	Pos         token.Pos
	Name        string
	Return      cc.Type
	CParameters []cc.Parameter
	Variadic    bool
	Declarator  *cc.Declarator
}

// func (d Declaration) Format(f fmt.State, c rune) {
// 	if !f.Flag('#') {
// 		fmt.Fprintf(f, "Declaration{%v}", d.Name)
// 		return
// 	}
// 	fmt.Fprintf(f, "func %v(", d.Name)
// 	for i, param := range d.Parameters() {
// 		fmt.Fprintf(f, "%v %v", param.Name(), param.Type())
// 		if i < len(d.CParameters) {
// 			fmt.Fprint(f, ", ")
// 		}
// 	}
// 	fmt.Fprintf(f, ") %v", d.Return)
// }

// Position returns the token position of the declaration.
func (d Declaration) Position() token.Position { return xc.FileSet.Position(d.Pos) }

// Parameters returns the declaration's CParameters converted to a []Parameter.
func (d *Declaration) Parameters() []Parameter {
	p := make([]Parameter, len(d.CParameters))
	for i, c := range d.CParameters {
		p[i] = Parameter{c, TypeDefOf(c.Type)}
	}
	return p
}

// Parameter is a C function parameter.
type Parameter struct {
	cc.Parameter
	TypeDefName string // can be empty
}

// Name returns the name of the parameter.
func (p *Parameter) Name() string { return string(xc.Dict.S(p.Parameter.Name)) }

// Type returns the C type of the parameter.
func (p *Parameter) Type() cc.Type { return p.Parameter.Type }

// Kind returns the C kind of the parameter.
func (p *Parameter) Kind() cc.Kind { return p.Parameter.Type.Kind() }

// Elem returns the pointer type of a pointer parameter or the element type of an
// array parameter.
func (p *Parameter) Elem() cc.Type { return p.Parameter.Type.Element() }

type byPosition []Declaration

func (d byPosition) Len() int { return len(d) }
func (d byPosition) Less(i, j int) bool {
	iPos := d[i].Position()
	jPos := d[j].Position()
	if iPos.Filename == jPos.Filename {
		return iPos.Line < jPos.Line
	}
	return iPos.Filename < jPos.Filename
}
func (d byPosition) Swap(i, j int) { d[i], d[j] = d[j], d[i] }

// GoTypeFor returns a string representation of the given type using a mapping in
// types. GoTypeFor will panic if no type mapping is found after searching the
// user-provided types mappings and then the following mapping:
//  {Kind: cc.Int}:                     "int",
//  {Kind: cc.Float}:                   "float32",
//  {Kind: cc.Float, IsPointer: true}:  "[]float32",
//  {Kind: cc.Double}:                  "float64",
//  {Kind: cc.Double, IsPointer: true}: "[]float64",
//  {Kind: cc.Bool}:                    "bool",
//  {Kind: cc.FloatComplex}:            "complex64",
//  {Kind: cc.DoubleComplex}:           "complex128",
func GoTypeFor(typ cc.Type, name string, types ...map[TypeKey]*template.Template) string {
	if typ == nil {
		return "<nil>"
	}
	k := typ.Kind()
	isPtr := typ.Kind() == cc.Ptr
	if isPtr {
		k = typ.Element().Kind()
	}
	var buf bytes.Buffer
	for _, t := range types {
		if s, ok := t[TypeKey{Kind: k, IsPointer: isPtr}]; ok {
			err := s.Execute(&buf, name)
			if err != nil {
				panic(err)
			}
			return buf.String()
		}
	}
	s, ok := goTypes[TypeKey{Kind: k, IsPointer: isPtr}]
	if ok {
		err := s.Execute(&buf, name)
		if err != nil {
			panic(err)
		}
		return buf.String()
	}
	log.Printf("%v", typ.Tag())
	panic(fmt.Sprintf("unknown type key: %v %+v", typ, TypeKey{Kind: k, IsPointer: isPtr}))
}

// GoTypeForEnum returns a string representation of the given enum type using a mapping
// in types. GoTypeForEnum will panic if no type mapping is found after searching the
// user-provided types mappings or the type is not an enum.
func GoTypeForEnum(typ cc.Type, name string, types ...map[string]*template.Template) string {
	if typ == nil {
		return "<nil>"
	}
	if typ.Kind() != cc.Enum {
		panic(fmt.Sprintf("invalid type: %v", typ))
	}
	tag := typ.Tag()
	if tag != 0 {
		n := string(xc.Dict.S(tag))
		for _, t := range types {
			if s, ok := t[n]; ok {
				var buf bytes.Buffer
				err := s.Execute(&buf, name)
				if err != nil {
					panic(err)
				}
				return buf.String()
			}
		}
	}
	log.Printf("%s", typ.Declarator())
	panic(fmt.Sprintf("unknown type: %+v", typ))
}

func IsConstType(a cc.Type) bool {
	return a.Specifier().IsConst()
}

func IsPointer(a cc.Type) bool {
	return a.RawDeclarator().PointerOpt != nil
}

func IsVoid(a cc.Type) bool {
	return a.String() == "void"
}
