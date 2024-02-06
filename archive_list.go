package mpq

import "bufio"

type archive_list struct {
	listfile *File
	scanner  *bufio.Scanner
}

func (archive *Archive) List() (list List, err error) {
	// create archive_list
	archive_list := new(archive_list)
	// open listfile
	archive_list.listfile, err = archive.Open("(listfile)")
	if err != nil {
		return
	}
	// create scanner
	archive_list.scanner = bufio.NewScanner(archive_list.listfile)
	// use custom scan function to process listfile
	archive_list.scanner.Split(scan_lines)
	list = archive_list
	return
}

func (list *archive_list) Next() bool {
	for {
		if !list.scanner.Scan() {
			return false
		}

		path := list.scanner.Text()
		if path == "" {
			continue
		}

		return true
	}

}

func (list *archive_list) Close() error {
	list.scanner = nil
	return list.listfile.Close()
}

func (list *archive_list) Path() string {
	return list.scanner.Text()
}
