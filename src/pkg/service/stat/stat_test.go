package stat

import (
	"context"
	"errors"
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/service/stat/mocks"
	"github.com/golang/mock/gomock"
)

//go:generate mockgen -package mocks -source ./stat.go -destination ./mocks/stat_mock.go

func TestCheckCounting(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		period   string
		prepare  func(repo *mocks.MockPatientsRepository)
		expected uint64
		wantErr  bool
	}{
		{
			"1. expected 0 for old period",
			"202001",
			func(repo *mocks.MockPatientsRepository) {
				repo.EXPECT().StatPatientsCountGet(ctx, int64(32503662000), int64(32503662000)).Return(uint64(0), nil)
			},
			0,
			false,
		},
		{
			"2. expected 31 for empty period",
			"",
			func(repo *mocks.MockPatientsRepository) {
				repo.EXPECT().StatPatientsCountGet(ctx, int64(0), int64(32503662000)).Return(uint64(31), nil)
			},
			31,
			false,
		},
		{
			"3. expected 31 for correct period",
			"202201",
			func(repo *mocks.MockPatientsRepository) {
				repo.EXPECT().StatPatientsCountGet(ctx, int64(1640995200), int64(1643673600)).Return(uint64(31), nil)
			},
			31,
			false,
		},
		{
			"4. expected 0 for period in future",
			"202301",
			func(repo *mocks.MockPatientsRepository) {
				repo.EXPECT().StatPatientsCountGet(ctx, int64(1672531200), int64(1675209600)).Return(uint64(0), nil)
			},
			0,
			false,
		},
		{
			"5. error on get data",
			"202301",
			func(repo *mocks.MockPatientsRepository) {
				repo.EXPECT().StatPatientsCountGet(ctx, int64(1672531200), int64(1675209600)).Return(uint64(0), errors.New("some error")) //nolint
			},
			0,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repoMock := mocks.NewMockPatientsRepository(ctrl)
			tt.prepare(repoMock)

			service := NewService(repoMock)
			count, err := service.GetPatientsCount(ctx, tt.period)
			if (err != nil) != tt.wantErr {
				t.Fatal(err)
			}

			if count != tt.expected {
				t.Fatalf("Expected %d, received %d", tt.expected, count)
			}
		})
	}
}
