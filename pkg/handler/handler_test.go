package handler

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/FluFFka/tarantool-kv-storage/pkg/repo"
	gomock "github.com/golang/mock/gomock"
	"github.com/labstack/echo"
)

var (
	ErrDB = errors.New("DB_ERROR")
)

type DataDB struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// go test -v -coverprofile="../../test/handler_cover.out"
// go tool cover -html="../../test/handler_cover.out" -o "../../test/handler_cover.html"

func checkUnmarshalled(val1, val2 string) bool {
	var res1 interface{}
	json.Unmarshal([]byte(val1), &res1)
	var res2 interface{}
	json.Unmarshal([]byte(val2), &res2)
	return reflect.DeepEqual(res1, res2)
}

type HandleFunc func(echo.Context) error

func MakeRequest(reqType, path string, body io.Reader, params map[string]string, checkFunc HandleFunc) (*http.Response, error) {
	r := httptest.NewRequest(reqType, path, body)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ctx := echo.New().NewContext(r, w)
	for k, v := range params {
		ctx.SetParamNames(k)
		ctx.SetParamValues(v)
	}
	err := checkFunc(ctx)
	resp := w.Result()
	return resp, err
}

func CheckResponse(t *testing.T, errGot error, codeGot, codeExp int, bodyGot io.Reader, bodyExp string) bool {
	if errGot != nil {
		t.Errorf("unexpected error: %v", errGot)
		return false
	}
	if codeGot != codeExp {
		t.Errorf("expected code %d, got %d", codeExp, codeGot)
		return false
	}
	if bodyGot != nil {
		body, err := ioutil.ReadAll(bodyGot)
		if err != nil {
			t.Errorf("Can't unmarshall body")
		}
		if !checkUnmarshalled(string(body), bodyExp) {
			t.Errorf("expected %s\ngot %s", bodyExp, string(body))
			return false
		}
	}
	return true
}

func TestGetByKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repoMock := NewMockRepositoryInterface(ctrl)
	handlerTest := Handler{Repo: repoMock}
	data := DataDB{
		Key:   "keyV",
		Value: `{"extra":"no","message":"ok"}`,
	}
	data1 := DataDB{
		Key:   "food",
		Value: `{"banana": ["banana1", "banana2"], "apple": "7", "grape": 9}`,
	}

	// Good requests
	repoMock.EXPECT().GetByKey(data.Key).Return(data.Value, nil)
	resp, err := MakeRequest("GET", "/kv/"+data.Key, nil, map[string]string{"key": data.Key}, handlerTest.GetByKey)
	if !CheckResponse(t, err, resp.StatusCode, 200, resp.Body, data.Value) {
		return
	}
	repoMock.EXPECT().GetByKey(data1.Key).Return(data1.Value, nil)
	resp, err = MakeRequest("GET", "/kv/"+data1.Key, nil, map[string]string{"key": data1.Key}, handlerTest.GetByKey)
	if !CheckResponse(t, err, resp.StatusCode, 200, resp.Body, data1.Value) {
		return
	}

	// Key not found
	repoMock.EXPECT().GetByKey(data.Key).Return("", repo.ErrNoContent)
	resp, err = MakeRequest("GET", "/kv/"+data.Key, nil, map[string]string{"key": data.Key}, handlerTest.GetByKey)
	if !CheckResponse(t, err, resp.StatusCode, 404, nil, "") {
		return
	}

	// DB error
	repoMock.EXPECT().GetByKey(data.Key).Return("", ErrDB)
	resp, err = MakeRequest("GET", "/kv/"+data.Key, nil, map[string]string{"key": data.Key}, handlerTest.GetByKey)
	if !CheckResponse(t, err, resp.StatusCode, 500, nil, "") {
		return
	}
}

func TestPostKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repoMock := NewMockRepositoryInterface(ctrl)
	handlerTest := Handler{Repo: repoMock}
	data := DataDB{
		Key:   "keyV",
		Value: `{"extra":"no","message":"ok"}`,
	}
	data1 := DataDB{
		Key:   "food",
		Value: `{"banana": ["banana1", "banana2"], "apple": "7", "grape": "9"}`,
	}

	// Good requests
	repoMock.EXPECT().InsertValue(data.Key, gomock.Any()).Return(nil)
	bodybyte, _ := json.Marshal(data)
	body := strings.NewReader(string(bodybyte))
	resp, err := MakeRequest("POST", "/kv", body, nil, handlerTest.InsertValue)
	if !CheckResponse(t, err, resp.StatusCode, 200, nil, "") {
		return
	}
	repoMock.EXPECT().InsertValue(data1.Key, gomock.Any()).Return(nil)
	bodybyte, _ = json.Marshal(data1)
	body = strings.NewReader(string(bodybyte))
	resp, err = MakeRequest("POST", "/kv", body, nil, handlerTest.InsertValue)
	if !CheckResponse(t, err, resp.StatusCode, 200, nil, "") {
		return
	}

	// incorrect body
	resp, err = MakeRequest("POST", "/kv", strings.NewReader(`"key":"keyV","value": "valueV"`), nil, handlerTest.InsertValue)
	if !CheckResponse(t, err, resp.StatusCode, 400, nil, "") {
		return
	}
	resp, err = MakeRequest("POST", "/kv", strings.NewReader(`{"NOTkey":"keyV","value": "valueV"}`), nil, handlerTest.InsertValue)
	if !CheckResponse(t, err, resp.StatusCode, 400, nil, "") {
		return
	}
	resp, err = MakeRequest("POST", "/kv", strings.NewReader(`{"key":123,"value": "valueV"}`), nil, handlerTest.InsertValue)
	if !CheckResponse(t, err, resp.StatusCode, 400, nil, "") {
		return
	}
	resp, err = MakeRequest("POST", "/kv", strings.NewReader(`{"key":"keyV","NOTvalue": "valueV"}`), nil, handlerTest.InsertValue)
	if !CheckResponse(t, err, resp.StatusCode, 400, nil, "") {
		return
	}

	// key already exists
	repoMock.EXPECT().InsertValue(data1.Key, gomock.Any()).Return(repo.ErrKeyFound)
	bodybyte, _ = json.Marshal(data1)
	body = strings.NewReader(string(bodybyte))
	resp, err = MakeRequest("POST", "/kv", body, nil, handlerTest.InsertValue)
	if !CheckResponse(t, err, resp.StatusCode, 409, nil, "") {
		return
	}

	// db error
	repoMock.EXPECT().InsertValue(data1.Key, gomock.Any()).Return(ErrDB)
	bodybyte, _ = json.Marshal(data1)
	body = strings.NewReader(string(bodybyte))
	resp, err = MakeRequest("POST", "/kv", body, nil, handlerTest.InsertValue)
	if !CheckResponse(t, err, resp.StatusCode, 500, nil, "") {
		return
	}
}

func TestChangeValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repoMock := NewMockRepositoryInterface(ctrl)
	handlerTest := Handler{Repo: repoMock}
	key := "keyV"
	params := map[string]string{"key": key}
	valueStr1 := `{"value":"valueV"}`
	valueStr := `{"value":{"banana":["banana1","banana2"],"apple":"7","grape": "9"}}`

	// Good requests
	repoMock.EXPECT().ChangeValue(key, gomock.Any()).Return(nil)
	body := strings.NewReader(valueStr1)
	resp, err := MakeRequest("PUT", "/kv"+params["key"], body, params, handlerTest.ChangeValue)
	if !CheckResponse(t, err, resp.StatusCode, 200, nil, "") {
		return
	}
	repoMock.EXPECT().ChangeValue(key, gomock.Any()).Return(nil)
	body = strings.NewReader(valueStr)
	resp, err = MakeRequest("PUT", "/kv"+params["key"], body, params, handlerTest.ChangeValue)
	if !CheckResponse(t, err, resp.StatusCode, 200, nil, "") {
		return
	}

	// Incorrect body
	resp, err = MakeRequest("PUT", "/kv"+params["key"], strings.NewReader(`[]`), params, handlerTest.ChangeValue)
	if !CheckResponse(t, err, resp.StatusCode, 400, nil, "") {
		return
	}
	resp, err = MakeRequest("PUT", "/kv"+params["key"], strings.NewReader(`{"NOTvalue":"valueV"}`), params, handlerTest.ChangeValue)
	if !CheckResponse(t, err, resp.StatusCode, 400, nil, "") {
		return
	}

	// key not found
	repoMock.EXPECT().ChangeValue(key, gomock.Any()).Return(repo.ErrNoContent)
	body = strings.NewReader(valueStr)
	resp, err = MakeRequest("PUT", "/kv"+params["key"], body, params, handlerTest.ChangeValue)
	if !CheckResponse(t, err, resp.StatusCode, 404, nil, "") {
		return
	}

	// db error
	repoMock.EXPECT().ChangeValue(key, gomock.Any()).Return(ErrDB)
	body = strings.NewReader(valueStr)
	resp, err = MakeRequest("PUT", "/kv"+params["key"], body, params, handlerTest.ChangeValue)
	if !CheckResponse(t, err, resp.StatusCode, 500, nil, "") {
		return
	}
}

func TestDeleteValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repoMock := NewMockRepositoryInterface(ctrl)
	handlerTest := Handler{Repo: repoMock}
	data := DataDB{
		Key:   "keyV",
		Value: `{"extra":"no","message":"ok"}`,
	}
	params := map[string]string{"key": data.Key}
	data1 := DataDB{
		Key:   "food",
		Value: `{"banana": ["banana1", "banana2"], "apple": "7", "grape": "9"}`,
	}
	params1 := map[string]string{"key": data1.Key}

	// Good requests
	repoMock.EXPECT().DeleteValue(data.Key).Return(nil)
	resp, err := MakeRequest("DELETE", "/kv/"+data.Key, nil, params, handlerTest.DeleteValue)
	if !CheckResponse(t, err, resp.StatusCode, 200, nil, "") {
		return
	}
	repoMock.EXPECT().DeleteValue(data1.Key).Return(nil)
	resp, err = MakeRequest("DELETE", "/kv/"+data1.Key, nil, params1, handlerTest.DeleteValue)
	if !CheckResponse(t, err, resp.StatusCode, 200, nil, "") {
		return
	}

	// key not found
	repoMock.EXPECT().DeleteValue(data.Key).Return(repo.ErrNoContent)
	resp, err = MakeRequest("DELETE", "/kv/"+data.Key, nil, params, handlerTest.DeleteValue)
	if !CheckResponse(t, err, resp.StatusCode, 404, nil, "") {
		return
	}
	// db error
	repoMock.EXPECT().DeleteValue(data.Key).Return(ErrDB)
	resp, err = MakeRequest("DELETE", "/kv/"+data.Key, nil, params, handlerTest.DeleteValue)
	if !CheckResponse(t, err, resp.StatusCode, 500, nil, "") {
		return
	}
}
