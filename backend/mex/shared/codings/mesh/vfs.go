package mesh

import (
	"bytes"
	"fmt"

	"github.com/psanford/sqlite3vfs"
)

const VFSName = "mem:ro"

// VFS type which only has the Open method implemented to work for our needs.
type inmemVFS struct {
	data []byte
}

func NewReadOnlyInMemoryVFS(data []byte) sqlite3vfs.VFS {
	return &inmemVFS{data: data}
}

func (x *inmemVFS) Open(_ string, _ sqlite3vfs.OpenFlag) (sqlite3vfs.File, sqlite3vfs.OpenFlag, error) {
	return vfsFile{bytes.NewReader(x.data)}, sqlite3vfs.OpenReadOnly | sqlite3vfs.OpenExclusive, nil
}

func (x *inmemVFS) Delete(_ string, _ bool) error {
	return fmt.Errorf("Delete: not impl")
}

func (x *inmemVFS) Access(_ string, _ sqlite3vfs.AccessFlag) (bool, error) {
	return false, fmt.Errorf("Access: not impl")
}

func (x *inmemVFS) FullPathname(name string) string {
	return name
}

type vfsFile struct {
	*bytes.Reader
}

// We define the methods on the value receiver (rather than a pointer) as the struct just embeds a pointer.
// Also, the method ReadAt is reused from the embedded *bytes.Reader.

func (b vfsFile) FileSize() (int64, error) {
	return int64(b.Reader.Len()), nil
}

func (b vfsFile) DeviceCharacteristics() sqlite3vfs.DeviceCharacteristic {
	return sqlite3vfs.IocapImmutable | sqlite3vfs.IocapUndeletableWhenOpen
}

// Unimplemented methods (used only for writing)

func (b vfsFile) Close() error {
	return nil
}

func (b vfsFile) WriteAt(p []byte, off int64) (int, error) {
	return -1, fmt.Errorf("WriteAt: not impl")
}

func (b vfsFile) Truncate(size int64) error {
	return fmt.Errorf("Truncate: not impl")
}

func (b vfsFile) Sync(flag sqlite3vfs.SyncType) error {
	return fmt.Errorf("Sync: not impl")
}

func (b vfsFile) Lock(elock sqlite3vfs.LockType) error {
	return fmt.Errorf("Lock: not impl")
}

func (b vfsFile) Unlock(elock sqlite3vfs.LockType) error {
	return fmt.Errorf("Unlock: not impl")
}

func (b vfsFile) CheckReservedLock() (bool, error) {
	return false, fmt.Errorf("CheckReservedLock: not impl")
}

func (b vfsFile) SectorSize() int64 {
	return -1
}
