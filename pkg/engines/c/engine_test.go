package C

import (
	"errors"
	"testing"

	"github.com/kampanosg/lazytest/pkg/models"
	"github.com/stretchr/testify/assert"
)

type mockRunner struct {
	runHandler func(cmd string) (string, error)
}

func (m *mockRunner) RunCmd(cmd string) (string, error) {
	return m.runHandler(cmd)
}

func TestCEngine_GetIcon(t *testing.T) {
	p := NewCEngine(nil)
	icon := p.GetIcon()
	if icon != "C" {
		t.Errorf("expected icont to be 'C, but got %s", icon)
	}
}

func TestCEngine_Load(t *testing.T) {
	type fields struct {
		runner *mockRunner
	}

	type args struct {
		dir string
	}

	tests := []struct {
		name    string
		fields  func() fields
		args    args
		wantErr bool
		wantNil bool
		want    *models.LazyTree
	}{
		{
			name: "runner returns error",
			fields: func() fields {
				return fields{
					runner: &mockRunner{
						runHandler: func(cmd string) (string, error) {
							return "", errors.New("an error")
						},
					},
				}
			},
			args: args{
				dir: ".",
			},
			wantErr: false,
			wantNil: true,
			want:    nil,
		},
		{
			name: "no tests in the project",
			fields: func() fields {
				return fields{
					runner: &mockRunner{
						runHandler: func(cmd string) (string, error) {
							return "", nil
						},
					},
				}
			},
			args: args{
				dir: ".",
			},
			wantErr: false,
			wantNil: true,
			want:    nil,
		},
		{
			name: "parse tests",
			fields: func() fields {
				return fields{
					runner: &mockRunner{
						runHandler: func(cmd string) (string, error) {
							return "/example/compare\n/example/rand\n/whatever/math\n", nil
						},
					},
				}
			},
			args: args{
				dir: ".",
			},
			wantErr: false,
			wantNil: false,
			want: &models.LazyTree{
				Root: &models.LazyNode{
					Name: "wallet_api",
					Children: []*models.LazyNode{
						{
							Name: "tests",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fields()
			p := NewCEngine(fields.runner)

			got, err := p.Load(tt.args.dir)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			if tt.wantNil {
				assert.Nil(t, got)
				return
			}

			assert.Equal(t, len(tt.want.Root.Children), len(got.Root.Children))
		})
	}
}
