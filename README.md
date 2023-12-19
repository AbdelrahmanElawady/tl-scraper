# tl-scraper

`tl-scraper` is a command line tool used for exporting TestLodge testing data using TestLodge provided [APIs](https://help.testlodge.com/hc/en-us/categories/203830188-TestLodge-API).

## Requirements

- TestLodge API key and account.

## Installation

- Download the archive from the [release](https://github.com/AbdelrahmanElawady/tl-scraper/releases) section.
- Unpack the archive and move the binary to any of `$PATH` directories. For example for linux x86:

```bash
tar xf ~/Downloads/tl-scraper_Linux_x86_64.tar.gz && sudo mv tl-scraper /usr/local/bin/
```

## Usage

`tl-scraper scrap <data> -e <env> -o <output-dir> -p <projects>...`

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

## Options

- e: key-value file with the user credentials.
- o: output directory where the scraped data will go.
- p: project IDs to include in the data scraped (not used when scraping projects data).

## Env File

`.env` file contains key-value pairs of user credentials, For example:

```env
key=<TestLodge API key>
accountID=<TestLodge account ID>
email=<TestLodge email>
```

You can get these data from TestLodge help [section](https://help.testlodge.com/hc/en-us/articles/226734768-API-Basics).

## Building

First you need Go (v1.21 or higher).

Run:

```bash
go build -o bin/tl-scraper main.go
```

Move binary to any `$PATH` directory (you might need to use sudo). For example:

```bash
mv bin/tl-scraper /usr/local/bin
```
