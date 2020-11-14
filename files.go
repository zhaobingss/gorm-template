package gorm_template

import (
	"io/ioutil"
	"os"
)

/// get all files from the dirPath
/// @param dirPath: the files dir
func GetFiles(dirPth string) (files []string, err error) {
	var dirs []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() {
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			fs, err := GetFiles(dirPth + PthSep + fi.Name())
			if err != nil {
				return nil, err
			} else {
				files = append(files, fs...)
			}
		} else {
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}

	return files, nil
}
