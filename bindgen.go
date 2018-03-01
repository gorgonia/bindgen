package bindgen

import (
	"go/token"
	"text/template"

	"github.com/cznic/cc"
	"github.com/cznic/xc"
)

// TypeKey is typically used as a representation of a C type that can  be used as a key in a map
type TypeKey struct {
	IsPointer bool
	Kind      cc.Kind
	Name      string
}

// Template represents a template of conversion. An optional InContext() function may be provided to check if the template needs to be executed
type Template struct {
	*template.Template
	InContext func() bool
}

// Declaration is anything with a position
type Declaration interface {
	Position() token.Position
}

type Namer interface {
	Name() string
}

// CSignature is a description of a C  declaration.
type CSignature struct {
	Pos         token.Pos
	Name        string
	Return      cc.Type
	CParameters []cc.Parameter
	Variadic    bool
	Declarator  *cc.Declarator
}

// Position returns the token position of the declaration.
func (d *CSignature) Position() token.Position { return xc.FileSet.Position(d.Pos) }

// Parameters returns the declaration's CParameters converted to a []Parameter.
func (d *CSignature) Parameters() []Parameter {
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
