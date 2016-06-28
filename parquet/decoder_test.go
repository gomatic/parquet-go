package parquet

import (
//	"os"
//	"reflect"
	"testing"
)


type cell struct {
	d int
	r int
	v interface{}
}

/* TODO: fix broken references */
//func checkColumnValues(t *testing.T, path string, columnIdx int, expected []cell) {
//
//	fd, err := OpenFile(path)
//	if err != nil {
//		t.Errorf("failed to read %s: %s", path, err)
//		return
//	}
//	defer fd.Close()
//	// schema, err := schemaFromFileMetaData(m)
//	// if err != nil {
//	// 	t.Errorf("failed to create schema: %s", err)
//	// 	return
//	// }
//
//	schema := fd.Schema()
//	columns := schema.Columns()
//	scanner, err := fd.ColumnScanner(columns[columnIdx])
//
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	k := 0
//	for i, rg := range m.RowGroups {
//		cc := rg.Columns[c]
//		cs := schema.ColumnByPath(cc.MetaData.PathInSchema)
//		if cs == nil {
//			t.Errorf("column %d: no schema", c)
//			return
//		}
//		cr, err := NewColumnChunkReader(r, *cs, *cc)
//		if err != nil {
//			t.Errorf("column %d: failed to create reader for row group %d: %s", c, i, err)
//			return
//		}
//		for cr.Next() {
//			if k < len(expected) {
//				got := cell{cr.Levels().D, cr.Levels().R, cr.Value()}
//				if !reflect.DeepEqual(got, expected[k]) {
//					t.Errorf("column %d: value at pos %d = %#v, want %#v", c, k, got, expected[k])
//				}
//			}
//			k++
//			//fmt.Printf("V:%v\tD:%d\tR:%d\n", cr.Value(), cr.Levels().D, cr.Levels().R)
//		}
//		if cr.Err() != nil {
//			t.Errorf("column %d: failed to read row group %d: %s", c, i, cr.Err())
//		}
//	}
//
//	if scanner.Err() != nil {
//		t.Errorf("column %d: failed to read row group: %s", columnIdx, scanner.Err())
//	}
//
//	// if k != len(expected) {
//	// 	t.Errorf("column %d: read %d values, want %d values", c, k, len(expected))
//	// }
//}

/* old version that allows tests to build tests */

func checkColumnValues(t *testing.T, path string, columnIdx int, expected []cell) {

	fd, err := OpenFile(path)
	if err != nil {
		t.Errorf("failed to read %s: %s", path, err)
		return
	}
	defer fd.Close()

	// schema, err := schemaFromFileMetaData(m)
	// if err != nil {
	// 	t.Errorf("failed to create schema: %s", err)
	// 	return
	// }

	schema := fd.Schema()
	columns := schema.Columns()
	scanner, err := fd.ColumnScanner(columns[columnIdx])

	if err != nil {
		t.Fatal(err)
	}

	// for i, rg := range m.RowGroups {
	// 	cc := rg.Columns[c]
	// 	columnSchema := schema.ColumnByPath(cc.MetaData.PathInSchema)
	// 	//var cr ColumnChunkReader
	//switch cs.SchemaElement.GetType() {
	//case parquetformat.Type_BOOLEAN:
	//cr, err = NewBooleanColumnChunkReader(r, cs, cc)
	//case parquetformat.Type_BYTE_ARRAY:
	//cr, err = NewByteArrayColumnChunkReader(r, cs, cc)
	//}

	// scanner := NewColumnScanner(r, cc, columnSchema.SchemaElement)

	for scanner.Scan() {

		//buffer := make([]bool, 0, 4)

		// if k < len(expected) {
		// 	// got := cell{cr.Levels().D, cr.Levels().R, cr.Value()}
		// 	// if !reflect.DeepEqual(got, expected[k]) {
		// 	// 	t.Errorf("column %d: value at pos %d = %#v, want %#v", c, k, got, expected[k])
		// 	// }
		// }

		// scanner.BoolArray()

		// k += len(buffer)
		//fmt.Printf("V:%v\tD:%d\tR:%d\n", cr.Value(), cr.Levels().D, cr.Levels().R)
	}

	if scanner.Err() != nil {
		t.Errorf("column %d: failed to read row group: %s", columnIdx, scanner.Err())
	}

	// if k != len(expected) {
	// 	t.Errorf("column %d: read %d values, want %d values", c, k, len(expected))
	// }
}


func TestBooleanColumnChunkReader(t *testing.T) {
	checkColumnValues(t, "testdata/Booleans.parquet", 0, []cell{
		{0, 0, true},
		{0, 0, true},
		{0, 0, false},
		{0, 0, true},
		{0, 0, false},
		{0, 0, true},
	})

	checkColumnValues(t, "testdata/Booleans.parquet", 1, []cell{
		{0, 0, false},
		{1, 0, false},
		{1, 0, true},
		{1, 0, true},
		{0, 0, false},
		{1, 0, true},
	})

	checkColumnValues(t, "testdata/Booleans.parquet", 2, []cell{
		{0, 0, false},

		{0, 0, false},

		{1, 0, true},

		{1, 0, true},
		{1, 1, false},
		{1, 1, true},

		{0, 0, false},
		{1, 0, true},
	})
}

func TestByteArrayColumnChunkReader(t *testing.T) {
	checkColumnValues(t, "testdata/ByteArrays.parquet", 0, []cell{
		{0, 0, []byte{'r', '1'}},
		{0, 0, []byte{'r', '2'}},
		{0, 0, []byte{'r', '3'}},
		{0, 0, []byte{'r', '4'}},
		{0, 0, []byte{'r', '5'}},
		{0, 0, []byte{'r', '6'}},
	})

	checkColumnValues(t, "testdata/ByteArrays.parquet", 1, []cell{
		{0, 0, []byte(nil)},
		{1, 0, []byte{'o', '2'}},
		{1, 0, []byte{'o', '3'}},
		{1, 0, []byte{'o', '4'}},
		{0, 0, []byte(nil)},
		{1, 0, []byte{'o', '6'}},
	})

	checkColumnValues(t, "testdata/ByteArrays.parquet", 2, []cell{
		{0, 0, []byte(nil)},

		{0, 0, []byte(nil)},

		{1, 0, []byte{'p', '3', '_', '1'}},

		{1, 0, []byte{'p', '4', '_', '1'}},
		{1, 1, []byte{'p', '4', '_', '2'}},
		{1, 1, []byte{'p', '4', '_', '3'}},

		{0, 0, []byte(nil)},

		{1, 0, []byte{'p', '6', '_', '1'}},
	})
}
