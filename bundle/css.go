package bundle

import (
	"strings"
    "net/http"
	"ibsi/utils"
)

const _IF = "##if"
const _ELSE = "##else"
const _ELSEIF = "##elseif"
const _ENDIF = "##endif"
const _CONST = "##const"
const _ENDCONST = "##endconst"
const _TEMPLATE = "##template"
const _ENDTEMPLATE = "##endtemplate"
const _INCLUDETEMPLATE = "##includetemplate"
const _ENDINCLUDE = "##endinclude"

type ParserFunc func(*FileReader)

type BundlerCss struct {
	Loader *CssLoader
	Variables utils.IKeyValue
	FileDatas map[int]FileData
	Actions map[string]ParserFunc
	FileFolder string
	Comment string
	// Comments []string

	filesCount int
}

func (b *BundlerCss) Init() {
	b.filesCount = 0
	b.Variables = utils.NewKeyValue()
	b.FileDatas = make(map[int]FileData)
	b.FileFolder = "css"
	b.Comment = "//,/*"
	// b.Comments = []string{"//", "/*"}

	b.Actions = make(map[string]ParserFunc)

	b.Actions[_IF] = b.Parse_IF;
	b.Actions[_CONST] = b.Parse_CONST;
	b.Actions[_TEMPLATE] = b.Parse_INCLUDE;
	b.Actions[_INCLUDETEMPLATE] = b.Parse_TEMPLATE;
}

func (b *BundlerCss) IsComment(comment string) bool {
	// for _, a := range b.Comments {
		// if a == comment {
			// return true
		// }
	// }

	// return false
	return strings.Index(b.Comment, comment) > -1
}

func (b BundlerCss) Compare(name string, value string, negate bool) bool {
	if b.Variables.ContainsKey(name) {
		if !negate {
			return b.Variables.AsString(name) == value
		} else {
			return b.Variables.AsString(name) != value
		}
	}
	
	return negate
}

func (b *BundlerCss) Parse_IF(reader *FileReader) {
	i := strings.Index(reader.Line(), "(") + 1
	j := strings.Index(reader.Line(), ")")
	expr := strings.Split(reader.Line()[i:j], "=")
	
	reader.Next()
	for !reader.LineTokenIs(_ENDIF) && !reader.LineTokenIs(_ELSE) && !reader.LineTokenIs(_ELSEIF) {
		if !reader.IsCommentLine() {
			if (b.Compare(expr[0], expr[1], false)) {
				reader.Container.WriteString(reader.Line()+"\n");
			}
		}
		reader.Next()
	}

	if reader.LineTokenIs(_ELSEIF) {
		b.Parse_IF(reader)
	}

	if reader.LineTokenIs(_ELSE) {
		reader.Next()
		for !reader.LineTokenIs(_ENDIF) {
			if !reader.IsCommentLine() {
				if (b.Compare(expr[0], expr[1], true)) {
					reader.Container.WriteString(reader.Line()+"\n");
				}
			}

			reader.Next()
		}
	}
}

func (b *BundlerCss) Parse_CONST(reader *FileReader) {
	reader.Next();
	for !reader.LineTokenIs(_ENDCONST) {
		if !reader.IsCommentLine() {
			reader.AddParameterEx(reader.Break(), false)
		}
		reader.Next()
	}
}

func (b *BundlerCss) Parse_INCLUDE(reader *FileReader) {
	// b.Loader.Output(reader.Line()+"\n")
	log := NewFileData("")

	if (reader.Variables != nil) {
		log.Keys.Append(reader.Variables)
	}

	log.Keys.Append(b.Variables)

	reader.Next();
	for !reader.LineTokenIs(_ENDTEMPLATE) {
		if !reader.IsCommentLine() {
			if values := reader.Break(); !reader.IsCommentLine() {
				if !log.Keys.ContainsKey(values[0]) {
					log.Keys.Set(values[0], values[1])
					// b.Loader.Output(values[0] +"="+ values[1] + "\n") // debug
				}
			}
		}
		reader.Next()
	}

	if reader.LineTokenIs(_ENDTEMPLATE) {
		// b.Loader.Output(reader.GetFileName(b.Loader.Extension)+"\n") // debug
		log.Keys.Set("file", reader.GetFileName(b.Loader.Extension))
		b.Log(log)
	}
}

