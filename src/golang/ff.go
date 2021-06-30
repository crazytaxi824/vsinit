package golang

import "local/src/util"

type foldersAndFiles struct {
	folders     []string
	files       []util.FileContent
	suggestions []*util.Suggestion // FIXME 全部改好之后这里不使用 * 类型
	cipath      string
}

func initFoldersAndFiles(folders []string, files []util.FileContent) foldersAndFiles {
	var ff foldersAndFiles
	ff.folders = folders
	ff.files = files

	return ff
}

func (ff *foldersAndFiles) _addFiles(files ...util.FileContent) {
	ff.files = append(ff.files, files...)
}

func (ff *foldersAndFiles) _addFolders(folders ...string) {
	ff.folders = append(ff.folders, folders...)
}

func (ff *foldersAndFiles) _addSuggestion(sug ...*util.Suggestion) {
	ff.suggestions = append(ff.suggestions, sug...)
}
