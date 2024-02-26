# runnr
runnr is an automated, cross-platform, open source project building cli tool. Simple yet powerfull.

>[!WARNING]
>This project is in alpha stages. If any one finds any bugs, it is recomended to report it to fix such bugs and issuses related to it.

>[!WARNING]
>Not tested for windows yet.

# Features

+ Faster and Efficient
+ runnr script to automate building and compiling project
+ Easy and Strong consistent rules
+ Better error checking

# Documentations

Will be added soon.

# How to install?

## Linux
### Maually:
```bash
go build -o runnr .
mkdir -p ~/.runnr
cp -r ./docs ~/.runnr
cp -r ./templates ~/runnr
cp ./runnr /usr/local/bin   #to install it gloabally
```

### Using runnr:
>[!TIP]
>runnr can be compiled using runnr itself. If you already had runnr installed then run the below command

```bash
runnr build
runnr build install #to install it gloabally
```

### Using makefile:
```bash
make
make install #to install it gloabally
```

# Future plans

+ Accessing OS's enviroment using $env() function.
+ if-else blocks

# Release Cycle

+ Patches: as required
+ Minor Version: every two months
+ Major Version: once an year

# Version

v0.0.1~alpha