func (b *BundlerCss) Parse_TEMPLATE(reader *FileReader) {
	variables := utils.NewKeyValue()

	reader.Next();
	for !reader.LineTokenIs(_ENDINCLUDE) {
		if values := reader.Break(); !reader.IsCommentLine() {
			if !variables.ContainsKey(values[0]) {
				variables.Set(values[0], values[1])
				// b.Loader.Output(values[0] +"="+ values[1] + "\n") // debug
			}
		}
		reader.Next()
	}

	if reader.LineTokenIs(_ENDINCLUDE) {
		// b.Loader.Output(reader.GetFileName(b.Loader.TemplateExtension)+"\n") // debug
		b.ParseFile(reader.GetFileName(b.Loader.TemplateExtension), variables);
	}
}

func (b *BundlerCss) Output(line string) {
	b.Loader.Output(line)
}

func (b *BundlerCss) Log(log FileData) {
	// if (log.ContainsKey("file")) {
		// UpdateModifiedDate(log.AsString["file"]);
	// }

	b.filesCount += 1
	b.FileDatas[b.filesCount] = log
}

func (b BundlerCss) ParseFile(fileName string, tkv utils.IKeyValue) {
	if tkv != nil {
		tkv.Append(b.Variables)
	}

	reader := NewFileReader(b, fileName, nil, tkv)

	for !reader.Eof() {
		if !reader.IsCommentLine() {
			// b.Loader.Output(reader.Line()+"\n")
			if reader.LineTokenIs("##") {
				// b.Loader.Output("roken ## found...")
				token := ""

				// if (reader.LineTokenIs("##if", false)) {
					// token = "##if";
				// } else if (reader.LineTokenIs("##elseif", false)) {
					// token = "##elseif";
				// } else
					token = reader.Break()[0];

				if _, ok := b.Actions[token]; ok {
					// b.Loader.Output(reader.Line()+"\n")
					b.Actions[token](&reader)
				// if (actions.ContainsKey(token)) {
					// actions[token].Invoke(reader);
				// } else {
					// actions["default"].Invoke(reader);
				}
			}
		}

		reader.Next()
	}

	reader.Close()
}

func (b BundlerCss) ProcessFile(r *http.Request, f FileData) {
	fn := f.Keys.AsString("file")
	reader := NewFileReader(b, fn, &f, nil)

	for !reader.Eof() {
		if !reader.IsCommentLine() {
			if reader.LineTokenIs("##") {
				token := ""

				if reader.LineTokenIs("##if") {
					token = "##if"
				// } else if reader.LineTokenIs("##elseif") {
					// token = "##elseif"
				// } else if reader.LineTokenIs("##else") {
					// token = "##else"
				} else {
					token = reader.Break()[0]
				}

				if _, ok := b.Actions[token]; ok {
					// reader.Container.WriteString("TOKEN: "+token+"\n")
					b.Actions[token](&reader)
				}
			} else {
				reader.Container.WriteString(reader.Line()+"\n")
			}
		}

		reader.Next()
	}

	s := reader.Container.String()
	for n := range f.Keys.GetCollection() {
		// kv.Set(n, src.Get(n))
		s = strings.Replace(s, "[" +n+ "]", f.Keys.AsString(n), -1)
	}
	
	b.Output("/* " + fn + " */\n")
	b.Output(s)

	reader.Close()
}

func (b BundlerCss) BundleFiles(r *http.Request) {
	for _, f := range b.FileDatas {
		b.ProcessFile(r, f)
	}
}

func NewBundlerCss(loader *CssLoader) BundlerCss {
	b := BundlerCss{
		Loader: loader,
	}

	b.Init()

	return b
}

func LoadCss(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	// w.Header().Set("Content-Type", "text/css; charset=utf-8")

	loader := NewCssLoader(w, r)
	loader.Load()
}
