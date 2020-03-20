package gocaptcha

import (
	"testing"
)

func TestNew(t *testing.T) {
	data,err := New(&Options{
		Curve:2,
	})
	if err != nil {
		t.Error(err)
	}

	t.Logf("%#v",data)
}

func BenchmarkNew(b *testing.B) {
	data,err := New(&Options{
		Curve:2,
	})
	if err != nil {
		b.Error(err)
	}
	b.Logf("%#v",data)
}
