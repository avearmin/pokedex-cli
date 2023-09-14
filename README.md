# pokedex-cli
 A command-line REPL Pokedex powered by PokéAPI.

The Pokedex CLI is a command-line tool that allows you to interact with Pokemon data. You can explore locations, catch Pokemon, inspect them, and view your Pokedex. This README provides an overview of the available commands and how to use them.

## Getting Started

1. Clone the repository:
```shell
git clone https://github.com/avearmin/pokedex-cli.git
cd pokedex-cli
```

2. Build the project:
```shell
go build
```

3. Run the Pokedex CLI:
```shell
./pokedex-cli
```

## Commands

-   **help**: Displays available commands and their descriptions.
    ```shell
    help
    ```
    
-   **exit**: Exits the Pokedex CLI.
    ```shell
    exit
    ```
    
-   **map**: Gets the next 20 locations.
    ```shell
    map
    ```
    
-   **mapb**: Gets the previous 20 locations.
    ```shell
    mapb
    ```
    
-   **explore**: Gets a list of Pokemon in a specific location.
    ```shell
    explore <area-name>
    ```
    
-   **catch**: Attempts to catch a Pokemon and saves it to your Pokedex.
    ```shell
    catch <pokemon>
    ```
    
-   **inspect**: Inspects a Pokemon in your Pokedex, providing information about its name, height, weight, stats, and types.
    ```shell
    inspect <pokemon>
    ```
    
-   **pokedex**: Displays a list of all Pokemon in your Pokedex.
    ```shell
    pokedex
    ```
