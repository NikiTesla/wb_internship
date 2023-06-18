package restserver

import (
	"fmt"
	"net/http/httptest"
	"testing"
	mock_natsserver "wb_internship/pkg/nats-server/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRestServer_order(t *testing.T) {
	type mockBehavior func(s *mock_natsserver.MockNats, id int)

	testTable := []struct {
		name                string
		inputUrl            string
		inputID             int
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:     "Not Found",
			inputUrl: "/order?id=123",
			inputID:  123,
			mockBehavior: func(s *mock_natsserver.MockNats, id int) {
				s.EXPECT().GetFromCache(id).Return(nil, fmt.Errorf("no such order"))
			},
			expectedStatusCode:  400,
			expectedRequestBody: "No such order saved",
		},
		{
			name:                "incorrect format of id",
			inputUrl:            "/order",
			mockBehavior:        func(s *mock_natsserver.MockNats, id int) {},
			expectedStatusCode:  400,
			expectedRequestBody: "incorrect format of id",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			// repo := mock_repository.NewMockRepo(c)
			nats := mock_natsserver.NewMockNats(c)
			testCase.mockBehavior(nats, testCase.inputID)

			handler := Handler{NatsServer: nats}
			rtr := handler.InitRouter()

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", testCase.inputUrl, nil)

			rtr.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
