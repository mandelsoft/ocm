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

package ctf

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gardener/ocm/pkg/errors"
	"github.com/gardener/ocm/pkg/ocm/accessmethods"
	"github.com/gardener/ocm/pkg/ocm/compdesc"
	metav1 "github.com/gardener/ocm/pkg/ocm/compdesc/meta/v1"
	"github.com/gardener/ocm/pkg/ocm/core"
	"github.com/mandelsoft/vfs/pkg/osfs"
	"github.com/mandelsoft/vfs/pkg/vfs"
	"github.com/opencontainers/go-digest"

	"github.com/gardener/component-spec/bindings-go/codec"
)

// ComponentDescriptorFileName is the name of the component-descriptor file.
const ComponentDescriptorFileName = "component-descriptor.yaml"

// BlobsDirectoryName is the name of the blob directory in the tar.
const BlobsDirectoryName = "blobs"

////////////////////////////////////////////////////////////////////////////////

type ComponentArchiveOptions struct {
	PathFileSystem      vfs.FileSystem
	ComponentFileSystem vfs.FileSystem
	ArchiveFormat       ArchiveFormat
	DecodeOptions       *compdesc.DecodeOptions
	EncodeOptions       *compdesc.EncodeOptions
}

var _ ComponentArchiveOption = &ComponentArchiveOptions{}

func (o *ComponentArchiveOptions) ApplyOption(options *ComponentArchiveOptions) {
	if o == nil {
		return
	}
	if o.PathFileSystem != nil {
		options.PathFileSystem = o.PathFileSystem
	}
	if o.ComponentFileSystem != nil {
		options.ComponentFileSystem = o.ComponentFileSystem
	}
	if o.ArchiveFormat != nil {
		options.ArchiveFormat = o.ArchiveFormat
	}
}

func (o *ComponentArchiveOptions) Default() *ComponentArchiveOptions {
	if o.PathFileSystem == nil {
		o.PathFileSystem = _osfs
	}
	if o.ArchiveFormat == nil {
		o.ArchiveFormat = CA_DIRECTORY
	}
	return o
}

// ApplyOptions applies the given list options on these options,
// and then returns itself (for convenient chaining).
func (o *ComponentArchiveOptions) ApplyOptions(opts []ComponentArchiveOption) *ComponentArchiveOptions {
	for _, opt := range opts {
		if opt != nil {
			opt.ApplyOption(o)
		}
	}
	return o
}

// ComponentArchiveOption is the interface to specify different archive options
type ComponentArchiveOption interface {
	ApplyOption(options *ComponentArchiveOptions)
}

// PathFileSystem set the evaltuation filesystem for the path name
func PathFileSystem(fs vfs.FileSystem) ComponentArchiveOption {
	return opt_PFS{fs}
}

type opt_PFS struct {
	vfs.FileSystem
}

// ApplyOption applies the configured path filesystem.
func (o opt_PFS) ApplyOption(options *ComponentArchiveOptions) {
	options.PathFileSystem = o.FileSystem
}

// ComponentFileSystem set the evaltuation filesystem for the path name
func ComponentFileSystem(fs vfs.FileSystem) ComponentArchiveOption {
	return opt_CFS{fs}
}

type opt_CFS struct {
	vfs.FileSystem
}

// ApplyOption applies the configured path filesystem.
func (o opt_CFS) ApplyOption(options *ComponentArchiveOptions) {
	options.ComponentFileSystem = o.FileSystem
}

type ArchiveFormat interface {
	String() string
	ComponentArchiveOption

	Open(ctx core.Context, path string, opts ComponentArchiveOptions) (*ComponentArchive, error)
	Create(ctx core.Context, path string, cd *compdesc.ComponentDescriptor, opts ComponentArchiveOptions) (*ComponentArchive, error)
	Write(ca *ComponentArchive, path string, opts ComponentArchiveOptions) error
}

type ComponentInfo struct {
	Path    string
	Options ComponentArchiveOptions
}

type ComponentCloser interface {
	Close(ca *ComponentArchive) error
}

type ComponentCloserFunction func(ca *ComponentArchive) error

func (f ComponentCloserFunction) Close(ca *ComponentArchive) error {
	return f(ca)
}

var formats = map[string]ArchiveFormat{}

func RegisterComponentArchiveFormat(f ArchiveFormat) {
	formats[f.String()] = f
}

func GetCompnentArchiveFormat(name string) ArchiveFormat {
	return formats[name]
}

////////////////////////////////////////////////////////////////////////////////

// ComponentArchive is the go representation for a component artefact
type ComponentArchive struct {
	ComponentDescriptor *compdesc.ComponentDescriptor
	repo                core.Repository
	fs                  vfs.FileSystem
	ctx                 core.Context
	closer              ComponentCloser
}

var _ core.ComponentAccess = &ComponentArchive{}

