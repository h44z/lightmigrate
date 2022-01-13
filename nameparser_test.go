package lightmigrate

import (
	"reflect"
	"testing"
)

func TestParseFileName(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     *migration
		wantErr  bool
	}{
		{
			name:     "1_foobar.up.sql",
			filename: "1_foobar.up.sql",
			wantErr:  false,
			want: &migration{
				Version:    1,
				Identifier: "foobar",
				Direction:  Up,
				Raw:        "1_foobar.up.sql",
			},
		},
		{
			name:     "1_foobar.down.sql",
			filename: "1_foobar.down.sql",
			wantErr:  false,
			want: &migration{
				Version:    1,
				Identifier: "foobar",
				Direction:  Down,
				Raw:        "1_foobar.down.sql",
			},
		},
		{
			name:     "1_f-o_ob+ar.up.sql",
			filename: "1_f-o_ob+ar.up.sql",
			wantErr:  false,
			want: &migration{
				Version:    1,
				Identifier: "f-o_ob+ar",
				Direction:  Up,
				Raw:        "1_f-o_ob+ar.up.sql",
			},
		},
		{
			name:     "1485385885_foobar.up.sql",
			filename: "1485385885_foobar.up.sql",
			wantErr:  false,
			want: &migration{
				Version:    1485385885,
				Identifier: "foobar",
				Direction:  Up,
				Raw:        "1485385885_foobar.up.sql",
			},
		},
		{
			name:     "20170412214116_date_foobar.up.sql",
			filename: "20170412214116_date_foobar.up.sql",
			wantErr:  false,
			want: &migration{
				Version:    20170412214116,
				Identifier: "date_foobar",
				Direction:  Up,
				Raw:        "20170412214116_date_foobar.up.sql",
			},
		},
		{
			name:     "20220412214116_date_foobar_2022.up.sql",
			filename: "20220412214116_date_foobar_2022.up.sql",
			wantErr:  false,
			want: &migration{
				Version:    20220412214116,
				Identifier: "date_foobar_2022",
				Direction:  Up,
				Raw:        "20220412214116_date_foobar_2022.up.sql",
			},
		},
		{
			name:     "18446744073709551616_date_foobar.up.sql", // uint64.max = 18446744073709551615
			filename: "18446744073709551616_date_foobar.up.sql",
			wantErr:  true,
			want:     nil,
		},
		{
			name:     "-1_foobar.up.sql",
			filename: "-1_foobar.up.sql",
			wantErr:  true,
			want:     nil,
		},
		{
			name:     "-1_foobar.up.sql",
			filename: "-1_foobar.up.sql",
			wantErr:  true,
			want:     nil,
		},
		{
			name:     "foobar.up.sql",
			filename: "foobar.up.sql",
			wantErr:  true,
			want:     nil,
		},
		{
			name:     "1.up.sql",
			filename: "1.up.sql",
			wantErr:  true,
			want:     nil,
		},
		{
			name:     "1_foobar.sql",
			filename: "1_foobar.sql",
			wantErr:  true,
			want:     nil,
		},
		{
			name:     "1_foobar.up",
			filename: "1_foobar.up",
			wantErr:  true,
			want:     nil,
		},
		{
			name:     "1_foobar.down",
			filename: "1_foobar.down",
			wantErr:  true,
			want:     nil,
		},
		{
			name:     "0_foobar.down.json",
			filename: "0_foobar.down.json",
			wantErr:  true,
			want:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFileName(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFileName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFileName() got = %v, want %v", got, tt.want)
			}
		})
	}
}
