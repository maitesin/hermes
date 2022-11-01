package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/maitesin/hermes/internal/app"
	appMock "github.com/maitesin/hermes/internal/app/mocks"
	"github.com/maitesin/hermes/pkg/comm"
	mMock "github.com/maitesin/hermes/pkg/comm/mocks"
	"github.com/maitesin/hermes/pkg/tracker"
	tmock "github.com/maitesin/hermes/pkg/tracker/mocks"
	"github.com/stretchr/testify/require"
)

type messengerChecks func(*mMock.MockMessenger)

type deliveriesRepositoryChecksWithContext func(context.Context, *appMock.MockDeliveriesRepository)

func TestChecker(t *testing.T) {
	tests := []struct {
		name                       string
		trackerChecks              trackerChecks
		deliveriesRepositoryChecks deliveriesRepositoryChecksWithContext
		messengerChecks            messengerChecks
		wantErr                    error
	}{
		{
			name: ``,
			trackerChecks: func(mt *tmock.MockTracker) {
				mt.
					EXPECT().
					Track("1234").
					Return([]tracker.DeliveryEvent{
						{
							Timestamp:   "666",
							Information: "Something",
						},
					}, false, nil)
			},
			deliveriesRepositoryChecks: func(ctx context.Context, mdr *appMock.MockDeliveriesRepository) {
				mdr.
					EXPECT().
					FindAllNotDelivered(ctx).
					Return([]app.Delivery{
						{
							TrackingID:     "1234",
							Log:            "Wololo",
							ConversationID: 9876,
						},
					}, nil)

				mdr.
					EXPECT().
					Update(
						ctx,
						app.NewDelivery("1234", "[{666 Something}]", 9876, false),
					).
					Return(nil)
			},
			messengerChecks: func(mm *mMock.MockMessenger) {
				mm.
					EXPECT().
					Message(comm.Message{
						Conversation: 9876,
						Text:         "[{666 Something}]",
					}).
					Return(nil)
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			tChannel := make(chan time.Time)
			ticker := &time.Ticker{C: tChannel}

			mockController := gomock.NewController(t)
			mockTracker := tmock.NewMockTracker(mockController)
			mockDeliveriesRepository := appMock.NewMockDeliveriesRepository(mockController)
			mockMessenger := mMock.NewMockMessenger(mockController)

			tt.trackerChecks(mockTracker)
			tt.deliveriesRepositoryChecks(ctx, mockDeliveriesRepository)
			tt.messengerChecks(mockMessenger)

			go func() {
				err := app.Checker(ctx, ticker, mockTracker, mockDeliveriesRepository, mockMessenger)
				if tt.wantErr != nil {
					require.ErrorAs(t, err, &tt.wantErr)
				} else {
					require.Nil(t, err)
				}
			}()
			tChannel <- time.Now()
			time.Sleep(time.Second)
			cancel()
		})
	}
}