// NewComponentArchive returns a new component descriptor with a filesystem
func NewComponentArchive(ctx core.Context, repo core.Repository, cd *compdesc.ComponentDescriptor, fs vfs.FileSystem, closer ComponentCloser) *ComponentArchive {
	ca := &ComponentArchive{
		repo:                repo,
		ComponentDescriptor: cd,
		fs:                  fs,
		ctx:                 ctx,
		closer:              closer,
	}
	if repo == nil {
		ca.repo = newPlainComponent(ca, ctx)
	}
	return ca
}

func evaluateOptions(opts []ComponentArchiveOption) ComponentArchiveOptions {
	return *(&ComponentArchiveOptions{}).ApplyOptions(opts).Default()
}

func CreateComponentArchive(ctx core.Context, cd *compdesc.ComponentDescriptor, path string, opts ...ComponentArchiveOption) (*ComponentArchive, error) {
	o := evaluateOptions(opts)
	return o.ArchiveFormat.Create(ctx, path, cd, o)
}

// OpenComponentArchive creates a new component archive from a file path.
func OpenComponentArchive(ctx core.Context, path string, opts ...ComponentArchiveOption) (*ComponentArchive, error) {
	o := evaluateOptions(opts)
	fi, err := o.PathFileSystem.Stat(path)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return CA_DIRECTORY.Open(ctx, path, o)
	}

	ca, err := CA_TGZ.Open(ctx, path, o)
	if err != nil {
		ca, err = CA_TAR.Open(ctx, path, o)
	}
	return ca, err
}

func (c *ComponentArchive) AccessMethod(a core.AccessSpec) (core.AccessMethod, error) {
	if a.GetName() == accessmethods.LocalBlobType {
		a, err := c.ctx.AccessSpecForSpec(a)
		if err != nil {
			return nil, err
		}
		if a.GetVersion() == "v1" {
			conv := &LocalFilesystemBlobAccessSpec{
				*a.(*accessmethods.LocalBlobAccessSpec),
			}
			return newLocalFilesystemBlobAccessMethod(conv, c)
		}
		return nil, errors.ErrNotSupported(errors.KIND_ACCESSMETHOD, a.GetType(), CTFRepositoryType)
	}
	// fall back to original version
	return a.AccessMethod(c)
}

func (c *ComponentArchive) GetContext() core.Context {
	return c.ctx
}

func (c *ComponentArchive) GetAccessType() string {
	return CTFRepositoryType
}

func (c *ComponentArchive) Close() error {
	if c.closer != nil {
		return c.closer.Close(c)
	}
	return nil
}

func (c *ComponentArchive) GetRepository() core.Repository {
	return c.repo
}

func (c *ComponentArchive) GetName() string {
	return c.ComponentDescriptor.GetName()
}

func (c *ComponentArchive) GetVersion() string {
	return c.ComponentDescriptor.GetVersion()
}

func (c *ComponentArchive) GetDescriptor() (*compdesc.ComponentDescriptor, error) {
	return c.ComponentDescriptor, nil
}

func (c *ComponentArchive) GetResource(meta metav1.Identity) (core.ResourceAccess, error) {
	// TODO
	return nil, fmt.Errorf("not implemented")
}
func (c *ComponentArchive) GetSource(meta metav1.Identity) (core.ResourceAccess, error) {
	// TODO
	return nil, fmt.Errorf("not implemented")
}

var _osfs = osfs.New()

func FileSystem(ofs []vfs.FileSystem) vfs.FileSystem {
	if len(ofs) > 0 {
		return ofs[0]
	}
	return _osfs
}

// Digest returns the digest of the component archive.
// The digest is computed serializing the included component descriptor into json and compute sha hash.
func (ca *ComponentArchive) Digest() (string, error) {
	data, err := codec.Encode(ca.ComponentDescriptor)
	if err != nil {
		return "", err
	}
	return digest.FromBytes(data).String(), nil
}

func (ca *ComponentArchive) writeCD() error {
	cdBytes, err := compdesc.Encode(ca.ComponentDescriptor)
	if err != nil {
		return fmt.Errorf("unable to encode component descriptor: %w", err)
	}
	if err := vfs.WriteFile(ca.fs, ComponentDescriptorFileName, cdBytes, os.ModePerm); err != nil {
		return fmt.Errorf("unable to copy component descritptor to %q: %w", ComponentDescriptorFileName, err)
	}
	return nil
}

func (ca *ComponentArchive) checkAccessSpec(acc compdesc.AccessSpec) error {
	spec, err := ca.GetContext().AccessSpecForSpec(acc)
	if err != nil {
		return err
	}
	if spec.ValidFor(ca.repo) {
		return nil
	}
	return errors.ErrInvalid(errors.KIND_ACCESSMETHOD, acc.GetName(), ca.repo.GetSpecification().GetName())
}

