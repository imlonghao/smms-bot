package smms

import (
	"io"
	"os"
	"testing"
)

func TestClient_Clear(t *testing.T) {
	type fields struct {
		Token string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *ClearJSON
		wantErr bool
	}{
		{
			"Clear test case",
			fields{},
			&ClearJSON{Success: true},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Token: tt.fields.Token,
			}
			got, err := c.Clear()
			t.Log(got, err)
			if (err != nil) != tt.wantErr {
				t.Errorf("Clear() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Success != tt.want.Success {
				t.Errorf("Clear() got = %v, want %v", got.Success, tt.want.Success)
			}
		})
	}
}

func TestClient_Delete(t *testing.T) {
	type fields struct {
		Token string
	}
	type args struct {
		hash string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *DeleteJSON
		wantErr bool
	}{
		{
			"Delete Test Case",
			fields{},
			args{hash: "not_exist"},
			&DeleteJSON{Success: false},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Token: tt.fields.Token,
			}
			got, err := c.Delete(tt.args.hash)
			t.Log(got, err)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Success != tt.want.Success {
				t.Errorf("Delete() got = %v, want %v", got.Success, tt.want.Success)
			}
		})
	}
}

func TestClient_History(t *testing.T) {
	type fields struct {
		Token string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *HistoryJSON
		wantErr bool
	}{
		{
			"History Test Case",
			fields{},
			&HistoryJSON{Success: true},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Token: tt.fields.Token,
			}
			got, err := c.History()
			t.Log(got, err)
			if (err != nil) != tt.wantErr {
				t.Errorf("History() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Success != tt.want.Success {
				t.Errorf("Delete() got = %v, want %v", got.Success, tt.want.Success)
			}
		})
	}
}

func TestClient_Upload(t *testing.T) {
	file, _ := os.Open("../test/avatar.png")
	type fields struct {
		Token string
	}
	type args struct {
		photo    io.Reader
		filename string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UploadJSON
		wantErr bool
	}{
		{
			"Upload Test Cast",
			fields{},
			args{
				photo:    file,
				filename: "avatar.png",
			},
			&UploadJSON{Success: true},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Token: tt.fields.Token,
			}
			got, err := c.Upload(tt.args.photo, tt.args.filename)
			t.Log(got, err)
			if (err != nil) != tt.wantErr {
				t.Errorf("Upload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Success != tt.want.Success {
				t.Errorf("Delete() got = %v, want %v", got.Success, tt.want.Success)
			}
			// Clearup
			c.Delete(got.Data.Hash)
		})
	}
}
