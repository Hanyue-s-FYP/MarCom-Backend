package utils

func ResponseJSON() {

}

// Solely for QoL development (like adding ternary to Go :))
// Ternary util
func If[T any](cond bool, valTrue, valFalse T) T {
    if cond {
        return valTrue
    } else {
        return valFalse
    }
}
