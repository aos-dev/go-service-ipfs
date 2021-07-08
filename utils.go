package ipfs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	ipfs "github.com/ipfs/go-ipfs-api"
	cmds "github.com/ipfs/go-ipfs-cmds"

	"github.com/beyondstorage/go-endpoint"
	ps "github.com/beyondstorage/go-storage/v4/pairs"
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

	var e string
	switch ep.Protocol() {
	case endpoint.ProtocolHTTP:
		e, _, _ = ep.HTTP()
	case endpoint.ProtocolHTTPS:
		e, _, _ = ep.HTTPS()
	default:
		return nil, services.PairUnsupportedError{Pair: ps.WithEndpoint(opt.Endpoint)}
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
	case os.ErrNotExist.Error():
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	}

	// ref: https://github.com/ipfs/go-ipfs-cmds/blob/4ade007405e5d3befb14184290576c63cc43a6a3/error.go#L31
	switch e.Code {
	case int(cmds.ErrRateLimited):
		return fmt.Errorf("%w: %v", services.ErrRequestThrottled, err)
	case int(cmds.ErrImplementation):
		return fmt.Errorf("%w: %v", services.ErrServiceInternal, err)
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
	if filepath.IsAbs(path) {
		return path
	}

	return s.workDir + path
}
