# boatswain

Learning blockchain technology project.

![logo](./logo.png)

### Requirements

* Ubuntu 20.04

### Setup

```shell
bash setup.sh # it will ask sudo
```

### Demo

```shell
rm -rf ./build/
mkdir -p ./build/
(
  cd ./build/
  cmake .. -G Ninja -DENABLE_CPPUNIT=yes
  ninja
  ninja test
)
./build/boatswain

```