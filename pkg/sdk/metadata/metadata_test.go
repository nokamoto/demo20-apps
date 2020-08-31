package metadata

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nokamoto/demo20-apis/cloud/api"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestAppendToOutgoingContext(t *testing.T) {
	expected := api.Metadata{
		Credential: &api.Metadata_MachineUserApiKey{
			MachineUserApiKey: "foo",
		},
		User: &api.Metadata_MachineUser{
			MachineUser: "bar",
		},
		Parent: "baz",
	}

	ctx, err := AppendToOutgoingContext(context.Background(), &expected)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := FromOutgoingContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(&expected, actual, protocmp.Transform()); len(diff) != 0 {
		t.Error(diff)
	}
}

func TestNewIncomingContext(t *testing.T) {
	expected := api.Metadata{
		Credential: &api.Metadata_MachineUserApiKey{
			MachineUserApiKey: "foo",
		},
		User: &api.Metadata_MachineUser{
			MachineUser: "bar",
		},
		Parent: "baz",
	}

	ctx, err := NewIncomingContext(context.Background(), &expected)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := FromIncomingContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(&expected, actual, protocmp.Transform()); len(diff) != 0 {
		t.Error(diff)
	}
}

func TestFromIncomingContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *api.Metadata
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromIncomingContext(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromIncomingContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromIncomingContext() = %v, want %v", got, tt.want)
			}
		})
	}
}
