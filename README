# tl-scraper

`tl-scraper` is a command line tool used for exporting TestLodge testing data using TestLodge provided [APIs](https://help.testlodge.com/hc/en-us/categories/203830188-TestLodge-API).

## Requirements

- Go (v1.21 or higher)
- TestLodge API key and account.

## Usage

`tl-scraper scrap <data> -e <email> -a <account-id> -k <api-key> -o <output-dir>`

Available data:

- all: Scrap all test data from TestLodge APIs
- plan-contents: Scrap test plan contents from TestLodge APIs
- plans: Scrap test plans from TestLodge APIs
- projects: Scrap projects from TestLodge APIs
- requirement-documents: Scrap requirement-documents from TestLodge APIs
- requirements: Scrap requirements from TestLodge APIs
- suite-sections: Scrap suite sections from TestLodge APIs
- suites: Scrap suites from TestLodge APIs
- test-cases: Scrap test cases from TestLodge APIs

## Building

Run:

```bash
go build -o bin/tl-scraper main.go
```

Move binary to any `$PATH` directory (you might need to use sudo). For example:

```bash
mv bin/tl-scraper /usr/local/bin
```
