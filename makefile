MAIN_PACKAGE_PATH := '.'
BINARY_NAME := Gigeii
 
 
 
build:
	# Include additional build steps, like TypeScript, SCSS or Tailwind compilation here...
	sqlc generate && templ generate && go build -o=/tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH} && golangci-lint run

run: build
	/tmp/bin/${BINARY_NAME}

run/live:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build" \
		--build.bin "/tmp/bin/${BINARY_NAME}" \
		--build.delay "100" \
		--build.exclude_dir "scripts,backups" \
		--build.exclude_regex ".*_templ.go|db/batch\.go|db/copyfrom\.go|db/db\.go|db/models\.go|.*\.sql\.go" \
		--build.include_ext "go, tpl, tmpl, templ, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
		--misc.clean_on_exit "true"