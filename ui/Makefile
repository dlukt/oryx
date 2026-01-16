.PHONY: default clean build
.PHONY: npm lint test build-cn build-en

default: build

build: npm build build-cn build-en

npm: lint
	pnpm install

# Don't run lint in GitHub Actions, because it crashes the build.
# See https://docs.github.com/en/actions/learn-github-actions/variables#default-environment-variables
lint:
	@if [ -z "$(GITHUB_ACTION)" ] && [ -z "$(SRS_NO_LINT)" ]; then pnpm run lint; fi

test: lint
	pnpm run test

build: build-cn npm

build-cn: npm
	env PUBLIC_URL=/mgmt REACT_APP_LOCALE=zh BUILD_PATH=build/zh pnpm run build

build-en: npm
	env PUBLIC_URL=/mgmt REACT_APP_LOCALE=en BUILD_PATH=build/en pnpm run build

clean:
	rm -rf ui/build
