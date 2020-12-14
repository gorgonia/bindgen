package bindgen

import (
	"bytes"
	"fmt"
	"strings"

	"modernc.org/cc"
)

func ExampleExplore() {
	functions := func(decl *cc.Declarator) bool {
		if !strings.HasPrefix(NameOf(decl), "func") {
			return false
		}
		if decl.Type.Kind() == cc.Function {
			return true
		}
		return false
	}
	enums := func(decl *cc.Declarator) bool {
		if decl.Type.Kind() == cc.Enum {
			return true
		}
		return false
	}
	others := func(decl *cc.Declarator) bool {
		if decl.Type.Kind() == cc.Ptr || decl.Type.Kind() == cc.Struct {
			return true
		}
		return false
	}
	tu, err := Parse(Model(), "testdata/dummy.h")
	if err != nil {
		panic(err)
	}

	if err = Explore(tu, functions, enums, others); err != nil {
		panic(err)
	}

	// Output:
	// func1i
	// func1f
	// func1fp
	// func2i
	// func2f
	// funcErr
	// funcCtx
	//
	// error
	// fntype_t
	//
	// dummy_t
	// dummy2_t
}

func ExampleGenIgnored() {
	functions := func(decl *cc.Declarator) bool {
		if !strings.HasPrefix(NameOf(decl), "func") {
			return false
		}
		if decl.Type.Kind() == cc.Function {
			return true
		}
		return false
	}
	tu, err := Parse(Model(), "testdata/dummy.h")
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if err := GenIgnored(&buf, tu, functions); err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
	// Output:
	// var ignored = map[string]struct{}{
	// "func1i":{},
	// "func1f":{},
	// "func1fp":{},
	// "func2i":{},
	// "func2f":{},
	// "funcErr":{},
	// "funcCtx":{},
	// }
}

func ExampleGenNameMap() {
	functions := func(decl *cc.Declarator) bool {
		if !strings.HasPrefix(NameOf(decl), "func") {
			return false
		}
		if decl.Type.Kind() == cc.Function {
			return true
		}
		return false
	}

	trans := func(a string) string {
		return strings.ToTitle(strings.TrimPrefix(a, "func"))
	}
	tu, err := Parse(Model(), "testdata/dummy.h")
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if err := GenNameMap(&buf, tu, "m", trans, functions, false); err != nil {
		panic(err)
	}
	fmt.Println(buf.String())

	// Output:
	// var m = map[string]string{
	// "func1i": "1I",
	// "func1f": "1F",
	// "func1fp": "1FP",
	// "func2i": "2I",
	// "func2f": "2F",
	// "funcErr": "ERR",
	// "funcCtx": "CTX",
	// }

}

func ExampleGenNameMap_2() {
	functions := func(decl *cc.Declarator) bool {
		if !strings.HasPrefix(NameOf(decl), "func") {
			return false
		}
		if decl.Type.Kind() == cc.Function {
			return true
		}
		return false
	}

	trans := func(a string) string {
		return strings.ToTitle(strings.TrimPrefix(a, "func"))
	}
	tu, err := Parse(Model(), "testdata/dummy.h")
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	fmt.Fprintln(&buf, "func init() {")
	if err := GenNameMap(&buf, tu, "m", trans, functions, true); err != nil {
		panic(err)
	}
	fmt.Fprintln(&buf, "}")
	fmt.Println(buf.String())

	// Output:
	// func init() {
	// m = map[string]string{
	// "func1i": "1I",
	// "func1f": "1F",
	// "func1fp": "1FP",
	// "func2i": "2I",
	// "func2f": "2F",
	// "funcErr": "ERR",
	// "funcCtx": "CTX",
	// }
	// }

}
