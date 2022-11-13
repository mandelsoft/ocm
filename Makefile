NAME      = r3transdemo
PROVIDER  ?= mandelsoft
COMPONENT = github.com/$(PROVIDER)/$(NAME)
OCMREPO   ?= ghcr.io/$(PROVIDER)/cnudie

REPO_ROOT                                      := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))/../..
VERSION                                        = $(shell git describe --tags --exact-match 2>/dev/null|| echo "$$(cat VERSION)-dev")
COMMIT                                         = $(shell git rev-parse HEAD)
EFFECTIVE_VERSION                              = $(VERSION)-$(COMMIT)

GEN = gen
ifneq ($(wildcard $(REPO_ROOT)/cmds/ocm),)
OCM = go run $(REPO_ROOT)/cmds/ocm
else
OCM = ocm
endif

REQ=BINK08154711
SYS=/tmp/r3trans/source/BIN

.PHONY: ctf
ctf: $(GEN)/ctf

$(GEN)/ctf: $(GEN)/ca
	$(OCM) transfer ca $(GEN)/ca $(GEN)/ctf
	touch $(GEN)/ctf

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: ca
ca: $(GEN)/ca

$(GEN)/ca: $(GEN)/.exists  $(SYS)/$(REQ) resources.yaml $(CHARTSRCS)
	$(OCM) create ca -f $(COMPONENT) "$(VERSION)" --provider $(PROVIDER) --file $(GEN)/ca
	$(OCM) add resources $(GEN)/ca VERSION="$(VERSION)" COMMIT="$(COMMIT)" REQUEST="$(REQ)" resources.yaml
	@touch $(GEN)/ca

.PHONY: push
push: $(GEN)/ctf $(GEN)/push.$(NAME)

$(GEN)/push.$(NAME): $(GEN)/ctf
	$(OCM) transfer ctf -f $(GEN)/ctf $(OCMREPO)
	@touch $(GEN)/push.$(NAME)

.PHONY: plain-push
plain-push: $(GEN)
	$(OCM) transfer ctf -f $(GEN)/ctf $(OCMREPO)
	@touch $(GEN)/push.$(NAME)

$(SYS)/$(REQ):
	mkdir -p "$(SYS)"
	echo "This is a demo transport request $(REQ)" >"$(SYS)/$(REQ)"

.PHONY: transport
transport:
ifneq ($(TARGETREPO),)
	$(OCM) transfer component -Vc  $(OCMREPO)//$(COMPONENT):$(VERSION) $(TARGETREPO)
endif

$(GEN)/.exists:
	@mkdir -p $(GEN)
	@mkdir -p local
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
	rm -r local/target local/transport.*
	rm -rf $(GEN)
	rm -rf $(SYS)

#
# demo commands
#

.PHONY: setup
setup:
	$(OCM) install plugin -f $(OCMREPO)//sap.com/r3trans/ocmplugin

.PHONY: local-download
local-download:
	$(OCM) download resource "$(GEN)"/ctf -O local/transport.local transport
	cat local/transport.local

.PHONY: local-transport
local-transport:
	@rm -rf local/target
	$(OCM) transfer ctf -V "$(GEN)"/ctf local/target

.PHONY: local-import
local-import:
	@rm -rf local/target
	@echo "target spec:"
	@echo '$(shell cat importtarget.yaml | yaml2json)'
	$(OCM) transfer ctf -V --uploader plugin/r3trans.sap.com=@importtarget.yaml '$(GEN)'/ctf local/target

.PHONY: target-download
target-download:
	$(OCM) download resource local/target -O local/transport.target transport
	cat local/transport.target

.PHONY: target-describe
target-describe:
	$(OCM) get resources local/target -o yaml