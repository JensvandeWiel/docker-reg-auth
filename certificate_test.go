package registry

import (
	"testing"
)

func TestLoadCertificateAndKey(t *testing.T) {

	pubKey, privKey, err := LoadCertificateAndKey("../.devcerts/RootCA.crt", "../.devcerts/RootCA.key")
	if err != nil {
		t.Errorf("LoadCertificateAndKey() error = %v", err)
		return
	}

	if pubKey == nil {
		t.Errorf("LoadCertificateAndKey() pubKey is nil")
	}

	if privKey == nil {
		t.Errorf("LoadCertificateAndKey() privKey is nil")
	}

	pubKey.MarshalJSON()
}
