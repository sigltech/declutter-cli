package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/tabwriter"
)

type Directory string

type FileExtension string

type DirInfo struct {
	Path  string
	Files []os.FileInfo
}

type NewDirsMap struct {
	Images      DirInfo
	Documents   DirInfo
	Pdf         DirInfo
	Media       DirInfo
	Zips        DirInfo
	Programming DirInfo
	Other       DirInfo
}

type FileHandler struct {
	DirectoryMap *NewDirsMap
	RawFiles     []os.FileInfo
	Excluded     []string
	WorkingDir   string
}

type DROutputColumns struct {
	FileName string
	NewDir   string
}

const (
	/* ------ Directories ------ */
	IMAGES      string = "images"
	DOCS        string = "documents"
	PDFS        string = "pdf"
	MEDIA       string = "media"
	ZIPS        string = "zips"
	PROGRAMMING string = "programming"
	OTHER       string = "other"

	/* ------ Image -------- */
	JPG FileExtension = ".jpg"
	PNG FileExtension = ".png"

	/* ------ Document-------- */
	PDF  FileExtension = ".pdf"
	DOC  FileExtension = ".doc"
	DOCX FileExtension = ".docx"
	XLS  FileExtension = ".xls"
	XLSX FileExtension = ".xlsx"
	PPT  FileExtension = ".ppt"
	PPTX FileExtension = ".pptx"
	TXT  FileExtension = ".txt"

	/* -------- Zips -------- */
	ZIP FileExtension = ".zip"

	/* -------- Media -------- */
	MP3 FileExtension = ".mp3"
	MP4 FileExtension = ".mp4"
	MPG FileExtension = ".mpg"

	/* -- Programming files--- */
	SQL  FileExtension = ".sql"
	DLL  FileExtension = ".dll"
	EXE  FileExtension = ".exe"
	JS   FileExtension = ".js"
	TS   FileExtension = ".ts"
	CSS  FileExtension = ".css"
	HTML FileExtension = ".html"
	TSX  FileExtension = ".tsx"
	JSX  FileExtension = ".jsx"
)

func InitializeFileHandler(pwd string, excluded []string) *FileHandler {
	absolutePath, err := filepath.Abs(pwd)
	if err != nil {
		panic(err)
	}
	newDirsMap := &NewDirsMap{
		Images: DirInfo{
			Path:  mustGetNewPath(absolutePath, IMAGES),
			Files: []os.FileInfo{},
		},
		Documents: DirInfo{
			Path:  mustGetNewPath(absolutePath, DOCS),
			Files: []os.FileInfo{},
		},
		Pdf: DirInfo{
			Path:  mustGetNewPath(absolutePath, PDFS),
			Files: []os.FileInfo{},
		},
		Media: DirInfo{
			Path:  mustGetNewPath(absolutePath, MEDIA),
			Files: []os.FileInfo{},
		},
		Zips: DirInfo{
			Path:  mustGetNewPath(absolutePath, ZIPS),
			Files: []os.FileInfo{},
		},
		Programming: DirInfo{
			Path:  mustGetNewPath(absolutePath, PROGRAMMING),
			Files: []os.FileInfo{},
		},
		Other: DirInfo{
			Path:  mustGetNewPath(absolutePath, OTHER),
			Files: []os.FileInfo{},
		},
	}
	fileHandler := &FileHandler{
		DirectoryMap: newDirsMap,
		WorkingDir:   absolutePath,
		RawFiles:     []os.FileInfo{},
	}

	// get the names of the excluded files
	for _, file := range excluded {
		fileHandler.Excluded = append(fileHandler.Excluded, file)
	}

	return fileHandler
}

// SortFiles sorts the file extensions to the directory it will be moved to.
// This adds the file to the appropriate key in the directoryMap
func (fileHandler *FileHandler) SortFiles() *FileHandler {
	files := fileHandler.RawFiles
	for _, file := range files {
		switch MustGetExtension(file.Name()) {
		case
			JPG,
			PNG:
			fileHandler.DirectoryMap.Images.Files = append(fileHandler.DirectoryMap.Images.Files, file)
		case DOC, DOCX, XLS, XLSX, PPT, PPTX, TXT, SQL:
			fileHandler.DirectoryMap.Documents.Files = append(fileHandler.DirectoryMap.Documents.Files, file)
		case ZIP:
			fileHandler.DirectoryMap.Zips.Files = append(fileHandler.DirectoryMap.Zips.Files, file)
		case MP3, MP4, MPG:
			fileHandler.DirectoryMap.Media.Files = append(fileHandler.DirectoryMap.Media.Files, file)
		case PDF:
			fileHandler.DirectoryMap.Pdf.Files = append(fileHandler.DirectoryMap.Pdf.Files, file)
		case DLL, EXE, JS, TS, CSS, HTML, TSX, JSX:
			fileHandler.DirectoryMap.Programming.Files = append(fileHandler.DirectoryMap.Programming.Files, file)
		default:
			fileHandler.DirectoryMap.Other.Files = append(fileHandler.DirectoryMap.Other.Files, file)
		}
	}

	return fileHandler
}

