package local

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kitops/pkg/lib/constants"
	"os"
	"path"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/content/oci"
	"oras.land/oras-go/v2/errdef"
	"oras.land/oras-go/v2/registry"
)

type LocalRepo interface {
	GetRepo() string
	GetIndex() (*ocispec.Index, error)
	getStorePath() string // TODO TODO TODO: We don't need this anymore
	oras.Target
	content.Deleter
	content.Untagger
}

type localRepo struct {
	storagePath string
	nameRef     string
	indexPath   string
	localIndex  *ocispec.Index
	*oci.Store
}

func NewLocalRepo(storagePath string, ref *registry.Reference) (LocalRepo, error) {
	repo := &localRepo{}
	repo.storagePath = storagePath
	repo.nameRef = path.Join(ref.Registry, ref.Repository)

	store, err := oci.New(storagePath)
	if err != nil {
		return nil, err
	}
	repo.Store = store

	// Initialize repo-specific index.json
	repo.indexPath = constants.IndexJsonPathForRepo(storagePath, repo.nameRef)
	localIndex, err := parseIndex(repo.indexPath)
	if err != nil {
		return nil, err
	}
	repo.localIndex = localIndex

	return repo, nil
}

// TODO: Temp implementations to make testing a little easier.
func (r *localRepo) GetIndex() (*ocispec.Index, error) {
	return r.localIndex, nil
}

// GetRepo returns the registry and repository for the current OCI store.
func (r *localRepo) GetRepo() string {
	return r.nameRef
}

// TODO: Is this still needed?
func (r *localRepo) getStorePath() string {
	return r.storagePath
}

// // Delete implements LocalRepo.
// func (l *localRepo) Delete(ctx context.Context, target ocispec.Descriptor) error {
// 	panic("unimplemented")
// }

// // Exists implements LocalRepo.
// func (l *localRepo) Exists(ctx context.Context, target ocispec.Descriptor) (bool, error) {
// 	panic("unimplemented")
// }

// // Fetch implements LocalRepo.
// func (l *localRepo) Fetch(ctx context.Context, target ocispec.Descriptor) (io.ReadCloser, error) {
// 	panic("unimplemented")
// }

func (l *localRepo) Push(ctx context.Context, expected ocispec.Descriptor, content io.Reader) error {
	if err := l.Store.Push(ctx, expected, content); err != nil {
		return err
	}
	if expected.MediaType == ocispec.MediaTypeImageManifest {
		l.addManifestToLocalIndex(expected)
		return l.saveLocalIndex()
	}
	return nil
}

// // Resolve implements LocalRepo.
// func (l *localRepo) Resolve(ctx context.Context, reference string) (ocispec.Descriptor, error) {
// 	panic("unimplemented")
// }

func (l *localRepo) Tag(ctx context.Context, desc ocispec.Descriptor, reference string) error {
	// TODO: should we tag it in the general index.json too?
	// TODO: should probably de-duplicate this (don't store a manifest without a tag)
	descExists := false
	for _, m := range l.localIndex.Manifests {
		tag := m.Annotations[ocispec.AnnotationRefName]
		if m.Digest == desc.Digest {
			if tag == reference {
				return nil
			}
			descExists = true
		}
	}
	if !descExists {
		return fmt.Errorf("%s: %s: %w", desc.Digest, desc.MediaType, errdef.ErrNotFound)
	}
	if desc.Annotations == nil {
		desc.Annotations = map[string]string{}
	}
	desc.Annotations[ocispec.AnnotationRefName] = reference
	l.addManifestToLocalIndex(desc)
	return l.saveLocalIndex()
}

// // Untag implements LocalRepo.
// func (l *localRepo) Untag(ctx context.Context, reference string) error {
// 	panic("unimplemented")
// }

// // getStorePath implements LocalRepo.
// func (l *localRepo) getStorePath() string {
// 	panic("unimplemented")
// }

func (l *localRepo) addManifestToLocalIndex(manifestDesc ocispec.Descriptor) {
	// TODO: consider using ORAS' tag resolver to make this a little cleaner
	curTag := manifestDesc.Annotations[ocispec.AnnotationRefName]
	for _, m := range l.localIndex.Manifests {
		manifestTag := m.Annotations[ocispec.AnnotationRefName]
		if m.Digest == manifestDesc.Digest && manifestTag == curTag {
			// Already included
			return
		}
	}
	l.localIndex.Manifests = append(l.localIndex.Manifests, manifestDesc)
}

func (l *localRepo) saveLocalIndex() error {
	indexJson, err := json.Marshal(l.localIndex)
	if err != nil {
		return fmt.Errorf("failed to marshal index: %w", err)
	}
	return os.WriteFile(l.indexPath, indexJson, 0666)
}

var _ LocalRepo = (*localRepo)(nil)

// parseIndexJson parses the OCI index.json stored in the OCI index at storageHome
func parseIndex(indexPath string) (*ocispec.Index, error) {
	indexBytes, err := os.ReadFile(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			index := &ocispec.Index{}
			index.SchemaVersion = 2
			return index, nil
		}
		return nil, fmt.Errorf("failed to read index: %w", err)
	}

	index := &ocispec.Index{}
	if err := json.Unmarshal(indexBytes, index); err != nil {
		return nil, fmt.Errorf("failed to parse index: %w", err)
	}

	return index, nil
}
