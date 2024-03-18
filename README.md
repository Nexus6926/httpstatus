# HTTP Status Checker
This tool is a command-line tool that checks the HTTP status codes of a list of URLs and outputs the results to the terminal and a set of files in an "outputs" folder.

### Requirements
- git
- go-lang

### Usage

```
httpstatus -dL <file.txt>

### Input File
The input file should be a text file with one URL per line.


The script will create a file in the "outputs" folder for each status code, and write the status code and URL to the corresponding file.

### Examples
Check the HTTP status codes of the URLs in urls.txt and write the results to `output.txt`:

httpstatus -dL urls.txt

```
### Installation

```
go install github.com/Nexus6926/httpstatus@latest
```