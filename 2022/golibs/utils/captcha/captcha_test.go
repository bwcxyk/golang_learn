package captcha

import "testing"

func TestGenerate(t *testing.T) {
	t.Log(Generate())
}

func TestVerify(t *testing.T) {
	t.Log(Verify("tcJ8dRZszmrHRFsYjx79", "58768"))
}
