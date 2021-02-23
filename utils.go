package bindgen

import (
	"fmt"
	"strings"
	"text/template"
	"unicode"

	"modernc.org/cc"
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

// byName is a sortable slice of Declarations
type byName []Declaration

func (ds byName) Len() int           { return len(ds) }
func (ds byName) Less(i, j int) bool { return NameOf(ds[i].Decl()) < NameOf(ds[j].Decl()) }
func (ds byName) Swap(i, j int)      { ds[i], ds[j] = ds[j], ds[i] }

// Snake2Camel converts snake case to camel case. It's not particularly performant. Rather it's a quick and dirty function.
func Snake2Camel(s string, exported bool) (retVal string) {
	nextUpper := exported
	for i, v := range s {
		switch {
		case unicode.IsNumber(v):
			retVal += string(v)
		case unicode.IsUpper(v):
			if i == 0 && !nextUpper {
				retVal += strings.ToLower(string(v))
			} else {
				retVal += string(v)
			}
		case unicode.IsLower(v):
			if nextUpper {
				retVal += strings.ToUpper(string(v))
			} else {
				retVal += string(v)
			}
		case v == '_':
			nextUpper = true
			continue
		default:
			retVal += string(v)
		}
		nextUpper = false
	}
	return
}

// LongestCommonPrefix takes a slice of strings, and finds the longest common prefix
//
// This function was taken from github.com/chewxy/lingo/corpus
func LongestCommonPrefix(strs ...string) string {
	switch len(strs) {
	case 0:
		return "" // idiots
	case 1:
		return strs[0]
	}

	min := strs[0]
	max := strs[0]

	for _, s := range strs[1:] {
		switch {
		case s < min:
			min = s
		case s > max:
			max = s
		}
	}

	for i := 0; i < len(min) && i < len(max); i++ {
		if min[i] != max[i] {
			return min[:i]
		}
	}

	// In the case where lengths are not equal but all bytes
	// are equal, min is the answer ("foo" < "foobar").
	return min
}
