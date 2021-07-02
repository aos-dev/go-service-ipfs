package ipfs

import (
	"errors"
	"fmt"

	ipfs "github.com/ipfs/go-ipfs-api"

	"github.com/beyondstorage/go-endpoint"
	"github.com/beyondstorage/go-storage/v4/services"
	"github.com/beyondstorage/go-storage/v4/types"
)

// Storage is the example client.
type Storage struct {
	ipfs *ipfs.Shell

	defaultPairs DefaultStoragePairs
	features     StorageFeatures

	workDir string

	types.UnimplementedStorager
	types.UnimplementedCopier
	types.UnimplementedMover
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf("Storager IPFS {WorkDir: %s}", s.workDir)
}

// NewStorager will create Storager only.
func NewStorager(pairs ...types.Pair) (types.Storager, error) {
	opt, err := parsePairStorageNew(pairs)
	if err != nil {
		return nil, err
	}

	st := &Storage{
		workDir: "/",
	}
	if opt.HasWorkDir {
		st.workDir = opt.WorkDir
	}

	ep, err := endpoint.Parse(opt.Endpoint)
	if err != nil {
		return nil, err
	}

	e := "http://localhost:5001"
	switch ep.Protocol() {
	case endpoint.ProtocolHTTP:
		e, _, _ = ep.HTTP()
	case endpoint.ProtocolHTTPS:
		e, _, _ = ep.HTTPS()
	}

	sh := ipfs.NewShell(e)
	if !sh.IsUp() {
		return nil, errors.New("ipfs not online")
	}
	st.ipfs = sh

	return st, nil
}

func formatError(err error) error {
	if _, ok := err.(services.InternalError); ok {
		return err
	}

	e, ok := err.(*ipfs.Error)
	if !ok {
		return fmt.Errorf("%w: %v", services.ErrUnexpected, err)
	}

	switch e.Message {
	case "file does not exist":
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	}

	return fmt.Errorf("%w: %v", services.ErrUnexpected, err)
}

func (s *Storage) formatError(op string, err error, path ...string) error {
	if err == nil {
		return nil
	}

	return services.StorageError{
		Op:       op,
		Err:      formatError(err),
		Storager: s,
		Path:     path,
	}
}

// getAbsPath will calculate object storage's abs path
func (s *Storage) getAbsPath(path string) string {
	return s.workDir + path
}
