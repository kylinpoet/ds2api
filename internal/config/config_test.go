package config

import (
	"strings"
	"testing"
)

func TestAccountIdentifierFallsBackToTokenHash(t *testing.T) {
	acc := Account{Token: "example-token-value"}
	id := acc.Identifier()
	if !strings.HasPrefix(id, "token:") {
		t.Fatalf("expected token-prefixed identifier, got %q", id)
	}
	if len(id) != len("token:")+16 {
		t.Fatalf("unexpected identifier length: %d (%q)", len(id), id)
	}
}

func TestStoreFindAccountWithTokenOnlyIdentifier(t *testing.T) {
	t.Setenv("DS2API_CONFIG_JSON", `{
		"keys":["k1"],
		"accounts":[{"token":"token-only-account"}]
	}`)

	store := LoadStore()
	accounts := store.Accounts()
	if len(accounts) != 1 {
		t.Fatalf("expected 1 account, got %d", len(accounts))
	}
	id := accounts[0].Identifier()
	if id == "" {
		t.Fatalf("expected synthetic identifier for token-only account")
	}
	found, ok := store.FindAccount(id)
	if !ok {
		t.Fatalf("expected FindAccount to locate token-only account by synthetic id")
	}
	if found.Token != "token-only-account" {
		t.Fatalf("unexpected token value: %q", found.Token)
	}
}

func TestStoreUpdateAccountTokenKeepsOldAndNewIdentifierResolvable(t *testing.T) {
	t.Setenv("DS2API_CONFIG_JSON", `{
		"accounts":[{"token":"old-token"}]
	}`)

	store := LoadStore()
	before := store.Accounts()
	if len(before) != 1 {
		t.Fatalf("expected 1 account, got %d", len(before))
	}
	oldID := before[0].Identifier()
	if oldID == "" {
		t.Fatal("expected old identifier")
	}
	if err := store.UpdateAccountToken(oldID, "new-token"); err != nil {
		t.Fatalf("update token failed: %v", err)
	}

	after := store.Accounts()
	newID := after[0].Identifier()
	if newID == "" || newID == oldID {
		t.Fatalf("expected changed identifier, old=%q new=%q", oldID, newID)
	}
	if got, ok := store.FindAccount(newID); !ok || got.Token != "new-token" {
		t.Fatalf("expected find by new identifier")
	}
	if got, ok := store.FindAccount(oldID); !ok || got.Token != "new-token" {
		t.Fatalf("expected find by old identifier alias")
	}
}
