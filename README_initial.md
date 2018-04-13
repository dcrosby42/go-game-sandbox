# getting up and running

## Setting up project w go dep

- Installing: https://golang.github.io/dep/docs/installation.html
  - I opted for the direct route: 
- Getting started doc: https://golang.github.io/dep/docs/introduction.html

```
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
cd $GOPATH/src/github.com/dcrosby/go-game-sandbox
dep init # makes Gopkg.toml and Gopkg.lock and vendor/
```


## Getting setup with game lib deps

```
# Get a main file with some imports that I lack
vim scratch1/main.go # copied from https://github.com/icexin/gocraft/blob/master/main.go, references several GL libs

dep ensure # downloads all the deps into vendor and updates Gopkg.toml
./vendor/github.com/go-gl/glfw/v3.2/glfw/glfw # <-- manually copy this subdir from a separately-cloned git repo https://github.com/go-gl/glfw.git
```
