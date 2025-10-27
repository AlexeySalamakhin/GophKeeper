// Package crypto содержит тесты для функций шифрования.
package crypto

import (
	"testing"
)

func TestEncryptDecryptPassword(t *testing.T) {
	password := "test_password_123"
	key := "test_encryption_key"

	// Шифрование пароля
	encrypted, err := EncryptPassword(password, key)
	if err != nil {
		t.Fatalf("Ошибка шифрования пароля: %v", err)
	}

	if encrypted == "" {
		t.Error("Зашифрованный пароль не должен быть пустым")
	}

	if encrypted == password {
		t.Error("Зашифрованный пароль не должен совпадать с исходным")
	}

	// Расшифровка пароля
	decrypted, err := DecryptPassword(encrypted, key)
	if err != nil {
		t.Fatalf("Ошибка расшифровки пароля: %v", err)
	}

	if decrypted != password {
		t.Errorf("Ожидался пароль %s, получен %s", password, decrypted)
	}
}

func TestEncryptDecryptPassword_DifferentKeys(t *testing.T) {
	password := "test_password_123"
	key1 := "key1"
	key2 := "key2"

	// Шифрование с первым ключом
	encrypted, err := EncryptPassword(password, key1)
	if err != nil {
		t.Fatalf("Ошибка шифрования пароля: %v", err)
	}

	// Попытка расшифровки с другим ключом должна вернуть ошибку
	_, err = DecryptPassword(encrypted, key2)
	if err == nil {
		t.Error("Расшифровка с неправильным ключом должна возвращать ошибку")
	}
}

func TestEncryptPassword_EmptyPassword(t *testing.T) {
	key := "test_key"

	encrypted, err := EncryptPassword("", key)
	if err != nil {
		t.Fatalf("Ошибка шифрования пустого пароля: %v", err)
	}

	decrypted, err := DecryptPassword(encrypted, key)
	if err != nil {
		t.Fatalf("Ошибка расшифровки пустого пароля: %v", err)
	}

	if decrypted != "" {
		t.Errorf("Ожидался пустой пароль, получен %s", decrypted)
	}
}
