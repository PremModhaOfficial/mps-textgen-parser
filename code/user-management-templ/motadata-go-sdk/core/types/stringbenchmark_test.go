package types

import (
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"
	"testing"
)

// Benchmark tests for String functions

func BenchmarkIsEmpty(b *testing.B) {
	testStr := "hello world"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		IsEmpty(testStr)
	}
}

func BenchmarkIsNotEmpty(b *testing.B) {
	testStr := "hello world"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		IsNotEmpty(testStr)
	}
}

func BenchmarkIsBlank(b *testing.B) {
	testStr := "   hello world   "

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		IsBlank(testStr)
	}
}

func BenchmarkIsNotBlank(b *testing.B) {
	testStr := "   hello world   "

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		IsNotBlank(testStr)
	}
}

func BenchmarkToUpper(b *testing.B) {
	testStr := "hello world"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ToUpper(testStr)
	}
}

func BenchmarkToLower(b *testing.B) {
	testStr := "HELLO WORLD"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ToLower(testStr)
	}
}

func BenchmarkTrim(b *testing.B) {
	testStr := "   hello world   "
	cutset := " "

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Trim(testStr, cutset)
	}
}

func BenchmarkTrimSpace(b *testing.B) {
	testStr := "   hello world   "

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		TrimSpace(testStr)
	}
}

func BenchmarkTrimLeft(b *testing.B) {
	testStr := "   hello world"
	cutset := " "

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		TrimLeft(testStr, cutset)
	}
}

func BenchmarkTrimRight(b *testing.B) {
	testStr := "hello world   "
	cutset := " "

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		TrimRight(testStr, cutset)
	}
}

func BenchmarkNormalizeSpace(b *testing.B) {
	testStr := "  hello   world  "

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		NormalizeSpace(testStr)
	}
}

func BenchmarkCollapseWhitespace(b *testing.B) {
	testStr := "hello\n\t world\r\n  test"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		CollapseWhitespace(testStr)
	}
}

func BenchmarkRemoveWhitespace(b *testing.B) {
	testStr := "hello world test"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		RemoveWhitespace(testStr)
	}
}

func BenchmarkToTitle(b *testing.B) {
	testStr := "hello world"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ToTitle(testStr)
	}
}

func BenchmarkHash32(b *testing.B) {
	testStr := "hello world"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Hash32(testStr)
	}
}

func BenchmarkHash64(b *testing.B) {
	testStr := "hello world"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Hash64(testStr)
	}
}

func BenchmarkChecksum(b *testing.B) {
	testStr := "hello world"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Checksum(testStr)
	}
}

func BenchmarkBase64Encode(b *testing.B) {
	testStr := "hello world"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Base64Encode(testStr)
	}
}

func BenchmarkBase64Decode(b *testing.B) {
	encoded := Base64Encode("hello world")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = Base64Decode(encoded)
	}
}

func BenchmarkIsEmail(b *testing.B) {
	email := "test@example.com"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		IsEmail(email)
	}
}

func BenchmarkIsURL(b *testing.B) {
	url := "https://example.com/path"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		IsURL(url)
	}
}

func BenchmarkIsNumeric(b *testing.B) {
	numeric := "12345"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		IsNumeric(numeric)
	}
}

func BenchmarkIsAlpha(b *testing.B) {
	alpha := "hello"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		IsAlpha(alpha)
	}
}

func BenchmarkIsAlphaNumeric(b *testing.B) {
	alphaNum := "hello123"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		IsAlphaNumeric(alphaNum)
	}
}

func BenchmarkAESEncrypt(b *testing.B) {
	plaintext := "hello world"
	key := "1234567890123456"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = AESEncrypt(plaintext, key)
	}
}

func BenchmarkAESDecrypt(b *testing.B) {
	plaintext := "hello world"
	key := "1234567890123456"
	encrypted, _ := AESEncrypt(plaintext, key)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = AESDecrypt(encrypted, key)
	}
}

func BenchmarkHashPassword(b *testing.B) {
	password := "mySecurePassword123"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = HashPassword(password)
	}
}

func BenchmarkVerifyPassword(b *testing.B) {
	password := "mySecurePassword123"
	hash, _ := HashPassword(password)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		VerifyPassword(hash, password)
	}
}

func BenchmarkToBytes(b *testing.B) {
	testStr := "hello world"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		utils.StringToBytes(testStr)
	}
}

func BenchmarkURLEncode(b *testing.B) {
	testStr := "hello world test"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		URLEncode(testStr)
	}
}

func BenchmarkURLDecode(b *testing.B) {
	encoded := URLEncode("hello world test")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = URLDecode(encoded)
	}
}

func BenchmarkIsASCII(b *testing.B) {
	testStr := "hello world 123"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		IsASCII(testStr)
	}
}

func BenchmarkIsLower(b *testing.B) {
	testStr := "hello world"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		IsLower(testStr)
	}
}

func BenchmarkIsUpper(b *testing.B) {
	testStr := "HELLO WORLD"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		IsUpper(testStr)
	}
}

func BenchmarkIsDigit(b *testing.B) {
	testStr := "123456"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		IsDigit(testStr)
	}
}

func BenchmarkIsHex(b *testing.B) {
	testStr := "deadbeef"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		IsHex(testStr)
	}
}
