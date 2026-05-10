package security

import (
	"crypto/rand"
	"fmt"
	"testing"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

// =============================================================================
// HELPERS
// =============================================================================

func randomPassword() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func randomSalt(n uint32) []byte {
	s := make([]byte, n)
	_, _ = rand.Read(s)
	return s
}

// =============================================================================
// BLOQUE 1 — Argon2id: comparación de perfiles de parámetros
//
// Objetivo: demostrar cómo varían memoria, tiempo e iteraciones e impactan
// la latencia. Permite justificar los parámetros elegidos en producción.
// =============================================================================

var argon2Profiles = []struct {
	name        string
	memory      uint32 // KiB
	iterations  uint32
	parallelism uint8
	keyLen      uint32
}{
	// Mínimo funcional (NO recomendado en producción)
	{"Min_m=8MB_t=1_p=1", 8 * 1024, 1, 1, 32},

	// OWASP mínimo recomendado para Argon2id (2023)
	{"OWASP_Min_m=19MB_t=2_p=1", 19 * 1024, 2, 1, 32},

	// Configuración del proyecto Alcatraz (equilibrio seguridad/rendimiento)
	{"Alcatraz_Default_m=64MB_t=3_p=2", 64 * 1024, 3, 2, 32},

	// Alta seguridad (servidores con más recursos)
	{"HighSec_m=128MB_t=4_p=4", 128 * 1024, 4, 4, 32},

	// Máximo práctico (para contextos de máxima seguridad)
	{"Max_m=256MB_t=5_p=4", 256 * 1024, 5, 4, 32},
}

// BenchmarkArgon2id_Profiles mide el tiempo de derivación de clave para cada
// perfil de parámetros. Ejecutar con:
//
//	go test -bench=BenchmarkArgon2id_Profiles -benchmem -benchtime=5s
func BenchmarkArgon2id_Profiles(b *testing.B) {
	password := []byte("CorrectHorseBatteryStaple1!")

	for _, p := range argon2Profiles {
		p := p
		b.Run(p.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				salt := randomSalt(16)
				_ = argon2.IDKey(password, salt, p.iterations, p.memory, p.parallelism, p.keyLen)
			}
		})
	}
}

// =============================================================================
// BLOQUE 2 — Argon2id vs Argon2i vs Argon2d
//
// Objetivo: comparar las tres variantes del algoritmo Argon2 en condiciones
// equivalentes. Argon2id es el recomendado (híbrido resistente a GPU y
// side-channel). Argon2d no está disponible directamente en el paquete de Go,
// así que comparamos id vs i.
// =============================================================================

// BenchmarkArgon2_Variants compara Argon2id vs Argon2i con los mismos
// parámetros base. El ganador en seguridad es siempre Argon2id.
//
//	go test -bench=BenchmarkArgon2_Variants -benchmem
func BenchmarkArgon2_Variants(b *testing.B) {
	password := []byte("CorrectHorseBatteryStaple1!")
	salt := randomSalt(16)

	// Parámetros iguales para comparación justa
	memory := uint32(64 * 1024)
	iterations := uint32(3)
	parallelism := uint8(2)
	keyLen := uint32(32)

	b.Run("Argon2id", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = argon2.IDKey(password, salt, iterations, memory, parallelism, keyLen)
		}
	})

	b.Run("Argon2i", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = argon2.Key(password, salt, iterations, memory, parallelism, keyLen)
		}
	})
}

// =============================================================================
// BLOQUE 3 — Argon2id vs bcrypt
//
// Objetivo: demostrar por qué Argon2id supera a bcrypt para gestores de
// contraseñas modernos. Bcrypt no permite configurar uso de memoria, lo que
// lo hace más vulnerable a ataques con hardware especializado (FPGAs/ASICs).
// =============================================================================

var bcryptCosts = []struct {
	name string
	cost int
}{
	{"bcrypt_cost=10", 10}, // Mínimo recomendado actualmente
	{"bcrypt_cost=12", 12}, // Recomendado por OWASP 2023
	{"bcrypt_cost=14", 14}, // Alta seguridad
}

