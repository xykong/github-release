ifeq ($(OS),Windows_NT)
    CCFLAGS += -D WIN32
    SWAGGER = ./scripts/swagger.exe
    ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
        CCFLAGS += -D AMD64
    else
        ifeq ($(PROCESSOR_ARCHITECTURE),AMD64)
            CCFLAGS += -D AMD64
        endif
        ifeq ($(PROCESSOR_ARCHITECTURE),x86)
            CCFLAGS += -D IA32
        endif
    endif
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        CCFLAGS += -D LINUX
        SWAGGER = ./scripts/swagger-Linux
    endif
    ifeq ($(UNAME_S),Darwin)
        CCFLAGS += -D OSX
        SWAGGER = ./scripts/swagger-Darwin
    endif
    UNAME_P := $(shell uname -p)
    ifeq ($(UNAME_P),x86_64)
        CCFLAGS += -D AMD64
    endif
    ifneq ($(filter %86,$(UNAME_P)),)
        CCFLAGS += -D IA32
    endif
    ifneq ($(filter arm%,$(UNAME_P)),)
        CCFLAGS += -D ARM
    endif
endif

SOURCE_DIR=.

FILES = $(foreach dir,$(DIRS),$(wildcard $(dir)/*.go))
SOURCES := $(shell find $(SOURCE_DIR) -name '*.go')

.PHONY : clean all version test docker

default : github-release

github-release : $(SOURCES)
	go build

version :
	./scripts/version.sh VERSION ./cmd/version.go README.md

test :
	go test `go list ./... | grep -v /vendor/`

clean :
	-rm github-release

publish : release
	git commit -a -m "publish version `cat VERSION`."
#	git push
	echo "publish version `cat VERSION` success."

release : version github-release
	@tag=$$(cat VERSION | awk '{print $$1}');\
	echo $$tag;\
	id=$$(./github-release create -s --tag_name $${tag} --name "github-release $${tag}" --body "Publish the release package.");\
	tar czvf github-release_$${tag}_darwin_amd64.tar.gz github-release;\
	./github-release upload -i $${id} github-release_$${tag}_darwin_amd64.tar.gz;\
    rm -f github-release_$${tag}_darwin_amd64.tar.gz
