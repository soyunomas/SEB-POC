package crypto

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"howett.net/plist"
)

// buildSEBJSON replica json_encode de PHP con JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE.
// Serialización recursiva que respeta el formato canónico que Moodle espera.
func buildSEBJSON(v interface{}) string {
	switch val := v.(type) {
	case map[string]interface{}:
		if len(val) == 0 {
			return "{}"
		}
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		// CRÍTICO: Moodle usa Collator('root')->asort() (case-insensitive).
		// sort.Strings() es case-sensitive y produce un hash diferente.
		sort.Slice(keys, func(i, j int) bool {
			li, lj := strings.ToLower(keys[i]), strings.ToLower(keys[j])
			if li != lj {
				return li < lj
			}
			return keys[i] < keys[j]
		})
		var parts []string
		for _, k := range keys {
			child := val[k]
			// Moodle/SEB omiten diccionarios vacíos del hash
			if d, ok := child.(map[string]interface{}); ok && len(d) == 0 {
				continue
			}
			parts = append(parts, fmt.Sprintf("%q:%s", k, buildSEBJSON(child)))
		}
		return "{" + strings.Join(parts, ",") + "}"
	case []interface{}:
		if len(val) == 0 {
			return "[]"
		}
		var parts []string
		for _, item := range val {
			parts = append(parts, buildSEBJSON(item))
		}
		return "[" + strings.Join(parts, ",") + "]"
	case bool:
		if val {
			return "true"
		}
		return "false"
	case uint64:
		return fmt.Sprintf("%d", val)
	case int64:
		return fmt.Sprintf("%d", val)
	case int:
		return fmt.Sprintf("%d", val)
	case float64:
		return fmt.Sprintf("%v", val)
	case string:
		b, _ := json.Marshal(val)
		s := string(b)
		// Go escapa HTML (&, <, >) pero PHP no (JSON_UNESCAPED_UNICODE)
		s = strings.ReplaceAll(s, "\\u003c", "<")
		s = strings.ReplaceAll(s, "\\u003e", ">")
		s = strings.ReplaceAll(s, "\\u0026", "&")
		return s
	case []byte:
		// plist <data> → base64 string
		return fmt.Sprintf("%q", base64.StdEncoding.EncodeToString(val))
	case time.Time:
		// plist <date> → ISO 8601 UTC
		return fmt.Sprintf("%q", val.UTC().Format("2006-01-02T15:04:05Z"))
	default:
		return "null"
	}
}

// DeriveConfigKey implementa la matemática de Moodle (generate_config_key)
func DeriveConfigKey(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var dict map[string]interface{}
	if _, err := plist.Unmarshal(data, &dict); err != nil {
		return "", err
	}

	// Moodle solo descarta originatorVersion antes de hashear
	delete(dict, "originatorVersion")

	jsonStr := buildSEBJSON(dict)

	fmt.Printf("\n[FORENSE] JSON Canonicalizado a inyectar:\n%s\n\n", jsonStr)

	hash := sha256.Sum256([]byte(jsonStr))
	return hex.EncodeToString(hash[:]), nil
}
