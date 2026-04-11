package security

import "testing"

func TestHashPasswordAndVerifyPassword(t *testing.T) {
	password := "correct horse battery staple"

	encoded1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}

	ok, err := VerifyPassword(password, encoded1)
	if err != nil {
		t.Fatalf("VerifyPassword error: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true for correct password")
	}

	ok, err = VerifyPassword("wrong-password", encoded1)
	if err != nil {
		t.Fatalf("VerifyPassword error (wrong password): %v", err)
	}
	if ok {
		t.Fatalf("expected ok=false for wrong password")
	}

	encoded2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword error (second): %v", err)
	}
	if encoded1 == encoded2 {
		t.Fatalf("expected different hashes due to random salt")
	}
}

func TestVerifyPasswordInvalidFormat(t *testing.T) {
	_, err := VerifyPassword("x", "not-a-hash")
	if err == nil {
		t.Fatalf("expected error for invalid hash format")
	}
}

func TestNeedsRehash(t *testing.T) {
	password := "pw"

	encoded, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}

	needs, err := NeedsRehash(encoded, DefaultParams)
	if err != nil {
		t.Fatalf("NeedsRehash error: %v", err)
	}
	if needs {
		t.Fatalf("expected needs=false with DefaultParams")
	}

	changed := DefaultParams
	changed.Memory = changed.Memory + 1

	needs, err = NeedsRehash(encoded, changed)
	if err != nil {
		t.Fatalf("NeedsRehash error (changed): %v", err)
	}
	if !needs {
		t.Fatalf("expected needs=true after changing params")
	}
}
