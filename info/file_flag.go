package info

import (
	"strings"
)

type FileFlag uint32

type FileCompression uint8

const (
	FileImplode      FileFlag = 0x00000100 // Implode method (By PKWARE Data Compression Library)
	FileCompress     FileFlag = 0x00000200 // Compress methods (By multiple methods)
	FileEncrypted    FileFlag = 0x00010000 // Indicates an encrypted file
	FileKey_v2       FileFlag = 0x00020000 // Indicates an encrypted file with key v2
	FilePatchFile    FileFlag = 0x00100000 // The file is a patch file. Raw file data begin with TPatchInfo structure
	FileSingleUnit   FileFlag = 0x01000000 // File is stored as a single unit, rather than split into sectors (Thx, Quantam)
	FileDeleteMarker FileFlag = 0x02000000 // File is a deletion marker. Used in MPQ patches, indicating that the file no longer exists.
	FileSectorCRC    FileFlag = 0x04000000 // File has checksums for each sector. Ignored if file is not compressed or imploded.
	FileSignature    FileFlag = 0x10000000 // Present on STANDARD.SNP\(signature). The only occurence ever observed
	FileExists       FileFlag = 0x80000000 // Set if file exists, reset when the file was deleted
	FileCompressMask FileFlag = 0x0000FF00 // Mask for a file being compressed
	FileFixKey       FileFlag = 0x00020000 // Obsolete, do not use
)

func (flag FileFlag) String() string {
	var readout []string
	if flag&FileImplode != 0 {
		readout = append(readout, "uses PKWARE compression method")
	}
	if flag&FileCompress != 0 {
		readout = append(readout, "uses multiple-compression method")
	}
	if flag&FileEncrypted != 0 {
		readout = append(readout, "an encrypted file")
	}
	if flag&FileKey_v2 != 0 {
		readout = append(readout, "an encrypted file (v2)")
	}
	if flag&FilePatchFile != 0 {
		readout = append(readout, "a patch file")
	}
	if flag&FileSingleUnit != 0 {
		readout = append(readout, "stored as a single unit")
	}
	if flag&FileDeleteMarker != 0 {
		readout = append(readout, "file deletion marker")
	}
	if flag&FileSectorCRC != 0 {
		readout = append(readout, "file contains checksums for each sector")
	}
	if flag&FileSignature != 0 {
		readout = append(readout, "signature flag (?)")
	}
	if flag&FileExists != 0 {
		readout = append(readout, "file exists (?)")
	}
	if flag&FileFixKey != 0 {
		readout = append(readout, "fix key (DEPRECATED!)")
	}
	return strings.Join(readout, ", ")
}
