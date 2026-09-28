package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"path/filepath"
	"testing"
	"time"

	"github.com/karamble/dcrgaming-orbitgolf/internal/bridgeconn"
)

func pemPair(t *testing.T) (string, string) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "dcrminigolf"},
		NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour)}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})),
		string(pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: raw}))
}

// Changing the address or network needs nothing pasted again, and a first save
// still needs all three credentials.
func TestSavingKeepsCredentialsLeftBlank(t *testing.T) {
	a := &app{dir: t.TempDir(), practice: -1, jobs: make(chan error, 1)}
	defer a.close()
	blank := &bridgeconn.Config{Network: "mainnet", Host: "127.0.0.1", Port: "8443"}
	if _, err := a.call(request{Method: "save_config", Config: blank}); err == nil {
		t.Fatal("a first save without credentials was accepted")
	}
	cert, key := pemPair(t)
	full := &bridgeconn.Config{Network: "mainnet", Host: "127.0.0.1", Port: "8443", ClientCert: cert, ClientKey: key, BridgeCert: cert}
	if _, err := a.call(request{Method: "save_config", Config: full}); err != nil {
		t.Fatal(err)
	}
	moved := &bridgeconn.Config{Network: "mainnet", Host: "95.216.110.66", Port: "8443"}
	if _, err := a.call(request{Method: "save_config", Config: moved}); err != nil {
		t.Fatalf("changing only the address needed the credentials again: %v", err)
	}
	got, err := bridgeconn.Load(filepath.Join(a.dir, "bridge.json"))
	if err != nil {
		t.Fatal(err)
	}
	if got.Host != "95.216.110.66" || got.ClientKey != key || got.ClientCert != cert || got.BridgeCert != cert {
		t.Fatalf("saved %+v", got.Host)
	}
}
