# shell-butler

This is simple tool to organize your frequently using shell command. 

## Features

* Autocomplete Search
* Auto resize in terminal


## Roadmap

* Feature : add new commands from command line 

## Prerequisite
- Go (>=1.14)

## How to install

Build the application using

```
go build -o out/ main.go
```

Create a alias to the run script 

```
alias j='source <path to package> /run'
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

## Tested on

* Ubuntu 20.04 (zsh)
