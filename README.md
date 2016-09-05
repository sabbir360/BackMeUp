# BackMeUp
### File backup system based on size.
This is a sample program which will copy files if new/modified based on config.json.

## How to Run? 
- See a sample JSON file named config.json.sample atteched with the project.
- Copy `config.json.sample` to a new file named `backmeup.config.json` on this directory. Modify source `source_dir` and destination `destination_dir` as required for you. You can add as much you required following this JSON format.
- Make sure you have go installed.
- From this directrory 
  - go run main.go or, 
  - go build main.go && ./main
