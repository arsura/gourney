package pgsql

import (
	"testing"

	"github.com/arsura/moonbase-service/pkg/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

func TestCurrencyModel_Insert(t *testing.T) {
	type fields struct {
		Pool *pgxpool.Pool
	}
	type args struct {
		p *models.Currency
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "foo",
			fields: fields{
				Pool: (func() {

				}),
			},
			args: args{
				p: &models.Currency{
					Name:       "string",
					Amount:     1000,
					Total:      1000,
					RiseRate:   1000,
					RiseFactor: 1000,
				},
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &CurrencyModel{
				Pool: tt.fields.Pool,
			}
			got, err := m.Insert(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("CurrencyModel.Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CurrencyModel.Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}
