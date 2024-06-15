build:
	cd cmd && go build
test:build
	cmd/cmd
%:examples/%
	cd $< && go build && ./$@