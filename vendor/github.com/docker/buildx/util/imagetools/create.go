package imagetools

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/platforms"
	"github.com/containerd/containerd/remotes"
	"github.com/distribution/reference"
	"github.com/docker/buildx/util/buildflags"
	"github.com/moby/buildkit/exporter/containerimage/exptypes"
	"github.com/moby/buildkit/util/contentutil"
	"github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type Source struct {
	Desc ocispec.Descriptor
	Ref  reference.Named
}

func (r *Resolver) Combine(ctx context.Context, srcs []*Source, ann []string) ([]byte, ocispec.Descriptor, error) {
	eg, ctx := errgroup.WithContext(ctx)

	dts := make([][]byte, len(srcs))
	for i := range dts {
		func(i int) {
			eg.Go(func() error {
				dt, err := r.GetDescriptor(ctx, srcs[i].Ref.String(), srcs[i].Desc)
				if err != nil {
					return err
				}
				dts[i] = dt

				if srcs[i].Desc.MediaType == "" {
					mt, err := detectMediaType(dt)
					if err != nil {
						return err
					}
					srcs[i].Desc.MediaType = mt
				}

				mt := srcs[i].Desc.MediaType

				switch mt {
				case images.MediaTypeDockerSchema2Manifest, ocispec.MediaTypeImageManifest:
					p := srcs[i].Desc.Platform
					if srcs[i].Desc.Platform == nil {
						p = &ocispec.Platform{}
					}
					if p.OS == "" || p.Architecture == "" {
						if err := r.loadPlatform(ctx, p, srcs[i].Ref.String(), dt); err != nil {
							return err
						}
					}
					srcs[i].Desc.Platform = p
				case images.MediaTypeDockerSchema1Manifest:
					return errors.Errorf("schema1 manifests are not allowed in manifest lists")
				}

				return nil
			})
		}(i)
	}

	if err := eg.Wait(); err != nil {
		return nil, ocispec.Descriptor{}, err
	}

	// on single source, return original bytes
	if len(srcs) == 1 && len(ann) == 0 {
		if mt := srcs[0].Desc.MediaType; mt == images.MediaTypeDockerSchema2ManifestList || mt == ocispec.MediaTypeImageIndex {
			return dts[0], srcs[0].Desc, nil
		}
	}

	m := map[digest.Digest]int{}
	newDescs := make([]ocispec.Descriptor, 0, len(srcs))

	addDesc := func(d ocispec.Descriptor) {
		idx, ok := m[d.Digest]
		if ok {
			old := newDescs[idx]
			if old.MediaType == "" {
				old.MediaType = d.MediaType
			}
			if d.Platform != nil {
				old.Platform = d.Platform
			}
			if old.Annotations == nil {
				old.Annotations = map[string]string{}
			}
			for k, v := range d.Annotations {
				old.Annotations[k] = v
			}
			newDescs[idx] = old
		} else {
			m[d.Digest] = len(newDescs)
			newDescs = append(newDescs, d)
		}
	}

	for i, src := range srcs {
		switch src.Desc.MediaType {
		case images.MediaTypeDockerSchema2ManifestList, ocispec.MediaTypeImageIndex:
			var mfst ocispec.Index
			if err := json.Unmarshal(dts[i], &mfst); err != nil {
				return nil, ocispec.Descriptor{}, errors.WithStack(err)
			}
			for _, d := range mfst.Manifests {
				addDesc(d)
			}
		default:
			addDesc(src.Desc)
		}
	}

	dockerMfsts := 0
	for _, desc := range newDescs {
		if strings.HasPrefix(desc.MediaType, "application/vnd.docker.") {
			dockerMfsts++
		}
	}

	var mt string
	if dockerMfsts == len(newDescs) {
		// all manifests are Docker types, use Docker manifest list
		mt = images.MediaTypeDockerSchema2ManifestList
	} else {
		// otherwise, use OCI index
		mt = ocispec.MediaTypeImageIndex
	}

	// annotations are only allowed on OCI indexes
	indexAnnotation := make(map[string]string)
	if mt == ocispec.MediaTypeImageIndex {
		annotations, err := buildflags.ParseAnnotations(ann)
		if err != nil {
			return nil, ocispec.Descriptor{}, err
		}
		for k, v := range annotations {
			switch k.Type {
			case exptypes.AnnotationIndex:
				indexAnnotation[k.Key] = v
			case exptypes.AnnotationManifestDescriptor:
				for i := 0; i < len(newDescs); i++ {
					if newDescs[i].Annotations == nil {
						newDescs[i].Annotations = map[string]string{}
					}
					if k.Platform == nil || k.PlatformString() == platforms.Format(*newDescs[i].Platform) {
						newDescs[i].Annotations[k.Key] = v
					}
				}
			case exptypes.AnnotationManifest, "":
				return nil, ocispec.Descriptor{}, errors.Errorf("%q annotations are not supported yet", k.Type)
			case exptypes.AnnotationIndexDescriptor:
				return nil, ocispec.Descriptor{}, errors.Errorf("%q annotations are invalid while creating an image", k.Type)
			}
		}
	}

	idxBytes, err := json.MarshalIndent(ocispec.Index{
		MediaType: mt,
		Versioned: specs.Versioned{
			SchemaVersion: 2,
		},
		Manifests:   newDescs,
		Annotations: indexAnnotation,
	}, "", "  ")
	if err != nil {
		return nil, ocispec.Descriptor{}, errors.Wrap(err, "failed to marshal index")
	}

	return idxBytes, ocispec.Descriptor{
		MediaType: mt,
		Size:      int64(len(idxBytes)),
		Digest:    digest.FromBytes(idxBytes),
	}, nil
}

