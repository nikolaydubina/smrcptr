package main

import (
	"github.com/nikolaydubina/smrcptr/analysis/smrcptr"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(smrcptr.Analyzer) }
