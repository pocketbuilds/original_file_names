package original_file_names

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/pocketbase/pocketbase/tests"
)

func TestPlugin(t *testing.T) {
	setupTestApp := func(t testing.TB) *tests.TestApp {
		testApp, err := tests.NewTestApp("./test/pb_data/")
		if err != nil {
			t.Fatal(err)
		}
		(&Plugin{
			// test config will go here
		}).Init(testApp)
		return testApp
	}

	var contentType string

	makeBody := func(filename string) io.Reader {
		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		_, err := w.CreateFormFile("file", filename)
		if err != nil {
			t.Fatal(err)
		}
		contentType = w.FormDataContentType()
		w.Close()
		return &b
	}

	scenarios := []tests.ApiScenario{
		{
			Name:           "create record",
			Method:         http.MethodPost,
			URL:            "/api/collections/test/records",
			TestAppFactory: setupTestApp,
			Body:           makeBody("test.txt"),
			Headers: map[string]string{
				"Content-Type": contentType,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"file":"test.txt"`,
			},
		},
		{
			Name:           "create record filename with space",
			Method:         http.MethodPost,
			URL:            "/api/collections/test/records",
			TestAppFactory: setupTestApp,
			Body:           makeBody("test 2.txt"),
			Headers: map[string]string{
				"Content-Type": contentType,
			},
			ExpectedStatus: http.StatusOK,
			ExpectedContent: []string{
				`"file":"test 2.txt"`,
			},
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
