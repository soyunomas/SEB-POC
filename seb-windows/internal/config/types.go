package config

// SEBConfig contiene los parámetros críticos derivados del archivo .seb.
// PRECEPTO #5: Struct Alignment (String, String, String, Bool) para evitar padding en memoria.
type SEBConfig struct {
	StartURL       string `plist:"startURL"`
	BrowserExamKey string `plist:"browserExamKey,omitempty"`
	ConfigKey      string `plist:"configKey,omitempty"`
	KioskMode      bool   `plist:"kioskMode,omitempty"`
}
