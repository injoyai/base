package oss

import (
	"encoding/hex"
	"strings"
)

var (
	MapFileType = map[string]string{
		"4D546864":                     "mid",
		"FFD8FF":                       "jpg",
		"89504E47":                     "png",
		"47494638":                     "gif",
		"49492A00":                     "tif",
		"424D":                         "bmp",
		"41433130":                     "dwg",
		"38425053":                     "psd",
		"7B5C727466":                   "rtf",
		"3C3F786D6C":                   "xml",
		"68746D6C3E":                   "html",
		"44656C69766572792D646174653A": "eml",
		"CFAD12FEC5FD746F":             "dbx",
		"2142444E":                     "pst",
		"D0CF11E0":                     "xls,doc",
		"5374616E64617264204A":         "mdb",
		"FF575043":                     "wpd",
		"252150532D41646F6265":         "eps,ps",
		"255044462D312E":               "pdf",
		"AC9EBD8F":                     "qdf",
		"E3828596":                     "pwl",
		"504B0304":                     "zip",
		"52617221":                     "rar",
		"57415645":                     "wav",
		"52494646":                     "avi",
		"2E7261FD":                     "ram",
		"2E524D46":                     "rm",
		"000001BA":                     "mpg",
		"000001B3":                     "mpg",
		"6D6F6F76":                     "mov",
		"3026B2758E66CF11":             "asf",
		"7F454C46":                     "elf", //可执行文件(linux)
		"EDABEEDB":                     "rpm",
	}
)

// FileType 文件类型
func FileType(bs []byte) string {
	if len(bs) > 32 {
		bs = bs[:32]
	}
	lower := hex.EncodeToString(bs)
	upper := strings.ToUpper(lower)
	for i, v := range MapFileType {
		if len(i) > 0 && strings.Index(upper, i) == 0 {
			return v
		}
	}
	return "unknown"
}
