package crypto

import (
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	type hash_test_entry struct {
		Type     HashType
		String   string
		Expected uint32
	}

	entries := []hash_test_entry{
		{
			Type:     HashEncryptKey,
			String:   "(hash table)",
			Expected: 0xc3af3770,
		},
	}

	for _, e := range entries {
		result := HashString(e.String, e.Type)
		if e.Expected != result {
			t.Fatal("expected", e.Expected, "got", result)
		}
		fmt.Printf("0x%08x", result)
	}
}
