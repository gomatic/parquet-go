package page

import "testing"

func TestEncodeDataPageHeader(t *testing.T) {
	values := make([]int32, 100)
	for i := 0; i < 100; i++ {
		values[i] = int32(i)
	}

	preferences := EncodingPreferences{
		CompressionCodec: "",
		Strategy:         "",
	}

	enc := NewPageEncoder(preferences)

	err := enc.WriteInt32(values)
	if err != nil {
		t.Fatalf("could not WriteInt32 %s", err)
	}

	//	page := enc.DataPage()
	// enc.DictionaryPage
	// enc.IndexPage()

	// var final bytes.Buffer

	// // DataPage
	// var b bytes.Buffer
	// w := bufio.NewWriter(&b)

	// enc := encoding.NewPlainEncoder(w)
	// err := enc.WriteInt32(values)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// enc.Flush()

	// enc := NewPageEncoder(w)
	// err := enc.WriteInt32(value)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// var compressed bytes.Buffer
	// wc := snappy.NewWriter(&compressed)
	// if _, err := io.Copy(wc, &b); err != nil {
	// 	t.Fatal(err)
	// }

	// // Page Header
	// header := thrift.NewPageHeader()
	// header.CompressedPageSize = int32(compressed.Len())
	// header.UncompressedPageSize = int32(b.Len())
	// header.Type = thrift.PageType_DATA_PAGE
	// header.DataPageHeader = thrift.NewDataPageHeader()
	// header.DataPageHeader.NumValues = int32(enc.NumValues())
	// header.DataPageHeader.Encoding = thrift.Encoding_PLAIN
	// header.DataPageHeader.DefinitionLevelEncoding = thrift.Encoding_BIT_PACKED
	// header.DataPageHeader.RepetitionLevelEncoding = thrift.Encoding_BIT_PACKED
	// header.Write(&final)
	// _, err := io.Copy(&final, &compressed)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// // ColumnChunk
	// offset := 0
	// name := "thisfile.parquet"
	// chunk := thrift.NewColumnChunk()
	// chunk.FileOffset = int64(offset)
	// chunk.FilePath = &name
	// chunk.MetaData = thrift.NewColumnMetaData()
	// chunk.MetaData.TotalCompressedSize = 0
	// chunk.MetaData.TotalUncompressedSize = 0
	// chunk.MetaData.Codec = thrift.CompressionCodec_SNAPPY

	// chunk.MetaData.DataPageOffset = 0
	// chunk.MetaData.DictionaryPageOffset = nil

	// chunk.MetaData.Type = thrift.Type_INT32
	// chunk.MetaData.PathInSchema = []string{"some"}

	// if _, err := chunk.Write(&final); err != nil {
	// 	t.Fatal(err)
	// }

	// // Schema Element
	// columnSchema := thrift.NewSchemaElement()
	// columnSchema.Name = name
	// columnSchema.NumChildren = nil
	// columnSchema.Type = typeInt32
	// columnSchema.RepetitionType = nil

	// // Encoder
	// eenc := NewEncoder([]*thrift.ColumnChunk{})
	// //eenc.AddRowGroup()
	// fd, err := os.Create("some.file.parquet")
	// if err != nil {
	// 	log.Println(err)
	// }

	// if err := eenc.Write(fd); err != nil {
	// 	log.Println(err)
	// }
	// fd.Close()
}
