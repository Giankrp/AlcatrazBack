// Alcatraz - Secure open source Password Manager and Storage System
// Copyright (C) 2026 Gian Carlo Ruiz Patiño
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
