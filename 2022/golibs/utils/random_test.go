package utils

import (
	rand2 "crypto/rand"
	"encoding/binary"
	"testing"
)

func TestRand(t *testing.T) {
	t.Run("常规测试", func(t *testing.T) {
		s := Rand(32, RandMix)
		t.Log(s)
		if len(s) != 32 {
			t.Error("预期不符")
		}
	})
}

func BenchmarkRand(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var s, s1 string
		for pb.Next() {
			s = Rand(16, RandMix)
			if s == s1 {
				b.Error("预期不符")
			}
			s1 = s
		}
	})
}

func BenchmarkRand1(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var s, s1 int32
		for pb.Next() {
			binary.Read(rand2.Reader, binary.LittleEndian, &s)
			//s, _ = rand2.Int(rand2.Reader, big.NewInt(100))
			b.Log(s)
			if s == s1 {
				b.Error("预期不符")
			}
			s1 = s
		}
	})
}
