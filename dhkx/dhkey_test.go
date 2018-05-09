package dhkx_test

import (
	"fmt"
	"testing"

	"github.com/lzutao/godhchat/dhkx"
)

func TestAll(t *testing.T) {
	a, _ := dhkx.NewDHKey(0)
	b, _ := dhkx.NewDHKey(0)

	ga := a.PublicKey()
	gb := b.PublicKey()

	gab, _ := a.SharedSecretKey(gb)
	gba, _ := b.SharedSecretKey(ga)

	if gab.Cmp(gba) == 0 {
		fmt.Println("Shared keys match.")
		fmt.Printf("Key: %v\n", gab)
	} else {
		fmt.Println("Shared secrets didn't match!")
		fmt.Println("Shared secret A: ", gab)
		fmt.Println("Shared secret B: ", gba)
	}
}
