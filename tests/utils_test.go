package tests

import (
	"os"
	"testing"

	ipfs "github.com/beyondstorage/go-service-ipfs"
	"github.com/beyondstorage/go-storage/v4/pairs"
	"github.com/beyondstorage/go-storage/v4/types"
	"github.com/google/uuid"
)

func setupTest(t *testing.T) types.Storager {
	t.Log("Setup test for IPFS")

	store, err := ipfs.NewStorager(
		pairs.WithName(os.Getenv("STORAGE_IPFS_NAME")),
		pairs.WithEndpoint("STORAGE_IPFS_ENDPOINT"),
		pairs.WithWorkDir("/"+uuid.New().String()+"/"),
	)
	if err != nil {
		t.Errorf("new storager: %v", err)
	}
	return store
}
