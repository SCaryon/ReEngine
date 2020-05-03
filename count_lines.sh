find . -path "./static" -prune -o -type f "(" -name "*.go" -or -name "*.html" ")" -print | xargs wc -l
