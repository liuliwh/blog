## How to find all the shell script files of a directory, and give them the x mode?
### How to tell the file is a shell script?
```shell script
[vagrant@localhost ~]$ file ./test_sh 
./test_sh: POSIX shell script, ASCII text executable
```
### Find all shell script files under ./
```shell script
[vagrant@localhost ~]$ find . -type f -exec file '{}' \; | awk -F: '/shell.script/ {print $1}'
./scripts/check_websites.sh
./test_sh
```
### Option 1
```shell script
find . -type f -exec file '{}' \; | awk -F: '/shell.script/ {print $1}' | xargs chmod +x
```
```shell script
[vagrant@localhost ~]$ find . -type f -exec file '{}' \; | awk -F: '/shell.script/ {print $1}' | xargs ls -l
-rw-rw-r--. 1 vagrant vagrant 197 Aug 28 21:31 ./scripts/check_websites.sh
-rw-rw-r--. 1 vagrant vagrant  22 Sep  2 21:22 ./test_sh
[vagrant@localhost ~]$ find . -type f -exec file '{}' \; | awk -F: '/shell.script/ {print $1}' | xargs chmod +x
[vagrant@localhost ~]$ find . -type f -exec file '{}' \; | awk -F: '/shell.script/ {print $1}' | xargs ls -l
-rwxrwxr-x. 1 vagrant vagrant 197 Aug 28 21:31 ./scripts/check_websites.sh
-rwxrwxr-x. 1 vagrant vagrant  22 Sep  2 21:22 ./test_sh
```
### Option 2
```shell script
chmod +x `find . -type f -exec file '{}' \; | awk -F: '/shell.script/ {print $1}'`
```
```shell script
[vagrant@localhost ~]$ find . -type f -exec file '{}' \; | awk -F: '/shell.script/ {print $1}' | xargs ls -l
-rw-rw-r--. 1 vagrant vagrant 197 Aug 28 21:31 ./scripts/check_websites.sh
-rw-rw-r--. 1 vagrant vagrant  22 Sep  2 21:22 ./test_sh
[vagrant@localhost ~]$ chmod +x `find . -type f -exec file '{}' \; | awk -F: '/shell.script/ {print $1}'`
[vagrant@localhost ~]$ find . -type f -exec file '{}' \; | awk -F: '/shell.script/ {print $1}' | xargs ls -l
-rwxrwxr-x. 1 vagrant vagrant 197 Aug 28 21:31 ./scripts/check_websites.sh
-rwxrwxr-x. 1 vagrant vagrant  22 Sep  2 21:22 ./test_sh
```
### Notes
1. file command. file — determine file type
2. find command with exec. 
```shell script
[vagrant@localhost ~]$ # If sometimes I forgot -exec cmd '{}' \;
[vagrant@localhost ~]$ # I would use pipeline line with xargs. for example
[vagrant@localhost ~]$ find . -type f | xargs file
./.bash_logout:              ASCII text
./.bash_profile:             ASCII text
./.bashrc:                   ASCII text
./.ssh/authorized_keys:      OpenSSH RSA public key
./.bash_history:             ASCII text
./scripts/oldboy.txt:        ASCII text
./scripts/check_websites.sh: Bourne-Again shell script, ASCII text executable
./random_number:             data
./test_sh:                   POSIX shell script, ASCII text executable
./test:                      ASCII text
./.viminfo:                  ASCII text
```
3. awk pattern match
