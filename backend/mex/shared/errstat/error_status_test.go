package errstat

import (
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorInfo struct {
	name                  string
	err                   error
	expectedMexErrorCode  string
	expectedGrpcCode      codes.Code
	expectedErrorCauseMsg string
	expectedErrorMsg      string
}

func CheckError(err error, tt ErrorInfo, t *testing.T) {
	st, _ := status.FromError(err)
	if st.Code() != tt.expectedGrpcCode {
		t.Errorf("Wrong GRPC status code: got '%s', expected '%s'", st.Code(), tt.expectedGrpcCode)
	}
	errorCauseMsg := ""
	for _, d := range st.Details() {
		if d, ok := d.(*ErrorDetailReason); ok {
			if d.Reason != tt.expectedErrorCauseMsg {
				t.Errorf("Wrong error cause message: got '%s', expected '%s'", d.Reason, tt.expectedErrorCauseMsg)
			}
			errorCauseMsg = d.Reason
		}
		if d, ok := d.(*ErrorDetailCode); ok {
			if d.Code != tt.expectedMexErrorCode {
				t.Errorf("Wrong MexErrorCode: got '%s', expected '%s'", d.Code, tt.expectedMexErrorCode)
			}
		}
	}

	if errorCauseMsg == "" && tt.expectedErrorCauseMsg != "" {
		t.Errorf("Expected an error cause message: got none, expected '%s'", tt.expectedErrorCauseMsg)
	}

	if st.Message() != tt.expectedErrorMsg {
		t.Errorf("Wrong error message: got '%s', expected '%s'", st.Message(), tt.expectedErrorMsg)
	}

}

func Test_MakeMexError(t *testing.T) {

	tests := []ErrorInfo{
		{
			name:                  "Creating an MexError returns correct status and error codes",
			err:                   MakeMexStatus(InvalidClientQuery, "first error message").Err(),
			expectedGrpcCode:      codeToStatus[InvalidClientQuery].GrpcCode,
			expectedMexErrorCode:  codeToStatus[InvalidClientQuery].MexErrorString,
			expectedErrorMsg:      "first error message",
			expectedErrorCauseMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckError(tt.err, tt, t)
		})
	}

}
