package app

import (
	"context"
	"net/url"
	"reflect"
	"testing"
)

func TestApp_Context(t *testing.T) {
	type fields struct {
		id        string
		name      string
		version   string
		metadata  map[string]string
		endpoints []string
		want      struct {
			id        string
			name      string
			version   string
			metadata  map[string]string
			endpoints []string
		}
	}
	tests := []fields{
		{
			id:        "1",
			name:      "toho-v1",
			version:   "v1.0.0",
			metadata:  map[string]string{},
			endpoints: []string{"localhost"},
			want: struct {
				id        string
				name      string
				version   string
				metadata  map[string]string
				endpoints []string
			}{
				id:        "1",
				name:      "toho-v1",
				version:   "v1.0.0",
				metadata:  map[string]string{},
				endpoints: []string{"localhost"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{
				opts: options{
					id:       tt.id,
					name:     tt.name,
					version:  tt.version,
					metadata: tt.metadata,
				},
			}

			for _, e := range tt.endpoints {
				eu, err := url.Parse(e)
				if err != nil {
					t.Errorf("parse endpoint = %v", e)
				}

				a.opts.endpoints = append(a.opts.endpoints, eu)
			}

			ctx := NewContext(context.Background(), a)

			if got, ok := FromContext(ctx); ok {
				if got.ID() != tt.want.id {
					t.Errorf("ID() = %v, want %v", got.ID(), tt.want.id)
				}
				if got.Name() != tt.want.name {
					t.Errorf("Name() = %v, want %v", got.Name(), tt.want.name)
				}
				if got.Version() != tt.want.version {
					t.Errorf("Version() = %v, want %v", got.Version(), tt.want.version)
				}
				if !reflect.DeepEqual(got.Metadata(), tt.want.metadata) {
					t.Errorf("Metadata() = %v, want %v", got.Metadata(), tt.want.metadata)
				}
				if !reflect.DeepEqual(got.Endpoint(), tt.want.endpoints) {
					t.Errorf("Endpoint() = %v, want %v", got.Endpoint(), tt.want.endpoints)
				}
			} else {
				t.Errorf("ok() = %v, want %v", ok, true)
			}
		})
	}
}
