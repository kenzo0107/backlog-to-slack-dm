package utility

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInSlice(t *testing.T) {
	expected := true
	slice := []string{"hoge", "moge"}
	exists, index, err := InSlice("hoge", slice)

	if diff := cmp.Diff(exists, expected); diff != "" {
		t.Errorf("wrong result : %s", diff)
	}
	if diff := cmp.Diff(index, 0); diff != "" {
		t.Errorf("wrong result : %s", diff)
	}
	if err != nil {
		t.Error("wrong result : err is not nil", err)
	}
}

func TestNotInSlice(t *testing.T) {
	expected := false
	slice := []string{"hoge", "moge"}
	exists, index, err := InSlice("bar", slice)

	if diff := cmp.Diff(exists, expected); diff != "" {
		t.Errorf("wrong result : %s", diff)
	}
	if diff := cmp.Diff(index, -1); diff != "" {
		t.Errorf("wrong result : %s", diff)
	}
	if err != nil {
		t.Error("wrong result : err is not nil", err)
	}
}

func TestInSliceInt(t *testing.T) {
	expected := true
	slice := []int{
		1,
		-1,
	}
	exists, index, err := InSlice(-1, slice)

	if diff := cmp.Diff(exists, expected); diff != "" {
		t.Errorf("wrong result : %s", diff)
	}
	if diff := cmp.Diff(index, 1); diff != "" {
		t.Errorf("wrong result : %s", diff)
	}
	if err != nil {
		t.Error("wrong result : err is not nil", err)
	}
}

func TestInSliceBool(t *testing.T) {
	expected := true
	slice := []bool{
		true,
		false,
		true,
		false,
	}
	exists, index, err := InSlice(false, slice)

	if diff := cmp.Diff(exists, expected); diff != "" {
		t.Errorf("wrong result : %s", diff)
	}
	if diff := cmp.Diff(index, 1); diff != "" {
		t.Errorf("wrong result : %s", diff)
	}
	if err != nil {
		t.Error("wrong result : err is not nil", err)
	}
}

func TestInNotSliceBool(t *testing.T) {
	expected := false
	slice := []int{
		1,
		2,
		3,
		4,
	}
	exists, index, err := InSlice(false, slice)

	if diff := cmp.Diff(exists, expected); diff != "" {
		t.Errorf("wrong result : %s", diff)
	}
	if diff := cmp.Diff(index, -1); diff != "" {
		t.Errorf("wrong result : %s", diff)
	}
	if err != nil {
		t.Error("wrong result : err is not nil", err)
	}
}

func TestInNotSlice(t *testing.T) {
	expected := false
	notSlice := "hoge"
	exists, index, err := InSlice("hoge", notSlice)

	if diff := cmp.Diff(exists, expected); diff != "" {
		t.Errorf("wrong result : %s", diff)
	}
	if diff := cmp.Diff(index, -1); diff != "" {
		t.Errorf("wrong result : %s", diff)
	}
	if err == nil {
		t.Error("wrong result : err is nil")
	}
}
