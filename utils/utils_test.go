package utils_test

import (
	"testing"

	"github.com/Hanyue-s-FYP/Marcom-Backend/utils"
)

func TestTernary(t *testing.T) {
    trueInt := utils.If(true, 5, 10)
    if trueInt != 5 {
        t.Errorf("expected 5, got %d", trueInt)
    }
    falseInt := utils.If(false, 5, 10)
    if falseInt != 10 {
        t.Errorf("expected 10, got %d", falseInt)
    }
}
