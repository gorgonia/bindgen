package bindgen

import (
	"fmt"
	"sort"

	"github.com/cznic/cc"
	"github.com/cznic/xc"
)

// FilterFunc is a function to filter types
type FilterFunc func(*cc.Declarator) bool

// Parse parses with the given model, as well as having some hard coded predefined definitions that are useful
// for translating C to Go code
func Parse(model *cc.Model, paths ...string) (*cc.TranslationUnit, error) {
	predefined, includePaths, sysIncludePaths, err := cc.HostConfig()
	if err != nil {
		return nil, fmt.Errorf("binding: failed to get host config: %v", err)
	}

	return cc.Parse(
		predefined+`
#define __const const
#define __attribute__(...)
#define __extension__
#define __inline
#define __restrict
unsigned __builtin_bswap32 (unsigned x);
unsigned long long __builtin_bswap64 (unsigned long long x);
`,
		paths,
		model,
		cc.IncludePaths(includePaths),
		cc.SysIncludePaths(sysIncludePaths),
	)
}

func Get(t *cc.TranslationUnit, filter FilterFunc) ([]Declaration, error) {
	var decls []Declaration
	for ; t != nil; t = t.TranslationUnit {
		if t.ExternalDeclaration.Case != 1 { /* Declaration */
			continue
		}

		d := t.ExternalDeclaration.Declaration
		if d.Case != 0 {
			// Other case is 1: StaticAssertDeclaration.
			continue
		}

		init := d.InitDeclaratorListOpt
		if init == nil {
			continue
		}
		idl := init.InitDeclaratorList
		if idl.InitDeclaratorList != nil {
			// We do not want comma-separated lists.
			continue
		}
		id := idl.InitDeclarator
		if id.Case != 0 {
			// We do not want assignments.
			continue
		}

		declarator := id.Declarator
		name := NameOf(declarator)

		if !filter(declarator) {
			continue
		}
		params, variadic := declarator.Type.Parameters()

		var retType cc.Type
		switch declarator.Type.Kind() {
		case cc.Function:
			// raw := declarator.Type.RawDeclarator()
			// extractParams(raw)
			// log.Println("SPECIFIER", raw.DirectDeclarator.ParameterTypeList.ParameterList.ParameterDeclaration.DeclarationSpecifiers.String())

			retType = declarator.Type.Result()
		case cc.Enum:
			// do nothing
		}
		decls = append(decls, Declaration{
			Pos:         declarator.Pos(),
			Name:        name,
			Return:      retType,
			CParameters: params,
			Variadic:    variadic,
			Declarator:  declarator,
		})
	}

	sort.Sort(byPosition(decls))

	return decls, nil
}

func NameOf(decl *cc.Declarator) (name string) {
	var id int
	id, _ = decl.Identifier()
	name = string(xc.Dict.S(id))
	// id, bindings := decl.Identifier()
	return
}

func TypeDefOf(t cc.Type) (name string) {
	id := t.Declarator().RawSpecifier().TypedefName()
	return string(xc.Dict.S(id))
}
