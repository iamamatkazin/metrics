package postgresql

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/iamamatkazin/metrics.git/internal/model"
)

func TestStorage_GetMetric(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name        string
		args        args
		setupMock   func(sqlmock.Sqlmock, string)
		expectError bool
		expected    *model.Metric
		want        *model.Metric
		wantErr     bool
	}{
		{
			name: "Test 1",
			args: args{id: "test"},
			setupMock: func(mock sqlmock.Sqlmock, id string) {
				mock.ExpectQuery(`SELECT id, mtype, val, delta FROM metrics WHERE id = \$1`).
					WithArgs(id).
					WillReturnRows(sqlmock.NewRows([]string{"id", "mtype", "val", "delta"}).
						AddRow("test-gauge", "gauge", 123.45, nil))
			},
			expected: &model.Metric{
				ID:    "test-gauge",
				MType: "gauge",
				Value: nil,
				Delta: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// s := &Storage{
			// 	// cfg: tt.fields.cfg,
			// 	// db:  tt.fields.db,
			// }
			// got, err := s.GetMetric(tt.args.ctx, tt.args.id)
			// if (err != nil) != tt.wantErr {
			// 	t.Errorf("Storage.GetMetric() error = %v, wantErr %v", err, tt.wantErr)
			// 	return
			// }
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("Storage.GetMetric() = %v, want %v", got, tt.want)
			// }
		})
	}
}
