package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/errwrap/errwrap"
	"github.com/sonatard/noctx"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpmux"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/slog"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stdversion"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/staticcheck"

	osexit "github.com/sodiqit/gophermart/cmd/staticlint/analysis"
)

const ConfigFileName = `config.json`

type Config struct {
	Staticcheck    []string
	Analysispasses []string
	Custom         []string
}

func main() {
	appfile, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(filepath.Dir(appfile), ConfigFileName))
	if err != nil {
		log.Fatal(err)
	}
	var cfg Config
	if err = json.Unmarshal(data, &cfg); err != nil {
		log.Fatal(err)
	}

	var mychecks []*analysis.Analyzer
	checks := make(map[string]bool)
	for _, v := range cfg.Staticcheck {
		checks[v] = true
	}

	for _, v := range staticcheck.Analyzers {
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	for _, v := range cfg.Analysispasses {
		switch v {
		case "appends":
			mychecks = append(mychecks, appends.Analyzer)
		case "asmdecl":
			mychecks = append(mychecks, asmdecl.Analyzer)
		case "assign":
			mychecks = append(mychecks, assign.Analyzer)
		case "atomic":
			mychecks = append(mychecks, atomic.Analyzer)
		case "atomicalign":
			mychecks = append(mychecks, atomicalign.Analyzer)
		case "bools":
			mychecks = append(mychecks, bools.Analyzer)
		case "buildssa":
			mychecks = append(mychecks, buildssa.Analyzer)
		case "buildtag":
			mychecks = append(mychecks, buildtag.Analyzer)
		case "cgocall":
			mychecks = append(mychecks, cgocall.Analyzer)
		case "composite":
			mychecks = append(mychecks, composite.Analyzer)
		case "copylock":
			mychecks = append(mychecks, copylock.Analyzer)
		case "ctrlflow":
			mychecks = append(mychecks, ctrlflow.Analyzer)
		case "deepequalerrors":
			mychecks = append(mychecks, deepequalerrors.Analyzer)
		case "defers":
			mychecks = append(mychecks, defers.Analyzer)
		case "directive":
			mychecks = append(mychecks, directive.Analyzer)
		case "errorsas":
			mychecks = append(mychecks, errorsas.Analyzer)
		case "fieldalignment":
			mychecks = append(mychecks, fieldalignment.Analyzer)
		case "findcall":
			mychecks = append(mychecks, findcall.Analyzer)
		case "framepointer":
			mychecks = append(mychecks, framepointer.Analyzer)
		case "httpmux":
			mychecks = append(mychecks, httpmux.Analyzer)
		case "httpresponse":
			mychecks = append(mychecks, httpresponse.Analyzer)
		case "ifaceassert":
			mychecks = append(mychecks, ifaceassert.Analyzer)
		case "inspect":
			mychecks = append(mychecks, inspect.Analyzer)
		case "loopclosure":
			mychecks = append(mychecks, loopclosure.Analyzer)
		case "lostcancel":
			mychecks = append(mychecks, lostcancel.Analyzer)
		case "nilfunc":
			mychecks = append(mychecks, nilfunc.Analyzer)
		case "nilness":
			mychecks = append(mychecks, nilness.Analyzer)
		case "pkgfact":
			mychecks = append(mychecks, pkgfact.Analyzer)
		case "printf":
			mychecks = append(mychecks, printf.Analyzer)
		case "reflectvaluecompare":
			mychecks = append(mychecks, reflectvaluecompare.Analyzer)
		case "shadow":
			mychecks = append(mychecks, shadow.Analyzer)
		case "shift":
			mychecks = append(mychecks, shift.Analyzer)
		case "sigchanyzer":
			mychecks = append(mychecks, sigchanyzer.Analyzer)
		case "slog":
			mychecks = append(mychecks, slog.Analyzer)
		case "sortslice":
			mychecks = append(mychecks, sortslice.Analyzer)
		case "stdmethods":
			mychecks = append(mychecks, stdmethods.Analyzer)
		case "stdversion":
			mychecks = append(mychecks, stdversion.Analyzer)
		case "stringintconv":
			mychecks = append(mychecks, stringintconv.Analyzer)
		case "structtag":
			mychecks = append(mychecks, structtag.Analyzer)
		case "testinggoroutine":
			mychecks = append(mychecks, testinggoroutine.Analyzer)
		case "tests":
			mychecks = append(mychecks, tests.Analyzer)
		case "timeformat":
			mychecks = append(mychecks, timeformat.Analyzer)
		case "unmarshal":
			mychecks = append(mychecks, unmarshal.Analyzer)
		case "unreachable":
			mychecks = append(mychecks, unreachable.Analyzer)
		case "unsafeptr":
			mychecks = append(mychecks, unsafeptr.Analyzer)
		case "unusedresult":
			mychecks = append(mychecks, unusedresult.Analyzer)
		case "unusedwrite":
			mychecks = append(mychecks, unusedwrite.Analyzer)
		case "usesgenerics":
			mychecks = append(mychecks, usesgenerics.Analyzer)
		}
	}
	// Добавляем публичные анализаторы, включачая свой кастомный osexit.
	for _, v := range cfg.Custom {
		switch v {
		case "errwrap":
			mychecks = append(mychecks, errwrap.Analyzer)
		case "noctx":
			mychecks = append(mychecks, noctx.Analyzer)
		case "osexit":
			mychecks = append(mychecks, osexit.ExitCheckAnalyzer)
		}
	}

	multichecker.Main(
		mychecks...,
	)
}
