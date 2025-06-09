# mvcommon

`mvcommon` tidies directories by detecting a shared prefix in file names and automatically moving those files into a folder named after that prefix. It exists so you don't need to keep writing throwaway scripts or manually sort files one by one.

## Usage

```
$ mvcommon              
Usage: mvcommon [-stopword=<stopword:` - `,`] `,`[`>] [-trim=<trim:-_ .>] [-min=3] [-dry-run] [-interactive] <file1> <file2> ...
```
## Features

- Automatically detects common filename prefixes
- Optionally remove stop words and trimming characters
- `-interactive` mode to confirm operations
- `-dry-run` shows what would change without modifying files

## Why use mvcommon?

`mvcommon` shines when you regularly download or create files that share a common prefix. Instead of manually creating folders and dragging files around, a single command sorts everything for you. Typical scenarios include:

- Grouping invoices or reports that start with the same project code
- Collecting export files that all begin with a date or build number
- Tidying up camera downloads so your holiday photos end up in their own directory

By automating these chores, `mvcommon` keeps folders organised and saves you the tedious clicks of moving files individually.


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

```bash
$ ls
mvcommon
$ touch "Report 234 - Draft1.txt" "Report 234 - Draft2.txt" "Report 234 - Final.txt"
$ ./mvcommon *.txt                                                                  
Creating folder: Report 234
Moved Report 234 - Draft1.txt -> Report 234/Report 234 - Draft1.txt
Moved Report 234 - Draft2.txt -> Report 234/Report 234 - Draft2.txt
Moved Report 234 - Final.txt -> Report 234/Report 234 - Final.txt
Operation completed successfully.
$ ls -lR                                                                            
.:
total 2384
-rwxr-xr-x 1 arran arran 2440869 Dec  9 13:45  mvcommon
drwxr-xr-x 1 arran arran     136 Dec  9 13:46 'Report 234'

'./Report 234':
total 0
-rw-r--r-- 1 arran arran 0 Dec  9 13:45 'Report 234 - Draft1.txt'
-rw-r--r-- 1 arran arran 0 Dec  9 13:45 'Report 234 - Draft2.txt'
-rw-r--r-- 1 arran arran 0 Dec  9 13:45 'Report 234 - Final.txt'
```


```bash
$ touch "[Draft] Report 234.txt" "[For Review a] Report 234 - Version 2.txt" "[Final] Report 234.txt"
$ ls -l
total 2384
-rw-r--r-- 1 arran arran       0 Dec  9 13:46 '[Draft] Report 234.txt'
-rw-r--r-- 1 arran arran       0 Dec  9 13:46 '[Final] Report 234.txt'
-rw-r--r-- 1 arran arran       0 Dec  9 13:46 '[For Review a] Report 234 - Version 2.txt'
-rwxr-xr-x 1 arran arran 2440869 Dec  9 13:45  mvcommon
$ ./mvcommon *.txt                                                                                   
Creating folder: Report 234
Moved [Draft] Report 234.txt -> Report 234/[Draft] Report 234.txt
Moved [Final] Report 234.txt -> Report 234/[Final] Report 234.txt
Moved [For Review a] Report 234 - Version 2.txt -> Report 234/[For Review a] Report 234 - Version 2.txt
Operation completed successfully.
$ find .
.
./mvcommon
./Report 234
./Report 234/[Draft] Report 234.txt
./Report 234/[Final] Report 234.txt
./Report 234/[For Review a] Report 234 - Version 2.txt
```


```bash
$ touch "file_one.txt" "file_two.txt" "file_three.txt"                                               
$ ls -l
total 2384
-rw-r--r-- 1 arran arran       0 Dec  9 13:47 file_one.txt
-rw-r--r-- 1 arran arran       0 Dec  9 13:47 file_three.txt
-rw-r--r-- 1 arran arran       0 Dec  9 13:47 file_two.txt
-rwxr-xr-x 1 arran arran 2440869 Dec  9 13:45 mvcommon
$ ./mvcommon *.txt                                    
Creating folder: file
Moved file_one.txt -> file/file_one.txt
Moved file_three.txt -> file/file_three.txt
Moved file_two.txt -> file/file_two.txt
Operation completed successfully.
$ tree .
.
├── file
│   ├── file_one.txt
│   ├── file_three.txt
│   └── file_two.txt
└── mvcommon

2 directories, 4 files
```

## Interactive mode

```bash
% touch "file_one.txt" "file_two.txt" "file_three.txt"
% ./mvcommon -interactive *.txt

Interactive Mode Enabled:
The following files are detected:
1. file_one.txt
2. file_three.txt
3. file_two.txt

Enter file numbers to include (e.g., 1,2,3 or 1-3,5-6) or press 'a' to confirm all:
Your choice: 1-2
Selected files:
 - file_one.txt
 - file_three.txt
Creating folder: file
Moved file_one.txt -> file/file_one.txt
Moved file_three.txt -> file/file_three.txt
Operation completed successfully.
% tree .
.
├── file
│   ├── file_one.txt
│   └── file_three.txt
├── file_two.txt
└── mvcommon

2 directories, 4 files
```

## Download

See Github releases here: https://github.com/arran4/mvcommon/releases

## Gentoo

See: https://github.com/arran4/arrans_overlay

```bash
$ eselect repository enable arrans-overlay
$ emerge -va app-misc/mvcommon-bin 
```

---
Licensed under the terms of the GPL-3.0.
