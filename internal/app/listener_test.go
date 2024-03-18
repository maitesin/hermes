package app_test

import (
	"context"
	"testing"

	"github.com/maitesin/hermes/pkg/tracker/group"

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
			name: `Given a working tracker, a working deliveries repository and a valid tracking ID,
					when the Listen function received a message containing the valid tracking ID,
					then the information retrieved from the tracker gets stored in the deliveries repository`,
			message: comm.Message{
				Text: "1234",
			},
			trackerChecks: func(mt *tmock.MockTracker) {
				mt.EXPECT().Track("1234")
				mt.EXPECT().Name().Return("correos")
			},
			deliveriesRepositoryChecks: func(mdr *appMock.MockDeliveriesRepository) {
				mdr.EXPECT().Insert(context.Background(), app.NewDelivery("correos", "1234", "[]", 0, false))
			},
			wantErr: nil,
		},
		{
			name: `Given a failing tracker, a working deliveries repository and a valid tracking ID,
					when the Listen function received a message containing the valid tracking ID,
					then the tracker fails to obtain the information and returns an error`,
			message: comm.Message{
				Text: "1234",
			},
			trackerChecks: func(mt *tmock.MockTracker) {
				mt.EXPECT().Track("1234").Return(nil, false, app.NewDeliveryNotFoundError("1234"))
			},
			deliveriesRepositoryChecks: func(mdr *appMock.MockDeliveriesRepository) {},
			wantErr:                    app.NewDeliveryNotFoundError("1234"),
		},
		{
			name: `Given a working tracker, a failing deliveries repository and a valid tracking ID,
					when the Listen function received a message containing the valid tracking ID,
					then the tracker obtains the information, but the deliveries repository fails to store it and returns an error`,
			message: comm.Message{
				Text: "1234",
			},
			trackerChecks: func(mt *tmock.MockTracker) {
				mt.EXPECT().Track("1234")
				mt.EXPECT().Name().Return("correos")
			},
			deliveriesRepositoryChecks: func(mdr *appMock.MockDeliveriesRepository) {
				mdr.EXPECT().
					Insert(
						context.Background(),
						app.NewDelivery("correos", "1234", "[]", 0, false),
					).
					Return(app.NewDeliveryNotFoundError("1234"))
			},
			wantErr: app.NewDeliveryNotFoundError("1234"),
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

			gotListen := app.Listen(ctx, group.NewGroup(mockTracker), mockDeliveriesRepository)
			err := gotListen(tt.message)
			if tt.wantErr != nil {
				require.ErrorAs(t, err, &tt.wantErr)
			} else {
				require.Nil(t, err)
			}
		})
	}
}
