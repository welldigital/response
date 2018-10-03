package response

import (
	"encoding/xml"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSON(t *testing.T) {
	tests := []struct {
		name         string
		v            interface{}
		status       int
		expectedBody string
		expectedErr  error
	}{
		{
			name: "json",
			v: struct {
				Name string `json:"name"`
			}{Name: "value"},
			status:       http.StatusOK,
			expectedBody: `{"name":"value"}`,
		},
		{
			name:         "json: unsupported",
			v:            unmarshallable{Name: "value"},
			status:       http.StatusInternalServerError,
			expectedBody: `{"error":{"statusCode":500,"errorCode":"json_marshal","message":"failed to marshal JSON"}}`,
			expectedErr:  errors.New("json: error calling MarshalJSON for type response.unmarshallable: failed to marshal"),
		},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		actualErr := JSON(test.v, w, test.status)
		if w.Code != test.status {
			t.Errorf("%s: expected status %v, got %v", test.name, test.status, w.Code)
		}
		if w.Body.String() != test.expectedBody {
			t.Errorf("%s: expected body '%v', got '%v'", test.name, test.expectedBody, w.Body.String())
		}
		if test.expectedErr == nil && actualErr != nil {
			t.Fatalf("%s: error not expected, got '%v'", test.name, actualErr)
		}
		if test.expectedErr != nil && actualErr == nil {
			t.Fatalf("%s: error expected '%v', got nil'", test.name, test.expectedErr)
		}
		if test.expectedErr != nil && test.expectedErr.Error() != actualErr.Error() {
			t.Fatalf("%s: expected error '%v', got '%v'", test.name, test.expectedErr, actualErr)
		}
	}
}

type unmarshallable struct {
	Name string
}

func (u unmarshallable) MarshalJSON() ([]byte, error) {
	return []byte{}, errors.New("failed to marshal")
}

func TestOK(t *testing.T) {
	tests := []struct {
		name         string
		ok           bool
		status       int
		expectedBody string
	}{
		{
			name:         "true",
			ok:           true,
			status:       http.StatusOK,
			expectedBody: `{"ok":true}`,
		},
		{
			name:         "true",
			ok:           false,
			status:       http.StatusBadGateway,
			expectedBody: `{"ok":false}`,
		},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		OK(test.ok, w, test.status)
		if w.Code != test.status {
			t.Errorf("%s: expected status %v, got %v", test.name, test.status, w.Code)
		}
		if w.Body.String() != test.expectedBody {
			t.Errorf("%s: expected body '%v', got '%v'", test.name, test.expectedBody, w.Body.String())
		}
	}
}

func TestErrorResponse(t *testing.T) {
	e := NewErrorResponse(http.StatusOK, "error_code", "message")
	if e.Error.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, e.Error.StatusCode)
	}
	if e.Error.ErrorCode != "error_code" {
		t.Errorf("expected ErrorCode '%s', got '%s'", "error_code", e.Error.ErrorCode)
	}
	if e.Error.Message != "message" {
		t.Errorf("expected Message '%s', got '%s'", "message", e.Error.Message)
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		data         map[string]interface{}
		status       int
		expectedBody string
	}{
		{
			name:         "test1",
			err:          errors.New("test2"),
			status:       http.StatusInternalServerError,
			expectedBody: `{"error":{"statusCode":500,"errorCode":"test1","message":"test2"}}`,
		},
		{
			name:         "test1",
			err:          errors.New("test2"),
			status:       http.StatusBadGateway,
			expectedBody: `{"error":{"statusCode":502,"errorCode":"test1","message":"test2"}}`,
		},
		{
			name: "test with data",
			err:  errors.New("test2"),
			data: map[string]interface{}{
				"additionalValue": 123,
			},
			status:       http.StatusBadGateway,
			expectedBody: `{"error":{"statusCode":502,"errorCode":"test with data","message":"test2","data":{"additionalValue":123}}}`,
		},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		if test.data == nil {
			Error(test.name, test.err, w, test.status)
		} else {
			ErrorWithData(test.name, test.err, test.data, w, test.status)
		}
		if w.Code != test.status {
			t.Errorf("%s: expected status %v, got %v", test.name, test.status, w.Code)
		}
		if w.Body.String() != test.expectedBody {
			t.Errorf("%s: expected body '%v', got '%v'", test.name, test.expectedBody, w.Body.String())
		}
	}
}

