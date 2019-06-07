# Introduction
This is just a cute little puzzle generator for DND or any other sand boxed tabletop RPG. 
It generates a series of operators, and then sets of problems that pair with them. The puzzle is
to discover the operators using the provided problems.

## Example
**Secret Operators**:
*, -
**Numbers**:
((2 ? 13) ? 13) = 13  
((16 ? 14) ? 6) = 218

The dungeon master gives the players the numbers and it is their job to discover the secret operators.

## Caveats
* Conventional order of operations is ignored. Just go left to right.  
* This is integer math, so don't round, just drop the decimal.
* I think these problems get really hard to solve once you go past 2-3 operators. You can go up
to ~10 right now, but I don't think the resulting problem is remotely solvable.

# Usage
## Building
To build the executable, run: `go build`.
This was written for go 1.12. The libraries aren't particularly new, so you can probably build with 
older versions, but I don't guarantee it.
## Execution
The base command is just
`./dnd_puzzle_1`
### Flags
`terms`: specify the number of terms (numbers) in each sub problem. Default: 3. Example: -terms 10
`min`: specify the minimum number any term can be. Default 2. Example: -min 0
`max`: specify the maximum number any term can be. Default 20. Example:  -max 12
`help`: shows documentation, then exits.
### Example
`./dnd_puzzle_1 -terms 8 -min 4 -max 16`