func (r *Resolver) Push(ctx context.Context, ref reference.Named, desc ocispec.Descriptor, dt []byte) error {
	ctx = remotes.WithMediaTypeKeyPrefix(ctx, "application/vnd.in-toto+json", "intoto")

	ref = reference.TagNameOnly(ref)
	p, err := r.resolver().Pusher(ctx, ref.String())
	if err != nil {
		return err
	}
	cw, err := p.Push(ctx, desc)
	if err != nil {
		if errdefs.IsAlreadyExists(err) {
			return nil
		}
		return err
	}

	err = content.Copy(ctx, cw, bytes.NewReader(dt), desc.Size, desc.Digest)
	if errdefs.IsAlreadyExists(err) {
		return nil
	}
	return err
}

func (r *Resolver) Copy(ctx context.Context, src *Source, dest reference.Named) error {
	ctx = remotes.WithMediaTypeKeyPrefix(ctx, "application/vnd.in-toto+json", "intoto")

	dest = reference.TagNameOnly(dest)
	p, err := r.resolver().Pusher(ctx, dest.String())
	if err != nil {
		return err
	}

	srcRef := reference.TagNameOnly(src.Ref)
	f, err := r.resolver().Fetcher(ctx, srcRef.String())
	if err != nil {
		return err
	}

	refspec := reference.TrimNamed(src.Ref).String()
	u, err := url.Parse("dummy://" + refspec)
	if err != nil {
		return err
	}
	source, repo := u.Hostname(), strings.TrimPrefix(u.Path, "/")
	if src.Desc.Annotations == nil {
		src.Desc.Annotations = make(map[string]string)
	}
	src.Desc.Annotations["containerd.io/distribution.source."+source] = repo

	err = contentutil.CopyChain(ctx, contentutil.FromPusher(p), contentutil.FromFetcher(f), src.Desc)
	if err != nil {
		return err
	}
	return nil
}

func (r *Resolver) loadPlatform(ctx context.Context, p2 *ocispec.Platform, in string, dt []byte) error {
	var manifest ocispec.Manifest
	if err := json.Unmarshal(dt, &manifest); err != nil {
		return errors.WithStack(err)
	}

	dt, err := r.GetDescriptor(ctx, in, manifest.Config)
	if err != nil {
		return err
	}

	var p ocispec.Platform
	if err := json.Unmarshal(dt, &p); err != nil {
		return errors.WithStack(err)
	}

	p = platforms.Normalize(p)

	if p2.Architecture == "" {
		p2.Architecture = p.Architecture
		if p2.Variant == "" {
			p2.Variant = p.Variant
		}
	}
	if p2.OS == "" {
		p2.OS = p.OS
	}

	return nil
}

func detectMediaType(dt []byte) (string, error) {
	var mfst struct {
		MediaType string          `json:"mediaType"`
		Config    json.RawMessage `json:"config"`
		FSLayers  []string        `json:"fsLayers"`
	}

	if err := json.Unmarshal(dt, &mfst); err != nil {
		return "", errors.WithStack(err)
	}

	if mfst.MediaType != "" {
		return mfst.MediaType, nil
	}
	if mfst.Config != nil {
		return images.MediaTypeDockerSchema2Manifest, nil
	}
	if len(mfst.FSLayers) > 0 {
		return images.MediaTypeDockerSchema1Manifest, nil
	}

	return images.MediaTypeDockerSchema2ManifestList, nil
}
