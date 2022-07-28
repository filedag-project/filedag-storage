package policy

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/s3action"
	"testing"
)

func TestStatement_IsAllowed(t *testing.T) {
	type fields struct {
		SID       ID
		Effect    Effect
		Principal Principal
		Actions   s3action.ActionSet
		Resources ResourceSet
	}
	type args struct {
		args auth.Args
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "test1",
			fields: fields{
				SID:       ID("test1"),
				Effect:    Allow,
				Principal: NewPrincipal("*"),
				Actions:   DefaultPolicies[0].Definition.Statements[0].Actions,
				Resources: NewResourceSet(NewResource("mybucket", "*")),
			},
			args: args{
				args: auth.Args{
					AccountName: "test1",
					Action:      "s3:GetObject",
					BucketName:  "mybucket",
					ObjectName:  "test1.txt",
				},
			},
			want: true,
		},
		{
			name: "test2",
			fields: fields{
				SID:       ID("test2"),
				Effect:    Allow,
				Principal: NewPrincipal("test2"),
				Actions:   DefaultPolicies[0].Definition.Statements[0].Actions,
				Resources: NewResourceSet(NewResource("mybucket", "*")),
			},
			args: args{
				args: auth.Args{
					AccountName: "test2",
					Action:      "s3:GetObject",
					BucketName:  "mybucket",
					ObjectName:  "test1.txt",
				},
			},
			want: true,
		},
		{
			name: "test3",
			fields: fields{
				SID:       ID("test3"),
				Effect:    Allow,
				Principal: NewPrincipal("test3"),
				Actions:   DefaultPolicies[0].Definition.Statements[0].Actions,
				Resources: NewResourceSet(NewResource("mybucket", "test1.txt")),
			},
			args: args{
				args: auth.Args{
					AccountName: "test3",
					Action:      "s3:GetObject",
					BucketName:  "mybucket",
					ObjectName:  "test1.txt",
				},
			},
			want: true,
		},
		{
			name: "test3",
			fields: fields{
				SID:       ID("test3"),
				Effect:    Allow,
				Principal: NewPrincipal("test3"),
				Actions:   DefaultPolicies[0].Definition.Statements[0].Actions,
				Resources: NewResourceSet(NewResource("mybucket", "test2.txt")),
			},
			args: args{
				args: auth.Args{
					AccountName: "test3",
					Action:      "s3:GetObject",
					BucketName:  "mybucket",
					ObjectName:  "test1.txt",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &Statement{
				SID:       tt.fields.SID,
				Effect:    tt.fields.Effect,
				Principal: tt.fields.Principal,
				Actions:   tt.fields.Actions,
				Resources: tt.fields.Resources,
			}
			if got := st.IsAllowed(tt.args.args); got != tt.want {
				t.Errorf("Statement.IsAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}
