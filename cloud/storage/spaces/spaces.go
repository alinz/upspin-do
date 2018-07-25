package spaces

import (
	"bytes"
	"fmt"
	"os"

	minio "github.com/minio/minio-go"

	"upspin.io/cloud/storage"
	"upspin.io/errors"
)

// Keys used for storing dial options.
const (
	regionName = "spacesRegion"
	spaceName  = "spaceName"
)

// spacesImpl is an implementation of Storage that connects to an Amazon Simple
// Storage (S3) backend.
type spacesImpl struct {
	client    *minio.Client
	spaceName string
	endpoint  string
}

// New initializes a Storage implementation that stores data to Spaces Simple
// Storage Service.
func New(opts *storage.Opts) (storage.Storage, error) {
	const op errors.Op = "cloud/storage/spaces.New"
	const ssl = true

	accessKey := os.Getenv("SPACES_KEY")
	if accessKey == "" {
		return nil, errors.E(op, errors.Invalid, errors.Errorf("SPACES_KEY env variable is required"))
	}

	secKey := os.Getenv("SPACES_SECRET")
	if secKey == "" {
		return nil, errors.E(op, errors.Invalid, errors.Errorf("SPACES_SECRET env variable is required"))
	}

	region, ok := opts.Opts[regionName]
	if !ok {
		return nil, errors.E(op, errors.Invalid, errors.Errorf("%q option is required", regionName))
	}

	name, ok := opts.Opts[spaceName]
	if !ok {
		return nil, errors.E(op, errors.Invalid, errors.Errorf("%q option is required", name))
	}

	endpoint := fmt.Sprintf("%s.digitaloceanspaces.com", region)

	// Initiate a client using DigitalOcean Spaces.
	client, err := minio.New(endpoint, accessKey, secKey, ssl)
	if err != nil {
		return nil, errors.E(op, errors.IO, errors.Errorf("unable to create minio session: %s", err))
	}

	return &spacesImpl{
		client:    client,
		spaceName: name,
		endpoint:  endpoint,
	}, nil
}

func init() {
	storage.Register("Spaces", New)
}

// Guarantee we implement the Storage interface.
var _ storage.Storage = (*spacesImpl)(nil)

// LinkBase implements Storage.
func (s *spacesImpl) LinkBase() (base string, err error) {
	return fmt.Sprintf("%s.%s", s.spaceName, s.endpoint), nil
}

// Download implements Storage.
func (s *spacesImpl) Download(ref string) ([]byte, error) {
	const op errors.Op = "cloud/storage/spaces.Download"

	obj, err := s.client.GetObject(s.spaceName, ref, minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.E(op, errors.IO, errors.Errorf(
			"unable to download ref %q from bucket %q: %s", ref, s.spaceName, err))
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(obj)
	return buf.Bytes(), nil
}

// Put implements Storage.
func (s *spacesImpl) Put(ref string, contents []byte) error {
	const op errors.Op = "cloud/storage/spaces.Put"

	_, err := s.client.PutObject(s.spaceName, ref, bytes.NewReader(contents), int64(len(contents)), minio.PutObjectOptions{})
	if err != nil {
		return errors.E(op, errors.IO, errors.Errorf(
			"unable to upload ref %q to bucket %q: %s", ref, s.spaceName, err))
	}

	return nil
}

// Delete implements Storage.
func (s *spacesImpl) Delete(ref string) error {
	const op errors.Op = "cloud/storage/spaces.Delete"

	err := s.client.RemoveObject(s.spaceName, ref)
	if err != nil {
		return errors.E(op, errors.IO, errors.Errorf(
			"unable to delete ref %q from bucket %q: %s", ref, s.spaceName, err))
	}

	return nil
}

// Close implements Storage.
func (s *spacesImpl) Close() {
	s.client = nil
	s.spaceName = ""
}
