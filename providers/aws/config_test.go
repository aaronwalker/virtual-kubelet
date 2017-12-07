package aws

import (
	"bytes"
	"strings"
	"testing"
)

const cfg = `
Region = "us-east-1"
Cluster = "mycluster"
CPU = "100"
Memory = "100Gi"
Pods = "20"`

func TestConfig(t *testing.T) {
	br := bytes.NewReader([]byte(cfg))
	var p ECSProvider
	err := p.loadConfig(br)
	if err != nil {
		t.Error(err)
	}
	wanted := "us-east-1"
	if p.region != wanted {
		t.Errorf("Wanted %s, got %s.", wanted, p.region)
	}

	wanted = "mycluster"
	if p.cluster != wanted {
		t.Errorf("Wanted %s, got %s.", wanted, p.region)
	}

	wanted = "100"
	if p.cpu != wanted {
		t.Errorf("Wanted %s, got %s.", wanted, p.cpu)
	}

	wanted = "100Gi"
	if p.memory != wanted {
		t.Errorf("Wanted %s, got %s.", wanted, p.memory)
	}

	wanted = "20"
	if p.pods != wanted {
		t.Errorf("Wanted %s, got %s.", wanted, p.pods)
	}
}

const cfgBad = `
Region = "us-east-1"
OperatingSystem = "noop"`

func TestBadConfig(t *testing.T) {
	br := bytes.NewReader([]byte(cfgBad))
	var p ECSProvider
	err := p.loadConfig(br)
	if err == nil {
		t.Fatal("expected loadConfig to fail with bad operating system option")
	}

	if !strings.Contains(err.Error(), "is not a valid operating system") {
		t.Fatalf("expected loadConfig to fail with 'is not a valid operating system' but got: %v", err)

	}
}

const defCfg = `
Region = "us-east-1"`

func TestDefaultedConfig(t *testing.T) {
	br := bytes.NewReader([]byte(defCfg))
	var p ECSProvider
	err := p.loadConfig(br)
	if err != nil {
		t.Error(err)
	}
	// Test that defaults work with no settings in config.
	wanted := "default"
	if p.cluster != wanted {
		t.Errorf("Wanted default %s, got %s.", wanted, p.cpu)
	}

	wanted = "20"
	if p.cpu != wanted {
		t.Errorf("Wanted default %s, got %s.", wanted, p.cpu)
	}

	wanted = "100Gi"
	if p.memory != wanted {
		t.Errorf("Wanted default %s, got %s.", wanted, p.memory)
	}

	wanted = "20"
	if p.pods != wanted {
		t.Errorf("Wanted default %s, got %s.", wanted, p.pods)
	}
}
