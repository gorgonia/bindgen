package bindgen

import (
	"fmt"
	"text/template"

	"github.com/cznic/cc"
)

// Pure "lifts" a string or *template.Template into a template
func Pure(any interface{}) Template {
	switch a := any.(type) {
	case string:
		return Template{Template: template.Must(template.New(a).Parse(a))}
	case *template.Template:
		return Template{Template: a}
	case Template:
		return a
	case struct {
		*template.Template
		InContext func() bool
	}:
		return Template(a)
	default:
		panic(fmt.Sprintf("%v of %T unhandled", any, any))
	}
}

// IsConstType returns true if the C-type is specified with a `const`
func IsConstType(a cc.Type) bool { return a.Specifier().IsConst() }

// IsPointer returns true if the C-type is specified as a pointer
func IsPointer(a cc.Type) bool { return a.RawDeclarator().PointerOpt != nil }

// IsVoid returns true if the C type is
func IsVoid(a cc.Type) bool { return a.Kind() == cc.Void }

// byPosition implements a sorting for a slice of Declaration
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
