# Go Connect Four: Human vs. Computer

### About

This is a command-line version of [the classic vertical boardgame](https://en.wikipedia.org/wiki/Connect_Four) in which players take turns dropping chips into columns in an attempt to form horizontal, vertical or diagonal lines of 4 or more.  It’s a human vs. computer challenge.  The computer’s decisions are governed by [the negamax algorithm](https://en.wikipedia.org/wiki/Negamax) with optimizations that include [alpha-beta pruning](https://en.wikipedia.org/wiki/Alpha%E2%80%93beta_pruning), move ordering and a [transposition table](https://en.wikipedia.org/wiki/Transposition_table).  On startup, the program prompts the user for the algorithm’s maximum search depth; difficulty increases with larger values.   

I wrote this program as an exercise while learning [the Go programming language](https://en.wikipedia.org/wiki/Go_(programming_language)).

### Sample Gameplay
```
Connect Four: Human vs. Computer     

Enter max search depth (4--20) [10]: 16

Human (X) plays first.

  1 2 3 4 5 6 7       
1 . . . . . . .       
2 . . . . . . .       
3 . . . . . . .       
4 . . . . . . .       
5 . . . . . . .       
6 . . . . . . .       

Enter column: 4

  1 2 3 4 5 6 7 
1 . . . . . . . 
2 . . . . . . . 
3 . . . . . . . 
4 . . . . . . . 
5 . . . . . . . 
6 . . . X . . . 

Computer drops O into column 4.

  1 2 3 4 5 6 7
1 . . . . . . .
2 . . . . . . .
3 . . . . . . .
4 . . . . . . .
5 . . . O . . .
6 . . . X . . .

Enter column: 4

  1 2 3 4 5 6 7
1 . . . . . . .
2 . . . . . . .
3 . . . . . . .
4 . . . X . . .
5 . . . O . . .
6 . . . X . . .

Computer drops O into column 4.

  1 2 3 4 5 6 7
1 . . . . . . .
2 . . . . . . .
3 . . . O . . .
4 . . . X . . .
5 . . . O . . .
6 . . . X . . .

Enter column: 4

  1 2 3 4 5 6 7
1 . . . . . . .
2 . . . X . . .
3 . . . O . . .
4 . . . X . . .
5 . . . O . . .
6 . . . X . . .

Computer drops O into column 4.

  1 2 3 4 5 6 7
1 . . . O . . .
2 . . . X . . .
3 . . . O . . .
4 . . . X . . .
5 . . . O . . .
6 . . . X . . .

Enter column: 3

  1 2 3 4 5 6 7
1 . . . O . . .
2 . . . X . . .
3 . . . O . . .
4 . . . X . . .
5 . . . O . . .
6 . . X X . . .

Computer drops O into column 2.

  1 2 3 4 5 6 7
1 . . . O . . .
2 . . . X . . .
3 . . . O . . .
4 . . . X . . .
5 . . . O . . .
6 . O X X . . .

Enter column: 5

  1 2 3 4 5 6 7
1 . . . O . . .
2 . . . X . . .
3 . . . O . . .
4 . . . X . . .
5 . . . O . . .
6 . O X X X . .

Computer drops O into column 6.

  1 2 3 4 5 6 7
1 . . . O . . .
2 . . . X . . .
3 . . . O . . .
4 . . . X . . .
5 . . . O . . .
6 . O X X X O .

Enter column: 5

  1 2 3 4 5 6 7
1 . . . O . . .
2 . . . X . . .
3 . . . O . . .
4 . . . X . . .
5 . . . O X . .
6 . O X X X O .

Computer drops O into column 5.

  1 2 3 4 5 6 7
1 . . . O . . .
2 . . . X . . .
3 . . . O . . .
4 . . . X O . .
5 . . . O X . .
6 . O X X X O .

Enter column: 3

  1 2 3 4 5 6 7
1 . . . O . . .
2 . . . X . . .
3 . . . O . . .
4 . . . X O . .
5 . . X O X . .
6 . O X X X O .

Computer drops O into column 3.

  1 2 3 4 5 6 7
1 . . . O . . .
2 . . . X . . .
3 . . . O . . .
4 . . O X O . .
5 . . X O X . .
6 . O X X X O .

Enter column: 5

  1 2 3 4 5 6 7
1 . . . O . . .
2 . . . X . . .
3 . . . O X . .
4 . . O X O . .
5 . . X O X . .
6 . O X X X O .

Computer drops O into column 2.

  1 2 3 4 5 6 7
1 . . . O . . .
2 . . . X . . .
3 . . . O X . .
4 . . O X O . .
5 . O X O X . .
6 . O X X X O .

Enter column: 1

  1 2 3 4 5 6 7
1 . . . O . . .
2 . . . X . . .
3 . . . O X . .
4 . . O X O . .
5 . O X O X . .
6 X O X X X O .

Computer drops O into column 5.

  1 2 3 4 5 6 7
1 . . . O . . .
2 . . . X O . .
3 . . . O X . .
4 . . O X O . .
5 . O X O X . .
6 X O X X X O .

Computer wins.

Play again (y/n)? n

Thanks for playing.

Goodbye.
```