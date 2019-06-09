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
### How to Build
To build the executable, run: `go build`.
This was written for go 1.12. The libraries aren't particularly new, so you can probably build with 
older versions, but I don't guarantee it.
### Current Build Status
[![Build Status](https://travis-ci.org/Luke-Sikina/dnd_puzzle_1.svg?branch=master)](https://travis-ci.org/Luke-Sikina/dnd_puzzle_1)
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

# Purpose
I did this little project for fun, but also because I was rusty with Go, and I thought this
would be a good opportunity to get back into the swing of things. This project does or will do a few
things that I think outline how comfortable I am with Go.

## Topics Currently Covered
### Types
- Make more descriptive types by typing already defined types to reflect their function
within the context of the program
- Using iota and constants to create enums
- Writing basic functions for types

### Structs
- Create structs to represent structured data
- Leverage structural typing and type inference to minimize type boilerplate

### Interfaces
- Implement the Formatter interface to customize how objects are printed

### Error Handling
- Implemented standard Go error handling using multiple returns
- Checked errors before continuing execution

### Logging
- Used Go's standard, simple logging to provide extra information for critical steps
- Logging steps that can take a while to show progress
- Logging significant errors

### Unit Testing
- Writing clean, readable tests with reasonable test coverage
- Minimizing test logic while still testing functions with randomness

## Topics To Be Covered
### Channels
- Basic Channel IO

### Goroutines
- Small scale goroutine creation to facilitate concurrent execution

### Iterators / Generators
- Implementing a generator to reduce memory profile
- Writing generators in an idiomatic manner so that they can be used with `range`