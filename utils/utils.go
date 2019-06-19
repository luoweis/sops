package utils


import (
	"fmt"
	"os"
	"flag"
	"crypto/sha256"
	"encoding/hex"
)

func FlagUsage() {
	fmt.Fprintf(os.Stderr,`sops version:1.0.0
Usage:./sops [-p port] [-r role]

Options:
`)
	flag.PrintDefaults()
}


// @Description： 计算哈希值
func CalculateHash(str string) string {
	hashInBytes := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hashInBytes[:])
}
