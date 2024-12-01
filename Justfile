compile:
    pnpm build

    rm -rf ./compile

    bun build --compile --target=bun-darwin-arm64 --minify --sourcemap --bytecode ./dist/index.js --outfile ./compile/cpenv-darwin-arm64

    bun build --compile --target=bun-darwin-x64 --minify --sourcemap --bytecode ./dist/index.js --outfile ./compile/cpenv-darwin-x64
