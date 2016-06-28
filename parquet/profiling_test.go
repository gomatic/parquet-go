package parquet

import "testing"
// import _ "net/http/pprof"

/*  This file is currently for exploration purposes.
    We repeat a test enough times to get a cpu profile, then
    utilize pprof tool to build a call graph.


Usage:
	$ go test -run=<Name of Test> -cpuprofile=cpu.prof
	$ go tool pprof --pdf parquet.test cpu.prof > <name to save the pdf of callgraph>

	* note: must have GraphViz installed ($ brew install graphviz)
*/


func TestProfileBooleanColumnChunkReader(t *testing.T) {
	for i := 0; i < 10000; i++ {
		TestBooleanColumnChunkReader(t)
	}
}

func TestProfileByteArrayColumnChunkReader(t *testing.T) {
	for i := 0; i < 10000; i++ {
		TestByteArrayColumnChunkReader(t)
	}
}

func TestProfileCreateSchema(t *testing.T) {
	for i := 0; i < 10000; i++ {
		TestCreateSchema(t)
	}
}

func TestProfileReadFileMetaData(t *testing.T) {
	for i := 0; i < 10000; i++ {
		TestReadFileMetaData(t)
	}
}

func TestProfileCreateSchemaFromFileMetaDataAndMarshal(t *testing.T) {
	for i := 0; i < 10000; i++ {
		TestCreateSchemaFromFileMetaDataAndMarshal(t)
	}
}

