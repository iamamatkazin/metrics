package handler

import (
	"compress/gzip"
	"net/http/httptest"
	"testing"
)

func Test_gzipWriter_Write(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Test 1",
			args: args{b: []byte("test")},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gz, err := gzip.NewWriterLevel(httptest.NewRecorder(), gzip.BestSpeed)
			if err != nil {
				t.Errorf("gzipWriter.Write()")
				return
			}
			defer gz.Close()

			w := &gzipWriter{
				ResponseWriter: &responseWriter{},
				Writer:         gz,
			}
			got, err := w.Write(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("gzipWriter.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("gzipWriter.Write() = %v, want %v", got, tt.want)
			}
		})
	}
}
