package models

import (
	"testing"
	"time"
)

func TestUser_Validate(t *testing.T) {
	type fields struct {
		ID        uint32
		Nickname  string
		Email     string
		Password  string
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	type args struct {
		action string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "ok",
			fields:  fields{
				ID:        0,
				Nickname:  "Baby",
				Email:     "babytigergmail.com",
				Password:  "ohbabybaby",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			args:    args{
					action: "update",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:        tt.fields.ID,
				Nickname:  tt.fields.Nickname,
				Email:     tt.fields.Email,
				Password:  tt.fields.Password,
				CreatedAt: tt.fields.CreatedAt,
				UpdatedAt: tt.fields.UpdatedAt,
			}
			if err := u.Validate(tt.args.action); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}