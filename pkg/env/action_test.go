/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package env

import (
	"context"
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/internal/types"
)

func TestAction_Run(t *testing.T) {
	tests := []struct {
		name string
		ctx context.Context
		setup func (context.Context) (int, error)
		expected int
		shouldFail bool
	}{
		{
			name : "single-step action",
			ctx : context.WithValue(context.TODO(), 0, 1),
			setup: func(ctx context.Context) (val int, err error) {
				funcs := []types.EnvFunc{
					func(ctx context.Context) (context.Context, error) {
						val = 12
						return ctx, nil
					},
				}
				_, err = action{role: roleSetup, funcs: funcs}.run(ctx)
				return
			},
			expected: 12,
		},
		{
			name : "multi-step action",
			ctx : context.WithValue(context.TODO(), 0, 1),
			setup: func(ctx context.Context) (val int, err error) {
				funcs := []types.EnvFunc{
					func(ctx context.Context) (context.Context, error) {
						val = 12
						return ctx, nil
					},
					func(ctx context.Context) (context.Context, error) {
						val = val * 2
						return ctx, nil
					},
				}
				_ , err = action{role: roleSetup, funcs: funcs}.run(ctx)
				return
			},
			expected: 24,
		},
		{
			name : "read from context",
			ctx : context.WithValue(context.TODO(), 0, 1),
			setup: func(ctx context.Context) (val int, err error) {
				funcs := []types.EnvFunc{
					func(ctx context.Context) (context.Context, error) {
						i := ctx.Value(0).(int) + 2
						val = i
						return ctx, nil
					},
					func(ctx context.Context) (context.Context, error) {
						val = val + 3
						return ctx, nil
					},
				}
				_, err = action{role: roleSetup, funcs: funcs}.run(ctx)
				return
			},
			expected: 6,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T){
			result, err := test.setup(test.ctx)
			if !test.shouldFail && err != nil{
				t.Fatalf("unexpected failure: %v",err)
			}
			if result != test.expected {
				t.Error("unexpected value:", result)
			}
		})
	}
}
