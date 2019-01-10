package redirect

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test(t *testing.T) {
	tests := []struct {
		name         string
		handler      *Handler
		input        *http.Request
		expectedCode Code
		expectedURL  string
	}{
		{
			name:         "if no URL matches, the default redirect is used",
			handler:      NewHandler("https://example.com/default", CodeTemporary),
			input:        httptest.NewRequest("GET", "/test", nil),
			expectedCode: CodeTemporary,
			expectedURL:  "https://example.com/default",
		},
		{
			name: "if a redirect matches, the expected URL changes",
			handler: NewHandler("https://example.com/default", CodeTemporary).
				Add("/redirectToA", "/targetA/", CodeTemporary),
			input:        httptest.NewRequest("GET", "/redirectToA", nil),
			expectedCode: CodeTemporary,
			expectedURL:  "/targetA/",
		},
		{
			name: "if a redirect matches, the expected URL changes, even if the case doesn't match",
			handler: NewHandler("https://example.com/default", CodeTemporary).
				Add("/redirectToA", "/targetA/", CodeTemporary),
			input:        httptest.NewRequest("GET", "/redirecttoa", nil),
			expectedCode: CodeTemporary,
			expectedURL:  "/targetA/",
		},
		{
			name: "if a redirect matches, the expected URL changes, even if there's a trailing slash",
			handler: NewHandler("https://example.com/default", CodeTemporary).
				Add("/redirectToA", "/targetA/", CodeTemporary),
			input:        httptest.NewRequest("GET", "/redirecttoa/", nil),
			expectedCode: CodeTemporary,
			expectedURL:  "/targetA/",
		},
		{
			name: "if a redirect matches, the expected URL changes, even if there's a trailing slash",
			handler: NewHandler("https://example.com/default", CodeTemporary).
				Add("/redirectToA", "/targetA/", CodeTemporary).
				Add("/redirectToB", "/targetB/", CodePermanent),
			input:        httptest.NewRequest("GET", "/redirecttob/", nil),
			expectedCode: CodePermanent,
			expectedURL:  "/targetB/",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			test.handler.ServeHTTP(w, test.input)
			result := w.Result()
			actualURL, err := result.Location()
			if err != nil {
				t.Fatalf("expected location '%s', got error: %v", test.expectedURL, err)
			}
			if actualURL.String() != test.expectedURL {
				t.Errorf("expected location '%s', got '%s'", test.expectedURL, actualURL)
			}
			if result.StatusCode != int(test.expectedCode) {
				t.Errorf("expected status code '%d', got '%d'", test.expectedCode, result.StatusCode)
			}
		})
	}
}

func parseOrPanic(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}
