package app_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/maitesin/hermes/internal/app"
	appMock "github.com/maitesin/hermes/internal/app/mocks"
	"github.com/maitesin/hermes/pkg/comm"
	tmock "github.com/maitesin/hermes/pkg/tracker/mocks"
	"github.com/stretchr/testify/require"
)

type trackerChecks func(*tmock.MockTracker)

type deliveriesRepositoryChecks func(*appMock.MockDeliveriesRepository)

func TestListen(t *testing.T) {
	tests := []struct {
		name                       string
		message                    comm.Message
		trackerChecks              trackerChecks
		deliveriesRepositoryChecks deliveriesRepositoryChecks
		wantErr                    error
	}{
		{
			name: ``,
			message: comm.Message{
				Text: "1234",
			},
			trackerChecks: func(mt *tmock.MockTracker) {
				mt.EXPECT().Track("1234")
			},
			deliveriesRepositoryChecks: func(mdr *appMock.MockDeliveriesRepository) {
				mdr.EXPECT().Insert(context.Background(), app.NewDelivery("1234", "[]", 0, false))
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			mockController := gomock.NewController(t)
			mockTracker := tmock.NewMockTracker(mockController)
			mockDeliveriesRepository := appMock.NewMockDeliveriesRepository(mockController)

			tt.trackerChecks(mockTracker)
			tt.deliveriesRepositoryChecks(mockDeliveriesRepository)

			gotListen := app.Listen(ctx, mockTracker, mockDeliveriesRepository)
			err := gotListen(tt.message)
			if tt.wantErr != nil {
				require.ErrorAs(t, err, &tt.wantErr)
			} else {
				require.Nil(t, err)
			}
		})
	}
}
