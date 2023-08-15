# Teams Presence Enumeration Tool

Teams Presence Enumeration Tool is a command-line application written in Go that allows users to query the presence information of Microsoft Teams users. Users can supply individual object IDs or a file containing a list of object IDs. The results can be displayed in the console or exported to a CSV file.

This tool is based on a presentation given by nyxgeek, available at [this link](https://media.defcon.org/DEF%20CON%2031/DEF%20CON%2031%20presentations/nyxgeek%20-%20Track%20the%20Planet%20Mapping%20Identities%20Monitoring%20Presence%20and%20Decoding%20Business%20Alliances%20in%20the%20Azure%20Ecosystem.pdf).

When available, I will provide a link here to his repository containing the original tool.
## Features

- Query presence information for a single object ID or a list from a file.
- Export results to a CSV file.

## Installation

To install the Teams Presence Enumeration Tool, you need to have Go installed on your machine. Then, you can clone the repository and build the project.

```bash
git clone https://github.com/username/teamsPresenceEnum.git
cd teamsPresenceEnum
go build
```

## Usage

Query a single object ID:

```
.\teamsPresenceEnum.exe -o d4957c9d-869e-4364-830c-d0c95be72738
```

Query a list of object IDs from a file:


```
.\teamsPresenceEnum.exe -f "C:/Users/objectIds.txt"
```


Export the results to a CSV file:


```
.\teamsPresenceEnum.exe -f "C:/Users/objectIds.txt" -path "path/to/export/file.csv"
```


### Options

- `-o` : Specify a single object ID.
- `-f` : Specify a file path containing object IDs (one per line).
- `-path` : Specify the path to export the CSV file.


## Acknowledgments

- Thanks to nyxgeek for the inspiration and technical guidance for this tool. 
- Thanks to defcon for hosting the presentation.