// BenchmarkArgon2id_vs_Bcrypt compara directamente ambos algoritmos.
// Todos los perfiles de argon2id deben superar en seguridad a bcrypt
// manteniendo tiempos aceptables (~100-500ms por operación en servidor).
//
//	go test -bench=BenchmarkArgon2id_vs_Bcrypt -benchmem -benchtime=3s
func BenchmarkArgon2id_vs_Bcrypt(b *testing.B) {
	password := []byte("CorrectHorseBatteryStaple1!")

	// Argon2id con parámetros de Alcatraz
	b.Run("Argon2id_Alcatraz_Default", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			salt := randomSalt(16)
			_ = argon2.IDKey(password, salt, 3, 64*1024, 2, 32)
		}
	})

	// bcrypt a distintos costs
	for _, c := range bcryptCosts {
		c := c
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = bcrypt.GenerateFromPassword(password, c.cost)
			}
		})
	}
}

// =============================================================================
// BLOQUE 4 — Impacto del paralelismo en Argon2id
//
// Objetivo: mostrar cómo varía el tiempo de derivación al escalar el número
// de threads. Relevante para servidores multi-core. Un mayor paralelismo no
// siempre implica mayor velocidad si el CPU ya está saturado.
// =============================================================================

// BenchmarkArgon2id_Parallelism fija memoria e iteraciones y varía p=1..8.
//
//	go test -bench=BenchmarkArgon2id_Parallelism -benchmem
func BenchmarkArgon2id_Parallelism(b *testing.B) {
	password := []byte("CorrectHorseBatteryStaple1!")
	salt := randomSalt(16)

	for _, p := range []uint8{1, 2, 4, 8} {
		p := p
		b.Run(fmt.Sprintf("p=%d", p), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = argon2.IDKey(password, salt, 3, 64*1024, p, 32)
			}
		})
	}
}

// =============================================================================
// BLOQUE 5 — Impacto de la longitud de la clave de salida
//
// Objetivo: verificar que la longitud de la clave derivada (keyLen) no es un
// factor determinante en el tiempo, ya que el costo real es el llenado de
// memoria. Útil para justificar usar 32 bytes (256 bits) como output.
// =============================================================================

// BenchmarkArgon2id_KeyLength muestra que variar keyLen tiene impacto mínimo.
//
//	go test -bench=BenchmarkArgon2id_KeyLength -benchmem
func BenchmarkArgon2id_KeyLength(b *testing.B) {
	password := []byte("CorrectHorseBatteryStaple1!")
	salt := randomSalt(16)

	for _, kl := range []uint32{16, 32, 64} {
		kl := kl
		b.Run(fmt.Sprintf("keyLen=%d", kl), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = argon2.IDKey(password, salt, 3, 64*1024, 2, kl)
			}
		})
	}
}

// =============================================================================
// BLOQUE 6 — HashPassword y VerifyPassword del proyecto (full pipeline)
//
// Objetivo: medir el tiempo de las funciones reales de Alcatraz, incluyendo
// generación de sal aleatoria, codificación base64 y parseo del hash.
// Estos son los números reales que el usuario experimenta en login/registro.
// =============================================================================

// BenchmarkHashPassword mide el tiempo completo de HashPassword tal como se
// llama en el handler de registro/cambio de contraseña.
//
//	go test -bench=BenchmarkHashPassword -benchmem
func BenchmarkHashPassword(b *testing.B) {
	password := "CorrectHorseBatteryStaple1!"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HashPassword(password)
	}
}

// BenchmarkVerifyPassword mide el tiempo de VerifyPassword tal como se llama
// en el handler de login. El tiempo debe ser idéntico al de HashPassword
// ya que el trabajo computacional es el mismo.
//
//	go test -bench=BenchmarkVerifyPassword -benchmem
func BenchmarkVerifyPassword(b *testing.B) {
	password := "CorrectHorseBatteryStaple1!"
	encoded, _ := HashPassword(password)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = VerifyPassword(password, encoded)
	}
}

// BenchmarkVerifyPassword_WrongPass mide si el tiempo de verificación es
// constante también para contraseñas incorrectas (resistencia a timing attacks).
//
//	go test -bench=BenchmarkVerifyPassword_WrongPass -benchmem
func BenchmarkVerifyPassword_WrongPass(b *testing.B) {
	encoded, _ := HashPassword("CorrectHorseBatteryStaple1!")
	wrong := "TotallyWrongPassword!"
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = VerifyPassword(wrong, encoded)
	}
}
