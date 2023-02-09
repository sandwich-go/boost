package xos

import (
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	// tmpPermissionForDirectory makes the destination directory writable,
	// so that stuff can be copied recursively even if any original directory is NOT writable.
	// See https://github.com/otiai10/copy/pull/9 for more information.
	tmpPermissionForDirectory = os.FileMode(0755)
)

// Copy copies src to dest, doesn't matter if src is a directory or a file.
func Copy(src, dest string, opt ...CopyOption) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}
	return switchboard(src, dest, info, NewCopyOptions(opt...))
}

// switchboard switches proper copy functions regarding file type, etc...
// If there would be anything else here, add a case to this switchboard.
func switchboard(src, dest string, info os.FileInfo, opt *CopyOptions) error {
	switch {
	case info.Mode()&os.ModeSymlink != 0:
		return onsymlink(src, dest, info, opt)
	case info.IsDir():
		return dcopy(src, dest, info, opt)
	default:
		return fcopy(src, dest, info, opt)
	}
}

// copy decide if this src should be copied or not.
// Because this "copy" could be called recursively,
// "info" MUST be given here, NOT nil.
func _copy(src, dest string, info os.FileInfo, opt *CopyOptions) error {
	skip, err := opt.Skip(src)
	if err != nil {
		return err
	}
	if skip {
		return nil
	}
	return switchboard(src, dest, info, opt)
}

// fcopy is for just a file,
// with considering existence of parent directory
// and file permission.
func fcopy(src, dest string, info os.FileInfo, opt *CopyOptions) error {
	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer fclose(f, &err)
	if err = os.Chmod(f.Name(), info.Mode()|opt.AddPermission); err != nil {
		return err
	}
	var s *os.File
	s, err = os.Open(src)
	if err != nil {
		return err
	}
	defer fclose(s, &err)
	if _, err = io.Copy(f, s); err != nil {
		return err
	}
	if opt.Sync {
		err = f.Sync()
	}
	return err
}

// dcopy is for a directory,
// with scanning contents inside the directory
// and pass everything to "copy" recursively.
func dcopy(srcdir, destdir string, info os.FileInfo, opt *CopyOptions) (err error) {
	originalMode := info.Mode()
	// Make dest dir with 0755 so that everything writable.
	if err = os.MkdirAll(destdir, tmpPermissionForDirectory); err != nil {
		return
	}
	// Recover dir mode with original one.
	defer chmod(destdir, originalMode|opt.AddPermission, &err)
	var contents []fs.FileInfo
	contents, err = ioutil.ReadDir(srcdir)
	if err != nil {
		return
	}
	for _, content := range contents {
		cs, cd := filepath.Join(srcdir, content.Name()), filepath.Join(destdir, content.Name())
		if err = _copy(cs, cd, content, opt); err != nil {
			// If any error, exit immediately
			return
		}
	}
	return
}

func onsymlink(src, dest string, info os.FileInfo, opt *CopyOptions) error {
	switch opt.OnSymlink(src) {
	case Shallow:
		return lcopy(src, dest)
	case Deep:
		orig, err := filepath.EvalSymlinks(src)
		if err != nil {
			return err
		}
		info, err = os.Lstat(orig)
		if err != nil {
			return err
		}
		return _copy(orig, dest, info, opt)
	case Skip:
		fallthrough
	default:
		return nil // do nothing
	}
}

// lcopy is for a symlink,
// with just creating a new symlink by replicating src symlink.
func lcopy(src, dest string) error {
	linkSrc, err := os.Readlink(src)
	if err != nil {
		return err
	}
	return os.Symlink(linkSrc, dest)
}

// fclose ANYHOW closes file,
// with asiging error raised during Close,
// BUT respecting the error already reported.
func fclose(f *os.File, reported *error) {
	if err := f.Close(); *reported == nil {
		*reported = err
	}
}

// chmod ANYHOW changes file mode,
// with asiging error raised during Chmod,
// BUT respecting the error already reported.
func chmod(dir string, mode os.FileMode, reported *error) {
	if err := os.Chmod(dir, mode); *reported == nil {
		*reported = err
	}
}
