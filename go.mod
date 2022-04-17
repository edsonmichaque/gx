module github.com/edsonmichaque/omni

go 1.17

replace github.com/edsonmichaque/omni => ./

replace github.com/edsonmichaque/libomni => ../libomni

require (
	github.com/edsonmichaque/libomni v0.0.0-20220417100936-aad246d01148 // indirect
	github.com/google/uuid v1.3.0 // indirect
)
