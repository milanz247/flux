// Package dto defines the data transfer objects exchanged with the Vue
// frontend. GORM models are never serialized directly — every page gets a
// DTO describing exactly the shape of its props, and every write endpoint
// gets a DTO describing (and validating) its input.
package dto
