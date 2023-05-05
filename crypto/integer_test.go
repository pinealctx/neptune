package crypto

import (
	"testing"
	"time"
)

func TestIntCypher_DecryptU32(t *testing.T) {
	var t1 = time.Now()
	var ic = NewIntCypher([]byte{54, 179, 221, 82, 230, 144, 168, 47, 124, 130, 37, 240, 255, 53, 121, 80},
		[]byte{115, 141, 121, 11, 126, 146, 20, 188, 225, 177, 134, 227, 184, 148, 105, 66})
	for i := uint32(10000); i < 10000000; i++ {
		cypherU32N0(t, ic, i)
	}
	var t2 = time.Now()
	var dur = t2.Sub(t1)
	t.Log("duration:", dur, "average:", dur/time.Duration(10000000-10000))
}

func TestIntCypher_DecryptU32V2(t *testing.T) {
	var t1 = time.Now()
	var ic = NewIntCypher([]byte{54, 179, 221, 82, 230, 144, 168, 47, 124, 130, 37, 240, 255, 53, 121, 80},
		[]byte{115, 141, 121, 11, 126, 146, 20, 188, 225, 177, 134, 227, 184, 148, 105, 66})
	for i := uint32(10000); i < 10000000; i++ {
		cypherU32N1(t, ic, i)
	}
	var t2 = time.Now()
	var dur = t2.Sub(t1)
	t.Log("duration:", dur, "average:", dur/time.Duration(10000000-10000))
}

func TestIntCypher_DecryptU32V3(t *testing.T) {
	var t1 = time.Now()
	var ic = NewIntCypher([]byte{54, 179, 221, 82, 230, 144, 168, 47, 124, 130, 37, 240, 255, 53, 121, 80},
		[]byte{115, 141, 121, 11, 126, 146, 20, 188, 225, 177, 134, 227, 184, 148, 105, 66})
	for i := uint32(10000); i < 10000000; i++ {
		cypherU32N2(t, ic, i)
	}
	var t2 = time.Now()
	var dur = t2.Sub(t1)
	t.Log("duration:", dur, "average:", dur/time.Duration(10000000-10000))
}

func TestIntCypher_DecryptU64(t *testing.T) {
	var t1 = time.Now()
	var ic = NewIntCypher([]byte{54, 179, 221, 82, 230, 144, 168, 47, 124, 130, 37, 240, 255, 53, 121, 80},
		[]byte{115, 141, 121, 11, 126, 146, 20, 188, 225, 177, 134, 227, 184, 148, 105, 66})
	for i := uint64(10000); i < 10000000; i++ {
		cypherU64(t, ic, i)
	}
	var t2 = time.Now()
	var dur = t2.Sub(t1)
	t.Log("duration:", dur, "average:", dur/time.Duration(10000000-10000))
}

func TestIntCypher_EncryptU32(t *testing.T) {
	var start = uint32(100081)
	var ic = NewIntCypher([]byte{54, 179, 221, 82, 230, 144, 168, 47, 124, 130, 37, 240, 255, 53, 121, 80},
		[]byte{115, 141, 121, 11, 126, 146, 20, 188, 225, 177, 134, 227, 184, 148, 105, 66})
	for k := byte(0); k < 24; k++ {
		t.Log("k:", k, "------------------------")
		for i := uint32(0); i < 100; i++ {
			t.Log(start+i, "-->", ic.EncU32[k](start+i))
		}
		start = uint32(200000)
		for i := uint32(0); i < 100; i++ {
			t.Log(start+i, "-->", ic.EncU32[k](start+i))
		}
		t.Log("")
		t.Log("")
	}
}

func TestIntCypher_EncryptComp(t *testing.T) {
	var start = uint32(100081)
	var ic = NewIntCypher([]byte{54, 179, 221, 82, 230, 144, 168, 47, 124, 130, 37, 240, 255, 53, 121, 80},
		[]byte{115, 141, 121, 11, 126, 146, 20, 188, 225, 177, 134, 227, 184, 148, 105, 66})

	for i := uint32(0); i < 100; i++ {
		t.Log(start+i, "-->", ic.EncU32[20](start+i), calculateOutput(start+i))
	}
	start = uint32(200000)
	for i := uint32(0); i < 100; i++ {
		t.Log(start+i, "-->", ic.EncU32[20](start+i), calculateOutput(start+i))
	}

}

