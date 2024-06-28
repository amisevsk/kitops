package local

import (
	"encoding/json"
	"fmt"
	"kitops/pkg/lib/constants"
	"os"
	"path"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/content/oci"
	"oras.land/oras-go/v2/registry"
)

type LocalRepo interface {
	GetRepo() string
	GetIndex() (*ocispec.Index, error)
	oras.Target
	content.Deleter
	content.Untagger
}

type localRepo struct {
	storagePath string
	nameRef     string
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
	localIndex, err := parseIndex(constants.IndexJsonPathForRepo(storagePath, repo.nameRef))
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

// // Push implements LocalRepo.
// func (l *localRepo) Push(ctx context.Context, expected ocispec.Descriptor, content io.Reader) error {
// 	panic("unimplemented")
// }

// // Resolve implements LocalRepo.
// func (l *localRepo) Resolve(ctx context.Context, reference string) (ocispec.Descriptor, error) {
// 	panic("unimplemented")
// }

// // Tag implements LocalRepo.
// func (l *localRepo) Tag(ctx context.Context, desc ocispec.Descriptor, reference string) error {
// 	panic("unimplemented")
// }

// // Untag implements LocalRepo.
// func (l *localRepo) Untag(ctx context.Context, reference string) error {
// 	panic("unimplemented")
// }

// // getStorePath implements LocalRepo.
// func (l *localRepo) getStorePath() string {
// 	panic("unimplemented")
// }

var _ LocalRepo = (*localRepo)(nil)

// parseIndexJson parses the OCI index.json stored in the OCI index at storageHome
func parseIndex(indexPath string) (*ocispec.Index, error) {
	indexBytes, err := os.ReadFile(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &ocispec.Index{}, nil
		}
		return nil, fmt.Errorf("failed to read index: %w", err)
	}

	index := &ocispec.Index{}
	if err := json.Unmarshal(indexBytes, index); err != nil {
		return nil, fmt.Errorf("failed to parse index: %w", err)
	}

	return index, nil
}
