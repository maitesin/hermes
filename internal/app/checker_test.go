package app_test

import (
	"context"
	"errors"
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
			name: `Given a working tracker checker, a working deliveries repository and a valid messenger checker,
					when the Checker function is called, a new event is called by the time ticker and a the list of undelivered shipments is retrieved,
					then each of the deliveries is checked and the new information found is stored in the DB`,
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
					Insert(
						ctx,
						app.NewDelivery("1234", "- 666:\n  Something\n", 9876, false),
					).
					Return(nil)
			},
			messengerChecks: func(mm *mMock.MockMessenger) {
				mm.
					EXPECT().
					Message(comm.Message{
						Conversation: 9876,
						Text:         "1234:\n- 666:\n  Something\n",
					}).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name: `Given a working tracker checker, a failing deliveries repository and a valid messenger checker,
					when the Checker function is called, a new event is called by the time ticker and a the list of undelivered shipments is retrieved,
					then it fails to retrieve the deliveries from the repository`,
			trackerChecks: func(mt *tmock.MockTracker) {},
			deliveriesRepositoryChecks: func(ctx context.Context, mdr *appMock.MockDeliveriesRepository) {
				mdr.
					EXPECT().
					FindAllNotDelivered(ctx).
					Return(nil, errors.New("something went wrong"))
			},
			messengerChecks: func(mm *mMock.MockMessenger) {},
			wantErr:         errors.New("something went wrong"),
		},
		{
			name: `Given a failing tracker checker, a working deliveries repository and a valid messenger checker,
					when the Checker function is called, a new event is called by the time ticker and a the list of undelivered shipments is retrieved,
					then the tracker returns a failure`,
			trackerChecks: func(mt *tmock.MockTracker) {
				mt.
					EXPECT().
					Track("1234").
					Return(nil, false, errors.New("something went wrong"))
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
			},
			messengerChecks: func(mm *mMock.MockMessenger) {},
			wantErr:         nil,
		},
		{
			name: `Given a working tracker checker, a working deliveries repository and a failing messenger checker,
					when the Checker function is called, a new event is called by the time ticker and a the list of undelivered shipments is retrieved,
					then the messenger fails to deliver the message`,
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
			},
			messengerChecks: func(mm *mMock.MockMessenger) {
				mm.
					EXPECT().
					Message(comm.Message{
						Conversation: 9876,
						Text:         "1234:\n- 666:\n  Something\n",
					}).
					Return(errors.New("something went wrong"))
			},
			wantErr: errors.New("something went wrong"),
		},
		{
			name: `Given a working tracker checker, a failing deliveries repository and a valid messenger checker,
					when the Checker function is called, a new event is called by the time ticker and a the list of undelivered shipments is retrieved,
					then the deliveries repository fails to update the information`,
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
					Insert(
						ctx,
						app.NewDelivery("1234", "- 666:\n  Something\n", 9876, false),
					).
					Return(errors.New("something went wrong"))
			},
			messengerChecks: func(mm *mMock.MockMessenger) {
				mm.
					EXPECT().
					Message(comm.Message{
						Conversation: 9876,
						Text:         "1234:\n- 666:\n  Something\n",
					}).
					Return(nil)
			},
			wantErr: errors.New("something went wrong"),
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
