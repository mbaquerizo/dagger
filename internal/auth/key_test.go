package auth

import "testing"

func TestHashKey(t *testing.T) {
	got := HashKey("")
	want := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	if got != want {
		t.Errorf("HashKey(\"\") = %q, want %q", got, want)
	}

	got2 := HashKey("dgr_test-key-123")

	if len(got2) != 64 {
		t.Errorf("HashKey() returned %d chars, want 64", len(got2))
	}
}

func TestPrefix(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "long key", input: "dgr_iliketurtles123", want: "dgr_ilik"},
		{name: "short key", input: "abc", want: "abc"},
		{name: "empty", input: "", want: ""},
		{name: "exactly 8", input: "ilikecat", want: "ilikecat"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Prefix(tt.input)
			want := tt.want

			if got != want {
				t.Errorf("Prefix() = %q, want %q", got, want)
			}
		})
	}
}

func TestGenerateKey(t *testing.T) {
	rawKey, hash, prefix, err := GenerateKey()

	if err != nil {
		t.Fatalf("GenerateKey() returned error: %v", err)
	}

	if len(rawKey) < 4 || rawKey[:4] != KeyPrefix {
		t.Errorf("GenerateKey() rawKey = %q, want prefix %q", rawKey, KeyPrefix)
	}

	if len(hash) != 64 {
		t.Errorf("GenerateKey() hash = %q, want 64 characters", hash)
	}

	if prefix != Prefix(rawKey) {
		t.Errorf("GenerateKey() prefix = %q, Prefix(rawKey) = %q", prefix, Prefix(rawKey))
	}

	if hash != HashKey(rawKey) {
		t.Errorf("GenerateKey() hash = %q, HashKey(rawKey) = %q", hash, HashKey(rawKey))
	}
}
