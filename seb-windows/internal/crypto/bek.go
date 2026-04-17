package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"sync"
)

// bufPool mantiene buffers de 32KB preasignados en RAM.
// CRÍTICO: Previene el trabajo del Garbage Collector durante operaciones I/O intensivas.
var bufPool = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 32*1024)
		return &buf
	},
}

// GenerateBEK calcula el SHA-256 del binario en ejecución usando streaming.
func GenerateBEK() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	f, err := os.Open(exePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hash := sha256.New()
	
	// Obtenemos el puntero al array de bytes preasignado
	bufPtr := bufPool.Get().(*[]byte)
	defer bufPool.Put(bufPtr)

	buf := *bufPtr

	// Bucle Zero-Alloc
	for {
		n, err := f.Read(buf)
		if n > 0 {
			hash.Write(buf[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
