package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-cantabular-dimension-api/service/mock"

	"github.com/ONSdigital/dp-cantabular-dimension-api/handler"
	"github.com/ONSdigital/dp-cantabular-dimension-api/service"
)

func newTestServer(path string, h func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	mux.HandleFunc(path, h)
	return server
}

func TestGetAreas(t *testing.T) {
	Convey("Successfully get a list of areas", func() {
		r, err := http.NewRequest("GET", "http://localhost:27200/areas?dataset=Example", nil)
		So(err, ShouldBeNil)

		w := httptest.NewRecorder()

		serverMock := &mock.HTTPServerMock{
			ListenAndServeFunc: func() error {
				return nil
			},
			ShutdownFunc: func(ctx context.Context) error {
				return nil
			},
		}
		service.GetHTTPServer = func(bindAddr string, router http.Handler) service.HTTPServer {
			return serverMock
		}

		svc := &service.Service{}

		areas := handler.NewAreas(svc.responder, svc.cantabularClient)
		areas.Get("/areas", r)
		So(w.Code, ShouldEqual, http.StatusOK)
		//validateBody(w.Body.Bytes(), expectedBodyFull())
	})

}
