![Screenshot from 2020-05-03 13-51-24](https://user-images.githubusercontent.com/462648/80909535-bd37cf80-8d46-11ea-8658-af760b359791.png)

# shell-butler
![Go](https://github.com/vinodsr/shell-butler/workflows/Go/badge.svg)

This is simple tool to organize your frequently using shell command. 

## Features

* Autocomplete Search
* Auto resize in terminal


## Roadmap

* Feature : add new commands from command line 

## Prerequisite
- Go (>=1.14)

## How to install

You can also download deb packages from [releases](https://github.com/vinodsr/shell-butler/releases/)

```bash
sudo apt install <shell-butler>.deb
```
>
>
> **Note** : 
> 
> If you want to get the shell features to the command execution then create an alias in your ~/.profile or ~/.bash_profile sourcing the binary
> 
> alias sb='source /usr/bin/shellbutler'
>
> This is helpful if you are running a command thay involves some shell theme color variable eg : LS_COLORS
>
>


## How to build

Build the application using

```
go build -o out/ main.go
```

Create a alias to the run script 

```
alias j='source <path to package>/alfred'
```


Now execute "j" from any folder. It will create a dummy settings.json. 

To add your commands edit the settings.json

```json
{
    "commands": [
        {
            "context" : "server:date",   #The command structure. Split levels by :
            "program" : "date"   # The command you need to execute.
        }
    ]
}
```

You can also run the add command

```
alfred add
```

Make sure that run is available in the PATH.

## Tested on

* Ubuntu 20.04 (zsh)
