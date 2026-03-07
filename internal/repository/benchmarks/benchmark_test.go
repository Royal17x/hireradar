package benchmarks

import (
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
	"testing"
)

func BenchmarkExistsPostgres(b *testing.B) {
	b.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vacancyRepo.Exists(ctx, "test-hh-id")
	}
}

func BenchmarkIsSeenRedis(b *testing.B) {
	b.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cacheRepo.IsSeen(ctx, "test-hh-id")
	}
}
