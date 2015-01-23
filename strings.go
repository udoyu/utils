package utils

import "strings"

func StringIndexAny(s, cutset string, n int) int {
    r := 0
    for n > 0 {
        n--
        t := strings.Index(s[r:], cutset)
        if t == -1 {
            return -1
        }
        r += t + 1
    }
    return r-1
}

func StringCutRight(s string, offset int) string {
    if offset < 0 || offset > len(s) {
        return s
    }
    return s[:offset]
}

func StringCutLeft(s string, offset int) string {
     if offset < 0 || offset > len(s) {
        return ""
    }
    return s[offset:]
}

func StringCutRightExp(s string, cutset string, n int) string {
    return StringCutRight(s, StringIndexAny(s, cutset, n))
}

func StringCutLeftExp(s string, cutset string, n int) string {
    return StringCutLeft(s, StringIndexAny(s, cutset, n)+1)
}
