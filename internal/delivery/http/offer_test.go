package http

// // Тест GetOffersHandler (успешный запрос)
// func TestGetOffersHandler_Success(t *testing.T) {
// 	req := httptest.NewRequest(http.MethodGet, "/offers", nil)
// 	rr := httptest.NewRecorder()

// 	GetOffersHandler(rr, req)

// 	if rr.Code != http.StatusOK {
// 		t.Errorf("Ожидался статус 200, получен %d", rr.Code)
// 	}

// 	// Проверяем заголовки
// 	if rr.Header().Get("Content-Type") != "application/json" {
// 		t.Errorf("Ожидался Content-Type 'application/json', получен '%s'", rr.Header().Get("Content-Type"))
// 	}
// }

// // Тест GetOffersHandler (неправильный метод)
// func TestGetOffersHandler_MethodNotAllowed(t *testing.T) {
// 	req := httptest.NewRequest(http.MethodPost, "/offers", nil)
// 	rr := httptest.NewRecorder()

// 	GetOffersHandler(rr, req)

// 	if rr.Code != http.StatusMethodNotAllowed {
// 		t.Errorf("Ожидался статус 405, получен %d", rr.Code)
// 	}
// }
