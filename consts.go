package bindgen

import (
	"unsafe"

	"modernc.org/cc"
)

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
