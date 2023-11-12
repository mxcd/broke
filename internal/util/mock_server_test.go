package util

import "testing"

func TestListLimitOffset(t *testing.T) {
	data := []string{"a", "b", "c", "d", "e", "f"}

	// Test offset >= len(data)
	result := ListLimitOffset(data, 1, 6)
	if len(result) != 0 {
		t.Errorf("Expected empty list, got %v", result)
	}

	// Test offset+limit > len(data)
	result = ListLimitOffset(data, 2, 4)
	if len(result) != 2 {
		t.Errorf("Expected list of length 2, got %v", result)
	}
	if result[0] != "e" || result[1] != "f" {
		t.Errorf("Expected list to be [e, f], got %v", result)
	}

	// Test offset+limit <= len(data)
	result = ListLimitOffset(data, 2, 2)
	if len(result) != 2 {
		t.Errorf("Expected list of length 2, got %v", result)
	}
	if result[0] != "c" || result[1] != "d" {
		t.Errorf("Expected list to be [c, d], got %v", result)
	}
}