func TestList(t *testing.T) {
	tests := []struct {
		name         string
		list         List
		status       int
		expectedBody string
	}{
		{
			name:         "empty",
			list:         NewList(0, 0, 0, []interface{}{}),
			status:       http.StatusOK,
			expectedBody: `{"startIndex":0,"pageSize":0,"data":[]}`,
		},
		{
			name:         "3 items",
			list:         NewList(0, 3, 0, []interface{}{1, 2, 3}),
			status:       http.StatusOK,
			expectedBody: `{"startIndex":0,"pageSize":3,"data":[1,2,3]}`,
		},
		{
			name:         "3 items, with more available",
			list:         NewList(0, 3, 3, []interface{}{1, 2, 3}),
			status:       http.StatusOK,
			expectedBody: `{"startIndex":0,"pageSize":3,"nextIndex":3,"data":[1,2,3]}`,
		},
		{
			name:         "start at index 1, page size of 2, next index 3",
			list:         NewList(1, 2, 3, []interface{}{2, 3}),
			status:       http.StatusOK,
			expectedBody: `{"startIndex":1,"pageSize":2,"nextIndex":3,"data":[2,3]}`,
		},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		JSON(test.list, w, test.status)
		if w.Code != test.status {
			t.Errorf("%s: expected status %v, got %v", test.name, test.status, w.Code)
		}
		if w.Body.String() != test.expectedBody {
			t.Errorf("%s: expected body '%v', got '%v'", test.name, test.expectedBody, w.Body.String())
		}
	}
}

func TestCSV(t *testing.T) {
	tests := []struct {
		name         string
		v            [][]string
		status       int
		expectedBody string
		expectedErr  error
	}{
		{
			name: "csv",
			v: [][]string{
				{
					"header1",
					"header2",
				},
				{
					"1",
					"2",
				},
			},
			status:       http.StatusOK,
			expectedBody: `header1,header2
1,2
`,
		},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		actualErr := CSV(test.v, w, test.status)
		if w.Code != test.status {
			t.Errorf("%s: expected status %v, got %v", test.name, test.status, w.Code)
		}
		if w.Body.String() != test.expectedBody {
			t.Errorf("%s: expected body '%v', got '%v'", test.name, test.expectedBody, w.Body.String())
		}
		if test.expectedErr == nil && actualErr != nil {
			t.Fatalf("%s: error not expected, got '%v'", test.name, actualErr)
		}
		if test.expectedErr != nil && actualErr == nil {
			t.Fatalf("%s: error expected '%v', got nil'", test.name, test.expectedErr)
		}
		if test.expectedErr != nil && test.expectedErr.Error() != actualErr.Error() {
			t.Fatalf("%s: expected error '%v', got '%v'", test.name, test.expectedErr, actualErr)
		}
	}
}

type XMLExample struct {
	XMLName xml.Name `xml:"Response"`
	Message string   `xml:"Message"`
}

func TestXML(t *testing.T) {
	tests := []struct {
		name         string
		v            interface{}
		status       int
		expectedBody string
		expectedErr  error
	}{
		{
			name: "xml",
			v: XMLExample{Message:"message"},
			status:       http.StatusOK,
			expectedBody: `<Response><Message>message</Message></Response>`,
		},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		actualErr := XML(test.v, w, test.status)
		if w.Code != test.status {
			t.Errorf("%s: expected status %v, got %v", test.name, test.status, w.Code)
		}
		if w.Body.String() != test.expectedBody {
			t.Errorf("%s: expected body '%v', got '%v'", test.name, test.expectedBody, w.Body.String())
		}
		if test.expectedErr == nil && actualErr != nil {
			t.Fatalf("%s: error not expected, got '%v'", test.name, actualErr)
		}
		if test.expectedErr != nil && actualErr == nil {
			t.Fatalf("%s: error expected '%v', got nil'", test.name, test.expectedErr)
		}
		if test.expectedErr != nil && test.expectedErr.Error() != actualErr.Error() {
			t.Fatalf("%s: expected error '%v', got '%v'", test.name, test.expectedErr, actualErr)
		}
	}
}