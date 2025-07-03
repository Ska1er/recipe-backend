package internal

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"path/filepath"
	"time"
)

const StaticFolder = "./public"

func GenerateUniqueFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	uuid := make([]byte, 16)
	rand.Read(uuid)
	return fmt.Sprintf("%x-%d%s", uuid, time.Now().UnixNano(), ext)
}

func ToNullString(str *string) sql.NullString {
	if str == nil {
		return sql.NullString{
			String: "",
			Valid:  false,
		}
	}

	return sql.NullString{
		String: *str,
		Valid:  true,
	}

}
