package bindgen

import (
	"html/template"
	"unsafe"

	"github.com/cznic/cc"
)

var goTypes = map[TypeKey]*template.Template{
	{Kind: cc.Undefined}: template.Must(template.New("<undefined>").Parse("<undefined>")),
	{Kind: cc.Int}:       template.Must(template.New("int").Parse("int")),
	{Kind: cc.Float}:     template.Must(template.New("float32").Parse("float32")),
	{Kind: cc.Float, IsPointer: true}: template.Must(template.New("[]float32").Parse(
		`{{if eq . "alpha" "beta" "cScalar" "sScalar" "result" "retVal"}}float32{{else}}[]float32{{end}}`)),
	{Kind: cc.Double}: template.Must(template.New("float64").Parse("float64")),
	{Kind: cc.Double, IsPointer: true}: template.Must(template.New("[]float64").Parse(
		`{{if eq . "alpha" "beta" "cScalar" "sScalar" "result" "retVal"}}float64{{else}}[]float64{{end}}`)),
	{Kind: cc.Bool}:          template.Must(template.New("bool").Parse("bool")),
	{Kind: cc.FloatComplex}:  template.Must(template.New("complex64").Parse("complex64")),
	{Kind: cc.DoubleComplex}: template.Must(template.New("complex128").Parse("complex128")),

	{Kind: cc.FloatComplex, IsPointer: true}: template.Must(template.New("cuComplex*").Parse(
		`{{if eq . "alpha" "beta" "cScalar" "sScalar" "result" "retVal"}}complex64{{else}}[]complex64{{end}}`,
	)),
	{Kind: cc.DoubleComplex, IsPointer: true}: template.Must(template.New("cuDoubleComplex*").Parse(
		`{{if eq . "alpha" "beta" "cScalar" "sScalar" "result" "retVal"}}complex128{{else}}[]complex128{{end}}`,
	)),
	{Kind: cc.Int, IsPointer: true}: template.Must(template.New("int*").Parse(
		`{{if eq . "alpha" "beta" "cScalar" "sScalar" "result" "retVal"}}int{{else}}[]int{{end}}`)),
}

func Model() *cc.Model {
	p := int(unsafe.Sizeof(uintptr(0)))
	i := int(unsafe.Sizeof(int(0)))
	return &cc.Model{
		Items: map[cc.Kind]cc.ModelItem{
			cc.Ptr:               {Size: p, Align: p, StructAlign: p},
			cc.UintPtr:           {Size: p, Align: p, StructAlign: p},
			cc.Void:              {Size: 0, Align: 1, StructAlign: 1},
			cc.Char:              {Size: 1, Align: 1, StructAlign: 1},
			cc.SChar:             {Size: 1, Align: 1, StructAlign: 1},
			cc.UChar:             {Size: 1, Align: 1, StructAlign: 1},
			cc.Short:             {Size: 2, Align: 2, StructAlign: 2},
			cc.UShort:            {Size: 2, Align: 2, StructAlign: 2},
			cc.Int:               {Size: 4, Align: 4, StructAlign: 4},
			cc.UInt:              {Size: 4, Align: 4, StructAlign: 4},
			cc.Long:              {Size: i, Align: i, StructAlign: i},
			cc.ULong:             {Size: i, Align: i, StructAlign: i},
			cc.LongLong:          {Size: 8, Align: 8, StructAlign: 8},
			cc.ULongLong:         {Size: 8, Align: 8, StructAlign: 8},
			cc.Float:             {Size: 4, Align: 4, StructAlign: 4},
			cc.Double:            {Size: 8, Align: 8, StructAlign: 8},
			cc.LongDouble:        {Size: 8, Align: 8, StructAlign: 8},
			cc.Bool:              {Size: 1, Align: 1, StructAlign: 1},
			cc.FloatComplex:      {Size: 8, Align: 8, StructAlign: 8},
			cc.DoubleComplex:     {Size: 16, Align: 16, StructAlign: 16},
			cc.LongDoubleComplex: {Size: 16, Align: 16, StructAlign: 16},
		},
	}
}
