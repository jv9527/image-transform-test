run-go:
	go run go/main.go

bench-resize-gmagick:
	go test ./go/graphicsmagick/... -bench=. -benchmem

bench-resize-vips:
	go test ./go/vips/... -bench=. -benchmem