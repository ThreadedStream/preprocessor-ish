# preprocessor-ish
"almost-c-preprocessor" written in Go

This project is an example of toying with a Golang's stdlib text/scanner package, which allows one to easily build parsers for any language. 

In order to run a program, you need to type in commands listed below:
```bash
  go build main.go
  ./main
```
By the end of the execution of a program, you will see a created directory titled "gen", and two files happily residing within the latter. These two files
represent the preprocessed versions of original sources located in eponymous directory.  

Well, that's it for now. Have a great day, pal (or night). 

## Known limitations

Currently, preprocessor only is capable of dealing with the simple define macro directives. I didn't want to complicate the program that much, as 
it should only serve as an example.  
