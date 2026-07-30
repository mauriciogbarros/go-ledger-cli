# Go Transaction Ledger CLI

![Go](https://img.shields.io/badge/Go-1.22+-blue.svg)
![Build](https://img.shields.io/badge/build-passing-brightgreen)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Contributions](https://img.shields.io/badge/contributions-welcome-orange.svg)

A command-line application written in Go that simulates a **financial transaction ledger system**, demonstrating backend engineering fundamentals such as data structures, domain modeling, and abstraction.

---

## Table of Contents
- [Problem Statement](#problem-statement)
- [Deliverable Definition](#deliverable-definition)
- [Concepts Used](#concepts-used)
- [Project Requirements](#project-requirements)
- [Core Features](#core-features)
- [Application Structure](#application-structure)
- [Development Plan](#development-plan)
- [Specifications](#specifications)
- [User Stories](#user-stories)
- [Getting Started](#getting-started)
- [Demo](#demo)
- [Portfolio Value](#portfolio-value)
- [Future Improvements](#future-improvements)
- [Contributing](#contributing)
- [License](#license)

---

## Problem Statement

Financial accounting requires a systematic way to record, organize, and review monetary transactions. Without a structured system, tracking debits, credits, account balances, and journal entries becomes error-prone and difficult to audit. This project addresses that by implementing a double-entry bookkeeping ledger as a CLI application.

---

## Deliverable Definition

A Go CLI application that implements core double-entry bookkeeping concepts: a Chart of Accounts, a General Journal, and a General Ledger. Users can create accounts, record and update journal entries, and view financial data through a set of commands backed by a SQLite database.

---

## Concepts Used
### Data Structures
- Slices for storing transactions
- Maps for fast account lookup
- Structs for domain modeling

### Types & Abstraction
- Custom types (e.g., Id, Currency)
- Encapsulation via unexported fields
- Struct embedding for reuse
- Interfaces for decoupling behavior

### Engineering Practices
- Package-based structure
- Separation of concerns
- CLI-driven architecture

---

## Project Requirements

### Functional Requirements

- Create and manage accounts organized by account type (Assets, Liabilities, Equities, Revenues, Expenses)
- Record double-entry journal entries with debit/credit accounts, amount, date, and explanation
- View the General Journal with filtering by date range and posting status
- View the Chart of Accounts and account details
- View the General Ledger with per-account transaction history
- Post journal entries to the ledger
- Unpost journal entries from the ledger
- Update or delete journal entries that are not posted

### Non-Functional Requirements
- Clear and maintainable code structure
- Domain-driven naming
- No direct mutation of critical state
- Error handling for invalid operations

---

## Core Features

- Chart of Accounts with 5 default account types (Assets, Liabilities, Equities, Revenues, Expenses), each with a numeric ref prefix (1xxx–9xxx)
- General Journal with double-entry recording, date filtering, and posted/not-posted views
- General Ledger with per-account T-account style display (Date | Debit | Credit | Balance)
- SQLite persistence via modernc.org/sqlite — no CGO required
- Currency stored as integer cents with banker's rounding on float conversion
- UUID-based identity for all domain entities

---

## Application Structure

```bash
src/
├── main.go
└── internal/
    ├── cli/          # Command parsing and dispatch
    ├── database/     # SQLite persistence layer
    ├── ledger/       # Root aggregate, orchestrates chart + journal
    ├── chart/        # Chart of Accounts
    ├── accountType/  # Account type (Assets, Liabilities, etc.)
    ├── account/      # Account and AccountEntry
    ├── journal/      # General Journal
    ├── entry/        # Journal Entry
    ├── currency/     # Currency type (int64 cents) with banker's rounding
    ├── id/           # UUID wrapper
    └── ui/           # CLI input helpers
```

### Structure Goals
- Separate domain logic from CLI handling
- Enable future extension (e.g., database, API)
- Keep business logic testable

---

## Development Plan

---

## Specifications

### Ledger
Root aggregate. Holds one Chart and one Journal. Provides all view and create operations used by the CLI.

### Chart (Chart of Accounts)
Holds a map of AccountTypes keyed by UUID. Accounts are looked up by ref (e.g. 1001) or by UUID.

### Account Type
Groups accounts by category. Has a refPrefix (1–9) that determines the thousands digit of account refs. Default types: Assets (1), Liabilities (2), Equities (3), Revenues (4), Expenses (9). Max 999 accounts per type.

### Account
Represents a ledger account. Has a UUID, a numeric ref (e.g. 1001), name (max 30 chars), description, and a collection of AccountEntries. Balance is calculated from entries.

### Account Entry
A single posting to an account. Stores date, amount (cents), explanation, side (debit=false / credit=true), and running balance.

### Journal (General Journal)
Holds a map of Entries keyed by UUID. Entries are returned sorted by date.

### Entry (A journal entry)
Has a UUID, date, debit account, credit account, amount (Currency), explanation, and a posted flag. Can only be updated or deleted when not posted.

---

## User Stories
### Views
#### General Journal (Journal)
- [x] As a user, I want to view all journal entries
- [x] As a user, I want to view all journal entries between two dates
- [x] As a user, I want to view all posted entries
- [x] As a user, I want to view all not-posted entries
- [x] As a user, I want to view an entry details

#### Chart of Accounts (Chart)
- [x] As a user, I want to view the "Chart of Accounts"
- [x] As a user, I want to view the details for an account type
- [x] As a user, I want to view the details for an account

#### Ledger
- [ ] As a user, I want to view the "General Ledger"
- [ ] As a user, I want to view the ledger of an account type
- [ ] As a user, I want to view the ledger of an account
- [ ] As a user, I want to view the "Trial Balance"

### Create
- [x] As a user, I want to create a new entry in the journal
- [x] As a user, I want to create a new account

### Update
- [ ] As a user, I want to update certain details for an entry that is not-posted: date, debit account, credit account, amount, and explanation.
- [ ] As a user, I want to change the status of an entry to posted
- [ ] As a user, I want to change the status of an entry to not-posted
- [ ] As a user, I want to change certain details of an account: name, description

### Delete
- [ ] As a user, I want to delete an entry that is not-posted

---

## Getting Started

```bash
git clone https://github.com/yourusername/go-ledger-cli.git
cd go-ledger-cli/src
go run main.go help
```

---

## Demo


---

## Portfolio Value
This project demonstrates:
- Backend engineering fundamentals
- Financial domain modeling
- Code organization and abstraction
- Readiness for:
  - FinTech
  - Banking Systems
  - Public sector applications

---

## Possible Future Improvements
- REST API interface
- Authentication layer
- Concurrency-safe processing
- Integration with external systems

---

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

1. Fork the repository.
2. Create a new branch for your changes.
3. Make your changes and commit them.
4. Push your changes to your forked repository.
5. Open a pull request.

---

## License

This project is licensed under the [MIT License](./LICENSE.md)
