# run `curl -fsSL https://raw.githubusercontent.com/wymli/makex/master/install.sh | INSTALL_TYPE=release sh -` to install from remote.


case $INSTALL_TYPE in
  build)
    # build from scrach
    git clone https://github.com/wymli/makex.git
    pushd makex
    go build
    sudo mv makex /usr/local/bin
    popd
    rm -rf makex
    ;;
  release)
    # download release
    url=https://github.com/wymli/makex/releases/latest/download/makex
    curl -fsSL url > makex
    chmod +x makex
    sudo mv makex /usr/local/bin
    ;;
  *)
    echo "not supported"
    ;;
esac

echo "makex installed: $(which makex)"
