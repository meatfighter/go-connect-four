# Go Connect Four: Human vs. Computer

## About

This is a command-line version of [the classic vertical boardgame](https://en.wikipedia.org/wiki/Connect_Four) in which players take turns dropping chips into columns in an attempt to form horizontal, vertical or diagonal lines of 4 or more.  It’s a human vs. computer challenge.  The computer’s decisions are governed by [the negamax algorithm](https://en.wikipedia.org/wiki/Negamax) with optimizations that include [alpha-beta pruning](https://en.wikipedia.org/wiki/Alpha%E2%80%93beta_pruning), move ordering and a [transposition table](https://en.wikipedia.org/wiki/Transposition_table).  On startup, the program prompts the user for the algorithm’s maximum search depth; difficulty increases with larger values.   

I wrote this program as an exercising in learning [the Go programming language](https://en.wikipedia.org/wiki/Go_(programming_language)).  Plus, this is my first experience using [GitHub](https://en.wikipedia.org/wiki/GitHub).
