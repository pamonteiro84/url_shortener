package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"url_shortener/internal/handlers"
	"url_shortener/internal/models"
	"url_shortener/internal/service"
)

type fakeStorage struct {
	urls map[string]*models.URL
}

func newFakeStorage() *fakeStorage {
	return &fakeStorage{urls: make(map[string]*models.URL)}
}

func (f *fakeStorage) Save(u *models.URL) error {
	f.urls[u.ShortCode] = u
	return nil
}

func (f *fakeStorage) GetByShortURL(shortURL string) (*models.URL, error) {
	u, ok := f.urls[shortURL]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return u, nil
}

func (f *fakeStorage) Delete(shortCode string) error {
	if _, ok := f.urls[shortCode]; !ok {
		return gorm.ErrRecordNotFound
	}
	delete(f.urls, shortCode)
	return nil
}

func newTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	storage := newFakeStorage()
	svc := service.NewService(storage)
	h := handlers.NewHandler(svc)

	r := gin.New()
	Setup(r, h)
	return r
}

func shorten(t *testing.T, r *gin.Engine, url string) *httptest.ResponseRecorder {
	t.Helper()

	body, _ := json.Marshal(map[string]string{"url": url})
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestShorten_NewURL(t *testing.T) {
	r := newTestRouter()

	w := shorten(t, r, "https://example.com")

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d, body = %s", w.Code, http.StatusCreated, w.Body.String())
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if resp["short_code"] == "" {
		t.Error("response missing short_code")
	}
}

func TestShorten_Repeated(t *testing.T) {
	r := newTestRouter()

	first := shorten(t, r, "https://example.com")
	second := shorten(t, r, "https://example.com")

	if second.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", second.Code, http.StatusOK, second.Body.String())
	}

	var firstResp, secondResp map[string]string
	json.Unmarshal(first.Body.Bytes(), &firstResp)
	json.Unmarshal(second.Body.Bytes(), &secondResp)

	if secondResp["short_code"] != firstResp["short_code"] {
		t.Errorf("short_code = %q, want %q (same as first)", secondResp["short_code"], firstResp["short_code"])
	}
}

func TestRedirect_Found(t *testing.T) {
	r := newTestRouter()

	shortenResp := shorten(t, r, "https://example.com")
	var resp map[string]string
	json.Unmarshal(shortenResp.Body.Bytes(), &resp)
	shortCode := resp["short_code"]

	req := httptest.NewRequest(http.MethodGet, "/"+shortCode, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusFound)
	}
	if got := w.Header().Get("Location"); got != "https://example.com" {
		t.Errorf("Location = %q, want %q", got, "https://example.com")
	}
}

func TestRedirect_NotFound(t *testing.T) {
	r := newTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/doesnotexist", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestDelete_Found(t *testing.T) {
	r := newTestRouter()

	shortenResp := shorten(t, r, "https://example.com")
	var resp map[string]string
	json.Unmarshal(shortenResp.Body.Bytes(), &resp)
	shortCode := resp["short_code"]

	req := httptest.NewRequest(http.MethodDelete, "/"+shortCode, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/"+shortCode, nil)
	getW := httptest.NewRecorder()
	r.ServeHTTP(getW, getReq)

	if getW.Code != http.StatusNotFound {
		t.Fatalf("after delete, GET status = %d, want %d", getW.Code, http.StatusNotFound)
	}
}

func TestDelete_NotFound(t *testing.T) {
	r := newTestRouter()

	req := httptest.NewRequest(http.MethodDelete, "/doesnotexist", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}
