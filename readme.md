# mvcommmon

This is a basic program where you specify a bunch of files or folders as arguments, then it will find the first common 
string in the filename, and create a folder and move the contents into that folder.

# Usage:

```
Usage: mvcommon [-stopword=<stopword: - >] [-trim=<trim:- _>] [-dry-run] <file1> <file2> ...
```

# Examples

```bash
$ touch 2024-02-34 2024-04-02 2024-04-03
$ ls -lh
total 2.4M
-rw-r--r-- 1 arran arran    0 Dec  9 10:35 2024-02-34
-rw-r--r-- 1 arran arran    0 Dec  9 10:35 2024-04-02
-rw-r--r-- 1 arran arran    0 Dec  9 10:35 2024-04-03
-rwxr-xr-x 1 arran arran 2.3M Dec  9 10:34 mvcommon
$ ./mvcommon -stopword=- 20* 
Creating folder: 2024
Moved 2024-02-34 -> 2024/2024-02-34
Moved 2024-04-02 -> 2024/2024-04-02
Moved 2024-04-03 -> 2024/2024-04-03
Operation completed successfully.
```

