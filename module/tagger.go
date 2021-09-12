package module

import (
	"go/parser"
	"go/printer"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/fatih/structtag"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type mod struct {
	*pgs.ModuleBase
	pgsgo.Context
}

func New() pgs.Module {
	return &mod{ModuleBase: &pgs.ModuleBase{}}
}

func (m *mod) InitContext(c pgs.BuildContext) {
	m.ModuleBase.InitContext(c)
	m.Context = pgsgo.InitContext(c.Parameters())
}

func (mod) Name() string {
	return "gotag"
}

func (m mod) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {

	xtv := m.Parameters().Str("xxx")

	xtv = strings.Replace(xtv, "+", ":", -1)

	xt, err := structtag.Parse(xtv)
	m.CheckErr(err)

	autoTag := m.Parameters().Str("auto")
	var autoTags []string
	if autoTag != "" {
		autoTags = strings.Split(autoTag, "+")
	}

	extractor := newTagExtractor(m, m.Context, autoTags)

	for _, f := range targets {
		tags := extractor.Extract(f)

		tags.AddTagsToXXXFields(xt)

		gfname := m.Context.OutputPath(f).SetExt(".go").String()
		module := m.Parameters().Str("module")
		outdir := m.Parameters().Str("outdir")
		outputFilename := gfname
		inputFilename := gfname

		if module != "" {
			outputFilename = strings.TrimPrefix(outputFilename, module+"/")
			inputFilename = strings.TrimPrefix(inputFilename, module+"/")
		}

		if outdir != "" {
			inputFilename = filepath.Join(outdir, inputFilename)
		} // else {
		// 	filename = "." + filename
		// }

		//panic("outdir=" + outdir + "\nmodule = " + module + "\ngfname = " + gfname + "\nFilename: " + filename)

		fs := token.NewFileSet()
		fn, err := parser.ParseFile(fs, inputFilename, nil, parser.ParseComments)
		m.CheckErr(err)
		m.CheckErr(Retag(fn, tags))

		var buf strings.Builder
		m.CheckErr(printer.Fprint(&buf, fs, fn))
		m.OverwriteGeneratorFile(outputFilename, buf.String())
	}

	return m.Artifacts()
}
