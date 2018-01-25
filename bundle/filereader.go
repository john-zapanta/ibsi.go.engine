package bundle

import (
	"os"
	"bufio"
    "bytes"
	"path/filepath"
	// "io"
	// "bufio"
    // "net/http"
    "strings"
    "log"
	"ibsi/utils"
	// "github.com/tdewolff/minify"
	// "github.com/tdewolff/minify/js"
)

type FileReader struct {
	Variables utils.IKeyValue
	Container bytes.Buffer
	filesParser BundlerCss
	fileData *FileData
	fileName string
	s string
	
	position int
	length int
	lines []string
	eof bool
}

func (fr *FileReader) Init() {
	fr.lines = make([]string, 0)
	fr.eof = false
	fr.length = 0
	fr.position = 0
	
	// fr.filesParser.Output("Test")
	
	// if _, err := os.Stat(fr.fileName); os.IsNotExist(err) {
		file, err := os.Open(fr.fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		
		// thanks to: https://stackoverflow.com/questions/21371673/reading-files-with-a-bom-in-go
		// bom := make([]byte, 3)
		// io.ReadFull(file, bom[:])
		
		// if bom[0] != 0xef || bom[1] != 0xbb || bom[1] != 0xbf {
			// file.Seek(0, 0) // Not a BOM -- seek back to the beginning
		// } else {
			// fr.lines = append(fr.lines, ".test2 {}")
		// }
		// _, err = io.ReadFull(fd, bom[:])
			// if err != nil {
				// log.Fatal(err)
			// }		
			
		i := 0
		scanner := bufio.NewScanner(file)
		for scanner.Scan() { 
			// fr.filesParser.Output(strings.TrimSpace(scanner.Text()))
			s := strings.TrimSpace(scanner.Text())
			if i == 0 {
				b := scanner.Bytes()
				_, bom, _ := bufio.ScanRunes(b, false)
				// if bom[0] != 0xef || bom[1] != 0xbb || bom[2] != 0xbf {
					
				if bom[0] == 0xef && bom[1] == 0xbb && bom[2] == 0xbf {
				// } else {
					s = s[3:len(s)]
					// s = "HELLO"
					
					// s[0] = 0x0
					// s[1] = 0x0
					// s[2] = 0x0
				}
				
				// fr.lines = append(fr.lines, string(s[0]))
				// s = strings.Replace(s, "\uFFFD", "", 1)
			}
			i++
			
			
			// if i == 0 {
				// fr.lines = append(fr.lines, string(s[0]))
				// s = strings.Replace(s, "\xef\xbf\xbd", "", 1)
			// }
			
			fr.lines = append(fr.lines, s)
		}
		
		if fr.length = len(fr.lines); fr.length > 0 {
			fr.s = fr.lines[fr.position];
		} else {
			fr.eof = true
		}
	// } else {
		// fr.filesParser.Output("Not exists daw") // debug
	// }
	
}

func (fr *FileReader) Line() string {
	return fr.s
}

func (fr *FileReader) Eof() bool {
	return fr.eof
}

func (fr *FileReader) Next() string {	
	if fr.eof = fr.position + 1 == fr.length; fr.eof {
		fr.s = "";
	} else {
		fr.position++
		fr.s = fr.lines[fr.position];
	}
	
	return fr.s
}

func (fr *FileReader) EmptyLine() bool {
	return fr.s == ""
}

func (fr *FileReader) IsCommentLine() bool {
	if !fr.EmptyLine() {
		if len(fr.s) > 1 {
			return fr.filesParser.IsComment(fr.s[0:2])
		} else {
			return false
		}
	} else {
		return true
	}
}

func (fr *FileReader) LineTokenIs(token string) bool {
	return strings.Index(fr.s, token) == 0
}

func (fr *FileReader) Break() []string {
	return strings.Split(fr.s, "=")
}

// func (fr *FileReader) KeyValue() struct{} {
	// values := strings.Split(fr.s, "=")
	// return struct{
		// key: values[0],
		// name: values[1],
	// }
// }

func (fr *FileReader) GetFileName(extension string) string {
	fullPath := ""
	parts := fr.Break()
	nameParts := strings.Split(parts[1], "\\")
	
	if nameParts[0] == "engine" || nameParts[0] == "app" {
		// fullPath = filepath.Join(ResourcePath, nameParts[0], fr.filesParser.FileFolder, nameParts[1]+fr.filesParser.Loader.Extension)
		fullPath = filepath.Join(ResourcePath, nameParts[0], fr.filesParser.FileFolder, nameParts[1])
	} else {
		// fullPath = filepath.Join(ResourcePath, fr.filesParser.FileFolder, parts[1]+fr.filesParser.Loader.Extension)
		fullPath = filepath.Join(ResourcePath, fr.filesParser.FileFolder, parts[1])
	}
	
	ext := filepath.Ext(fullPath)
	// fullPath = fullPath[0:len(fullPath)-len(ext)] + fr.filesParser.Loader.Extension
	
	// return fullPath[0:len(fullPath)-len(ext)] + fr.filesParser.Loader.Extension
	return fullPath[0:len(fullPath)-len(ext)] + extension
}

func (fr *FileReader) Close() {
	
}

func (fr *FileReader) AddParameter(name string, value interface{}, overwrite bool) {
	if fr.fileData != nil {
		fr.fileData.AddParameter(name, value, overwrite)
	}
}

func (fr *FileReader) AddParameterEx(values []string, overwrite bool) {
	if fr.fileData != nil {
		fr.fileData.AddParameterEx(values, overwrite)
	}
}

func NewFileReader(filesParser BundlerCss, fileName string, fileData *FileData, variables utils.IKeyValue) FileReader {
	fr := FileReader{
		filesParser: filesParser,
		fileName: fileName,
		fileData: fileData,
		Variables: variables,
	}

	fr.Init()
	return fr
}
