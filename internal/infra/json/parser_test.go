package json

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func TestParseInput_ValidOperations(t *testing.T) {
	inputJSON := `[{"operation":"buy","unit-cost":10.00,"quantity":100},{"operation":"sell","unit-cost":20.00,"quantity":50}]`
	inputBytes := []byte(inputJSON)
	expected := []Operation{
		{Operation: "buy", UnitCost: 10.00, Quantity: 100},
		{Operation: "sell", UnitCost: 20.00, Quantity: 50},
	}

	result, err := ParseInput(inputBytes)

	if err != nil {
		t.Fatalf("Assertion failed: expected no error, but got: %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Assertion failed: result = %v, want %v", result, expected)
	}
}

func TestParseInput_EmptyArray(t *testing.T) {
	inputJSON := `[]`
	inputBytes := []byte(inputJSON)

	result, err := ParseInput(inputBytes)

	if err != nil {
		t.Fatalf("Assertion failed: expected no error for empty array, but got: %v", err)
	}
	if result == nil {
		t.Errorf("Assertion failed: result for empty array was nil, want non-nil empty slice")
	} else if len(result) != 0 {
		t.Errorf("Assertion failed: result for empty array should have len 0, but got %d", len(result))
	}
}

func TestParseInput_InvalidJsonSyntax(t *testing.T) {
	inputJSON := `[{"operation":"buy","unit-cost":10.00,"quantity":100},`
	inputBytes := []byte(inputJSON)
	var syntaxError *json.SyntaxError

	_, err := ParseInput(inputBytes)

	if err == nil {
		t.Fatalf("Assertion failed: expected an error for invalid JSON syntax, but got nil")
	}
	if !errors.As(err, &syntaxError) {
		t.Errorf("Assertion failed: expected error of type *json.SyntaxError, but got type %T: %v", err, err)
	}
}

func TestParseInput_WrongJsonStructure_ObjectInsteadOfArray(t *testing.T) {
	inputJSON := `{"operation":"buy","unit-cost":10.00,"quantity":100}`
	inputBytes := []byte(inputJSON)
	var unmarshalTypeError *json.UnmarshalTypeError

	_, err := ParseInput(inputBytes)

	if err == nil {
		t.Fatalf("Assertion failed: expected an error for object-instead-of-array mismatch, but got nil")
	}
	if !errors.As(err, &unmarshalTypeError) {
		t.Errorf("Assertion failed: expected error of type *json.UnmarshalTypeError, but got type %T: %v", err, err)
	}
}

func TestParseInput_WrongJsonStructure_ArrayOfStrings(t *testing.T) {
	inputJSON := `["buy", "sell"]`
	inputBytes := []byte(inputJSON)
	var unmarshalTypeError *json.UnmarshalTypeError

	_, err := ParseInput(inputBytes)

	if err == nil {
		t.Fatalf("Assertion failed: expected an error for array-of-strings mismatch, but got nil")
	}
	if !errors.As(err, &unmarshalTypeError) {
		t.Errorf("Assertion failed: expected error of type *json.UnmarshalTypeError, but got type %T: %v", err, err)
	}
}

func TestParseInput_MissingField(t *testing.T) {
	inputJSON := `[{"operation":"buy","quantity":100}]`
	inputBytes := []byte(inputJSON)
	expected := []Operation{
		{Operation: "buy", UnitCost: 0.0, Quantity: 100},
	}

	result, err := ParseInput(inputBytes)

	if err != nil {
		t.Fatalf("Assertion failed: expected no error for missing field, but got: %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Assertion failed: result for missing field = %v, want %v", result, expected)
	}
}

func TestParseInput_WrongFieldType(t *testing.T) {
	inputJSON := `[{"operation":"buy","unit-cost":"expensive","quantity":100}]`
	inputBytes := []byte(inputJSON)
	var unmarshalTypeError *json.UnmarshalTypeError

	_, err := ParseInput(inputBytes)

	if err == nil {
		t.Fatalf("Assertion failed: expected an error for wrong field type, but got nil")
	}
	if !errors.As(err, &unmarshalTypeError) {
		t.Errorf("Assertion failed: expected error of type *json.UnmarshalTypeError, but got type %T: %v", err, err)
	}
}
