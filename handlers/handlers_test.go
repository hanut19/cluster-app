package handlers

import (
	"cluster-app/db"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// --- LOGIN TESTS ---

func TestLogin_Success(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	db.DB = mockDB

	mock.ExpectQuery(`SELECT password, role FROM users WHERE username=\$1`).
		WithArgs("admin").
		WillReturnRows(sqlmock.NewRows([]string{"password", "role"}).AddRow("admin", "admin"))

	form := url.Values{}
	form.Add("username", "admin")
	form.Add("password", "admin")

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	Login(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("got %d, want %d", resp.StatusCode, http.StatusSeeOther)
	}
	if loc := resp.Header.Get("Location"); loc != "/portal" {
		t.Errorf("redirect location: got %s, want /portal", loc)
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	db.DB = mockDB

	mock.ExpectQuery(`SELECT password, role FROM users WHERE username=\$1`).
		WithArgs("wronguser").
		WillReturnError(sqlmock.ErrCancelled)

	form := url.Values{}
	form.Add("username", "wronguser")
	form.Add("password", "wrongpass")

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	Login(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, want 200", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Invaild username and password") {
		t.Errorf("expected login failure message")
	}
}

// --- UPDATE TESTS ---

func TestUpdate_Success(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	db.DB = mockDB

	mock.ExpectExec(`UPDATE clusters SET nodes=\$1 WHERE name=\$2`).
		WithArgs(20, "cluster-a").
		WillReturnResult(sqlmock.NewResult(1, 1))

	form := url.Values{}
	form.Add("cluster-a", "20")

	req := httptest.NewRequest(http.MethodPost, "/update", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "session", Value: "admin"})

	w := httptest.NewRecorder()
	Update(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("got %d, want %d", resp.StatusCode, http.StatusSeeOther)
	}
	if loc := resp.Header.Get("Location"); loc != "/portal" {
		t.Errorf("redirect location: got %s, want /portal", loc)
	}
}

func TestUpdate_Unauthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/update", nil)
	w := httptest.NewRecorder()

	Update(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("got %d, want %d", resp.StatusCode, http.StatusSeeOther)
	}
	if loc := resp.Header.Get("Location"); loc != "/login" {
		t.Errorf("redirect location: got %s, want /login", loc)
	}
}

// --- LOGOUT TEST ---

func TestLogout(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "admin"})
	w := httptest.NewRecorder()

	Logout(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("got %d, want %d", resp.StatusCode, http.StatusSeeOther)
	}
	if loc := resp.Header.Get("Location"); loc != "/login" {
		t.Errorf("redirect location: got %s, want /login", loc)
	}

	// Check session cookie cleared
	cookies := resp.Cookies()
	var sessionCleared bool
	for _, c := range cookies {
		if c.Name == "session" && c.Value == "" && c.MaxAge == -1 {
			sessionCleared = true
		}
	}
	if !sessionCleared {
		t.Error("session cookie was not cleared")
	}
}
