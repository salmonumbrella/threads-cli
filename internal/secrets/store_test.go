package secrets

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/99designs/keyring"
)

// mockKeyring implements keyring.Keyring for testing
type mockKeyring struct {
	items      map[string]keyring.Item
	failGet    bool
	failSet    bool
	failRemove bool
	failKeys   bool
}

func newMockKeyring() *mockKeyring {
	return &mockKeyring{
		items: make(map[string]keyring.Item),
	}
}

func (m *mockKeyring) Get(key string) (keyring.Item, error) {
	if m.failGet {
		return keyring.Item{}, errors.New("mock get error")
	}
	item, ok := m.items[key]
	if !ok {
		return keyring.Item{}, keyring.ErrKeyNotFound
	}
	return item, nil
}

func (m *mockKeyring) Set(item keyring.Item) error {
	if m.failSet {
		return errors.New("mock set error")
	}
	m.items[item.Key] = item
	return nil
}

func (m *mockKeyring) Remove(key string) error {
	if m.failRemove {
		return errors.New("mock remove error")
	}
	delete(m.items, key)
	return nil
}

func (m *mockKeyring) Keys() ([]string, error) {
	if m.failKeys {
		return nil, errors.New("mock keys error")
	}
	var keys []string
	for key := range m.items {
		keys = append(keys, key)
	}
	return keys, nil
}

func (m *mockKeyring) GetMetadata(key string) (keyring.Metadata, error) {
	return keyring.Metadata{}, nil
}

// Tests for Credentials methods

