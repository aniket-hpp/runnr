# runnr
runnr is an automated, cross-platform, open source project building cli tool. Simple yet powerfull.

>[!WARNING]
>This project is in alpha stages. If any one finds any bugs, it is recomended to report it to fix such bugs and issuses related to it.

# Features

+ Faster and Efficient
+ runnr script to automate building and compiling project
+ Easy and Strong consistent rules
+ Better error checking
+ Backwards compatibility with older 'build.runnr' scripts.

# Documentations

Will be added soon with version v0.3.0~alpha.

# How to install?

>[!IMPORTANT]
>You need to install golang first, if you haven't install it from [here](https://go.dev/doc/install).

## Clone the project:
```bash
git clone https://github.com/aniket-hpp/runnr.git
cd runnr
```

OR

## Download the project:
Go to this [link](https://github.com/aniket-hpp/runnr.git) and download latest release '.zip' file and extract it.

```bash
cd runnr
```

## Linux
### Maually:
```bash
go build -o runnr .
mkdir -p ~/.runnr
cp -r ./docs ~/.runnr
cp -r ./templates ~/runnr
sudo cp ./runnr /usr/local/bin   #to install it globally
```

### Using runnr:
>[!TIP]
>runnr can be compiled using runnr itself. If you already have runnr installed then run the below command.

First copy the unix.runnr file as build.runnr from "build" directory to root of the project.

```bash
cp ./build/unix.runnr ./build.runnr
```

Then run these below commands:

```bash
runnr build
runnr build install #to install it globally
```

### Using makefile:
```bash
make
make install #to install it globally
```

## Windows
### Manually:
Open poweshell with Administrative privileges and run these following commands:

```powershell
go build -o runnr.exe .
rm -r $HOME\.runnr
mkdir $HOME\.runnr
cp -r .\docs $HOME\.runnr
cp -r .\templates $HOME\.runnr
cp .\runnr.exe $HOME\.runnr
```

### Using runnr:
>[!TIP]
>runnr can be compiled using runnr itself. If you already have runnr installed then run the below command.

First copy the windows.runnr file as build.runnr from "build" directory to root of the project.

```powershell
cp .\build\windows.runnr .\build.runnr
```

Then run these below commands in powershell with Administrative privileges:

```powershell
runnr build
runnr build install #to install it globally
```

>[!IMPORTANT]
>After choosing one of the above options for installing "runnr" in windows, you have to add the "%USERPROFILE%\.runnr" to your enviorment variable by following the steps below:

1. Search for "View advanced system setting" and open it.
2. Go to "Advanced" tab.
3. Click "Environment Variables".
4. Under "User Variables", find the `PATH` or `Path` variable, select it, and click "Edit".
5. Then click "New" and add "%USERPROFILE%\.runnr".
6. Click "OK".
7. Restart your terminal.

# Changes

### v0.2.0~alpha:

+ Added accessing OS's enviroment using $env([name]) function.

### v0.2.5~alpha:

+ Bug Fixes: config-path
+ Working on support for windows.

### v0.2.7~alpha:

+ Added support for Microsoft's Windows OS.

### v0.2.8~alpha:

+ Added option to call $env([name]) in .modin directive.
+ Added enhanced error handling for pre-processor with file-path stack.

# Release Cycle

>Will be followed after Stable Release:
+ Patches: as required
+ Minor Version: every two months
+ Major Version: once an year

# Version

v0.2.8~alpha
