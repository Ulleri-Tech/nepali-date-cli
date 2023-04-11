# nepali-date-cli
Command Line Tool (CLI) for BS Date (Nepali Date)

## How to install Linux (arm64)
Download Latest Binary using Wget from https://github.com/Ulleri-Tech/nepali-date-cli/releases/tag/latest

```shell
wget https://github.com/Ulleri-Tech/nepali-date-cli/releases/download/v0.8/nepali-date-cli_0.8_linux_arm64.tar.gz

tar -xzf nepali-date-cli_0.8_linux_arm64.tar.gz convertdate                      

sudo mv convertdate /usr/bin

```
## How to install using brew
```shell
brew tap Ulleri-tech/tap-ulleritech
brew install convertdate
```

## Usage
```shell
convertdate -t #Today's Date in BS

convertdate -A 1993-06-29 # Convert date from AD to BS

convertdate -B 2050-03-15 # Convert date from BS to AD

convertdate -h # Help command

```
#### Run Locally 
```shell
go run main.go
```

#### Build GO Binary
```shell
go build -o name-of-binary    
```

#### Tagging and Push
```shell
git tag -a v0.1 -m "message"
git push origin v0.1
```

#### Misc 
```shell
go mod vendor 
```