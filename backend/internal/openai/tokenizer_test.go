package openai

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountTokens(t *testing.T) {
	tokens := CountTokens("Hello, world!", "gpt-4o-mini")
	assert.Greater(t, tokens, 0)
	assert.Less(t, tokens, 10)
}

func TestCountTokens_EmptyString(t *testing.T) {
	tokens := CountTokens("", "gpt-4o-mini")
	assert.Equal(t, 0, tokens)
}

func TestCountTokens_UnknownModel(t *testing.T) {
	tokens := CountTokens("Hello", "unknown-model")
	assert.Greater(t, tokens, 0)
}

func TestPreprocessDiff_RemovesTestFiles(t *testing.T) {
	diff := `diff --git a/main.go b/main.go
--- a/main.go
+++ b/main.go
@@ -1,3 +1,4 @@
+func main() {}
diff --git a/main_test.go b/main_test.go
--- a/main_test.go
+++ b/main_test.go
@@ -1,3 +1,4 @@
+func TestMain() {}
diff --git a/app.spec.ts b/app.spec.ts
--- a/app.spec.ts
+++ b/app.spec.ts
@@ -1,3 +1,4 @@
+test('app', () => {});`

	result := PreprocessDiff(diff)
	assert.Contains(t, result, "main.go")
	assert.NotContains(t, result, "main_test.go")
	assert.NotContains(t, result, "app.spec.ts")
}

func TestPreprocessDiff_RemovesBinaries(t *testing.T) {
	diff := `diff --git a/main.go b/main.go
--- a/main.go
+++ b/main.go
@@ -1,3 +1,4 @@
+func main() {}
diff --git a/image.png b/image.png
Binary files /dev/null and b/image.png differ`

	result := PreprocessDiff(diff)
	assert.Contains(t, result, "main.go")
	assert.NotContains(t, result, "image.png")
}

func TestPreprocessDiff_RemovesLockFiles(t *testing.T) {
	diff := `diff --git a/main.go b/main.go
--- a/main.go
+++ b/main.go
@@ -1,3 +1,4 @@
+func main() {}
diff --git a/package-lock.json b/package-lock.json
--- a/package-lock.json
+++ b/package-lock.json
@@ -1,100 +1,200 @@
+lots of lock content
diff --git a/go.sum b/go.sum
--- a/go.sum
+++ b/go.sum
@@ -1,50 +1,100 @@
+hash content`

	result := PreprocessDiff(diff)
	assert.Contains(t, result, "main.go")
	assert.NotContains(t, result, "package-lock.json")
	assert.NotContains(t, result, "go.sum")
}

func TestPreprocessDiff_RemovesVendorAndDist(t *testing.T) {
	diff := `diff --git a/src/app.go b/src/app.go
--- a/src/app.go
+++ b/src/app.go
@@ -1,3 +1,4 @@
+func app() {}
diff --git a/vendor/lib/lib.go b/vendor/lib/lib.go
--- a/vendor/lib/lib.go
+++ b/vendor/lib/lib.go
@@ -1,3 +1,4 @@
+vendor code
diff --git a/dist/bundle.js b/dist/bundle.js
--- a/dist/bundle.js
+++ b/dist/bundle.js
@@ -1,3 +1,4 @@
+bundled code`

	result := PreprocessDiff(diff)
	assert.Contains(t, result, "src/app.go")
	assert.NotContains(t, result, "vendor/lib")
	assert.NotContains(t, result, "dist/bundle")
}

func TestPreprocessDiff_ReducesContext(t *testing.T) {
	diff := `diff --git a/main.go b/main.go
--- a/main.go
+++ b/main.go
@@ -1,5 +1,6 @@
 package main

 import "fmt"
+import "os"

 func main() {`

	result := PreprocessDiff(diff)
	assert.Contains(t, result, "+import \"os\"")
	assert.NotContains(t, result, " package main")
	assert.NotContains(t, result, " import \"fmt\"")
}

func TestPreprocessDiff_EmptyDiff(t *testing.T) {
	result := PreprocessDiff("")
	assert.Equal(t, "", result)
}

func TestPreprocessDiffForLevel_QADetailed_KeepsContextLines(t *testing.T) {
	diff := `diff --git a/main.go b/main.go
--- a/main.go
+++ b/main.go
@@ -1,10 +1,11 @@
 package main

 import "fmt"

 func hello() {
+	fmt.Println("hello")
 }

 func goodbye() {
 	fmt.Println("bye")
 }`

	result := PreprocessDiffForLevel(diff, "qa_detailed")
	// Should keep the change line
	assert.Contains(t, result, "+\tfmt.Println(\"hello\")")
	// Should keep context lines around the change (within 5 lines)
	assert.Contains(t, result, " package main")
	assert.Contains(t, result, " func hello() {")
	assert.Contains(t, result, " }")
	// Should keep headers
	assert.Contains(t, result, "diff --git a/main.go b/main.go")
	assert.Contains(t, result, "--- a/main.go")
	assert.Contains(t, result, "+++ b/main.go")
}

func TestPreprocessDiffForLevel_Functional_RemovesContext(t *testing.T) {
	diff := `diff --git a/main.go b/main.go
--- a/main.go
+++ b/main.go
@@ -1,5 +1,6 @@
 package main

 import "fmt"
+import "os"

 func main() {`

	result := PreprocessDiffForLevel(diff, "functional")
	assert.Contains(t, result, "+import \"os\"")
	// Context lines should be removed for functional level
	assert.NotContains(t, result, " package main")
	assert.NotContains(t, result, " import \"fmt\"")
}

func TestChunkDiff_SmallDiff(t *testing.T) {
	diff := `diff --git a/main.go b/main.go
+func main() {}`

	chunks := ChunkDiff(diff, 1000)
	assert.Len(t, chunks, 1)
}

func TestChunkDiff_RespectsTokenLimit(t *testing.T) {
	var sections []string
	for i := 0; i < 20; i++ {
		section := "diff --git a/file" + strings.Repeat("x", 100) + ".go b/file.go\n" +
			"+line " + strings.Repeat("content ", 50)
		sections = append(sections, section)
	}
	diff := strings.Join(sections, "\n")

	chunks := ChunkDiff(diff, 200)
	assert.Greater(t, len(chunks), 1)
}

func TestChunkDiff_EmptyDiff(t *testing.T) {
	chunks := ChunkDiff("", 1000)
	assert.Nil(t, chunks)
}
