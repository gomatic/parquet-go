package datatypes

/*
type byteArrayPlainDecoder struct {
	data []byte
	pos  int
}

func newByteArrayPlainDecoder() *byteArrayPlainDecoder {
	return &byteArrayPlainDecoder{}
}

func (d *byteArrayPlainDecoder) init(data []byte) {
	d.data = data
	d.pos = 0
}

func (d *byteArrayPlainDecoder) next() (value []byte, err error) {
	if d.pos > len(d.data)-4 {
		return nil, fmt.Errorf("bytearray/plain: no more data")
	}
	size := int(binary.LittleEndian.Uint32(d.data[d.pos:])) // TODO: think about int overflow here
	d.pos += 4
	if d.pos > len(d.data)-size {
		return nil, fmt.Errorf("bytearray/plain: not enough data")
	}
	// TODO: configure copy or not
	value = make([]byte, size)
	copy(value, d.data[d.pos:d.pos+size])
	d.pos += size
	return
}

// TODO: shorter name?
type ByteArrayColumnChunkReader struct {
	// TODO: consider using a separate reader for each column chunk (no seeking)
	r *countingReader

	// initialized once
	maxLevels Levels
	totalSize int64

	// changing state
	err            error
	curLevels      Levels
	curValue       []byte
	pageValuesRead int
	atStartOfPage  bool
	atLastPage     bool
	valuesRead     int64
	dataPageOffset int64
	header         *thrift.PageHeader

	// decoders
	decoder  *byteArrayPlainDecoder
	dDecoder *rle32Decoder
	rDecoder *rle32Decoder
}

// NewByteArrayColumnChunkReader
func NewByteArrayColumnChunkReader(r io.ReadSeeker, cs *ColumnSchema, chunk *thrift.ColumnChunk) (*ByteArrayColumnChunkReader, error) {
	if chunk.FilePath != nil {
		return nil, fmt.Errorf("data in another file: '%s'", *chunk.FilePath)
	}

	// chunk.FileOffset is useless
	// see https://issues.apache.org/jira/browse/PARQUET-291
	meta := chunk.MetaData
	if meta == nil {
		return nil, fmt.Errorf("missing ColumnMetaData")
	}

	if meta.Type != thrift.Type_BYTE_ARRAY {
		return nil, fmt.Errorf("wrong type, expected BYTE_ARRAY was %s", meta.Type)
	}

	schemaElement := cs.SchemaElement
	if schemaElement.RepetitionType == nil {
		return nil, fmt.Errorf("nil RepetitionType (root SchemaElement?)")
	}

	// uncompress
	if meta.Codec != thrift.CompressionCodec_UNCOMPRESSED {
		return nil, fmt.Errorf("unsupported compression codec: %s", meta.Codec)
	}

	// only REQUIRED non-nested columns are supported for now
	// so definitionLevel = 1 and repetitionLevel = 1
	cr := ByteArrayColumnChunkReader{
		r:              &countingReader{rs: r},
		totalSize:      meta.TotalCompressedSize,
		dataPageOffset: meta.DataPageOffset,
		maxLevels:      cs.MaxLevels,
		decoder:        newByteArrayPlainDecoder(),
	}

	if *schemaElement.RepetitionType == thrift.FieldRepetitionType_REQUIRED {
		// TODO: also check that len(Path) = maxD
		// For data that is required, the definition levels are not encoded and
		// always have the value of the max definition level.
		cr.curLevels.D = cr.maxLevels.D
		// TODO: document level ranges
	} else {
		cr.dDecoder = newRLE32Decoder(bitWidth(cr.maxLevels.D))
	}
	if cr.curLevels.D == 0 && *schemaElement.RepetitionType != thrift.FieldRepetitionType_REPEATED {
		// TODO: I think we need to check all schemaElements in the path
		cr.curLevels.R = 0
		// TODO: clarify the following comment from parquet-format/README:
		// If the column is not nested the repetition levels are not encoded and
		// always have the value of 1
	} else {
		cr.rDecoder = newRLE32Decoder(bitWidth(cr.maxLevels.R))
	}

	return &cr, nil
}

// AtStartOfPage returns true if the reader is positioned at the first value of a page.
func (cr *ByteArrayColumnChunkReader) AtStartOfPage() bool {
	return cr.atStartOfPage
}

// PageHeader returns page header of the current page.
func (cr *ByteArrayColumnChunkReader) PageHeader() *thrift.PageHeader {
	return cr.header
}

// readDataPage returns the content of the current DataPage
func (cr *ByteArrayColumnChunkReader) readDataPage() ([]byte, error) {
	var err error
	n := cr.r.n

	if _, err := cr.r.rs.Seek(cr.dataPageOffset, 0); err != nil {
		return nil, err
	}
	ph := thrift.PageHeader{}
	if err := ph.Read(cr.r); err != nil {
		return nil, err
	}
	if ph.Type != thrift.PageType_DATA_PAGE {
		return nil, fmt.Errorf("DATA_PAGE type expected, but was %s", ph.Type)
	}
	dph := ph.DataPageHeader
	if dph == nil {
		return nil, fmt.Errorf("null DataPageHeader in %+v", ph)
	}
	if dph.Encoding != thrift.Encoding_PLAIN {
		return nil, fmt.Errorf("unsupported encoding %s", dph.Encoding)
	}

	size := int64(ph.CompressedPageSize)
	data, err := ioutil.ReadAll(io.LimitReader(cr.r, size))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) != size {
		return nil, fmt.Errorf("unable to read page fully: got %d bytes, expected %d", len(data), size)
	}
	if cr.r.n > cr.totalSize {
		return nil, fmt.Errorf("over-read")
	}

	cr.header = &ph
	if cr.r.n == cr.totalSize {
		cr.atLastPage = true
	}
	cr.dataPageOffset += (cr.r.n - n)

	return data, nil
}

// Next advances the reader to the next value, which then will be available
// through ByteArray() method. It returns false when the reading stops, either by
// reaching the end of the input or an error.
func (cr *ByteArrayColumnChunkReader) Next() bool {
	if cr.err != nil {
		return false
	}

	if cr.decoder.data == nil {
		if cr.atLastPage {
			// TODO: may be set error if trying to read past the end of the column chunk
			return false
		}

		// read next page
		data, err := cr.readDataPage()
		if err != nil {
			// TODO: handle EOF
			cr.err = err
			return false
		}
		//fmt.Printf("%v\n", data)
		start := 0
		// TODO: it looks like thrift README is incorrect
		// first R then D
		if cr.rDecoder != nil {
			// decode definition levels data
			// TODO: uint32 -> int overflow
			// TODO: error handing
			n := int(binary.LittleEndian.Uint32(data[:4]))
			cr.rDecoder.init(data[4 : n+4])
			start = n + 4
		}
		if cr.dDecoder != nil {
			// decode repetition levels data
			// TODO: uint32 -> int overflow
			// TODO: error handing
			n := int(binary.LittleEndian.Uint32(data[start : start+4]))
			cr.dDecoder.init(data[start+4 : start+n+4])
			start += n + 4
		}
		cr.decoder.init(data[start:])
		cr.atStartOfPage = true
	} else {
		cr.atStartOfPage = false
	}

	pageNumValues := int(cr.header.DataPageHeader.NumValues)

	// TODO: hasNext and error checking
	if cr.dDecoder != nil {
		d, _ := cr.dDecoder.next()
		cr.curLevels.D = int(d)
	}
	if cr.rDecoder != nil {
		r, _ := cr.rDecoder.next()
		cr.curLevels.R = int(r)
	}
	if cr.curLevels.D == cr.maxLevels.D {
		cr.curValue, cr.err = cr.decoder.next()
		if cr.err != nil {
			return false
		}
	} else {
		cr.curValue = nil // just to be deterministic
	}

	cr.valuesRead++
	cr.pageValuesRead++
	if cr.pageValuesRead >= pageNumValues {
		cr.pageValuesRead = 0
		cr.decoder.data = nil
	}

	return true
}

// Skip behaves like Next but it may skip decoding the current value.
// TODO: think about use case and implement if necessary
func (cr *ByteArrayColumnChunkReader) Skip() bool {
	panic("nyi")
}

// NextPage advances the reader to the first value of the next page. It behaves as Next.
// TODO: think about how this method can be used to skip unreadable/corrupted pages
func (cr *ByteArrayColumnChunkReader) NextPage() bool {
	if cr.err != nil {
		return false
	}
	cr.atStartOfPage = true
	// TODO: implement
	return false
}

// Levels returns repetition and definition levels of the most recent value
// generated by a call to Next or NextPage.
func (cr *ByteArrayColumnChunkReader) Levels() Levels {
	return cr.curLevels
}

// ByteArray returns the most recent value generated by a call to Next or NextPage.
func (cr *ByteArrayColumnChunkReader) ByteArray() []byte {
	return cr.curValue
}

// Value returns ByteArray()
func (cr *ByteArrayColumnChunkReader) Value() interface{} {
	return cr.ByteArray()
}

// Err returns the first non-EOF error that was encountered by the reader.
func (cr *ByteArrayColumnChunkReader) Err() error {
	return cr.err
}
*/
