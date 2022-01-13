package lightmigrate

import (
	"errors"
	"os"
	"testing"
)

func TestErrDuplicateMigration_Error(t *testing.T) {
	sampleFile := os.NewFile(0, "TheFileName.txt")
	stat, _ := sampleFile.Stat()
	e := ErrDuplicateMigration{
		FileInfo: stat,
	}
	wantMsg := "duplicate migration file: TheFileName.txt"
	if gotMsg := e.Error(); gotMsg != wantMsg {
		t.Errorf("Error() = %v, want %v", gotMsg, wantMsg)
	}
}

func TestDriverError_Error_NoMsg(t *testing.T) {
	e := DriverError{
		Line:    0,
		Query:   []byte("the db query"),
		Msg:     "",
		OrigErr: errors.New("suberr"),
	}
	wantMsg := "suberr in line 0: the db query"
	if gotMsg := e.Error(); gotMsg != wantMsg {
		t.Errorf("Error() = %v, want %v", gotMsg, wantMsg)
	}
}

func TestDriverError_Error_Msg(t *testing.T) {
	e := DriverError{
		Line:    0,
		Query:   []byte("the db query"),
		Msg:     "error message",
		OrigErr: errors.New("suberr"),
	}
	wantMsg := "error message in line 0: the db query (details: suberr)"
	if gotMsg := e.Error(); gotMsg != wantMsg {
		t.Errorf("Error() = %v, want %v", gotMsg, wantMsg)
	}
}

func TestDriverError_Unwrap_NoSubErr(t *testing.T) {
	var wantSubErr error = nil
	e := DriverError{
		Line:    0,
		Query:   []byte("the db query"),
		Msg:     "error message",
		OrigErr: wantSubErr,
	}

	if got := e.Unwrap(); got != wantSubErr {
		t.Errorf("Unwrap() = %v, want %v", got, wantSubErr)
	}
}

func TestDriverError_Unwrap_SubErr(t *testing.T) {
	wantSubErr := errors.New("suberr")
	e := DriverError{
		Line:    0,
		Query:   []byte("the db query"),
		Msg:     "error message",
		OrigErr: wantSubErr,
	}

	if got := e.Unwrap(); got != wantSubErr {
		t.Errorf("Unwrap() = %v, want %v", got, wantSubErr)
	}
}
