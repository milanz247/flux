// TypeScript mirrors of the Go DTOs (app/dto). Keep these in sync — every
// req.View() call serializes one of the Go DTOs into these shapes.

// Note: UserDTO is still used in shared auth props but since it comes from
// framework.AuthUser (flux/types.ts), there is no standalone DTO type needed
// on the frontend. The shared AuthUser type in flux/types.ts covers it.
