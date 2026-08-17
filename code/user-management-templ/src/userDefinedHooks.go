// userDefinedHooks.go
// AUTO-GENERATED STUBS — safe to edit.
// New hooks are appended automatically. Your implementations are preserved.
package main

import (
	"context"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/tracer"
)

func (s *RolesHandler) preCreateNotifyAdmin(ctx context.Context, span tracer.Span, event *RolesCreatedEvent) {
ctx, hookSpan := tracer.Start(ctx, "hook.pre.Roles.create.NotifyAdmin")
defer hookSpan.End()
	// TODO: implement
hookSpan.SetOK()
}

func (s *RolesPermissionsHandler) preAssignHook(ctx context.Context, span tracer.Span, event *RolesPermissionsAssignedEvent) error {
 return nil
}

func (s *RolesPermissionsHandler) postAssignHook(ctx context.Context, span tracer.Span, event *RolesPermissionsAssignedEvent, data []byte) []byte {
 return data
}

func (s *RolesPermissionsHandler) preListHook(ctx context.Context, span tracer.Span, event *RolesPermissionsListRequest) error {
 return nil
}

func (s *RolesPermissionsHandler) postListHook(ctx context.Context, span tracer.Span, event *RolesPermissionsListRequest, data []byte) []byte {
 return data
}

func (s *RolesPermissionsHandler) preRemoveHook(ctx context.Context, span tracer.Span, event *RolesPermissionsRemovedEvent) error {
 return nil
}

func (s *RolesPermissionsHandler) postRemoveHook(ctx context.Context, span tracer.Span, event *RolesPermissionsRemovedEvent, data []byte) []byte {
 return data
}

func (s *PermissionsHandler) preCreateNotifyManager(ctx context.Context, span tracer.Span, event *PermissionsCreatedEvent) error {
ctx, hookSpan := tracer.Start(ctx, "hook.pre.Permissions.create.NotifyManager")
defer hookSpan.End()
	// TODO: implement
hookSpan.SetOK()
   return nil
}

func (s *PermissionsHandler) preCreateNotifyAdmin(ctx context.Context, span tracer.Span, event *PermissionsCreatedEvent) {
ctx, hookSpan := tracer.Start(ctx, "hook.pre.Permissions.create.NotifyAdmin")
defer hookSpan.End()
	// TODO: implement
hookSpan.SetOK()
}

func (s *PermissionsHandler) postCreateAudit(ctx context.Context, span tracer.Span, event *PermissionsCreatedEvent, data []byte) {
ctx, hookSpan := tracer.Start(ctx, "hook.post.Permissions.create.Audit")
defer hookSpan.End()
	// TODO: implement
hookSpan.SetOK()
}

func (s *UserHandler) preCreateD(ctx context.Context, span tracer.Span, event *UserCreatedEvent) error {
ctx, hookSpan := tracer.Start(ctx, "hook.pre.User.create.D")
defer hookSpan.End()
	// TODO: implement
hookSpan.SetOK()
   return nil
}

func (s *UserHandler) preCreateC(ctx context.Context, span tracer.Span, event *UserCreatedEvent) {
ctx, hookSpan := tracer.Start(ctx, "hook.pre.User.create.C")
defer hookSpan.End()
	// TODO: implement
hookSpan.SetOK()
}

func (s *UserHandler) preCreateB(ctx context.Context, span tracer.Span, event *UserCreatedEvent) error {
ctx, hookSpan := tracer.Start(ctx, "hook.pre.User.create.B")
defer hookSpan.End()
	// TODO: implement
hookSpan.SetOK()
   return nil
}

func (s *UserHandler) preCreateM(ctx context.Context, span tracer.Span, event *UserCreatedEvent) {
ctx, hookSpan := tracer.Start(ctx, "hook.pre.User.create.M")
defer hookSpan.End()
	// TODO: implement
hookSpan.SetOK()
}

func (s *UserHandler) preCreateA(ctx context.Context, span tracer.Span, event *UserCreatedEvent) {
ctx, hookSpan := tracer.Start(ctx, "hook.pre.User.create.A")
defer hookSpan.End()
	// TODO: implement
hookSpan.SetOK()
}

func (s *UserHandler) preDeleteXyz(ctx context.Context, span tracer.Span, event *UserDeletedEvent) error {
ctx, hookSpan := tracer.Start(ctx, "hook.pre.User.delete.Xyz")
defer hookSpan.End()
	// TODO: implement
hookSpan.SetOK()
   return nil
}

func (s *UserHandler) preDeleteAbc(ctx context.Context, span tracer.Span, event *UserDeletedEvent) error {
ctx, hookSpan := tracer.Start(ctx, "hook.pre.User.delete.Abc")
defer hookSpan.End()
	// TODO: implement
hookSpan.SetOK()
   return nil
}

