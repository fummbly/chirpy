package auth

import "testing"

func TestHashingPassword(t *testing.T) {

	password := "test"

	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("Failed to hash the password: %v", err)
	}

	t.Logf("Got %s out of password %s", hash, password)
}

func TestCompareHash(t *testing.T) {
	password_hash, err := HashPassword("test")
	if err != nil {
		t.Errorf("Failed to create hash: %v", err)
	}

	err = CheckPasswordHash("test", password_hash)
	if err != nil {
		t.Errorf("Failed: passwords are supposed to match got: %v", err)
	}

}

func TestWrongPassword(t *testing.T) {
	password_hash, err := HashPassword("test")
	if err != nil {
		t.Errorf("Failed to create hash: %v", err)
	}

	err = CheckPasswordHash("another", password_hash)
	if err == nil {
		t.Error("Failed: passwords are not supposed to match")
	}
}
