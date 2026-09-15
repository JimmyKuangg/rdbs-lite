# RDBSLite

A minimal remake of a relational database in Go that emphasizes education over covering all its existing features.

RDBSLite features only a subset of what a real RDBMS has to offer. The goal wasn't to replace Postgres or MySQL. It was to understand the underlying mechanisms and design decisions behind a table-based database with a query language.

## Disclaimer

RDBSLite is intentionally designed to be a learning project. It does not feature all of the robust features that a production database contains, but prioritizes education and understanding over having a highly scalable database. For design decisions made, read my notes about how my brain tried to figure out what to do [here](https://jimmy-kuang.com/notes).

## Installation

- Ensure you have the [latest version of Go installed](https://go.dev/)
- Clone the repository and run commands using Go.

```
git clone https://github.com/JimmyKuangg/rdbslite.git
cd rdbslite

# Running the project
go run .
RDBSLite started
```

## Features

- A REPL to interact with the database directly from your terminal
- Commands such as `CREATE TABLE`, `INSERT INTO`, and `SELECT` that allow you to interact with the database (full list below)
- Typed columns (`INT`, `TEXT`, `BOOL`) with validation on insert
- Data persistence and hydration via an append-only file (AOF) and checkpointed snapshots

## Why Rebuild a Database?

Because I literally had no idea how it worked under the hood. So, what better way to learn than to jump into the deep end and try?

## Usage

### Table of Commands

| Command                 | Args | Example                             | Response      | Explanation                                                                                                                                    |
| ----------------------- | ---- | ----------------------------------- | ------------- | ---------------------------------------------------------------------------------------------------------------------------------------------- |
| CREATE TABLE            | >= 2 | CREATE TABLE users id int name text | OK            | Creates a table inside of the database. Use a command, table name, column name, column type... format                                          |
| INSERT INTO             | >= 2 | INSERT INTO users id 1 name Bob     | OK            | Inserts a value into a table. Use a command, table name, column name, value... format                                                          |
| SELECT... FROM          | >= 3 | SELECT \* FROM users                | Table of data | Selects a certain row of data. Use the \* to select all available data from the table. Use a command, col or wildcard, from, table name format |
| SELECT... FROM... WHERE | >= 6 | SELECT name FROM users WHERE id > 1 | Table of data | Selects data based on a boolean check. Use the appropriate column names and a command, col, from, table name, where, boolean format            |
| SAVE                    | 0    | SAVE                                | OK            | Forces a checkpoint of the current database state to disk                                                                                      |
| PRINT                   | 0    | PRINT                               | Table of data | Prints out the current state of every table in your database                                                                                   |

### Example of Usage

```
go run .
RDBSLite started

# Creates a table
CREATE TABLE users id int name text
OK

# Inserts a row
INSERT INTO users id 1 name Bob
OK

# Selects all data from the table
SELECT * FROM users
+----+------+
| id | name |
+----+------+
| 1  | Bob  |
+----+------+

# Selects data based on a condition
SELECT name FROM users WHERE id > 0
+------+
| name |
+------+
| Bob  |
+------+

# Prints out every table in the database
PRINT
TABLE users
+----+------+
| id | name |
+----+------+
| 1  | Bob  |
+----+------+

# Forces a checkpoint to disk
SAVE
OK

# To exit out and stop the REPL, press Ctrl + C
# This applies for both Windows and Mac
```

## How it Works

Overall structure of the project -

```
                    RDBSLite

              +------------------+
              |      REPL        |
              | Parse & Execute  |
              +--------+---------+
                       |
                       ▼
              +------------------+
              |     Commands     |
              | CREATE / INSERT  |
              |  SELECT / PRINT  |
              +--------+---------+
                       |
                       ▼
              +------------------+
              |   In-Memory DB   |
              | map[string]Table |
              +---+----------+---+
                  |          |
        Append AOF|          |Checkpoint
                  |          |
                  ▼          ▼
           +----------+  +----------+
           |   AOF    |  | Snapshot |
           +----------+  +----------+
```

Multiple things happen upon startup -

```
                  go run .
                      │
                      ▼
        +----------------------------+
        |   Initialize .rdbslite/    |
        |  (storage.db & append.aof) |
        +----------------------------+
                      │
                      ▼
        +----------------------------+
        |      Load Snapshot         |
        |    (.rdbslite/storage.db)  |
        +----------------------------+
                      │
                      ▼
        +----------------------------+
        |       Restore Data         |
        |    Into Memory (Tables)    |
        +----------------------------+
                      │
                      ▼
        +----------------------------+
        |        Replay AOF          |
        |    (.rdbslite/append.aof)  |
        +----------------------------+
                      │
                      ▼
        +----------------------------+
        |     Start the REPL         |
        |   Wait for Client Input    |
        +----------------------------+
                      │
                      ▼
              Ready to Accept Commands
```

Whenever you enter a CLI command, this is the flow of the commands -

```
        Command entered into REPL
                │
                ▼
        +----------------+
        | Parse Command  |
        +----------------+
                │
                ▼
        +----------------+
        | Execute Logic  |
        +----------------+
                │
                ▼
        +----------------+
        | Database (RAM) |
        +----------------+
                │
         ┌──────┴────────┐
         │               │
         ▼               ▼
 Return Response    Append to AOF
   To Client         (if modified)
```

The checkpoint logic flow -

```
        After Every Write Command
                  │
                  ▼
        +----------------------+
        | Check AOF File Size  |
        +----------------------+
                  │
                  ▼
        +----------------------+
        | Above Threshold?     |
        +----------------------+
           │             │
        Yes│             │No
           ▼             ▼
     Save Snapshot    Do Nothing
     & Truncate AOF
```

## Project Structure

`.rdbslite/`

- Where the storage files live. Take your time to watch the two files as you use commands in the CLI! It'll write to them automatically.
- Try NOT to edit them. If you have issues with starting the REPL using `go run .`, delete this folder and run it again. Note that this will delete your data.

`commands/`

- The execution logic behind each CLI command (`CREATE`, `INSERT`, `SELECT`, `PRINT`, `SAVE`) lives here.

`data/`

- Holds the core database logic: tables, rows, column types, and persistence (save/load) of the database.

`repl/`

- All the logic for the REPL sits here. Parsing user input, executing commands, appending to the AOF, replaying the AOF, and checkpointing sit here.

## Limitations

RDBSLite intentionally does not contain all the robust features that the brilliant engineers behind production databases spent years developing, such as:

- A real SQL parser (I just parse whitespace-separated tokens)
- Indexes beyond a simple column-to-position map
- Joins across multiple tables
- Transactions and rollback support
- Concurrent client connections

These features are beyond the scope of an educational project, but can be explored down the line.

## Takeaways and Learnings

Learned a lot about:

- Parsing and validating a small command language
- Designing around failure (What happens if a write fails? A crash mid-checkpoint? etc.)
- A general idea of how a table-based database with an append-only log works

RDBSLite answered many of my original questions about how databases store and query data, but it also raised new ones. I didn't implement joins, transactions, or concurrent access, but understanding why those features exist was just as valuable as building the pieces I did. Those questions are what I'll be exploring next.
