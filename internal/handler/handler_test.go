package handler_test

import (
	"encoding/json"
	"go-checklist-pet/internal/handler"
	"go-checklist-pet/internal/repository"
	"go-checklist-pet/internal/service"
	"go-checklist-pet/internal/task"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupRouter() http.Handler {
	repo := repository.NewMemoryRepository()
	svc := service.NewTaskService(repo)
	h := handler.NewHandler(svc)

	router := http.NewServeMux()
	h.RegisterRoutes(router)

	return router
}

func TestCreateTask(t *testing.T) {
	router := setupRouter()

	body := strings.NewReader(
		`{"title":"Learn Go","description":"Learn tests"}`,
	)

	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		body,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf(
			"expected status 201, got %d",
			rec.Code,
		)
	}

	var got task.Task

	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}

	if got.ID != 1 {
		t.Errorf("expected ID 1, got %d", got.ID)
	}

	if got.Done {
		t.Error("new task must not be completed")
	}
}

func TestHealth(t *testing.T) {
	router := setupRouter()

	req := httptest.NewRequest(
		http.MethodGet,
		"/health",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			rec.Code,
		)
	}

	var got map[string]string

	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}

	if got["status"] != "ok" {
		t.Errorf("expected status 'ok', got '%s'", got["status"])
	}
}

func TestListTasks(t *testing.T) {
	router := setupRouter()

	req := httptest.NewRequest(
		http.MethodGet,
		"/tasks",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			rec.Code,
		)
	}

	var got []task.Task

	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}

	if len(got) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(got))
	}

	body1 := strings.NewReader(
		`{"title":"Learn Go","description":"Learn tests"}`,
	)

	req1 := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		body1,
	)

	body2 := strings.NewReader(
		`{"title":"Learn Go","description":"Learn tests"}`,
	)

	req2 := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		body2,
	)

	rec1 := httptest.NewRecorder()
	rec2 := httptest.NewRecorder()

	router.ServeHTTP(rec1, req1)
	router.ServeHTTP(rec2, req2)

	if rec1.Code != http.StatusCreated || rec2.Code != http.StatusCreated {
		t.Fatalf(
			"failed to create tasks: got statuses %d and %d",
			rec1.Code,
			rec2.Code,
		)
	}

	req = httptest.NewRequest(
		http.MethodGet,
		"/tasks",
		nil,
	)

	rec = httptest.NewRecorder()

	var gotAfter []task.Task

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &gotAfter); err != nil {
		t.Fatal(err)
	}

	if len(gotAfter) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(gotAfter))
	}

	if gotAfter[0].ID != 1 || gotAfter[1].ID != 2 {
		t.Errorf("expected IDs 1 and 2, got %d and %d", gotAfter[0].ID, gotAfter[1].ID)
	}
}

func TestDeleteTask(t *testing.T) {
	router := setupRouter()

	body := strings.NewReader(
		`{"title":"Learn Go","description":"Learn tests"}`,
	)

	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		body,
	)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf(
			"expected status 201, got %d",
			rec.Code,
		)
	}

	req = httptest.NewRequest(
		http.MethodDelete,
		"/tasks/1",
		nil,
	)
	rec = httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf(
			"expected status 204, got %d",
			rec.Code,
		)
	}

	req = httptest.NewRequest(
		http.MethodGet,
		"/tasks",
		nil,
	)
	rec = httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var got []task.Task

	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}

	if len(got) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(got))
	}
}

func TestMarkTaskDone(t *testing.T) {
	router := setupRouter()

	body := strings.NewReader(
		`{"title":"Learn Go","description":"Learn tests"}`,
	)

	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		body,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf(
			"expected status 201, got %d",
			rec.Code,
		)
	}

	req = httptest.NewRequest(
		http.MethodPut,
		"/tasks/1/done",
		nil,
	)

	rec = httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			rec.Code,
		)
	}

	var got task.Task

	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}

	if !got.Done {
		t.Error("task should be marked as done")
	}

	if got.ID != 1 {
		t.Errorf("expected ID 1, got %d", got.ID)
	}

	req = httptest.NewRequest(
		http.MethodPut,
		"/tasks/1/done",
		nil,
	)
	rec = httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			rec.Code,
		)
	}
}

func TestNotFoundTask(t *testing.T) {
	router := setupRouter()

	req := httptest.NewRequest(
		http.MethodDelete,
		"/tasks/999",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status 404, got %d",
			rec.Code,
		)
	}
}

func TestInvalidTask(t *testing.T) {
	router := setupRouter()

	body := strings.NewReader(
		`{"title":"","description":""}`,
	)

	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		body,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status 400, got %d",
			rec.Code,
		)
	}
}

func TestDeleteTaskInvalidIDs(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{
			name:       "zero ID",
			path:       "/tasks/0",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "negative ID",
			path:       "/tasks/-1",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "text ID",
			path:       "/tasks/abc",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			path:       "/tasks/999",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter()

			req := httptest.NewRequest(
				http.MethodDelete,
				tt.path,
				nil,
			)

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf(
					"expected %d, got %d",
					tt.wantStatus,
					rec.Code,
				)
			}
		})
	}
}

func TestPutDoneInvalidIDs(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{
			name:       "zero ID",
			path:       "/tasks/0/done",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "negative ID",
			path:       "/tasks/-1/done",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "text ID",
			path:       "/tasks/abc/done",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not found",
			path:       "/tasks/999/done",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter()

			req := httptest.NewRequest(
				http.MethodPut,
				tt.path,
				nil,
			)

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf(
					"expected %d, got %d",
					tt.wantStatus,
					rec.Code,
				)
			}
		})
	}
}
