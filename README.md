# Patchi

## Diff & migration tool like Git but for databases
Patchi connects to 2 of your databases and shows you the differences between them. Useful for
migrating database environments. It also generates a migration SQL script for you.


## Requirements
- Go 1.18 or higher

## Build
```bash
go build -o patchi
```

## Usage

### Commands

#### 1. `add`
Adds a new connection to a local config file. You can add as many connections as you want. It prompts you to enter the connection details.
All database information are stored locally in a config file. The config file is located at the equivalent of `~/.patchi/config.json` on your OS.
```bash
./patchi add
```

#### 2. `list`
Lists all the connections in the config file.
```bash
./patchi list
```

#### 3. `compare`
Shows the differences between 2 connections. It prompts you to select the connections you want to compare.
```bash
./patchi compare
```

#### 4. `rm`
Removes a connection from the config file. It prompts you to select the connection you want to remove.
It takes an optional argument which is the name of the connection you want to remove. Otherwise it prompts you to select the connection you want to remove.
```bash
./patchi rm [optional-connection-name]
```


## Contributing
Pull requests are always welcomed and encouraged. For major changes, please open an issue first to discuss what you would like to change.

## License
This project is licensed under the [MIT License](https://choosealicense.com/licenses/mit/). See the [LICENSE](LICENSE) file for details.