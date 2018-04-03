# Sqlfmt

Format the SQL of the Go code.

In the Go language, tools that allow you to adjust the format and import of the code are provided with standard tools.
sqlfmt rewrites and format the SQL in the code like those tools.
supports only MySQL SQL statements.

## Install

```
go get -u github.com/sgswtky/sqlfmt
```

## Usage

### Replace file

The SQL is formatted if both of the following conditions are satisfied.
  - Variable name 'sql'
  - If '// sqlfmt' is commented one line before that variable

Please be aware that all comments will be deleted after formatting.


For example, if the following file exists,

```sample.go
package main

import "fmt"

func main() {
	// sqlfmt
	sql := `select * from example where example_id = ? and example_name like '%example%'`
	
	fmt.Println(sql)
}

```

`sqlfmt` specifies the example file with the` -f` option and adds the `-w` option.

```
./sqlfmt -f sample.go -w
```

The SQL statement in the `sample.go` file is changed as follows.

```sample.go
package main

import "fmt"

func main() {
	// sqlfmt
	sql := `
SELECT
  *
FROM
  example
WHERE
  example_id = ?
  AND example_name like "%example%"`
	fmt.Println(sql)
}

```

### Action mode

| option | detail |
|---|---|
| -f | File mode. Please specify the file name |
| -w | Rewrite the target file. It works only when combined with file mode |
| -i | Interactive mode. |
|    | Pipe mode. See below. |

#### Pipe mode

```
⇒  echo "select * from example where example_id = ? and example_name LIKE '%sgsw%'" | sqlfmt
SELECT
  *
FROM
  example
WHERE
  example_id = ?
  AND example_name like "%sgsw%"
```

### Detailed Behavior

 - Comments are not supported and will be deleted
 - Since it only supports static SQL
 - do not combine strings and variables in the format target variable.

## Thanks

- github.com/vitessio/vitess
- github.com/xwb1989/sqlparser

## License

apache license 2.0