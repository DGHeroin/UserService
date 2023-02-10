package utils

import (
    "math"
    "strings"
)

const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Base62Encode(num int64) string {
    bytes := []byte{}
    for num > 0 {
        bytes = append(bytes, chars[num%62])
        num = num / 62
    }
    _Base62reverse(bytes)
    return string(bytes)
}

func Base62Decode(str string) int64 {
    var num int64
    n := len(str)
    for i := 0; i < n; i++ {
        num += int64(math.Pow(62, float64(n-i-1))) * int64(strings.IndexByte(chars, str[i]))
    }
    return num
}

func _Base62reverse(a []byte) {
    for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
        a[left], a[right] = a[right], a[left]
    }
}
