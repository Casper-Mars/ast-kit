package parserz

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPkgBuilder_Build(t *testing.T) {
	pkg, err := NewPkgBuilder().WithPkgPath("./tests/mgr").Build()
	assert.NoError(t, err)
	assert.Equal(t, "mgr", pkg.Name())
	assert.Equal(t, 2, len(pkg.Struct))
	assert.Equal(t, 5, len(pkg.Func))
	assert.Equal(t, "UserMgr", pkg.Func[0].Rec.Type.Name())
	assert.Equal(t, "\"tests/store\"", pkg.Func[0].Results[0].ImportPaths()[0])
	assert.Equal(t, "*store.User", pkg.Func[0].Results[0].String())
}