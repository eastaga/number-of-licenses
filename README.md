# Number of licenses required for purchasing

Some applications from vendors are allowed to be installed on multiple computers per user with specific restrictions. 
In our scenario, each copy of the application (ID 374) allows the user to install the application on to two computers if at least one of them is a laptop.
Given the provided data, create a utility that would calculate the minimum number of copies of the application the company must purchase.

# Using Utility

Set app id and csv file as env vars. Then run the binary.
```shell
    export APP_ID="374"
    export CSV_FILE="sample-large.csv"
    go run main.go
```

# Testing Utility

Run below command to run unit tests with coverage
```shell
    make test
```