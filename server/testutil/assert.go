package testutil

import "testing"

func AssertInterface(t *testing.T, expected, actual map[string]interface{}) {
	t.Helper()
	for k, v := range expected {
		if vStr, ok := v.(string); ok {
			if actualStr, ok := actual[k].(string); ok {
				if actualStr != vStr {
					t.Errorf("Expected %s to be %s, got %s", k, v, actualStr)
				}
			}
		} else if vBool, ok := v.(bool); ok {
			if actualBool, ok := actual[k].(bool); ok {
				if actualBool != vBool {
					t.Errorf("Expected %s to be %t, got %t", k, v, actualBool)
				}
			}
		} else if vinterface, ok := v.(map[string]interface{}); ok {
			if actualMap, ok := actual[k].(map[string]interface{}); ok {
				AssertInterface(t, vinterface, actualMap)
			} else {
				t.Errorf("Expected %s to be map[string]interface{}, got %T", k, actual[k])
			}
		} else if v == nil {
			if actual[k] != nil {
				t.Errorf("Expected %s to be nil, got %v", k, actual[k])
			}
		} else {
			t.Errorf("Expected %s to be string or map[string]interface{}, got %T", k, v)
		}
	}
}
