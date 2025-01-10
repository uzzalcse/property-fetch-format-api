package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	_ "property-fetch-format-api/routers"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/smartystreets/goconvey/convey"

)

func init() {
    _, file, _, _ := runtime.Caller(0)
    apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
    beego.TestBeegoInit(apppath)
}

func TestPropertyDetailsEndpoint(t *testing.T) {
    mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        propertyID := r.URL.Query().Get("propertyId")
        if propertyID == "HG-123456" {
            response := map[string]interface{}{
                "S3": map[string]interface{}{
                    "id": propertyID,
                    "feed": 1,
                    "published": true,
                },
            }
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(response)
            return
        }
        w.WriteHeader(http.StatusNotFound)
    }))
    defer mockServer.Close()

    beego.AppConfig.Set("baseurl", mockServer.URL)

    convey.Convey("Test Property Details API", t, func() {
        convey.Convey("When requesting valid property", func() {
            r := httptest.NewRequest("GET", "/v1/api/property/details/HG-123456", nil)
            w := httptest.NewRecorder()
            beego.BeeApp.Handlers.ServeHTTP(w, r)
            convey.So(w.Code, convey.ShouldEqual, http.StatusOK)
        })

        convey.Convey("When property not found", func() {
            r := httptest.NewRequest("GET", "/v1/api/property/details/HG-999999", nil)
            w := httptest.NewRecorder()
            beego.BeeApp.Handlers.ServeHTTP(w, r)
            convey.So(w.Code, convey.ShouldEqual, 500)
        })
    })
}

