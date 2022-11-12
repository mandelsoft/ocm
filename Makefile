NAME      = r3trans
PROVIDER  ?= sap.com
COMPONENT = $(PROVIDER)/r3trans/ocmplugin
OCMREPO   ?= ghcr.io/mandelsoft/cnudie
PLATFORMS = linux/amd64 linux/arm64

REPO_ROOT                                      := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))/../..
VERSION                                        = $(shell git describe --tags --exact-match 2>/dev/null|| echo "$$(cat VERSION)-dev")
COMMIT                                         = $(shell git rev-parse HEAD)
EFFECTIVE_VERSION                              = $(VERSION)-$(COMMIT)

CMDSRCS=$(shell find . -type f -name "*.go") $(REPO_ROOT)/go.*
OCMSRCS=$(shell find $(REPO_ROOT)/pkg -type f -name "*.go") $(REPO_ROOT)/go.*

GEN = $(REPO_ROOT)/local/$(NAME)/gen
OCM = go run $(REPO_ROOT)/cmds/ocm

NOW         := $(shell date --rfc-3339=seconds | sed 's/ /T/')
BUILD_FLAGS := "-s -w \
 -X github.com/open-component-model/ocm/pkg/version.gitVersion=$(EFFECTIVE_VERSION) \
 -X github.com/open-component-model/ocm/pkg/version.gitTreeState=$(GIT_TREE_STATE) \
 -X github.com/open-component-model/ocm/pkg/version.gitCommit=$(COMMIT) \
 -X github.com/open-component-model/ocm/pkg/version.buildDate=$(NOW)"


.PHONY: build
build: $(GEN)/build

$(GEN)/build: $(CMDSRCS) $(OCMSRCS)
	@for i in $(PLATFORMS); do \
    tag=$$(echo $$i | sed -e s:/:-:g); \
    echo GOARCH=$$(basename $$i) GOOS=$$(dirname $$i) go build -ldflags $(BUILD_FLAGS) -o $(GEN)/$(NAME).$$tag ./plugin; \
    GOARCH=$$(basename $$i) GOOS=$$(dirname $$i) go build -ldflags $(BUILD_FLAGS) -o $(GEN)/$(NAME).$$tag ./plugin; \
    done
	@touch $(GEN)/build


.PHONY: ctf
ctf: $(GEN)/ctf

$(GEN)/ctf: $(GEN)/ca.done
	$(OCM) transfer ca $(GEN)/ca $(GEN)/ctf
	touch $(GEN)/ctf

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: ca
ca: $(GEN)/ca.done

$(GEN)/ca.done: $(GEN)/.exists $(GEN)/build component/resources.yaml $(CHARTSRCS)
	$(OCM) create ca -f $(COMPONENT) "$(VERSION)" --provider $(PROVIDER) --file $(GEN)/ca
	$(OCM) add resources --templater=spiff $(GEN)/ca NAME="$(NAME)" VERSION="$(VERSION)" COMMIT="$(COMMIT)" GEN="$(GEN)" PLATFORMS="$(PLATFORMS)" component/resources.yaml
	@touch $(GEN)/ca.done

.PHONY: push
push: $(GEN)/ctf $(GEN)/push.$(NAME)

$(GEN)/push.$(NAME): $(GEN)/ctf
	$(OCM) transfer ctf -f $(GEN)/ctf $(OCMREPO)
	@touch $(GEN)/push.$(NAME)

.PHONY: plain-push
plain-push: $(GEN)
	$(OCM) transfer ctf -f $(GEN)/ctf $(OCMREPO)
	@touch $(GEN)/push.$(NAME)

.PHONY: transport
transport:
ifneq ($(TARGETREPO),)
	$(OCM) transfer component -Vc  $(OCMREPO)//$(COMPONENT):$(VERSION) $(TARGETREPO)
endif

$(GEN)/.exists:
	@mkdir -p $(GEN)
	@touch $@

.PHONY: info
info:
	@echo "ROOT:     $(REPO_ROOT)"
	@echo "VERSION:  $(VERSION)"
	@echo "COMMIT;   $(COMMIT)"

.PHONY: describe
describe: $(GEN)/ctf
	ocm get resources --lookup $(OCMREPO) -c -o treewide $(GEN)/ctf

.PHONY: descriptor
descriptor: $(GEN)/ctf
	ocm get component -S v3alpha1 -o yaml $(GEN)/ctf

.PHONY: clean
clean:
	rm -rf $(GEN)
