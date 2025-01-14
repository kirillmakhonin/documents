package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServer(t *testing.T) {
	Convey("Error start server if have bad config", t, func() {
		err := Server(Config{
			Address: ":)",
		})
		So(err, ShouldBeError)
	})
}

func TestApi(t *testing.T) {
	srv := httptest.NewServer(NewHandler(Config{
		TemplatesPath: "../testdata/goodTemplates",
		PublicPath:    "../public",
	}))

	Convey("Html to PDF", t, func() {
		Convey("Params is require", func() {
			resp, err := http.Get(srv.URL + "/html-to-pdf")
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		})

		Convey("Got PDF by content", func() {
			resp, err := http.Get(srv.URL + "/html-to-pdf?content=test")
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			So(resp.Header.Get("Content-Type"), ShouldEqual, "application/pdf")
		})

		Convey("Got PDF by request body", func() {
			data := []byte(`<html><body><h1>test</h2></body></html>`)
			bodyReader := bytes.NewReader(data)
			resp, err := http.Post(srv.URL+"/body-to-pdf", "text/html", bodyReader)
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			So(resp.Header.Get("Content-Type"), ShouldEqual, "application/pdf")
		})
	})

	Convey("PDF from template", t, func() {
		Convey("Params is valid json", func() {
			data := []byte(`1`)
			resp, err := http.Post(srv.URL+"/render-template", "application/json", bytes.NewReader(data))
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusBadRequest)
		})

		Convey("May be exist template", func() {
			data := []byte(`{"templateName": "baz"}`)
			resp, err := http.Post(srv.URL+"/render-template", "application/json", bytes.NewReader(data))
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusInternalServerError)
		})

		Convey("Got PDF by template", func() {
			data := []byte(`{"templateName": "foo"}`)
			resp, err := http.Post(srv.URL+"/render-template", "application/json", bytes.NewReader(data))
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, http.StatusOK)
			So(resp.Header.Get("Content-Type"), ShouldEqual, "application/pdf")
		})
	})

	srv.Close()
}
