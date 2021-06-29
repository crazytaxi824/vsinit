package golang

import "local/src/util"

type foldersAndFiles struct {
	folders []string
	files   []util.FileContent
	cipath  string
}

func initFoldersAndFiles(folders []string, files []util.FileContent) foldersAndFiles {
	var ff foldersAndFiles
	ff.folders = folders
	ff.files = files

	return ff
}

func (ff *foldersAndFiles) addFiles(files ...util.FileContent) {
	ff.files = append(ff.files, files...)
}

func (ff *foldersAndFiles) addFolders(folders ...string) {
	ff.folders = append(ff.folders, folders...)
}
