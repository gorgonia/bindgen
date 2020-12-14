package bindgen_test

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gorgonia/bindgen"
	"modernc.org/cc"
)

// genEnums represents a list of enums we want to generate
var genEnums = map[bindgen.TypeKey]struct{}{
	{Kind: cc.Enum, Name: "error"}: {},
}

var enumMappings = map[bindgen.TypeKey]string{
	{Kind: cc.Enum, Name: "error"}: "Status",
}

// This is an example of how to convert enums.
func Example_convertingEnums() {
	t, err := bindgen.Parse(bindgen.Model(), "testdata/dummy.h")
	if err != nil {
		panic(err)
	}
	enums := func(decl *cc.Declarator) bool {
		name := bindgen.NameOf(decl)
		kind := decl.Type.Kind()
		tk := bindgen.TypeKey{Kind: kind, Name: name}
		if _, ok := genEnums[tk]; ok {
			return true
		}
		return false
	}
	decls, err := bindgen.Get(t, enums)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	for _, d := range decls {
		// first write the type
		//	 type ___ int
		// This is possible because cznic/cc parses all enums as int.
		//
		// you are clearly free to add your own mapping.
		e := d.(*bindgen.Enum)
		tk := bindgen.TypeKey{Kind: cc.Enum, Name: e.Name}
		fmt.Fprintf(&buf, "type %v int\nconst (\n", enumMappings[tk])

		// then write the const definitions:
		// 	const(...)

		for _, a := range e.Type.EnumeratorList() {
			// this is a straightforwards mapping of the C defined name. The name is kept exactly the same, with a lowecase mapping
			// in real life, you might not want this, (for example, you may not want to export the names, which are typically in all caps),
			// or you might want different names

			enumName := string(a.DefTok.S())
			goName := strings.ToLower(enumName)
			fmt.Fprintf(&buf, "%v %v = C.%v\n", goName, enumMappings[tk], enumName)
		}
		buf.Write([]byte(")\n"))
	}
	fmt.Println(buf.String())

	// Output:
	// type Status int
	// const (
	// success Status = C.SUCCESS
	// failure Status = C.FAILURE
	// )
}
