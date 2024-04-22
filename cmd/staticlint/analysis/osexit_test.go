package analysis

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func TestOSExitAnalyzer(t *testing.T) {
	// функция analysistest.Run применяет тестируемый анализатор ExitCheckAnalyzer
	// к пакетам из папки testdata и проверяет ожидания
	// ./... — проверка всех поддиректорий в testdata
	// можно указать ./pkg1 для проверки только pkg1
	analysistest.Run(t, analysistest.TestData(), ExitCheckAnalyzer, "./...")
}