func TestCredentials_IsExpired(t *testing.T) {
	tests := []struct {
		name     string
		creds    Credentials
		expected bool
	}{
		{
			name:     "zero expiry - not expired",
			creds:    Credentials{ExpiresAt: time.Time{}},
			expected: false,
		},
		{
			name:     "future expiry - not expired",
			creds:    Credentials{ExpiresAt: time.Now().Add(time.Hour)},
			expected: false,
		},
		{
			name:     "past expiry - expired",
			creds:    Credentials{ExpiresAt: time.Now().Add(-time.Hour)},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.creds.IsExpired()
			if result != tt.expected {
				t.Errorf("expected IsExpired() = %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCredentials_IsExpiringSoon(t *testing.T) {
	tests := []struct {
		name     string
		creds    Credentials
		within   time.Duration
		expected bool
	}{
		{
			name:     "zero expiry - not expiring soon",
			creds:    Credentials{ExpiresAt: time.Time{}},
			within:   time.Hour,
			expected: false,
		},
		{
			name:     "far future - not expiring soon",
			creds:    Credentials{ExpiresAt: time.Now().Add(24 * time.Hour)},
			within:   time.Hour,
			expected: false,
		},
		{
			name:     "within window - expiring soon",
			creds:    Credentials{ExpiresAt: time.Now().Add(30 * time.Minute)},
			within:   time.Hour,
			expected: true,
		},
		{
			name:     "already expired - expiring soon",
			creds:    Credentials{ExpiresAt: time.Now().Add(-time.Minute)},
			within:   time.Hour,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.creds.IsExpiringSoon(tt.within)
			if result != tt.expected {
				t.Errorf("expected IsExpiringSoon(%v) = %v, got %v", tt.within, tt.expected, result)
			}
		})
	}
}

func TestCredentials_DaysUntilExpiry(t *testing.T) {
	tests := []struct {
		name      string
		creds     Credentials
		checkFunc func(float64) bool
	}{
		{
			name:  "zero expiry returns -1",
			creds: Credentials{ExpiresAt: time.Time{}},
			checkFunc: func(d float64) bool {
				return d == -1
			},
		},
		{
			name:  "24 hours returns ~1 day",
			creds: Credentials{ExpiresAt: time.Now().Add(24 * time.Hour)},
			checkFunc: func(d float64) bool {
				return d > 0.9 && d < 1.1
			},
		},
		{
			name:  "expired returns negative",
			creds: Credentials{ExpiresAt: time.Now().Add(-48 * time.Hour)},
			checkFunc: func(d float64) bool {
				return d < -1.9 && d > -2.1
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.creds.DaysUntilExpiry()
			if !tt.checkFunc(result) {
				t.Errorf("DaysUntilExpiry() = %v, unexpected value", result)
			}
		})
	}
}

// Tests for normalizeName

func TestNormalizeName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"test", "test"},
		{"TEST", "test"},
		{"Test", "test"},
		{"  test  ", "test"},
		{"  TEST  ", "test"},
		{"", ""},
		{"  ", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeName(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeName(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

// Tests for KeyringStore

func TestKeyringStore_Set(t *testing.T) {
	t.Run("successful set", func(t *testing.T) {
		mock := newMockKeyring()
		store := &KeyringStore{ring: mock, warnedAccounts: make(map[string]bool)}

		creds := Credentials{
			Name:        "test",
			AccessToken: "token123",
			UserID:      "user123",
		}

		err := store.Set("test", creds)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify item was stored
		item, err := mock.Get("account:test")
		if err != nil {
			t.Fatalf("item not found in keyring: %v", err)
		}

		var stored storedCredentials
		if err := json.Unmarshal(item.Data, &stored); err != nil {
			t.Fatalf("failed to unmarshal stored data: %v", err)
		}

		if stored.AccessToken != "token123" {
			t.Errorf("expected AccessToken 'token123', got %q", stored.AccessToken)
		}
	})

	t.Run("empty name error", func(t *testing.T) {
		mock := newMockKeyring()
		store := &KeyringStore{ring: mock, warnedAccounts: make(map[string]bool)}

		creds := Credentials{AccessToken: "token123"}
		err := store.Set("", creds)
		if err == nil {
			t.Error("expected error for empty name")
		}
	})

	t.Run("empty token error", func(t *testing.T) {
		mock := newMockKeyring()
		store := &KeyringStore{ring: mock, warnedAccounts: make(map[string]bool)}

		creds := Credentials{AccessToken: ""}
		err := store.Set("test", creds)
		if err == nil {
			t.Error("expected error for empty token")
		}
	})

	t.Run("normalizes name", func(t *testing.T) {
		mock := newMockKeyring()
		store := &KeyringStore{ring: mock, warnedAccounts: make(map[string]bool)}

		creds := Credentials{AccessToken: "token123"}
		err := store.Set("  TEST  ", creds)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should be stored with normalized name
		_, err = mock.Get("account:test")
		if err != nil {
			t.Error("expected item stored with normalized name 'test'")
		}
	})

	t.Run("sets CreatedAt if zero", func(t *testing.T) {
		mock := newMockKeyring()
		store := &KeyringStore{ring: mock, warnedAccounts: make(map[string]bool)}

		creds := Credentials{
			AccessToken: "token123",
			CreatedAt:   time.Time{}, // zero value
		}

		err := store.Set("test", creds)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		item, _ := mock.Get("account:test")
		var stored storedCredentials
		_ = json.Unmarshal(item.Data, &stored)

		if stored.CreatedAt.IsZero() {
			t.Error("expected CreatedAt to be set")
		}
	})
}

func TestKeyringStore_Get(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		mock := newMockKeyring()
		store := &KeyringStore{ring: mock, warnedAccounts: make(map[string]bool)}

		// Pre-store credentials
		stored := storedCredentials{
			AccessToken: "token123",
			UserID:      "user123",
			Username:    "testuser",
		}
		data, _ := json.Marshal(stored)
		_ = mock.Set(keyring.Item{Key: "account:test", Data: data})

		creds, err := store.Get("test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if creds.AccessToken != "token123" {
			t.Errorf("expected AccessToken 'token123', got %q", creds.AccessToken)
		}
		if creds.Name != "test" {
			t.Errorf("expected Name 'test', got %q", creds.Name)
		}
	})

	t.Run("not found error", func(t *testing.T) {
		mock := newMockKeyring()
		store := &KeyringStore{ring: mock, warnedAccounts: make(map[string]bool)}

		_, err := store.Get("nonexistent")
		if err == nil {
			t.Error("expected error for nonexistent account")
		}
	})

	t.Run("normalizes name", func(t *testing.T) {
		mock := newMockKeyring()
		store := &KeyringStore{ring: mock, warnedAccounts: make(map[string]bool)}

		// Store with lowercase
		stored := storedCredentials{AccessToken: "token123"}
		data, _ := json.Marshal(stored)
		_ = mock.Set(keyring.Item{Key: "account:test", Data: data})

		// Get with uppercase
		creds, err := store.Get("TEST")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if creds.AccessToken != "token123" {
			t.Error("expected to find credentials with normalized name")
		}
	})

	t.Run("get error", func(t *testing.T) {
		mock := newMockKeyring()
		mock.failGet = true
		store := &KeyringStore{ring: mock, warnedAccounts: make(map[string]bool)}

		_, err := store.Get("test")
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestKeyringStore_Delete(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		mock := newMockKeyring()
		store := &KeyringStore{ring: mock, warnedAccounts: make(map[string]bool)}

		// Pre-store an item
		_ = mock.Set(keyring.Item{Key: "account:test", Data: []byte("{}")})

		err := store.Delete("test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify deletion
		_, err = mock.Get("account:test")
		if err != keyring.ErrKeyNotFound {
			t.Error("expected item to be deleted")
		}
	})

	t.Run("normalizes name", func(t *testing.T) {
		mock := newMockKeyring()
		store := &KeyringStore{ring: mock, warnedAccounts: make(map[string]bool)}

		_ = mock.Set(keyring.Item{Key: "account:test", Data: []byte("{}")})

		err := store.Delete("  TEST  ")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		_, err = mock.Get("account:test")
		if err != keyring.ErrKeyNotFound {
			t.Error("expected item to be deleted with normalized name")
		}
	})
}

func TestKeyringStore_Keys(t *testing.T) {
	t.Run("returns account names", func(t *testing.T) {
		mock := newMockKeyring()
		store := &KeyringStore{ring: mock, warnedAccounts: make(map[string]bool)}

		// Pre-store some items
		_ = mock.Set(keyring.Item{Key: "account:test1", Data: []byte("{}")})
		_ = mock.Set(keyring.Item{Key: "account:test2", Data: []byte("{}")})
		_ = mock.Set(keyring.Item{Key: "other:key", Data: []byte("{}")}) // should be filtered

		keys, err := store.Keys()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(keys) != 2 {
			t.Errorf("expected 2 keys, got %d", len(keys))
		}

		// Check that account names are returned without prefix
		found := make(map[string]bool)
		for _, key := range keys {
			found[key] = true
		}

		if !found["test1"] {
			t.Error("expected 'test1' in keys")
		}
		if !found["test2"] {
			t.Error("expected 'test2' in keys")
		}
	})

	t.Run("empty list", func(t *testing.T) {
		mock := newMockKeyring()
		store := &KeyringStore{ring: mock, warnedAccounts: make(map[string]bool)}

		keys, err := store.Keys()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(keys) != 0 {
			t.Errorf("expected empty list, got %d keys", len(keys))
		}
	})

	t.Run("keys error", func(t *testing.T) {
		mock := newMockKeyring()
		mock.failKeys = true
		store := &KeyringStore{ring: mock, warnedAccounts: make(map[string]bool)}

		_, err := store.Keys()
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestKeyringStore_List(t *testing.T) {
	mock := newMockKeyring()
	store := &KeyringStore{ring: mock, warnedAccounts: make(map[string]bool)}

	_ = mock.Set(keyring.Item{Key: "account:test", Data: []byte("{}")})

	list, err := store.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(list) != 1 || list[0] != "test" {
		t.Errorf("expected ['test'], got %v", list)
	}
}

// Test constants

func TestConstants(t *testing.T) {
	if serviceName != "threads-cli" {
		t.Errorf("expected serviceName 'threads-cli', got %q", serviceName)
	}

	if accountPrefix != "account:" {
		t.Errorf("expected accountPrefix 'account:', got %q", accountPrefix)
	}

	if rotationDays != 55 {
		t.Errorf("expected rotationDays 55, got %d", rotationDays)
	}
}

// Test struct fields

func TestCredentials_Fields(t *testing.T) {
	now := time.Now()
	creds := Credentials{
		Name:         "test",
		AccessToken:  "token123",
		UserID:       "user123",
		Username:     "testuser",
		ExpiresAt:    now.Add(time.Hour),
		CreatedAt:    now,
		ClientID:     "client123",
		ClientSecret: "secret123",
		RedirectURI:  "https://example.com/callback",
	}

	if creds.Name != "test" {
		t.Errorf("expected Name 'test', got %q", creds.Name)
	}
	if creds.AccessToken != "token123" {
		t.Errorf("expected AccessToken 'token123', got %q", creds.AccessToken)
	}
	if creds.ClientSecret != "secret123" {
		t.Errorf("expected ClientSecret 'secret123', got %q", creds.ClientSecret)
	}
}

func TestStoredCredentials_Fields(t *testing.T) {
	now := time.Now()
	stored := storedCredentials{
		AccessToken:  "token123",
		UserID:       "user123",
		Username:     "testuser",
		ExpiresAt:    now.Add(time.Hour),
		CreatedAt:    now,
		ClientID:     "client123",
		ClientSecret: "secret123",
		RedirectURI:  "https://example.com/callback",
	}

	if stored.AccessToken != "token123" {
		t.Errorf("expected AccessToken 'token123', got %q", stored.AccessToken)
	}
	if stored.ClientSecret != "secret123" {
		t.Errorf("expected ClientSecret 'secret123', got %q", stored.ClientSecret)
	}
}

func TestKeyringStore_Set_RingError(t *testing.T) {
	mock := newMockKeyring()
	mock.failSet = true
	store := &KeyringStore{ring: mock, warnedAccounts: make(map[string]bool)}

	creds := Credentials{AccessToken: "token123"}
	err := store.Set("test", creds)
	if err == nil {
		t.Error("expected error when ring.Set fails")
	}
}

func TestKeyringStore_Get_UnmarshalError(t *testing.T) {
	mock := newMockKeyring()
	store := &KeyringStore{ring: mock, warnedAccounts: make(map[string]bool)}

	// Store invalid JSON
	_ = mock.Set(keyring.Item{Key: "account:test", Data: []byte("not json")})

	_, err := store.Get("test")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestKeyringStore_Get_ExpiringWarning(t *testing.T) {
	mock := newMockKeyring()
	store := &KeyringStore{ring: mock, warnedAccounts: make(map[string]bool)}

	// Store credentials expiring in negative time (expired)
	stored := storedCredentials{
		AccessToken: "token123",
		ExpiresAt:   time.Now().Add(-24 * time.Hour), // Already expired
	}
	data, _ := json.Marshal(stored)
	_ = mock.Set(keyring.Item{Key: "account:expired", Data: data})

	// First get - should trigger warning logic (but not panic)
	creds, err := store.Get("expired")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if creds == nil {
		t.Fatal("expected credentials")
	}

	// Store credentials expiring soon (within rotation window)
	stored2 := storedCredentials{
		AccessToken: "token456",
		ExpiresAt:   time.Now().Add(2 * 24 * time.Hour), // 2 days
	}
	data2, _ := json.Marshal(stored2)
	_ = mock.Set(keyring.Item{Key: "account:expiringsoon", Data: data2})

	creds2, err := store.Get("expiringsoon")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if creds2 == nil {
		t.Fatal("expected credentials")
	}

	// Verify the warned flag is set after warning
	// (The warning happens when daysUntilExpiry < 0 for expired tokens)
}

func TestKeyringStore_Delete_Error(t *testing.T) {
	mock := newMockKeyring()
	mock.failRemove = true
	store := &KeyringStore{ring: mock, warnedAccounts: make(map[string]bool)}

	err := store.Delete("test")
	if err == nil {
		t.Error("expected error when ring.Remove fails")
	}
}
