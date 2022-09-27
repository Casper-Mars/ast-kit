package parserz

import (
	"fmt"
	"testing"
)

func TestPkgBuilder_Build(t *testing.T) {

	build, err := NewPkgBuilder().WithPkgPath("./tests/mgr").Build()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v\n", build)
}
