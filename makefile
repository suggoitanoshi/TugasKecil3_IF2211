WASM_EXEC ?= "$(go env GOROOT)/misc/wasm/wasm_exec.js"

OUT=out
MAIN=main/main.go
MAIN_OUT=$(OUT)/main.wasm
IDX_HTML_SRC=html/index.html
IDX_HTML_OUT=$(OUT)/index.html
NODEMAKER_SRC=html/nodemaker.html
NODEMAKER_OUT=$(OUT)/nodemaker.html
WASM_OUT=$(OUT)/wasm_exec.js

$(OUT):
	mkdir -p $@

$(MAIN_OUT): $(MAIN) $(OUT)
	GOOS=js GOARCH=wasm go build -o $@ $< 

$(IDX_HTML_OUT): $(IDX_HTML_SRC) $(OUT)
	cp $< $@

$(NODEMAKER_OUT): $(NODEMAKER_SRC) $(OUT)
	cp $< $@

$(WASM_OUT): $(WASM_EXEC)
	cp $< $@

run: server.go $(IDX_HTML_OUT) $(NODEMAKER_OUT) $(WASM_OUT) $(MAIN_OUT)
	go run $< -dir=$(OUT)

clean:
	rm -rf $(OUT)
.PHONY: run clean