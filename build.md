
udo apt update

# linux/amd64 (本机通常自带)
sudo apt install gcc

# linux/arm (32位 ARM)
sudo apt install gcc-arm-linux-gnueabihf

# linux/arm64 (64位 ARM)
sudo apt install gcc-aarch64-linux-gnu

# windows/amd64 (MinGW 交叉编译)
sudo apt install gcc-mingw-w64-x86-64

# musl-gcc (x86_64)
sudo apt install musl-tools

./build.sh -r