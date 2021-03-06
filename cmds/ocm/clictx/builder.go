// Copyright 2022 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package clictx

import (
	"context"
	"io"

	"github.com/mandelsoft/vfs/pkg/vfs"

	"github.com/open-component-model/ocm/cmds/ocm/clictx/core"
	"github.com/open-component-model/ocm/pkg/contexts/datacontext"
	"github.com/open-component-model/ocm/pkg/contexts/ocm"
)

func WithContext(ctx context.Context) core.Builder {
	return core.Builder{}.WithContext(ctx)
}

func WithSharedAttributes(ctx datacontext.AttributesContext) core.Builder {
	return core.Builder{}.WithSharedAttributes(ctx)
}

func WithOCM(ctx ocm.Context) core.Builder {
	return core.Builder{}.WithOCM(ctx)
}

func WithFileSystem(fs vfs.FileSystem) core.Builder {
	return core.Builder{}.WithFileSystem(fs)
}

func WithOutput(w io.Writer) core.Builder {
	return core.Builder{}.WithOutput(w)
}

func WithErrorOutput(w io.Writer) core.Builder {
	return core.Builder{}.WithErrorOutput(w)
}

func WithInput(r io.Reader) core.Builder {
	return core.Builder{}.WithInput(r)
}

func New() core.Context {
	return core.Builder{}.New()
}
