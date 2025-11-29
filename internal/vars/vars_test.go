package vars

import (
	"testing"
)

func TestExpand(t *testing.T) {
	vars := map[string]string{
		"NAME": "flux",
		"VER":  "1.0",
	}

	result := Expand("${NAME}-${VER}", vars)
	expected := "flux-1.0"

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestNestedExpand(t *testing.T) {
	vars := map[string]string{
		"A": "${B}",
		"B": "value",
	}

	err := ResolveVars(vars)
	if err != nil {
		t.Fatalf("ResolveVars error: %v", err)
	}

	if vars["A"] != "value" {
		t.Errorf("Expected value, got %s", vars["A"])
	}
}

func TestMergeVars(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	overlay := map[string]string{"B": "3", "C": "4"}

	result := MergeVars(base, overlay)

	if result["A"] != "1" {
		t.Error("Expected A=1")
	}
	if result["B"] != "3" {
		t.Error("Expected B=3 (overlay)")
	}
	if result["C"] != "4" {
		t.Error("Expected C=4")
	}
}
