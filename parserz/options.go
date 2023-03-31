package parserz

import "os"

type Option func(o *PkgBuilder)

func WithPath(path string) Option {
	return func(o *PkgBuilder) {
		o.path = path
	}
}

func WithFilter(filter func(info os.FileInfo) bool) Option {
	return func(o *PkgBuilder) {
		o.filter = filter
	}
}
