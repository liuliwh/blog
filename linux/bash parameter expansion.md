## Use bash parameter expansion for string manipulation.
1. '#' Remove match prefix pattern. 
2. '%' Remove match postfix pattern.
## Rename *.html to *.txt
```shell script
[vagrant@localhost ~]$ touch {00..10}.html
[vagrant@localhost ~]$ for name in `ls *.html`; do mv $name ${name%.html}.txt ;done
```