# rdbs-lite

Currently under construction! Read about my shenanigans on my [portfolio notes](https://jimmy-kuang.com/notes)!

## Example Usage

| Command                 | Args | Example                             | Response      | Explanation                                                                                                                                    |
| ----------------------- | ---- | ----------------------------------- | ------------- | ---------------------------------------------------------------------------------------------------------------------------------------------- |
| CREATE TABLE            | >= 2 | CREATE TABLE users id int name text | OK            | Creates a table inside of the database. Use a command, table name, column name, column type... format                                          |
| INSERT INTO             | >= 2 | INSERT INTO users id 1 name Bob     | OK            | Inserts a value into a table. Use a command, table name, column name, value... format                                                          |
| SELECT... FROM          | >= 3 | SELECT \* FROM users                | Table of data | Selects a certain row of data. Use the \* to select all available data from the table. Use a command, col or wildcard, from, table name format |
| SELECT... FROM... WHERE | >= 6 | SELECT name FROM users WHERE id > 1 | Table of data | Selects data based on a boolean check. Use the appropriate column names and a command, col, from, table name, where, boolean format            |