// FilterDirectories filters the objects in the current directory
// to only include files and not folders
func (fileHandler *FileHandler) FilterDirectories() *FileHandler {
	// get the filesAndFolders in the current directory
	filesAndFolders, err := os.ReadDir(fileHandler.WorkingDir)
	if err != nil {
		panic(err)
	}
	files := make([]os.FileInfo, 0)
	for _, obj := range filesAndFolders {
		if !obj.IsDir() {
			info, err := obj.Info()
			if err != nil {
				panic(err)
			}
			files = append(files, info)
		}
	}
	fileHandler.RawFiles = files
	return fileHandler
}

// DryRun prints the files and their new directories
// in a table format, without moving the files. This allows the user to see
// what will be moved before they confirm.
func (fileHandler *FileHandler) DryRun() {
	fmt.Println("Dry run mode")
	fmt.Println("")
	// ANSI escape codes for bold text
	bold := "\033[1m"
	reset := "\033[0m"
	// Use reflection to iterate over the struct fields
	v := reflect.ValueOf(*fileHandler.DirectoryMap)
	var files []DROutputColumns
	var excludedFiles []struct {
		FileName string
	}

	for i := 0; i < v.NumField(); i++ {
		dirInfo := v.Field(i).Interface().(DirInfo)
		for _, file := range dirInfo.Files {
			if isExcluded(file, fileHandler.Excluded) {
				excludedFiles = append(excludedFiles, struct {
					FileName string
				}{
					FileName: file.Name(),
				})
				continue
			}
			files = append(files, DROutputColumns{
				FileName: file.Name(),
				NewDir:   dirInfo.Path,
			})
		}
	}

	// Create a tab writer
	writer := tabwriter.NewWriter(os.Stdout, 2, 0, 4, ' ', 0)

	// Write headers
	_, _ = fmt.Fprintln(writer, " File Name \t New Directory \t")
	_, _ = fmt.Fprintln(writer, "-----------\t---------------\t")

	// Write rows
	for _, file := range files {
		// Check if the filename is too long
		ellipsisFilename := ellipsis(file.FileName, 60)
		// Dynamically create the arrow with the required number of spaces
		arrowTail := constructArrow(ellipsisFilename)

		// Print the row with the arrow
		_, err := fmt.Fprintf(writer, "%s %s \t%s\t\n", ellipsisFilename, arrowTail, file.NewDir)
		if err != nil {
			panic(err)
		}
	}

	// Write headers
	_, _ = fmt.Fprintf(writer, "\n%sExcluded Files%s\t\t\n\n", bold, reset)
	_, _ = fmt.Fprintln(writer, " File Name")
	_, _ = fmt.Fprintln(writer, "-----------")
	// Write excluded file rows
	for _, file := range excludedFiles {
		// Check if the filename is too long
		ellipsisFilename := ellipsis(file.FileName, 60)
		// Dynamically create the arrow with the required number of spaces

		// Print the row with the arrow
		_, err := fmt.Fprintf(writer, "%s\n", ellipsisFilename)
		if err != nil {
			panic(err)
		}
	}

	// Flush the writer
	if err := writer.Flush(); err != nil {
		return
	}
}

func (fileHandler *FileHandler) MoveFiles() {
	// Reflect on DirectoryMap
	v := reflect.ValueOf(fileHandler.DirectoryMap)

	// Dereference the pointer
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Iterate over the fields dynamically
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		// Get the DirInfo value
		dirInfo := field.Interface().(DirInfo)

		// Ensure the directory exists
		if _, err := os.Stat(dirInfo.Path); os.IsNotExist(err) {
			fmt.Printf("Creating directory: %s\n", dirInfo.Path)
			err := os.Mkdir(dirInfo.Path, 0755)
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Printf("Directory already exists: %s\n", dirInfo.Path)
		}

		// Move files to the directory
		for _, file := range dirInfo.Files {
			if isExcluded(file, fileHandler.Excluded) {
				continue
			}
			newPath := dirInfo.Path + "/" + file.Name()
			oldPath := fileHandler.WorkingDir + "/" + file.Name()

			fmt.Printf("Moving file: %s --> %s\n", oldPath, newPath)
			err := os.Rename(oldPath, newPath)
			if err != nil {
				panic(err)
			}
		}
	}
}

func MustGetExtension(filename string) FileExtension {
	return FileExtension("." + filename[strings.LastIndex(filename, ".")+1:])
}

func mustGetNewPath(pwd string, newDirName string) string {
	return pwd + "/" + newDirName
}

func ellipsis(s string, maxLength int) string {
	if len(s) > maxLength {
		return s[:maxLength] + "..."
	}
	return s
}

func constructArrow(filename string) string {
	arrow := "-->"
	spacesNeeded := 70 - len(filename) - len(arrow)
	if spacesNeeded < 0 {
		spacesNeeded = 0
	}
	return strings.Repeat("-", spacesNeeded) + arrow
}

func isExcluded(file os.FileInfo, excluded []string) bool {
	for _, excludedFile := range excluded {
		if file.Name() == excludedFile {
			return true
		}
	}
	return false
}
