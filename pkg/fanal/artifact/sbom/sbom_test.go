package sbom_test

import (
	"context"
	"errors"
	"github.com/package-url/packageurl-go"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aquasecurity/trivy/pkg/fanal/artifact"
	"github.com/aquasecurity/trivy/pkg/fanal/artifact/sbom"
	"github.com/aquasecurity/trivy/pkg/fanal/cache"
	"github.com/aquasecurity/trivy/pkg/fanal/types"
)

func TestArtifact_Inspect(t *testing.T) {
	tests := []struct {
		name               string
		filePath           string
		putBlobExpectation cache.ArtifactCachePutBlobExpectation
		want               types.ArtifactReference
		wantErr            []string
	}{
		{
			name:     "happy path",
			filePath: filepath.Join("testdata", "bom.json"),
			putBlobExpectation: cache.ArtifactCachePutBlobExpectation{
				Args: cache.ArtifactCachePutBlobArgs{
					BlobID: "sha256:c1cc58e08422fd7606a8e9ee2b42bf722b7af8b703b895461c23b83956f33227",
					BlobInfo: types.BlobInfo{
						SchemaVersion: types.BlobJSONSchemaVersion,
						OS: types.OS{
							Family: "alpine",
							Name:   "3.16.0",
						},
						PackageInfos: []types.PackageInfo{
							{
								Packages: types.Packages{
									{
										Name:       "musl",
										Version:    "1.2.3-r0",
										SrcName:    "musl",
										SrcVersion: "1.2.3-r0",
										Licenses:   []string{"MIT"},
										Ref:        "pkg:apk/alpine/musl@1.2.3-r0?distro=3.16.0",
										Layer: types.Layer{
											DiffID: "sha256:dd565ff850e7003356e2b252758f9bdc1ff2803f61e995e24c7844f6297f8fc3",
										},
										Identifier: types.PkgIdentifier{
											PURL: &types.PackageURL{
												PackageURL: packageurl.PackageURL{
													Type:      packageurl.TypeApk,
													Namespace: "alpine",
													Name:      "musl",
													Version:   "1.2.3-r0",
													Qualifiers: packageurl.Qualifiers{
														{
															Key:   "distro",
															Value: "3.16.0",
														},
													},
												},
											},
										},
									},
								},
							},
						},
						Applications: []types.Application{
							{
								Type:     "composer",
								FilePath: "app/composer/composer.lock",
								Libraries: types.Packages{
									{
										Name:    "pear/log",
										Version: "1.13.1",
										Ref:     "pkg:composer/pear/log@1.13.1",
										Layer: types.Layer{
											DiffID: "sha256:3c79e832b1b4891a1cb4a326ef8524e0bd14a2537150ac0e203a5677176c1ca1",
										},
										Identifier: types.PkgIdentifier{
											PURL: &types.PackageURL{
												PackageURL: packageurl.PackageURL{
													Type:      packageurl.TypeComposer,
													Namespace: "pear",
													Name:      "log",
													Version:   "1.13.1",
												},
											},
										},
									},
									{

										Name:    "pear/pear_exception",
										Version: "v1.0.0",
										Ref:     "pkg:composer/pear/pear_exception@v1.0.0",
										Layer: types.Layer{
											DiffID: "sha256:3c79e832b1b4891a1cb4a326ef8524e0bd14a2537150ac0e203a5677176c1ca1",
										},
										Identifier: types.PkgIdentifier{
											PURL: &types.PackageURL{
												PackageURL: packageurl.PackageURL{
													Type:      packageurl.TypeComposer,
													Namespace: "pear",
													Name:      "pear_exception",
													Version:   "v1.0.0",
												},
											},
										},
									},
								},
							},
							{
								Type:     "gobinary",
								FilePath: "app/gobinary/gobinary",
								Libraries: types.Packages{
									{
										Name:    "github.com/package-url/packageurl-go",
										Version: "v0.1.1-0.20220203205134-d70459300c8a",
										Ref:     "pkg:golang/github.com/package-url/packageurl-go@v0.1.1-0.20220203205134-d70459300c8a",
										Layer: types.Layer{
											DiffID: "sha256:3c79e832b1b4891a1cb4a326ef8524e0bd14a2537150ac0e203a5677176c1ca1",
										},
										Identifier: types.PkgIdentifier{
											PURL: &types.PackageURL{
												PackageURL: packageurl.PackageURL{
													Type:      packageurl.TypeGolang,
													Namespace: "github.com/package-url",
													Name:      "packageurl-go",
													Version:   "v0.1.1-0.20220203205134-d70459300c8a",
												},
											},
										},
									},
								},
							},
							{
								Type:     "jar",
								FilePath: "",
								Libraries: types.Packages{
									{
										Name:    "org.codehaus.mojo:child-project",
										Ref:     "pkg:maven/org.codehaus.mojo/child-project@1.0?file_path=app%2Fmaven%2Ftarget%2Fchild-project-1.0.jar",
										Version: "1.0",
										Layer: types.Layer{
											DiffID: "sha256:3c79e832b1b4891a1cb4a326ef8524e0bd14a2537150ac0e203a5677176c1ca1",
										},
										FilePath: "app/maven/target/child-project-1.0.jar",
										Identifier: types.PkgIdentifier{
											PURL: &types.PackageURL{
												PackageURL: packageurl.PackageURL{
													Type:      packageurl.TypeMaven,
													Namespace: "org.codehaus.mojo",
													Name:      "child-project",
													Version:   "1.0",
												},
												FilePath: "app/maven/target/child-project-1.0.jar",
											},
										},
									},
								},
							},
							{
								Type:     "node-pkg",
								FilePath: "",
								Libraries: types.Packages{
									{
										Name:     "bootstrap",
										Version:  "5.0.2",
										Ref:      "pkg:npm/bootstrap@5.0.2?file_path=app%2Fapp%2Fpackage.json",
										Licenses: []string{"MIT"},
										Layer: types.Layer{
											DiffID: "sha256:3c79e832b1b4891a1cb4a326ef8524e0bd14a2537150ac0e203a5677176c1ca1",
										},
										FilePath: "app/app/package.json",
										Identifier: types.PkgIdentifier{
											PURL: &types.PackageURL{
												PackageURL: packageurl.PackageURL{
													Type:    packageurl.TypeNPM,
													Name:    "bootstrap",
													Version: "5.0.2",
												},
												FilePath: "app/app/package.json",
											},
										},
									},
								},
							},
						},
					},
				},
				Returns: cache.ArtifactCachePutBlobReturns{},
			},
			want: types.ArtifactReference{
				Name: filepath.Join("testdata", "bom.json"),
				Type: types.ArtifactCycloneDX,
				ID:   "sha256:c1cc58e08422fd7606a8e9ee2b42bf722b7af8b703b895461c23b83956f33227",
				BlobIDs: []string{
					"sha256:c1cc58e08422fd7606a8e9ee2b42bf722b7af8b703b895461c23b83956f33227",
				},
			},
		},
		{
			name:     "happy path for sbom attestation",
			filePath: filepath.Join("testdata", "sbom.cdx.intoto.jsonl"),
			putBlobExpectation: cache.ArtifactCachePutBlobExpectation{
				Args: cache.ArtifactCachePutBlobArgs{
					BlobID: "sha256:c1cc58e08422fd7606a8e9ee2b42bf722b7af8b703b895461c23b83956f33227",
					BlobInfo: types.BlobInfo{
						SchemaVersion: types.BlobJSONSchemaVersion,
						OS: types.OS{
							Family: "alpine",
							Name:   "3.16.0",
						},
						PackageInfos: []types.PackageInfo{
							{
								Packages: types.Packages{
									{
										Name:       "musl",
										Version:    "1.2.3-r0",
										SrcName:    "musl",
										SrcVersion: "1.2.3-r0",
										Licenses:   []string{"MIT"},
										Identifier: types.PkgIdentifier{
											PURL: &types.PackageURL{
												PackageURL: packageurl.PackageURL{
													Type:      packageurl.TypeApk,
													Namespace: "alpine",
													Name:      "musl",
													Version:   "1.2.3-r0",
													Qualifiers: packageurl.Qualifiers{
														{
															Key:   "distro",
															Value: "3.16.0",
														},
													},
												},
											},
										},
										Ref: "pkg:apk/alpine/musl@1.2.3-r0?distro=3.16.0",
										Layer: types.Layer{
											DiffID: "sha256:dd565ff850e7003356e2b252758f9bdc1ff2803f61e995e24c7844f6297f8fc3",
										},
									},
								},
							},
						},
						Applications: []types.Application{
							{
								Type:     "composer",
								FilePath: "app/composer/composer.lock",
								Libraries: types.Packages{
									{
										Name:    "pear/log",
										Version: "1.13.1",
										Identifier: types.PkgIdentifier{
											PURL: &types.PackageURL{
												PackageURL: packageurl.PackageURL{
													Type:      packageurl.TypeComposer,
													Namespace: "pear",
													Name:      "log",
													Version:   "1.13.1",
												},
											},
										},
										Ref: "pkg:composer/pear/log@1.13.1",
										Layer: types.Layer{
											DiffID: "sha256:3c79e832b1b4891a1cb4a326ef8524e0bd14a2537150ac0e203a5677176c1ca1",
										},
									},
									{

										Name:    "pear/pear_exception",
										Version: "v1.0.0",
										Identifier: types.PkgIdentifier{
											PURL: &types.PackageURL{
												PackageURL: packageurl.PackageURL{
													Type:      packageurl.TypeComposer,
													Namespace: "pear",
													Name:      "pear_exception",
													Version:   "v1.0.0",
												},
											},
										},
										Ref: "pkg:composer/pear/pear_exception@v1.0.0",
										Layer: types.Layer{
											DiffID: "sha256:3c79e832b1b4891a1cb4a326ef8524e0bd14a2537150ac0e203a5677176c1ca1",
										},
									},
								},
							},
							{
								Type:     "gobinary",
								FilePath: "app/gobinary/gobinary",
								Libraries: types.Packages{
									{
										Name:    "github.com/package-url/packageurl-go",
										Version: "v0.1.1-0.20220203205134-d70459300c8a",
										Identifier: types.PkgIdentifier{
											PURL: &types.PackageURL{
												PackageURL: packageurl.PackageURL{
													Type:      packageurl.TypeGolang,
													Namespace: "github.com/package-url",
													Name:      "packageurl-go",
													Version:   "v0.1.1-0.20220203205134-d70459300c8a",
												},
											},
										},
										Ref: "pkg:golang/github.com/package-url/packageurl-go@v0.1.1-0.20220203205134-d70459300c8a",
										Layer: types.Layer{
											DiffID: "sha256:3c79e832b1b4891a1cb4a326ef8524e0bd14a2537150ac0e203a5677176c1ca1",
										},
									},
								},
							},
							{
								Type:     "jar",
								FilePath: "",
								Libraries: types.Packages{
									{
										Name:    "org.codehaus.mojo:child-project",
										Version: "1.0",
										Identifier: types.PkgIdentifier{
											PURL: &types.PackageURL{
												PackageURL: packageurl.PackageURL{
													Type:      packageurl.TypeMaven,
													Namespace: "org.codehaus.mojo",
													Name:      "child-project",
													Version:   "1.0",
												},
												FilePath: "app/maven/target/child-project-1.0.jar",
											},
										},
										Ref: "pkg:maven/org.codehaus.mojo/child-project@1.0?file_path=app%2Fmaven%2Ftarget%2Fchild-project-1.0.jar",
										Layer: types.Layer{
											DiffID: "sha256:3c79e832b1b4891a1cb4a326ef8524e0bd14a2537150ac0e203a5677176c1ca1",
										},
										FilePath: "app/maven/target/child-project-1.0.jar",
									},
								},
							},
							{
								Type:     "node-pkg",
								FilePath: "",
								Libraries: types.Packages{
									{
										Name:    "bootstrap",
										Version: "5.0.2",
										Identifier: types.PkgIdentifier{
											PURL: &types.PackageURL{
												PackageURL: packageurl.PackageURL{
													Type:    packageurl.TypeNPM,
													Name:    "bootstrap",
													Version: "5.0.2",
												},
												FilePath: "app/app/package.json",
											},
										},
										Ref:      "pkg:npm/bootstrap@5.0.2?file_path=app%2Fapp%2Fpackage.json",
										Licenses: []string{"MIT"},
										Layer: types.Layer{
											DiffID: "sha256:3c79e832b1b4891a1cb4a326ef8524e0bd14a2537150ac0e203a5677176c1ca1",
										},
										FilePath: "app/app/package.json",
									},
								},
							},
						},
					},
				},
				Returns: cache.ArtifactCachePutBlobReturns{},
			},
			want: types.ArtifactReference{
				Name: filepath.Join("testdata", "sbom.cdx.intoto.jsonl"),
				Type: types.ArtifactCycloneDX,
				ID:   "sha256:c1cc58e08422fd7606a8e9ee2b42bf722b7af8b703b895461c23b83956f33227",
				BlobIDs: []string{
					"sha256:c1cc58e08422fd7606a8e9ee2b42bf722b7af8b703b895461c23b83956f33227",
				},
			},
		},
		{
			name:     "sad path with no such directory",
			filePath: filepath.Join("testdata", "unknown.json"),
			wantErr: []string{
				"no such file or directory",
				"The system cannot find the file specified",
			},
		},
		{
			name:     "sad path PutBlob returns an error",
			filePath: filepath.Join("testdata", "os-only-bom.json"),
			putBlobExpectation: cache.ArtifactCachePutBlobExpectation{
				Args: cache.ArtifactCachePutBlobArgs{
					BlobID: "sha256:033dc76e6daf7d8ba439d678dc7e33400687098f3e9f563f6975adf4eb440eee",
					BlobInfo: types.BlobInfo{
						SchemaVersion: types.BlobJSONSchemaVersion,
						OS: types.OS{
							Family: "alpine",
							Name:   "3.16.0",
						},
						PackageInfos: []types.PackageInfo{
							{},
						},
					},
				},
				Returns: cache.ArtifactCachePutBlobReturns{
					Err: errors.New("error"),
				},
			},
			wantErr: []string{"failed to store blob"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := new(cache.MockArtifactCache)
			c.ApplyPutBlobExpectation(tt.putBlobExpectation)

			a, err := sbom.NewArtifact(tt.filePath, c, artifact.Option{})
			require.NoError(t, err)

			got, err := a.Inspect(context.Background())
			if len(tt.wantErr) > 0 {
				require.NotNil(t, err)
				found := false
				for _, wantErr := range tt.wantErr {
					if strings.Contains(err.Error(), wantErr) {
						found = true
						break
					}
				}
				assert.True(t, found)
				return
			}

			// Not compare the original CycloneDX report
			got.CycloneDX = nil

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