func (ca *ComponentArchive) AddSource(meta *core.SourceMeta, acc compdesc.AccessSpec) error {
	if err := ca.checkAccessSpec(acc); err != nil {
		return err
	}
	src := &compdesc.Source{
		SourceMeta: *meta.Copy(),
		Access:     acc,
	}

	if idx := ca.ComponentDescriptor.GetSourceIndex(meta); idx == -1 {
		ca.ComponentDescriptor.Sources = append(ca.ComponentDescriptor.Sources, *src)
	} else {
		ca.ComponentDescriptor.Sources[idx] = *src
	}
	return ca.writeCD()
}

// AddSource adds a blob source to the current archive.
// If the specified source already exists it will be overwritten.
func (ca *ComponentArchive) AddSourceBlob(meta *core.SourceMeta, acc core.BlobAccess) error {
	if acc == nil {
		return errors.New("a source has to be defined")
	}
	if err := ca.ensureBlobsPath(); err != nil {
		return err
	}

	digest, err := core.Digest(acc)
	if err != nil {
		return err
	}

	blobpath := BlobPath(string(digest))
	if _, err := ca.fs.Stat(blobpath); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("unable to get file info for %s", blobpath)
		}
		file, err := ca.fs.OpenFile(blobpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return fmt.Errorf("unable to open file %s: %w", blobpath, err)
		}
		src, err := acc.Reader()
		if err != nil {
			return err
		}
		defer src.Close()
		if _, err := io.Copy(file, src); err != nil {
			return fmt.Errorf("unable to write blob %q to file: %w", blobpath, err)
		}
		if err := file.Close(); err != nil {
			return fmt.Errorf("unable to close file: %w", err)
		}
	}

	return ca.AddSource(meta, NewLocalFilesystemBlobAccessSpecV1(string(digest), acc.MimeType()))
}

func (ca *ComponentArchive) AddResource(meta *core.ResourceMeta, acc compdesc.AccessSpec) error {
	if err := ca.checkAccessSpec(acc); err != nil {
		return err
	}
	res := &compdesc.Resource{
		ResourceMeta: *meta.Copy(),
		Access:       acc,
	}

	if idx := ca.ComponentDescriptor.GetResourceIndex(meta); idx == -1 {
		ca.ComponentDescriptor.Resources = append(ca.ComponentDescriptor.Resources, *res)
	} else {
		ca.ComponentDescriptor.Resources[idx] = *res
	}
	return ca.writeCD()
}

// AddResource adds a blob resource to the current archive.
func (ca *ComponentArchive) AddResourceBlob(meta *core.ResourceMeta, acc core.BlobAccess) error {
	if acc == nil {
		return errors.New("a resource has to be defined")
	}
	if err := ca.ensureBlobsPath(); err != nil {
		return err
	}

	digest, err := core.Digest(acc)
	if err != nil {
		return err
	}

	blobpath := BlobPath(string(digest))
	if _, err := ca.fs.Stat(blobpath); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("unable to get file info for %s", blobpath)
		}
		file, err := ca.fs.OpenFile(blobpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return fmt.Errorf("unable to open file %s: %w", blobpath, err)
		}
		src, err := acc.Reader()
		if err != nil {
			return err
		}
		defer src.Close()

		_, err = io.Copy(file, src)
		if err != nil {
			return fmt.Errorf("unable to write blob to file %q: %w", blobpath, err)
		}
		if err := file.Close(); err != nil {
			return fmt.Errorf("unable to close file: %w", err)
		}
		return nil
	}

	return ca.AddResource(meta, NewLocalFilesystemBlobAccessSpecV1(string(digest), acc.MimeType()))
}

// ensureBlobsPath ensures that the blob directory exists
func (ca *ComponentArchive) ensureBlobsPath() error {
	if _, err := ca.fs.Stat(BlobsDirectoryName); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("unable to get file info for blob directory: %w", err)
		}
		return ca.fs.Mkdir(BlobsDirectoryName, os.ModePerm)
	}
	return nil
}

// BlobPath returns the path to the blob for a given name.
func BlobPath(name string) string {
	return filepath.Join(BlobsDirectoryName, name)
}

// ExtractTarToFs writes a tar stream to a filesystem.
func ExtractTarToFs(fs vfs.FileSystem, in io.Reader) error {
	tr := tar.NewReader(in)
	for {
		header, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := fs.MkdirAll(header.Name, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("unable to create directory %s: %w", header.Name, err)
			}
		case tar.TypeReg:
			file, err := fs.OpenFile(header.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("unable to open file %s: %w", header.Name, err)
			}
			if _, err := io.Copy(file, tr); err != nil {
				return fmt.Errorf("unable to copy tar file to filesystem: %w", err)
			}
			if err := file.Close(); err != nil {
				return fmt.Errorf("unable to close file %s: %w", header.Name, err)
			}
		}
	}
}
