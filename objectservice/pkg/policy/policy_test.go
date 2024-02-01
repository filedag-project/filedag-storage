package policy

import (
	"github.com/filedag-project/filedag-storage/objectservice/pkg/auth"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/policy/condition"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/s3action"
	"testing"
)

func TestPolicy_IsAllowed(t *testing.T) {
	type fields struct {
		SID        ID
		Effect     Effect
		Principal  Principal
		Actions    s3action.ActionSet
		Resources  ResourceSet
		Conditions condition.Conditions
	}
	c, _ := condition.NewStringEqualsFunc("", condition.S3Prefix.ToKey(), "object.txt", "")
	cf1 := condition.NewConFunctions(c)
	c, _ = condition.NewStringEqualsFunc("", condition.S3Prefix.ToKey(), "", "")
	cf2 := condition.NewConFunctions(c)
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
				SID:        ID("test1"),
				Effect:     Allow,
				Principal:  NewPrincipal("*"),
				Actions:    s3action.SupportedActions,
				Resources:  NewResourceSet(NewResource("mybucket", "*")),
				Conditions: cf1,
			},
			args: args{
				args: auth.Args{
					AccountName: "test1",
					Action:      "s3:GetObject",
					BucketName:  "mybucket",
					ObjectName:  "object.txt",
				},
			},
			want: true,
		},
		{
			name: "test1",
			fields: fields{
				SID:        ID("test1"),
				Effect:     Allow,
				Principal:  NewPrincipal("*"),
				Actions:    s3action.SupportedActions,
				Resources:  NewResourceSet(NewResource("mybucket", "")),
				Conditions: cf2,
			},
			args: args{
				args: auth.Args{
					AccountName: "test1",
					Action:      "s3:GetObject",
					BucketName:  "mybucket",
					ObjectName:  "object1.txt",
				},
			},
			want: false,
		},
		{
			name: "test2",
			fields: fields{
				SID:        ID("test2"),
				Effect:     Allow,
				Principal:  NewPrincipal("test2"),
				Actions:    s3action.SupportedActions,
				Resources:  NewResourceSet(NewResource("mybucket", "*")),
				Conditions: cf1,
			},
			args: args{
				args: auth.Args{
					AccountName: "test2",
					Action:      "s3:GetObject",
					BucketName:  "mybucket",
					ObjectName:  "object.txt",
				},
			},
			want: true,
		},
		{
			name: "test2",
			fields: fields{
				SID:        ID("test2"),
				Effect:     Allow,
				Principal:  NewPrincipal("test2"),
				Actions:    s3action.SupportedActions,
				Resources:  NewResourceSet(NewResource("mybucket", "")),
				Conditions: cf2,
			},
			args: args{
				args: auth.Args{
					AccountName: "test2",
					Action:      "s3:GetObject",
					BucketName:  "mybucket",
					ObjectName:  "object.txt",
				},
			},
			want: false,
		},
		{
			name: "test2",
			fields: fields{
				SID:        ID("test2"),
				Effect:     Allow,
				Principal:  NewPrincipal("test1"),
				Actions:    s3action.SupportedActions,
				Resources:  NewResourceSet(NewResource("mybucket", "*")),
				Conditions: cf1,
			},
			args: args{
				args: auth.Args{
					AccountName: "test2",
					Action:      "s3:GetObject",
					BucketName:  "mybucket",
					ObjectName:  "object.txt",
				},
			},
			want: false,
		},
		{
			name: "test3",
			fields: fields{
				SID:        ID("test3"),
				Effect:     Allow,
				Principal:  NewPrincipal("*"),
				Actions:    s3action.NewActionSet(s3action.GetObjectAction),
				Resources:  NewResourceSet(NewResource("mybucket", "*")),
				Conditions: cf1,
			},
			args: args{
				args: auth.Args{
					AccountName: "test1",
					Action:      "s3:GetObject",
					BucketName:  "mybucket",
					ObjectName:  "object.txt",
				},
			},
			want: true,
		},
		{
			name: "test3",
			fields: fields{
				SID:        ID("test3"),
				Effect:     Allow,
				Principal:  NewPrincipal("*"),
				Actions:    s3action.NewActionSet(s3action.GetObjectAction),
				Resources:  NewResourceSet(NewResource("mybucket", "")),
				Conditions: cf2,
			},
			args: args{
				args: auth.Args{
					AccountName: "test1",
					Action:      "s3:GetObject",
					BucketName:  "mybucket",
					ObjectName:  "object.txt",
				},
			},
			want: false,
		},
		{
			name: "test3",
			fields: fields{
				SID:       ID("test1"),
				Effect:    Allow,
				Principal: NewPrincipal("*"),
				Actions:   s3action.NewActionSet(s3action.CreateBucketAction),
				Resources: NewResourceSet(NewResource("mybucket", "*")),
			},
			args: args{
				args: auth.Args{
					AccountName: "test1",
					Action:      "s3:GetObject",
					BucketName:  "mybucket",
					ObjectName:  "object.txt",
				},
			},
			want: false,
		},

		{
			name: "test4",
			fields: fields{
				SID:        ID("test4"),
				Effect:     Allow,
				Principal:  NewPrincipal("test4"),
				Actions:    s3action.SupportedActions,
				Resources:  NewResourceSet(NewResource("mybucket", "test4.txt")),
				Conditions: cf1,
			},
			args: args{
				args: auth.Args{
					AccountName: "test4",
					Action:      "s3:GetObject",
					BucketName:  "mybucket",
					ObjectName:  "test4.txt",
				},
			},
			want: true,
		},
		{
			name: "test4",
			fields: fields{
				SID:        ID("test4"),
				Effect:     Allow,
				Principal:  NewPrincipal("test4"),
				Actions:    s3action.SupportedActions,
				Resources:  NewResourceSet(NewResource("mybucket", "")),
				Conditions: cf2,
			},
			args: args{
				args: auth.Args{
					AccountName: "test4",
					Action:      "s3:GetObject",
					BucketName:  "mybucket",
					ObjectName:  "test4.txt",
				},
			},
			want: false,
		},
		{
			name: "test4",
			fields: fields{
				SID:       ID("test4"),
				Effect:    Allow,
				Principal: NewPrincipal("test4"),
				Actions:   s3action.SupportedActions,
				Resources: NewResourceSet(NewResource("mybucket", "test2.txt")),
			},
			args: args{
				args: auth.Args{
					AccountName: "test4",
					Action:      "s3:GetObject",
					BucketName:  "mybucket",
					ObjectName:  "test4.txt",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &Policy{
				ID:      "test",
				Version: "2012-07-05",
				Statements: []Statement{
					{
						SID:       tt.fields.SID,
						Effect:    tt.fields.Effect,
						Principal: tt.fields.Principal,
						Actions:   tt.fields.Actions,
						Resources: tt.fields.Resources,
					},
				},
			}
			if got := st.IsAllowed(tt.args.args); got != tt.want {
				t.Errorf("Statement.IsAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}
