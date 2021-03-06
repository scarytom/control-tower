package terraform_test

import (
	"bytes"
	"reflect"
	"testing"

	. "github.com/EngineerBetter/control-tower/terraform"
)

func TestGCPInputVars_ConfigureTerraform(t *testing.T) {
	type FakeInputVars struct {
		Zone               string
		Tags               string
		Project            string
		GCPCredentialsJSON string
		ExternalIP         string
	}
	type args struct {
		terraformContents string
	}
	tests := []struct {
		name          string
		fakeInputVars FakeInputVars
		args          string
		want          string
		wantErr       bool
	}{
		{name: "Success",
			fakeInputVars: FakeInputVars{
				Zone:               "",
				Tags:               "",
				Project:            "",
				GCPCredentialsJSON: "",
				ExternalIP:         "",
			},
			args:    "",
			want:    "",
			wantErr: false,
		},
		{name: "Failure",
			fakeInputVars: FakeInputVars{},
			args:          "{{ .FakeKey }} \n",
			want:          "",
			wantErr:       true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := &GCPInputVars{
				Zone:               test.fakeInputVars.Zone,
				Tags:               test.fakeInputVars.Tags,
				Project:            test.fakeInputVars.Project,
				GCPCredentialsJSON: test.fakeInputVars.GCPCredentialsJSON,
				ExternalIP:         test.fakeInputVars.ExternalIP,
			}
			got, err := v.ConfigureTerraform(test.args)
			if (err != nil) != test.wantErr {
				t.Errorf("InputVars.ConfigureTerraform() test case \"%s\" failed\nReturned error %v\nExpects an error: %v", test.name, err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("InputVars.ConfigureTerraform() test case \"%s\" failed\nReturned value \"%v\"\nExpected value \"%v\"", test.name, got, test.want)
			}
		})
	}
}

func TestGCPMetadata_Get(t *testing.T) {
	type fields struct {
		Network MetadataStringValue
	}
	tests := []struct {
		name    string
		fields  fields
		fakeKey string
		want    string
		wantErr bool
	}{{
		name: "Success",
		fields: fields{
			Network: MetadataStringValue{
				Value: "fakeNetwork",
			},
		},
		want:    "fakeNetwork",
		fakeKey: "Network",
	},
		{
			name: "Failure",
			fields: fields{
				Network: MetadataStringValue{
					Value: "fakeMetadataStringValue",
				},
			},
			want:    "",
			fakeKey: "FakeKey",
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			outputs := &GCPOutputs{
				Network: test.fields.Network,
			}
			got, err := outputs.Get(test.fakeKey)
			if (err != nil) != test.wantErr {
				t.Errorf("Metadata.Get() test case  \"%s\" failed\nReturned error %v\nExpects an error: %v", test.name, err, test.wantErr)
			}
			if got != test.want {
				t.Errorf("Metadata.Get() test case \"%s\" failed\nReturned value \"%v\"\nExpected value \"%v\"", test.name, got, test.want)
			}
		})
	}
}

func TestGCPMetadata_Init(t *testing.T) {
	tests := []struct {
		name          string
		buffer        *bytes.Buffer
		data          string
		keyToSet      string
		expectedValue string
		wantErr       bool
	}{
		{
			name:          "Success",
			data:          `{"director_public_ip":{"sensitive":false,"type": "string","value": "fakeIP"}}`,
			keyToSet:      "DirectorPublicIP",
			expectedValue: "fakeIP",
		},
		{
			name:          "Failure",
			keyToSet:      "DirectorPublicIP",
			expectedValue: "",
			wantErr:       true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			outputs := &GCPOutputs{}
			buffer := bytes.NewBuffer(nil)
			buffer.WriteString(test.data)
			if err := outputs.Init(buffer); (err != nil) != test.wantErr {
				t.Errorf("Metadata.Init() error = %v, wantErr %v", err, test.wantErr)
			}
			mm := reflect.ValueOf(outputs)
			m := mm.Elem()
			mv := m.FieldByName(test.keyToSet)
			value := mv.FieldByName("Value").String()
			if value != test.expectedValue {
				t.Errorf("Metadata.Init() test case %s\nfailed testing key %s\nexpected value %s\nreceived value %s\n", test.name, test.keyToSet, value, test.expectedValue)
			}

		})
	}
}
