package auth

import "context"

type contextKey string

const keyWorkspaceID contextKey = "workspace_id"
const keyProjectID contextKey = "project_id"
const keyKeyID contextKey = "key_id"

func WithWorkspaceID(ctx context.Context, workspaceID int) context.Context {
	return context.WithValue(ctx, keyWorkspaceID, workspaceID)
}

func WorkspaceIDFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(keyWorkspaceID).(int)

	return id, ok
}

func WithProjectID(ctx context.Context, projectID int) context.Context {
	return context.WithValue(ctx, keyProjectID, projectID)
}

func ProjectIDFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(keyProjectID).(int)

	return id, ok
}

func WithKeyID(ctx context.Context, keyID int) context.Context {
	return context.WithValue(ctx, keyKeyID, keyID)
}

func KeyIDFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(keyKeyID).(int)

	return id, ok
}