/*
func TestIntCypher_EncryptU32V2C1(t *testing.T) {
	var start = uint32(100081)
	var ic = NewIntCypher([]byte{54, 179, 221, 82, 230, 144, 168, 47, 124, 130, 37, 240, 255, 53, 121, 80},
		[]byte{115, 141, 121, 11, 126, 146, 20, 188, 225, 177, 134, 227, 184, 148, 105, 66})
	for i := uint32(0); i < 100; i++ {
		t.Log(start+i, "-->", ic.EncryptU32V2(start+i))
	}
	start = uint32(200000)
	for i := uint32(0); i < 100; i++ {
		t.Log(start+i, "-->", ic.EncryptU32V2(start+i))
	}
}

func TestIntCypher_EncryptU32V3C1(t *testing.T) {
	var start = uint32(100081)
	var ic = NewIntCypher([]byte{54, 179, 221, 82, 230, 144, 168, 47, 124, 130, 37, 240, 255, 53, 121, 80},
		[]byte{115, 141, 121, 11, 126, 146, 20, 188, 225, 177, 134, 227, 184, 148, 105, 66})
	for i := uint32(0); i < 100; i++ {
		t.Log(start+i, "-->", ic.EncryptU32V3(start+i))
	}
	start = uint32(200000)
	for i := uint32(0); i < 100; i++ {
		t.Log(start+i, "-->", ic.EncryptU32V3(start+i))
	}
}

func TestIntCypher_EncryptU32C2(t *testing.T) {
	var start = uint32(100081)
	var ic = NewIntCypher([]byte{115, 141, 121, 11, 126, 146, 20, 188, 225, 177, 134, 227, 184, 148, 105, 66},
		[]byte{54, 179, 221, 82, 230, 144, 168, 47, 124, 130, 37, 240, 255, 53, 121, 80})
	for i := uint32(0); i < 100; i++ {
		t.Log(start+i, "-->", ic.EncryptU32(start+i))
	}
	start = uint32(200000)
	for i := uint32(0); i < 100; i++ {
		t.Log(start+i, "-->", ic.EncryptU32(start+i))
	}
}

func TestIntCypher_EncryptU32C2V3(t *testing.T) {
	var start = uint32(100081)
	var ic = NewIntCypher([]byte{115, 141, 121, 11, 126, 146, 20, 188, 225, 177, 134, 227, 184, 148, 105, 66},
		[]byte{54, 179, 221, 82, 230, 144, 168, 47, 124, 130, 37, 240, 255, 53, 121, 80})
	for i := uint32(0); i < 100; i++ {
		t.Log(start+i, "-->", ic.EncryptU32V3(start+i))
	}
	start = uint32(200000)
	for i := uint32(0); i < 100; i++ {
		t.Log(start+i, "-->", ic.EncryptU32V3(start+i))
	}
}

func TestIntCypher_EncryptU32C3(t *testing.T) {
	var start = uint32(100081)
	var ic = NewIntCypher(make([]byte, 16),
		make([]byte, 16))
	for i := uint32(0); i < 100; i++ {
		t.Log(start+i, "-->", ic.EncryptU32(start+i))
	}
	start = uint32(200000)
	for i := uint32(0); i < 100; i++ {
		t.Log(start+i, "-->", ic.EncryptU32(start+i))
	}
}
*/

func cypherU32N0(t *testing.T, x *IntCypher, number uint32) {
	var y = x.EncU32[0](number)
	var z = x.DecU32[0](y)
	if z != number {
		t.Fatal("failed to cypher number:", number, "cyphered:", y, "decyphered:", z)
	}
}

func cypherU32N1(t *testing.T, x *IntCypher, number uint32) {
	var y = x.EncU32[1](number)
	var z = x.DecU32[1](y)
	if z != number {
		t.Fatal("failed to cypher number:", number, "cyphered:", y, "decyphered:", z)
	}
}

func cypherU32N2(t *testing.T, x *IntCypher, number uint32) {
	var y = x.EncU32[2](number)
	var z = x.DecU32[2](y)
	if z != number {
		t.Fatal("failed to cypher number:", number, "cyphered:", y, "decyphered:", z)
	}
}

func cypherU64(t *testing.T, x *IntCypher, number uint64) {
	var y = x.EncryptU64(number)
	var z = x.DecryptU64(y)
	if z != number {
		t.Fatal("failed to cypher number:", number, "cyphered:", y, "decyphered:", z)
	}
}

func calculateOutput(input uint32) uint32 {
	base := uint32(200000)
	mod := int(input-base) % 8
	diff := uint32(0)

	switch mod {
	case 0, 3, 5, 6:
		diff = uint32(16777216)
	case 1, 2, 4, 7:
		diff = uint32(33947648)
	}

	if mod == 0 || mod > 4 {
		return input*diff + 1992294400
	} else {
		return input*diff - 1992294400
	}
}
