// SPDX-FileCopyrightText: 2022 SAP SE or an SAP affiliate company and Open Component Model contributors.
//
// SPDX-License-Identifier: Apache-2.0

package tag_test

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/open-component-model/ocm/cmds/ocm/testhelper"
	"github.com/open-component-model/ocm/pkg/contexts/oci/repositories/ctf"
	. "github.com/open-component-model/ocm/pkg/testutils"

	"github.com/open-component-model/ocm/pkg/common/accessio"
	"github.com/open-component-model/ocm/pkg/mime"
)

const ARCH = "/tmp/ctf"
const VERSION1 = "v1"
const VERSION2 = "v2"
const NS1 = "mandelsoft/test"
const NS2 = "mandelsoft/index"

var _ = Describe("Test Environment", func() {
	var env *TestEnv

	Context("without attached artifacts", func() {
		BeforeEach(func() {
			env = NewTestEnv()
			env.OCICommonTransport(ARCH, accessio.FormatDirectory, func() {
				env.Namespace(NS1, func() {
					env.Manifest(VERSION1, func() {
						env.Config(func() {
							env.BlobStringData(mime.MIME_JSON, "{}")
						})
						env.Layer(func() {
							env.BlobStringData(mime.MIME_TEXT, "testdata")
						})
					})
					env.Manifest(VERSION2, func() {
						env.Config(func() {
							env.BlobStringData(mime.MIME_JSON, "{}")
						})
						env.Layer(func() {
							env.BlobStringData(mime.MIME_TEXT, "otherdata")
						})
					})
				})

				env.Namespace(NS2, func() {
					env.Index(VERSION1, func() {
						env.Manifest("", func() {
							env.Config(func() {
								env.BlobStringData(mime.MIME_JSON, "{}")
							})
							env.Layer(func() {
								env.BlobStringData(mime.MIME_TEXT, "testdata")
							})
						})
						env.Manifest("", func() {
							env.Config(func() {
								env.BlobStringData(mime.MIME_JSON, "{}")
							})
							env.Layer(func() {
								env.BlobStringData(mime.MIME_TEXT, "otherdata")
							})
						})
					})
					env.Manifest(VERSION2, func() {
						env.Config(func() {
							env.BlobStringData(mime.MIME_JSON, "{}")
						})
						env.Layer(func() {
							env.BlobStringData(mime.MIME_TEXT, "yetanotherdata")
						})
					})
				})
			})
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("tag single artifact", func() {

			buf := bytes.NewBuffer(nil)
			Expect(env.CatchOutput(buf).Execute("tag", "artifact", "latest", ARCH+"//"+NS1+":"+VERSION1)).To(Succeed())
			Expect(buf.String()).To(StringEqualTrimmedWithContext(
				`
tagged ` + NS1 + `@sha256:2c3e2c59e0ac9c99864bf0a9f9727c09f21a66080f9f9b03b36a2dad3cce6ff9
`))

			repo := Must(env.OCIContext().RepositoryForSpec(Must(ctf.NewRepositorySpec(ctf.ACC_READONLY, ARCH, env))))
			defer Close(repo, "repo")
			a := Must(repo.LookupArtifact(NS1, "latest"))
			defer Close(a, "artifact")
		})
	})
})
