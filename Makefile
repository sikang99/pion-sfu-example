#
# Makefile for testing turn example
#
CLIENT=
SERVER=sfu-server
#----------------------------------------------------------------------------------
.PHONY: usage edit build clean docker compose git
usage:
	@echo "usage: make [edit|build|run|list|docker|compose]"
#----------------------------------------------------------------------------------
edit e:
	@echo "> make (edit:e) [client|server|history]"
edit-client ec:
	@echo "no action for $@"
edit-server es:
	vi $(SERVER)/main.go
edit-history eh:
	vi HISTORY.md
#----------------------------------------------------------------------------------
build b:
	CGO_ENABLED=0 GO111MODULE=on go build -o bin/$(SERVER) $(SERVER)/main.go
	@cp bin/* $(GOPATH)/bin
	@ls -al $(GOPATH)/bin/$(SERVER)
list l:
	go list -m all
clean:
	rm -f bin/*
	docker system prune -f
#----------------------------------------------------------------------------------
run r:
	@echo "> make (run:r) [client|server|web]"
run-client rc:
	@echo "curl http://localhost:8080/sdp -d [Base64 SDP]"
	curl http://localhost:8080/pub http://localhost:8080/sub http://localhost:8080/mon
run-server rs:
	#@-pkill $(SERVER)
	$(SERVER) &
run-web rw:
	open http://localhost:8080/static

kill k:
	@-pkill $(SERVER)
version v:
	go version -m bin/$(SERVER)
#----------------------------------------------------------------------------------
TAG=0.0.1
NAME=teamgrit/pion-sfu
IMAGE=$(NAME):$(TAG)

docker d:
	@echo "> make (docker:d) [build|run|push]"
docker-build db:
	docker build -t $(IMAGE) . -f Dockerfile
	docker images $(NAME)
docker-build-multi dbm:
	docker build -t $(IMAGE) . -f Dockerfile.multi
	docker images $(NAME)
docker-run dr:
	docker run -i -t $(IMAGE) 
docker-push dp:
	docker push $(IMAGE)
#----------------------------------------------------------------------------------
compose c:
	@echo "> make (compose:c) [up|down|status]"
compose-up cu:
	docker-compose up -d
	docker-compose ps
compose-down cd:
	docker-compose down
compose-status cs:
	docker-compose ps
#----------------------------------------------------------------------------------
git g:
	@echo "> make (git:g) [update|store]"
git-update gu:
	git add .
	git commit -a -m "create internal/util package"
	git push
git-store gs:
	git config credential.helper store
#----------------------------------------------------------------------------------

