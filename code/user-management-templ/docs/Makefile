ASCIIDOCTOR      := asciidoctor
ASCIIDOCTOR_PDF  := asciidoctor-pdf
DIAGRAMS_FLAGS   := -r asciidoctor-diagram
BUILD_DIR        := build
BLOG_DIR         := blog

GUIDE_SRC        := guide.adoc
BLOG_SRCS        := $(wildcard $(BLOG_DIR)/*.adoc)
BLOG_HTML        := $(patsubst $(BLOG_DIR)/%.adoc,$(BUILD_DIR)/blog/%.html,$(BLOG_SRCS))

.PHONY: all html pdf blog clean watch

all: html pdf

## HTML guide (for review)
html: $(BUILD_DIR)
	$(ASCIIDOCTOR) $(DIAGRAMS_FLAGS) \
	  -a source-highlighter=rouge \
	  -a rouge-style=github \
	  -a toc=left \
	  -D $(BUILD_DIR) \
	  $(GUIDE_SRC)

## PDF guide (for distribution)
pdf: $(BUILD_DIR)
	$(ASCIIDOCTOR_PDF) $(DIAGRAMS_FLAGS) \
	  -a source-highlighter=rouge \
	  -D $(BUILD_DIR) \
	  $(GUIDE_SRC)

## Blog post HTML files
blog: $(BUILD_DIR)/blog $(BLOG_HTML)

$(BUILD_DIR)/blog/%.html: $(BLOG_DIR)/%.adoc
	$(ASCIIDOCTOR) $(DIAGRAMS_FLAGS) \
	  -a source-highlighter=rouge \
	  -a rouge-style=github \
	  -a toc \
	  -D $(BUILD_DIR)/blog \
	  $<

## Live preview (requires: apt install entr)
watch:
	find . -name "*.adoc" | entr -r make html

clean:
	rm -rf $(BUILD_DIR)

$(BUILD_DIR) $(BUILD_DIR)/blog:
	mkdir -p $@
