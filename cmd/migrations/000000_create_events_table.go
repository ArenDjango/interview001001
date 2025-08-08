package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		_, err := db.ExecContext(ctx, `
			CREATE TABLE "events" (
				"id" UUID NOT NULL PRIMARY KEY,
				"title" TEXT NOT NULL,
				"description" TEXT NOT NULL,
				"start_time" TIMESTAMPTZ,
				"end_time" TIMESTAMPTZ,
				"created_at" TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
			)
		`)
		if err != nil {
			return err
		}

		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		_, err := db.ExecContext(ctx, `DROP TABLE IF EXISTS "events"`)
		return err
	})
}