func (s *UserHandler) preDeleteValidateTen(ctx context.Context, span tracer.Span, event *UserDeletedEvent) error {
ctx, hookSpan := tracer.Start(ctx, "hook.pre.User.delete.ValidateTen")
defer hookSpan.End()
	// TODO: implement
hookSpan.SetOK()
   return nil
}

func (s *UserHandler) preUpdateNotifyCache(ctx context.Context, span tracer.Span, event *UserUpdatedEvent) {
ctx, hookSpan := tracer.Start(ctx, "hook.pre.User.update.NotifyCache")
defer hookSpan.End()
	// TODO: implement
hookSpan.SetOK()
}

func (s *UserHandler) preUpdateAuditLog(ctx context.Context, span tracer.Span, event *UserUpdatedEvent) {
ctx, hookSpan := tracer.Start(ctx, "hook.pre.User.update.AuditLog")
defer hookSpan.End()
	// TODO: implement
hookSpan.SetOK()
}

func (s *UserHandler) preGetTraceAcces(ctx context.Context, span tracer.Span, event *UserGetRequest) {
ctx, hookSpan := tracer.Start(ctx, "hook.pre.User.get.TraceAcces")
defer hookSpan.End()
	// TODO: implement
hookSpan.SetOK()
}

func (s *UserHandler) postUpdateInvalidateCache(ctx context.Context, span tracer.Span, event *UserUpdatedEvent, data []byte) {
ctx, hookSpan := tracer.Start(ctx, "hook.post.User.update.InvalidateCache")
defer hookSpan.End()
	// TODO: implement
hookSpan.SetOK()
}

func (s *UserHandler) postGetFilterPII(ctx context.Context, span tracer.Span, event *UserGetRequest, data []byte) []byte {
ctx, hookSpan := tracer.Start(ctx, "hook.post.User.get.FilterPII")
defer hookSpan.End()
	// TODO: implement
hookSpan.SetOK()
   return data
}

func (s *UserHandler) postGetCacheResult(ctx context.Context, span tracer.Span, event *UserGetRequest, data []byte) {
ctx, hookSpan := tracer.Start(ctx, "hook.post.User.get.CacheResult")
defer hookSpan.End()
	// TODO: implement
hookSpan.SetOK()
}

func (s *UserHandler) postGetEnrichData(ctx context.Context, span tracer.Span, event *UserGetRequest, data []byte) []byte {
ctx, hookSpan := tracer.Start(ctx, "hook.post.User.get.EnrichData")
defer hookSpan.End()
	// TODO: implement
hookSpan.SetOK()
   return data
}

func (s *UserHandler) postCreateNotifyAdmin(ctx context.Context, span tracer.Span, event *UserCreatedEvent, data []byte) {
ctx, hookSpan := tracer.Start(ctx, "hook.post.User.create.NotifyAdmin")
defer hookSpan.End()
	// TODO: implement
hookSpan.SetOK()
}

func (s *UserHandler) postCreateSendWelcome(ctx context.Context, span tracer.Span, event *UserCreatedEvent, data []byte) {
ctx, hookSpan := tracer.Start(ctx, "hook.post.User.create.SendWelcome")
defer hookSpan.End()
	// TODO: implement
hookSpan.SetOK()
}

func (s *UserHandler) postDeleteMustDelete(ctx context.Context, span tracer.Span, event *UserDeletedEvent, data []byte) []byte {
ctx, hookSpan := tracer.Start(ctx, "hook.post.User.delete.MustDelete")
defer hookSpan.End()
	// TODO: implement
hookSpan.SetOK()
   return data
}

func (s *UserRolesHandler) preAssignHook(ctx context.Context, span tracer.Span, event *UserRolesAssignedEvent) error {
 return nil
}

func (s *UserRolesHandler) postAssignHook(ctx context.Context, span tracer.Span, event *UserRolesAssignedEvent, data []byte) []byte {
 return data
}

func (s *UserRolesHandler) preListHook(ctx context.Context, span tracer.Span, event *UserRolesListRequest) error {
 return nil
}

func (s *UserRolesHandler) postListHook(ctx context.Context, span tracer.Span, event *UserRolesListRequest, data []byte) []byte {
 return data
}

func (s *UserRolesHandler) preRemoveHook(ctx context.Context, span tracer.Span, event *UserRolesRemovedEvent) error {
 return nil
}

func (s *UserRolesHandler) postRemoveHook(ctx context.Context, span tracer.Span, event *UserRolesRemovedEvent, data []byte) []byte {
 return data
}
