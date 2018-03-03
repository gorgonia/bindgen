package bindgen

import "testing"

var snaketests = []struct {
	input, output string
	exported      bool
}{
	{"hello_world", "helloWorld", false},
	{"Hello_World", "HelloWorld", true},
	{"Hellow_Sekai_World", "hellowSekaiWorld", false},
	{"Hellow_Sekai_World", "HellowSekaiWorld", true},
	{"Hellow_Sekai_World_123", "HellowSekaiWorld123", true},
	{"Héllow_Sekai_World", "héllowSekaiWorld", false},
	{"_trailing_under||scores_", "TrailingUnder||scores", true}, // this is not a valud function or name, but added for completeness sake
}

func TestSnake2Camel(t *testing.T) {
	for _, st := range snaketests {
		out := Snake2Camel(st.input, st.exported)
		if out != st.output {
			t.Fatalf("Failed on Input %q. Wanted %q. Got %q", st.input, st.output, out)
		}
	}
}
