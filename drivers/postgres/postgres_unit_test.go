package postgres_test

import (
	"testing"

	"github.com/xo/usql/drivers"
	_ "github.com/xo/usql/drivers/postgres"
)

func TestOpenGaussAliasRegistered(t *testing.T) {
	if !drivers.Registered("opengauss") {
		t.Fatal("expected opengauss driver alias to be registered")
	}

	available := drivers.Available()
	pg, ok := available["postgres"]
	if !ok {
		t.Fatal("expected postgres driver to be registered")
	}
	og, ok := available["opengauss"]
	if !ok {
		t.Fatal("expected opengauss alias to resolve to a driver")
	}

	if pg.Name != "pq" {
		t.Fatalf("expected postgres go driver name to stay pq, got %q", pg.Name)
	}
	if og.Name != pg.Name {
		t.Fatalf("expected opengauss alias to reuse postgres driver, got %q want %q", og.Name, pg.Name)
	}
	if og.LexerName != "postgres" {
		t.Fatalf("expected opengauss alias to reuse postgres lexer, got %q", og.LexerName)
	}
}
